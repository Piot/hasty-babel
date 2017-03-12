package serializers

import (
	"encoding/binary"
	"fmt"
	"log"
	"reflect"

	"github.com/piot/hasty-babel/definition"
	"github.com/piot/hasty-babel/reflector"
)

func argumentToOctet(arg definition.Argument, data []byte, v *reflect.Value) (int, error) {
	switch arg.ArgumentType {
	case "u32":
		binary.BigEndian.PutUint32(data, uint32(v.Uint()))
		return 4, nil
	case "u16":
		binary.BigEndian.PutUint16(data, uint16(v.Uint()))
		return 2, nil
	case "u8":
		data[0] = uint8(v.Uint())
		return 1, nil
	case "bool":
		a := v.Bool()
		v := byte(0)
		if a {
			v = 1
		}
		data[0] = v
		return 1, nil
	case "string":
		n := copyStringToOctets(v.String(), data)
		return n, nil
	case "asciistring":
		n := copyStringToOctets(v.String(), data)
		return n, nil
	}

	return 0, fmt.Errorf("Illegal argument type '%s'", arg.ArgumentType)
}

func copyStringToOctets(s string, data []byte) int {
	n := len(s)
	data[0] = byte(n)
	stringOctets := []byte(s)
	copy(data[1:], stringOctets)
	return n + 1
}

func ObjectToOctets(in definition.ProtocolDefinition, foundCommand definition.Command, source interface{}) ([]byte, error) {
	s, reflectErr := reflector.NewStructure(source)
	if reflectErr != nil {
		log.Fatalf("Reflect:%v", reflectErr)
		return nil, reflectErr
	}
	var tempBuf = make([]byte, 256)
	pos := 0
	tempBuf[pos] = foundCommand.ID
	pos++
	for _, a := range foundCommand.Arguments {
		target, fieldErr := s.FindField(a.Name)
		if fieldErr != nil {
			return nil, fmt.Errorf("Field err:%s", fieldErr)
		}
		octetCount, parseErr := argumentToOctet(a, tempBuf[pos:], target)
		if parseErr != nil {
			return nil, fmt.Errorf("Argument to Octet:%s", parseErr)
		}
		pos += octetCount
	}
	return tempBuf[:pos], nil
}
