package demodcore

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type JsonFloat32 []float32

func (u JsonFloat32) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%f", u)), ",")
	}
	return []byte(result), nil
}

func (u JsonFloat32) MarshalByteArray() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, u)
	return buf.Bytes()
}
