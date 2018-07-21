package demodcore

import (
	"strings"
	"fmt"
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
