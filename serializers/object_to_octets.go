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

func argumentValueToOctet(arg definition.Argument, data []byte, v interface{}) (int, error) {
	switch arg.ArgumentType {
	case "u64":
		binary.BigEndian.PutUint64(data, uint64(v.(int)))
		return 8, nil
	case "u32":
		binary.BigEndian.PutUint32(data, uint32(v.(int)))
		return 4, nil
	case "u16":
		binary.BigEndian.PutUint16(data, uint16(v.(int)))
		return 2, nil
	case "u8":
		data[0] = uint8(uint8(v.(int)))
		return 1, nil
	case "bool":
		a := v.(bool)
		v := byte(0)
		if a {
			v = 1
		}
		data[0] = v
		return 1, nil
	case "string":
		n := copyStringToOctets(v.(string), data)
		return n, nil
	case "asciistring":
		n := copyStringToOctets(v.(string), data)
		return n, nil
	}

	return 0, fmt.Errorf("Illegal argument type '%s'", arg.ArgumentType)
}

func copyStringToOctets(s string, data []byte) int {
	n := len(s)
	lengthEncodedOctets, err := SmallLengthToOctets(uint16(n))
	if err != nil {
		fmt.Printf("Problem:%s", err)
	}
	lengthEncodingSize := len(lengthEncodedOctets)
	pos := 0
	copy(data, lengthEncodedOctets)
	pos += lengthEncodingSize
	stringOctets := []byte(s)
	copy(data[pos:], stringOctets)
	return n + lengthEncodingSize
}

func usingFields(in definition.ProtocolDefinition, foundCommand definition.Command, source interface{}) ([]byte, error) {
	s, reflectErr := reflector.NewStructure(source)
	if reflectErr != nil {
		log.Fatalf("Reflect:%v", reflectErr)
		return nil, reflectErr
	}
	var tempBuf = make([]byte, 32*1024)
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

func usingMap(in definition.ProtocolDefinition, foundCommand definition.Command, source *map[interface{}]interface{}) ([]byte, error) {
	var tempBuf = make([]byte, 32*1024)
	pos := 0
	tempBuf[pos] = foundCommand.ID
	pos++
	log.Printf("SOURCE: %+v", source)
	for _, a := range foundCommand.Arguments {
		uppercased := a.Name // strings.Title(a.Name)
		target := (*source)[uppercased]
		log.Printf("Converting %s to target: %v", a, target)
		octetCount, parseErr := argumentValueToOctet(a, tempBuf[pos:], target)
		if parseErr != nil {
			return nil, fmt.Errorf("Argument to Octet:%s", parseErr)
		}
		pos += octetCount
	}
	return tempBuf[:pos], nil
}

func ObjectToOctets(in definition.ProtocolDefinition, foundCommand definition.Command, source interface{}) ([]byte, error) {
	mapValue, isMap := source.(*map[interface{}]interface{})
	if isMap {
		return usingMap(in, foundCommand, mapValue)
	} else {
		return usingFields(in, foundCommand, source)
	}
}
