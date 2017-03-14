package serializers

import (
	"bytes"
	"encoding/binary"
)

// SmallLengthToOctets : serialize length
func SmallLengthToOctets(length uint16) ([]byte, error) {

	if length <= 127 {
		buf := make([]byte, 1)
		buf[0] = uint8(length)
		return buf, nil
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, length|0x8000)
	return buf.Bytes(), nil
}
