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
