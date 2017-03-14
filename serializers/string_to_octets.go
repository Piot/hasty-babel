package serializers

import (
	"log"

	yaml "gopkg.in/yaml.v2"

	"github.com/piot/hasty-babel/definition"
)

func StringToOctets(in definition.ProtocolDefinition, foundCommand definition.Command, values string) ([]byte, error) {
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(values), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
		return []byte{}, err
	}
	log.Printf("--- m:\n%v\n\n", m)

	return ObjectToOctets(in, foundCommand, &m)
}

func ValueStringToOctets(in definition.ProtocolDefinition, values string) ([]byte, error) {
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(values), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
		return []byte{}, err
	}
	log.Printf("m: %+v\n", m)
	valuesArray := m["values"].([]interface{})
	log.Printf("valuesArray: %+v\n", valuesArray)
	first := valuesArray[0].(map[interface{}]interface{})
	commandName := first["command"].(string)
	log.Printf("First: %+v\n", first)
	log.Printf("Command: %+v\n", commandName)
	foundCommand := in.FindCommandUsingName(commandName)

	return ObjectToOctets(in, *foundCommand, &first)
}
