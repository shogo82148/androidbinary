package androidbinary

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"
)

type injector interface {
	inject(table *TableFile, config *ResTableConfig)
}

var injectorType = reflect.TypeOf((*injector)(nil)).Elem()

func inject(val reflect.Value, table *TableFile, config *ResTableConfig) {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}
	if val.CanInterface() && val.Type().Implements(injectorType) {
		val.Interface().(injector).inject(table, config)
		return
	}
	if val.CanAddr() {
		pv := val.Addr()
		if pv.CanInterface() && pv.Type().Implements(injectorType) {
			pv.Interface().(injector).inject(table, config)
			return
		}
	}

	switch val.Kind() {
	default:
		// ignore other types
		return
	case reflect.Slice, reflect.Array:
		l := val.Len()
		for i := 0; i < l; i++ {
			inject(val.Index(i), table, config)
		}
		return
	case reflect.Struct:
		l := val.NumField()
		for i := 0; i < l; i++ {
			inject(val.Field(i), table, config)
		}
	}
}

// Bool is a boolean value in XML file.
// It may be an immediate value or a reference.
type Bool struct {
	value  string
	table  *TableFile
	config *ResTableConfig
}

// WithTableFile ties TableFile to the Bool.
func (v Bool) WithTableFile(table *TableFile) Bool {
	return Bool{
		value:  v.value,
		table:  table,
		config: v.config,
	}
}

// WithResTableConfig ties ResTableConfig to the Bool.
func (v Bool) WithResTableConfig(config *ResTableConfig) Bool {
	return Bool{
		value:  v.value,
		table:  v.table,
		config: config,
	}
}

func (v *Bool) inject(table *TableFile, config *ResTableConfig) {
	v.table = table
	v.config = config
}

// SetBool sets a boolean value.
func (v *Bool) SetBool(value bool) {
	v.value = strconv.FormatBool(value)
}

// SetResID sets a boolean value with the resource id.
func (v *Bool) SetResID(resID ResID) {
	v.value = resID.String()
}

// UnmarshalXMLAttr implements xml.UnmarshalerAttr.
func (v *Bool) UnmarshalXMLAttr(attr xml.Attr) error {
	v.value = attr.Value
	return nil
}

// Bool returns the boolean value.
// It resolves the reference if needed.
func (v Bool) Bool() (bool, error) {
	if v.value == "" {
		return false, nil
	}
	if !IsResID(v.value) {
		return strconv.ParseBool(v.value)
	}
	id, err := ParseResID(v.value)
	if err != nil {
		return false, err
	}
	value, err := v.table.GetResource(id, v.config)
	if err != nil {
		return false, err
	}
	ret, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("invalid type: %T", value)
	}
	return ret, nil
}
