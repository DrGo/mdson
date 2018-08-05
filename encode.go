package mdson

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
)

const lineBreak = "\n"

//Marshal encodes the contents of a pointer to an interface{} v into io.Writer
// Marshal traverses the value v recursively.
// Struct values encode as JSON objects.
//FIXME: change to producing msdon
func Marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := NewEncoder(&b)
	// enc.Indent(prefix, indent)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
	// return json.MarshalIndent(v, "", " ")
}

type fieldInfo struct {
	mdsonTagValues
	sf *reflect.StructField
}

// The encoding of each struct field can be customized by the format string
// stored under the "mdson" key in the struct field's tag.
// The format string gives the name of the field, possibly followed by a
// comma-separated list of options. The name may be empty in order to
// specify options without overriding the default field name.
//
// The "omitempty" option specifies that the field should be omitted
// from the encoding if the field has an empty value, defined as
// false, 0, a nil pointer, a nil interface value, and any empty array,
// slice, map, or string.
//
// As a special case, if the field tag is "-", the field is always omitted.
// Note that a field with name "-" can still be generated using the tag "-,".
//
// Examples of struct field tags and their meanings:
//
//   // Field appears in MDSon as key "myName".
//   Field int `mdson:"myName"`
//
//   // Field appears in MDSon as key "myName" and
//   // the field is omitted from the object if its value is empty,
//   // as defined above.
//   Field int `mdson:"myName,omitempty"`
//
//   // Field appears in MDSon as key "Field" (the default), but
//   // the field is skipped if empty.
//   // Note the leading comma.
//   Field int `mdson:",omitempty"`
//
//   // Field is ignored by this package.
//   Field int `mdson:"-"`
//
//   // Field appears in MDSon as key "-".
//   Field int `mdson:"-,"`

func newFieldInfo(sf reflect.StructField) *fieldInfo {
	tv := getMDsonTagValues(sf)
	return &fieldInfo{mdsonTagValues: tv, sf: &sf}
}

func (fi *fieldInfo) name() string {
	// Precedence for the field  name is:
	// 0. tag name
	// 1. filed name
	if fi.tagName != "" {
		return fi.tagName
	}
	return fi.sf.Name
}

func (fi *fieldInfo) typ() reflect.Type {
	return fi.sf.Type
}

//Encoder decodes
type Encoder struct {
	Debug int
	depth int
	*bufio.Writer
}

// NewEncoder initializes a new Encoder
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		Writer: bufio.NewWriter(w),
		Debug:  debug,
	}
}

// Encode writes the XML encoding of v to the stream.
// Encode calls Flush before returning.
func (enc *Encoder) Encode(v interface{}) error {
	st := reflect.ValueOf(v)
	//find a struct if v is an interface or pointer
	// This can turn into an infinite loop given a cyclic chain,
	for st.Kind() == reflect.Interface || st.Kind() == reflect.Ptr {
		if st.IsNil() {
			return nil
		}
		st = st.Elem()
	}
	blockName := getStructName(st, 0)
	err := enc.encodeStruct(st, blockName)
	if err != nil {
		return err
	}
	return enc.Flush()
}

func (enc *Encoder) encodeStruct(st reflect.Value, blockName string) error {
	var err error
	if !st.IsValid() {
		return nil
	}
	if enc.depth == 0 && st.Kind() != reflect.Struct { //root must be a struct
		return fmt.Errorf("root value is not a struct")
	}
	stType := st.Type()
	// Create an MDSon block element.
	enc.depth++ //got a root block
	enc.log("encodeStruct:", blockName, stType, enc.depth)
	//write the block name
	enc.WriteString(lineBreak + strings.Repeat("#", enc.depth) + " " + blockName + lineBreak)
	// loop:
	for i := 0; i < stType.NumField(); i++ {
		fld := st.Field(i)
		fldInfo := newFieldInfo((stType.Field(i)))
		k := ToMarshableKind(fld.Type().Kind())
		if k == MarshablePtr { //iface or ptr follow-it
			fld = Dereference(fld, false /*do not create element if nil*/)
			if !fld.IsValid() { //nil iface or ptr
				continue
			}
		}
		k = ToMarshableKind(fld.Type().Kind()) //refresh fld kind in case fld was updated above
		//TODO: handle omitEmpty
		// if finfo.flags&fOmitEmpty != 0 && isEmptyValue(fv) {
		// 	continue
		// }
		enc.log(blockName, "fld :", fldInfo.name(), fldInfo.typ, enc.depth)
		switch k {
		case MarshablePtr:
			panic("decodeStruct(): pointer not allowed at this point")
		case MarshableSlice:
			enc.encodeSlice(fld, fldInfo)
		case MarshableStruct:
			err = enc.encodeStruct(fld, fldInfo.name())
		case UnMarshable:
			continue
		case MarshablePrimitive:
			s := encodeSimple(fld)
			enc.WriteString(fldInfo.name() + " :" + s + lineBreak)
		default:
			panic("internal error")
		}
		if err != nil {
			return err
		}
	}
	fmt.Println("about to decrement depth from", enc.depth)
	enc.depth--
	return nil //p.cachedWriteError()
}

