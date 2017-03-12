package serializers

import (
	"bytes"
	"encoding/hex"
	"log"
	"testing"

	"github.com/piot/hasty-babel/definition"
)

var testProtocol = `
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

var testData = `
id: 1324
shorter: 13
flag: 144
path: Hello World
isactive: false
`

type testObject struct {
	Itemid   uint32
	Shorter  uint16
	Flag     uint8
	Path     string
	Isactive bool
}

func TestObjectToOctets(t *testing.T) {
	definition, _ := definition.NewProtocolDefinitionFromOctets([]byte(testProtocol))
	ox := testObject{Itemid: 0x1234, Shorter: 13, Flag: 0x33, Path: "this/is/a/path", Isactive: true}
	cmd := definition.FindCommand(1)
	octets, toOctetsErr := ObjectToOctets(definition, *cmd, &ox)
	if toOctetsErr != nil {
		t.Fatal(toOctetsErr)
	}
	correctAnswer := make([]byte, 150)
	length, _ := hex.Decode(correctAnswer, []byte("0100001234000D330E746869732F69732F612F7061746801"))
	correctAnswer = correctAnswer[:length]
	if bytes.Compare(octets, correctAnswer) != 0 {
		t.Fatal("Octets do not match")
	}

	humanReadableString, octetCountConsumed, toStringErr := OctetsToString(definition, octets)
	if toStringErr != nil {
		t.Fatal(toStringErr)
	}
	if octetCountConsumed != len(correctAnswer) {
		t.Fatalf("Wrong number of octets consumed:%d", octetCountConsumed)
	}

	log.Printf("Parsed back:%s", humanReadableString)

	newObject := testObject{}
	octetsConsumed, toObjectErr := OctetsToObject(definition, octets, &newObject)
	if toObjectErr != nil {
		t.Fatal(toObjectErr)
	}
	if octetsConsumed != len(correctAnswer) {
		t.Fatalf("Wrong number of octets consumed:%d", octetCountConsumed)
	}
	log.Printf("New Object:%v", newObject)
	if newObject.Flag != 0x33 {
		t.Fatalf("Wrong flag number:%d", newObject.Flag)
	}
	if newObject.Path != "this/is/a/path" {
		t.Fatalf("Wrong path:%s", newObject.Path)
	}
}
