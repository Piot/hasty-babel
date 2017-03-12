package serializers

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/piot/hasty-babel/definition"
)

func parseType(arg definition.Argument, data []byte) (string, byte, error) {
	switch arg.ArgumentType {
	case "u32":
		v := binary.BigEndian.Uint32(data)
		return fmt.Sprintf("0x%08X", v), 4, nil
	case "u16":
		v := binary.BigEndian.Uint16(data)
		return fmt.Sprintf("0x%04X", v), 2, nil
	case "u8":
		v := data[0]
		return fmt.Sprintf("0x%02X", v), 1, nil
	case "bool":
		v := data[0]
		var s string
		if v == 1 {
			s = "true"
		} else if v == 0 {
			s = "false"
		} else {
			return "", 0, fmt.Errorf("Bool must be either 0 or 1")
		}
		return s, 1, nil
	case "string":
		n := data[0]
		s := data[1 : 1+n]
		return "'" + string(s) + "'", n + 1, nil
	case "asciistring":
		n := data[0]
		s := data[1 : 1+n]
		return "'" + string(s) + "'", n + 1, nil
	}

	return "", 0, fmt.Errorf("Illegal argument type '%s'", arg.ArgumentType)
}

func parseArgument(in definition.ProtocolDefinition, arg definition.Argument, data []byte) (string, byte, error) {
	s, n, err := parseType(arg, data)
	if err != nil {
		return "", 0, err
	}
	return arg.Name + "=" + s, n, nil
}

func OctetsToString(in definition.ProtocolDefinition, data []byte) (string, int, error) {
	pos := 0
	var str []string
	commandID := data[pos]
	pos++

	foundCommand := in.FindCommand(commandID)
	if foundCommand == nil {
		return "", pos, fmt.Errorf("Couldn't find command 0x%02X", commandID)
	}
	for _, a := range foundCommand.Arguments {
		parsedString, octetCount, parseErr := parseArgument(in, a, data[pos:])
		if parseErr != nil {
			return "", pos, parseErr
		}
		pos += int(octetCount)
		str = append(str, parsedString)
	}
	return strings.Join(str, ", "), pos, nil
}