func (enc *Encoder) log(a ...interface{}) {
	if enc.Debug >= DebugAll {
		fmt.Printf(strings.Repeat("  ", enc.depth))
		fmt.Println(a...)
	}
}

func (enc *Encoder) warn(a ...interface{}) {
	if enc.Debug >= DebugWarning {
		fmt.Printf(strings.Repeat("  ", enc.depth) + "warning: ")
		fmt.Println(a...)
	}
}

// Slices and arrays of structs iterate over the elements. They do not have an enclosing tag.
func (enc *Encoder) encodeSlice(sl reflect.Value, fldInfo *fieldInfo) error {
	if sl.Kind() != reflect.Slice && sl.Kind() != reflect.Array {
		panic("encodeSlice: received non-slice/array as an arg")
	}
	el := sl.Type().Elem()
	//scenario 1: a slice/array of structs or ptrs to a struct, write each element as a block
	if el.Kind() == reflect.Struct || (el.Kind() == reflect.Ptr && el.Elem().Kind() == reflect.Struct) {
		enc.depth++
		//FIXME: what name to write
		enc.WriteString(lineBreak + strings.Repeat("#", enc.depth) + " " + fldInfo.name() + " List" + lineBreak)
		for i := 0; i < sl.Len(); i++ {
			st := Dereference(sl.Index(i), false /*do not create element if nil*/)
			if !st.IsValid() { //nil iface or ptr
				continue
			}
			blockName := getStructName(st, i+1)
			if err := enc.encodeStruct(st, blockName); err != nil {
				return err
			}
		}
		enc.depth--
		return nil
	}
	if el.Kind() == reflect.Uint8 { //byte array/slice: marshal as a string
		var bytes []byte
		if sl.Type().Kind() == reflect.Array {
			if sl.CanAddr() {
				bytes = sl.Slice(0, sl.Len()).Bytes()
			} else {
				bytes = make([]byte, sl.Len())
				reflect.Copy(reflect.ValueOf(bytes), sl)
			}
		}
		//must be a byte slice
		bytes = sl.Bytes()
		enc.WriteString(fldInfo.name() + " :" + string(bytes) + lineBreak)
		return nil
	}

	//scenario 3: a slice/array of primitive type or pointers to a primitive type, write as an MDSon list
	enc.log("primitve type:", primitiveType(el.Kind()))
	if primitiveType(el.Kind()) || (el.Kind() == reflect.Ptr && primitiveType(el.Elem().Kind())) {
		enc.WriteString(fldInfo.name() + " :" + lineBreak)
		for i := 0; i < sl.Len(); i++ {
			li := Dereference(sl.Index(i), false /*do not create element if nil*/)
			if !li.IsValid() { //nil iface or ptr
				continue
			}
			s := encodeSimple(sl.Index(i))
			enc.WriteString("- " + s + lineBreak)
		}
		return nil
	}
	return nil
}

//MarshableKind describes grouped reflect.type values
type MarshableKind int

const (
	// UnMarshable cannot be marshalled as MDson
	UnMarshable MarshableKind = iota
	// MarshableStruct marshalled as MDson block
	MarshableStruct
	//MarshablePrimitive marshalled as MDson primitive value
	MarshablePrimitive
	//MarshableSlice marshalled as MDson list
	MarshableSlice
	//MarshablePtr marshalled as what it points to
	MarshablePtr
)

// ToMarshableKind converts from reflect.Kind to MarshableKind
func ToMarshableKind(k reflect.Kind) MarshableKind {
	switch {
	case k == reflect.Struct:
		return MarshableStruct
	case primitiveType(k):
		return MarshablePrimitive
	case k == reflect.Array || k == reflect.Slice:
		return MarshableSlice
	case k == reflect.Interface || k == reflect.Ptr:
		return MarshablePtr
	default:
		return UnMarshable //unmarshable value
	}
}
