package reflector

import (
	"fmt"
	"reflect"
	"strings"
)

// Structure to reflect
type Structure struct {
	value reflect.Value
}

// NewStructure : Creates a new structure
func NewStructure(target interface{}) (Structure, error) {
	ps := reflect.ValueOf(target)
	s := ps.Elem()
	if s.Kind() != reflect.Struct {
		return Structure{}, fmt.Errorf("Can only set to structs")
	}

	st := Structure{value: s}
	return st, nil
}

// FindField : Finds a field by name
func (in Structure) FindField(name string) (*reflect.Value, error) {
	uppercased := strings.Title(name)
	f := in.value.FieldByName(uppercased)
	if !f.IsValid() {
		return nil, fmt.Errorf("Field '%s' is not valid", name)
	}
	if !f.CanSet() {
		return nil, fmt.Errorf("Field '%s' can not be set", name)
	}

	return &f, nil
}
