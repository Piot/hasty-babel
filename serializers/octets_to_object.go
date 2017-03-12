package serializers

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/piot/hasty-babel/definition"
	"github.com/piot/hasty-babel/reflector"
)

func parseTypeToValue(arg definition.Argument, data []byte, target *reflect.Value) (string, byte, error) {
	switch arg.ArgumentType {
	case "u32":
		v := binary.BigEndian.Uint32(data)
		if target.Kind() != reflect.Uint32 {
			return "", 0, fmt.Errorf("Wrong kind: %v expected uint32 but was %s", arg, target)
		}
		target.SetUint(uint64(v))
		return "", 4, nil
	case "u16":
		v := binary.BigEndian.Uint16(data)
		if target.Kind() != reflect.Uint16 {
			return "", 0, fmt.Errorf("Wrong kind: %v expected uint16 but was %s", arg, target)
		}
		target.SetUint(uint64(v))

		return "", 2, nil

	case "u8":
		v := data[0]
		if target.Kind() != reflect.Uint8 {
			return "", 0, fmt.Errorf("Wrong kind: %v expected uint8 but was %s", arg, target)
		}
		target.SetUint(uint64(v))
		return "", 1, nil

	case "bool":
		v := data[0] != 0x00
		if target.Kind() != reflect.Bool {
			return "", 0, fmt.Errorf("Wrong kind: %v expected bool but was %s", arg, target)
		}
		target.SetBool(v)
		return "", 1, nil
	case "string":
		n := data[0]
		s := data[1 : 1+n]
		str := string(s)
		if target.Kind() != reflect.String {
			return "", 0, fmt.Errorf("Wrong kind: %v expected string but was %s", arg, target)
		}
		target.SetString(str)
		return "", n + 1, nil
	case "asciistring":
		n := data[0]
		s := data[1 : 1+n]
		str := string(s)
		if target.Kind() != reflect.String {
			return "", 0, fmt.Errorf("Wrong kind: %v expected string but was %s", arg, target)
		}
		target.SetString(str)
		return "", n + 1, nil
	}

	return "", 0, fmt.Errorf("Illegal argument type '%s'", arg.ArgumentType)
}

func OctetsToObject(in definition.ProtocolDefinition, data []byte, target interface{}) (int, error) {
	pos := 0
	s, reflectErr := reflector.NewStructure(target)
	if reflectErr != nil {
		return pos, reflectErr
	}
	commandID := data[pos]
	pos++
	foundCommand := in.FindCommand(commandID)
	if foundCommand == nil {
		return pos, fmt.Errorf("Couldn't find command 0x%02X", commandID)
	}
	for _, a := range foundCommand.Arguments {
		target, err := s.FindField(a.Name)
		if err != nil {
			return pos, err
		}
		_, octetCount, parseErr := parseTypeToValue(a, data[pos:], target)
		if parseErr != nil {
			return pos, parseErr
		}
		pos += int(octetCount)
	}
	return pos, nil
}
