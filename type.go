package androidbinary

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

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
