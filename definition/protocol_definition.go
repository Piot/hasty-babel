package definition

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Argument struct {
	Name         string `yaml:"name"`
	ArgumentType string `yaml:"type"`
}

type Command struct {
	ID        byte       `yaml:"id"`
	Name      string     `yaml:"name"`
	Arguments []Argument `yaml:"args"`
}

func (in *Command) String() string {
	return fmt.Sprintf("[command '%s' id:%d]", in.Name, in.ID)
}

type ProtocolDefinition struct {
	Commands []Command
}

func NewProtocolDefinitionFromOctets(data []byte) (ProtocolDefinition, error) {
	protocolDefinition := ProtocolDefinition{}
	err := yaml.Unmarshal([]byte(data), &protocolDefinition)
	if err != nil {
		return ProtocolDefinition{}, err
	}
	return protocolDefinition, nil
}

func NewProtocolDefinitionFromFilePath(path string) (ProtocolDefinition, error) {
	octets, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		return ProtocolDefinition{}, readErr
	}

	return NewProtocolDefinitionFromOctets(octets)
}

func (in ProtocolDefinition) FindCommand(id byte) *Command {
	for _, c := range in.Commands {
		if c.ID == id {
			return &c
		}
	}

	return nil
}

// FindCommandUsingName :
func (in ProtocolDefinition) FindCommandUsingName(id string) *Command {
	for _, c := range in.Commands {
		if c.Name == id {
			return &c
		}
	}

	return nil
}

func (in ProtocolDefinition) String() string {
	s := ""
	for _, v := range in.Commands {
		s += fmt.Sprintf("\n Command: %d %s", v.ID, v.Name)
		for _, a := range v.Arguments {
			s += fmt.Sprintf("\n - %s : %s", a.Name, a.ArgumentType)
		}
	}

	return s
}
