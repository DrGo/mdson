package mdson

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func setValue(rv reflect.Value, value string) error {
	//FIXME: check that fld is settable
	switch kind := rv.Kind(); kind {
	case reflect.String:
		rv.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		rv.SetInt(int64(i))
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		rv.SetFloat(float64(f))
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		rv.SetBool(b)
	default:
		return fmt.Errorf("unsupported field type %s", kind)
	}
	return nil
}

func isValidBlock(block *ttBlock, rv reflect.Value) error {
	if block.isArray() { //if list of block, assert rv is a slice of struct
		if rv.Kind() != reflect.Slice {
			return ESyntaxError{-1, block.Name() + ": A blocks array has no corresponding slice"}
		}
		eT := rv.Type().Elem()           //get slice element type
		if eT.Kind() != reflect.Struct { //must be a slice of struct
			return ESyntaxError{-1, block.Name() + ": corresponding slice is not of structs"}
		}
	}

	if !block.isArray() && rv.Kind() != reflect.Struct { //if scalar block, assert rv is struct
		return ESyntaxError{-1, block.Name() + ": has no corresponding struct"}
	}
	return nil
}

func isArrayValue(rv reflect.Value) bool {
	return rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array
}

func getFieldByNameOrTag(st reflect.Value, fieldName string) reflect.Value {
	fld := st.FieldByNameFunc(func(name string) bool {
		found := strings.ToLower(name) == fieldName
		//TODO: handle names with spaces etc and extract into resusable func
		/*DEBUG*/ // fmt.Println("inside FieldByNameFunc:", name, n.name, found)
		return found
	})
	return fld
}

func getSettableField(st reflect.Value, fieldName string) reflect.Value {
	//TODO: handle tag search
	fld := getFieldByNameOrTag(st, fieldName)
	// fmt.Println(fld)
	if !fld.IsValid() || !fld.CanSet() {
		return reflect.Value{} //empty value
	}
	return fld
}

func primitiveType(k reflect.Kind) bool {
	// fmt.Println("primitiveType:", k)
	return k == reflect.String || (k >= reflect.Bool && k <= reflect.Float64)
}

// Dereference follows an interface or pointer until it finds a non-interface non-pointer element and returns it
// if it reaches a nil interface it return an invalid reflect.value unless createElement is true, in that case
// it allocates an element and returns that element.
func Dereference(rv reflect.Value, createElement bool) reflect.Value {
	for rv.Kind() == reflect.Interface || rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			if !createElement {
				return reflect.Value{}
			}
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	return rv
}

func encodeSimple(val reflect.Value) string {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(val.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'g', -1, val.Type().Bits())
	case reflect.String:
		return val.String()
	case reflect.Bool:
		return strconv.FormatBool(val.Bool())
	}
	return "unsupported type: " + val.Type().Name()
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

type mdsonTagValues struct {
	name string
	omit bool
}

func getMDsonTagValues(sf reflect.StructField) mdsonTagValues {
	fi := mdsonTagValues{}
	val, found := sf.Tag.Lookup("mdson")
	if !found {
		return fi //name=""; omit= false
	}
	if trimLower(val) == "-" {
		return mdsonTagValues{name: "-"} //name=""; omit= false
	}
	opts := strings.Split(val, ",")
	if len(opts) > 0 { //first options is always name
		fi.name = opts[0]
	}
	if len(opts) > 1 { //second options is omitempty
		fi.omit = trimLower(opts[1]) == "omitempty"
	}
	if !isValidTag(fi.name) {
		fi.name = ""
	}
	return fi
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		default:
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				return false
			}
		}
	}
	return true
}

func GetValidVarName(s string) string {
	ns := []rune{}
	for _, c := range s {
		if unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' {
			ns = append(ns, c)
		}
	}
	return string(ns)
}

func getStructName(st reflect.Value, index int) string {
	id := getFieldByNameOrTag(st, "id")                                        //find a variable with name ID
	if id.IsValid() && id.Type().String() == "mdson.ID" && id.String() != "" { //and type mdson.ID
		return strings.Title(id.String())
	}
	//No ID Field, use struct type name
	stName := st.Type().Name()
	if index > 0 {
		stName = stName + strconv.Itoa(index)
	}
	return stName
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

//TODO: centralize reflect.value type checking
// func mustBe(expected Kind) {
// 	if f.kind() != expected {
// 		panic(&ValueError{methodName(), f.kind()})
// 	}
// }

// // defer func(){ //catch all unhandled errors
// 	if r := recover(); r != nil {
// 		if _, ok := r.(runtime.Error); ok {
// 			panic(r)
// 		}
// 		err = r.(error)
// 	}
// }()
