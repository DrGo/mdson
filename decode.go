package mdson

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

//Unmarshal unmarshales the contents of io.Reader into a pointer to an interface{}
func Unmarshal(r io.Reader, out interface{}) (err error) {
	//TODO: define the defer recover func below
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(rv)}
	}
	//we have a valid pointer, get what it points to
	rv = rv.Elem()
	if !rv.IsValid() || rv.Kind() != reflect.Struct { //root must be struct
		return &InvalidUnmarshalError{reflect.TypeOf(rv)}
	}
	root, err := Parse(r)
	if err != nil {
		return err
	}
	blk, ok := root.(*ttBlock)
	if !ok {
		return fmt.Errorf("parser returned unexpected type: root is not *ttBlock")
	}
	dec := NewDecoder()
	dec.log("*********Decoding started")
	err = dec.decodeBlock(blk, rv)
	dec.depth = 0
	dec.log("***********decoding ended. Err= ", err)
	return err
}

//Decoder decodes
type Decoder struct {
	depth int
	Debug int
}

// NewDecoder initializes a new decoder
func NewDecoder() *Decoder {
	return &Decoder{
		Debug: debug,
	}
}

func (dec *Decoder) log(a ...interface{}) {
	if dec.Debug >= DebugAll {
		fmt.Printf(strings.Repeat("  ", dec.depth))
		fmt.Println(a...)
	}
}

func (dec *Decoder) warn(a ...interface{}) {
	if dec.Debug >= DebugWarning {
		fmt.Printf(strings.Repeat("  ", dec.depth) + "warning: ")
		fmt.Println(a...)
	}
}

//dereference for a pointer returns the value it points to (if nil it creates a new one)
// for a non-pointer it returns the passed value as is
//FIXME: replace with Dereference
func dereference(rv reflect.Value) reflect.Value {
	//TODO: repeat in a loop until you get a valid value
	if rv.Kind() != reflect.Ptr {
		//		fmt.Println("not a pointer")
		return rv
	}
	// if rv.Elem().Kind() != reflect.Ptr && rv.Elem().CanSet() {
	// 	fmt.Println("not pointing to a  pointer")
	// 	return rv.Elem()
	// }
	if rv.IsNil() {
		//		fmt.Println("nil pointer")
		rv.Set(reflect.New(rv.Type().Elem()))
	}
	return rv.Elem()
}

func (dec *Decoder) decodeBlock(block *ttBlock, rv reflect.Value) error {
	//	dec.log(block.Name(), "rv:", rv.String())
	rv = dereference(rv)
	//	dec.log(block.Name(), "rv:", rv.String())
	if err := isValidBlock(block, rv); err != nil {
		return err
	}
	dec.log("block:", block.Name(), ",struct type:", rv.Type().String())
	dec.depth++
	//set the ID field if one exists in the struct
	id := getSettableField(rv, "id")                      //find a variable with name ID
	if id.IsValid() && id.Type().String() == "mdson.ID" { //and type mdson.ID
		id.SetString(block.Name())
		dec.log("found ID field. Set to", block.Name())
	}
	for _, n := range block.children {
		/*DEBUG*/ dec.log("element:", n.Kind(), n.Name())
		switch n := n.(type) {
		case *ttBlock: //either array of blocks or a scalar block
			fld := getSettableField(rv, n.Name()) //find a variable with this array name
			if !fld.IsValid() {
				dec.warn(n.Name(), "no suitable field in corresponding struct")
				return nil
			}
			if n.isArray() { //array of blocks; fld is []struct or []*struct
				dec.log("array:", n.Name(), fld.Type().String())
				sElem := fld.Type().Elem() //get the type of the slice element
				if sElem.Kind() == reflect.Ptr {
					sElem = sElem.Elem() //get the struct that it points to
				}
				if sElem.Kind() != reflect.Struct { //must be a slice of struct
					return ESyntaxError{-1, block.Name() + "." + n.Name() + ": corresponding slice is not of structs"}
				}
				//FIXME: replace with fixed size slice to avoid append below
				fld.SetLen(0) //empty the slice
				// newv := reflect.MakeSlice(v.Type(), v.Len(), newcap)
				dec.depth++
				for _, c := range n.children {
					b, ok := c.(*ttBlock)
					if !ok {
						return ESyntaxError{-1, "A blocks array can only include blocks"}
					}
					// dec.log("fld.Type().Elem()", fld.Type().Elem().String())
					BlockPtr := reflect.New(sElem) //allocate struct of the slice's type
					// dec.log("struct type", BlockPtr.String())
					if err := dec.decodeBlock(b, BlockPtr.Elem()); err != nil {
						return err
					}
					if fld.Type().Elem().Kind() == reflect.Ptr {
						fld.Set(reflect.Append(fld, BlockPtr))
					} else {
						fld.Set(reflect.Append(fld, BlockPtr.Elem()))
					}
				}
				dec.depth--
				return nil
			}
			//a struct in a non-array block
			if err := dec.decodeBlock(n, fld); err != nil {
				return err
			}
		case *ttKVPair:
			if err := dec.decodeKVPair(n, rv); err != nil {
				return err
			}
		case *ttLiteralString:
			if err := dec.decodeKVPair(n, rv); err != nil {
				return err
			}
		case *ttList:
			if err := dec.decodeList(n, rv); err != nil {
				return err
			}
		}
	}
	dec.depth--
	return nil
}

func (dec *Decoder) decodeKVPair(n Node, rv reflect.Value) error {
	//FIXME: ensure rv is a valid struct?!
	fld := getSettableField(rv, n.Name()) //find a variable with this array name
	if !fld.IsValid() {
		dec.warn(n.Name(), "no suitable field in corresponding struct")
		return nil
	}
	//if this is an ID field do not set it b/c it was already filled
	if fld.Type().String() == "mdson.ID" {
		return nil
	}
	value := ""
	switch n := n.(type) {
	case *ttKVPair:
		value = n.value
	case *ttLiteralString:
		value = n.value
	default:
		return &InvalidUnmarshalError{reflect.TypeOf(n)}
	}
	return setValue(fld, value)
}

func (dec *Decoder) decodeList(list *ttList, rv reflect.Value) error {
	fld := getSettableField(rv, list.Name()) //find a field with this array name
	if !fld.IsValid() || fld.Kind() != reflect.Slice {
		dec.warn(list.Name(), "no suitable field in corresponding struct")
		return nil
	}
	if !primitiveType(fld.Type().Elem().Kind()) {
		dec.warn(list.Name(), "list with unsupported element type")
		return nil
	}
	//
	//TODO: add support for array of pointers eg []*string
	et := fld.Type().Elem()
	dec.log("list:", list.Name(), fld.Type().String(), "el-type", et.String())
	for _, item := range list.items {
		value := strings.TrimSpace(item.key)
		if value == "" {
			continue
		}
		newItem := reflect.New(et) //allocate an element of the slice's type
		if err := setValue(newItem.Elem(), value); err != nil {
			return err
		}
		fld.Set(reflect.Append(fld, newItem.Elem()))
	}
	return nil
}
