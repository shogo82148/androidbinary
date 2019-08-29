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

// MarshalXMLAttr implements xml.MarshalerAttr.
func (v Bool) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if v.value == "" {
		// return the zero value of bool
		return xml.Attr{
			Name:  name,
			Value: "false",
		}, nil
	}
	return xml.Attr{
		Name:  name,
		Value: v.value,
	}, nil
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

// MustBool is same as Bool, but it panics if it fails to parse the value.
func (v Bool) MustBool() bool {
	ret, err := v.Bool()
	if err != nil {
		panic(err)
	}
	return ret
}

// Int32 is an integer value in XML file.
// It may be an immediate value or a reference.
type Int32 struct {
	value  string
	table  *TableFile
	config *ResTableConfig
}

// WithTableFile ties TableFile to the Bool.
func (v Int32) WithTableFile(table *TableFile) Int32 {
	return Int32{
		value:  v.value,
		table:  table,
		config: v.config,
	}
}

// WithResTableConfig ties ResTableConfig to the Bool.
func (v Int32) WithResTableConfig(config *ResTableConfig) Bool {
	return Bool{
		value:  v.value,
		table:  v.table,
		config: config,
	}
}

func (v *Int32) inject(table *TableFile, config *ResTableConfig) {
	v.table = table
	v.config = config
}

// SetInt32 sets an integer value.
func (v *Int32) SetInt32(value int32) {
	v.value = strconv.FormatInt(int64(value), 10)
}

// SetResID sets a boolean value with the resource id.
func (v *Int32) SetResID(resID ResID) {
	v.value = resID.String()
}

// UnmarshalXMLAttr implements xml.UnmarshalerAttr.
func (v *Int32) UnmarshalXMLAttr(attr xml.Attr) error {
	v.value = attr.Value
	return nil
}

// MarshalXMLAttr implements xml.MarshalerAttr.
func (v Int32) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if v.value == "" {
		// return the zero value of int32
		return xml.Attr{
			Name:  name,
			Value: "0",
		}, nil
	}
	return xml.Attr{
		Name:  name,
		Value: v.value,
	}, nil
}

// Int32 returns the integer value.
// It resolves the reference if needed.
func (v Int32) Int32() (int32, error) {
	if v.value == "" {
		return 0, nil
	}
	if !IsResID(v.value) {
		v, err := strconv.ParseInt(v.value, 10, 32)
		return int32(v), err
	}
	id, err := ParseResID(v.value)
	if err != nil {
		return 0, err
	}
	value, err := v.table.GetResource(id, v.config)
	if err != nil {
		return 0, err
	}
	ret, ok := value.(uint32)
	if !ok {
		return 0, fmt.Errorf("invalid type: %T", value)
	}
	return int32(ret), nil
}

// MustInt32 is same as Int32, but it panics if it fails to parse the value.
func (v Int32) MustInt32() int32 {
	ret, err := v.Int32()
	if err != nil {
		panic(err)
	}
	return ret
}

// String is a boolean value in XML file.
// It may be an immediate value or a reference.
type String struct {
	value  string
	table  *TableFile
	config *ResTableConfig
}

// WithTableFile ties TableFile to the Bool.
func (v String) WithTableFile(table *TableFile) String {
	return String{
		value:  v.value,
		table:  table,
		config: v.config,
	}
}

// WithResTableConfig ties ResTableConfig to the Bool.
func (v String) WithResTableConfig(config *ResTableConfig) String {
	return String{
		value:  v.value,
		table:  v.table,
		config: config,
	}
}

func (v *String) inject(table *TableFile, config *ResTableConfig) {
	v.table = table
	v.config = config
}

// SetString sets a string value.
func (v *String) SetString(value string) {
	v.value = value
}

// SetResID sets a boolean value with the resource id.
func (v *String) SetResID(resID ResID) {
	v.value = resID.String()
}

// UnmarshalXMLAttr implements xml.UnmarshalerAttr.
func (v *String) UnmarshalXMLAttr(attr xml.Attr) error {
	v.value = attr.Value
	return nil
}

// MarshalXMLAttr implements xml.MarshalerAttr.
func (v String) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  name,
		Value: v.value,
	}, nil
}

// String returns the string value.
// It resolves the reference if needed.
func (v String) String() (string, error) {
	if !IsResID(v.value) {
		return v.value, nil
	}
	id, err := ParseResID(v.value)
	if err != nil {
		return "", err
	}
	value, err := v.table.GetResource(id, v.config)
	if err != nil {
		return "", err
	}
	ret, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("invalid type: %T", value)
	}
	return ret, nil
}

// MustString is same as String, but it panics if it fails to parse the value.
func (v String) MustString() string {
	ret, err := v.String()
	if err != nil {
		panic(err)
	}
	return ret
}
