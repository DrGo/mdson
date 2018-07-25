package mds

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func parseAttribs(s string) (TokenBlock, error) {
	entries := strings.Split(s, ",")
	attribs := make(TokenBlock, 10)
	for _, entry := range entries {
		parts := strings.SplitN(entry, ":", 2) //split on the first colon
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid value for '%s' likely missing ':'", entry)
		}
		attribs[trimLower(parts[0])] = trimLower(parts[1])
	}
	return attribs, nil
}

//decode fills a struct (passed an an interface) from a map of key-value pairs
func decode(in interface{}, block TokenBlock) error {
	if len(block) == 0 {
		return nil
	}
	//FIXME: validate in is not null reflect.IsNull
	st := reflect.ValueOf(in).Elem()
	//build a map of struc fields keyed by either tag name or field name
	//TODO: test that tagging struct works
	flds := make(map[string]reflect.Value, 10)
	for i := 0; i < st.NumField(); i++ {
		fldInfo := st.Type().Field(i)
		fldName := trimLower(fldInfo.Tag.Get("mds"))
		if fldName == "" { //no tag
			fldName = trimLower(fldInfo.Name)
		}
		flds[fldName] = st.Field(i)
		// /*DEBUG*/ fmt.Println("inside decode: fldName:", fldName)
	}
	// /*DEBUG*/ fmt.Println("inside decode: attribs", attribs)
	for name, value := range block {
		fld := flds[name]
		if !fld.IsValid() || !fld.CanSet() {
			continue
		}
		if err := setValue(fld, value); err != nil {
			return fmt.Errorf("could not set value: %s", err)
		}
		// /*DEBUG*/ fmt.Println("inside decode: new fld", fld)
	}
	return nil
}

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
