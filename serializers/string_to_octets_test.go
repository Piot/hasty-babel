package serializers

import (
	"log"
	"testing"

	"github.com/piot/hasty-babel/definition"
)

var testStringProtocol = `
---
commands:
  - id: 1
    name: Item
    args:
    - name: itemid
      type: u32
    - name: shorter
      type: u16
    - name: flag
      type: u8
    - name: path
      type: asciistring
    - name: isactive
      type: bool
`

var valueString = `
Itemid: 1
Shorter: 32
Flag: 33
Path: path/to/something
Isactive: true
`

func TestStringToOctets(t *testing.T) {
	definition, _ := definition.NewProtocolDefinitionFromOctets([]byte(testStringProtocol))
	cmd := definition.FindCommand(1)
	octets, toOctetsErr := StringToOctets(definition, *cmd, valueString)
	if toOctetsErr != nil {
		t.Fatal(toOctetsErr)
	}

	log.Printf("octets:%X", octets)
	output, _, toStringErr := OctetsToString(definition, octets)
	if toStringErr != nil {
		t.Fatal(toStringErr)
	}
	log.Printf("We are back '%s'", output)
}
