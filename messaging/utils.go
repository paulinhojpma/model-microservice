package messaging

import (
	"encoding/json"
	"math/rand"
)

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func interfaceToByte(param interface{}) ([]byte, error) {
	byt, errByte := json.Marshal(param)
	if errByte != nil {
		return nil, errByte
	}
	return byt, nil
}

func byteToInterface(b []byte) (interface{}, error) {
	var i interface{}
	errJSON := json.Unmarshal(b, i)
	if errJSON != nil {
		return nil, errJSON
	}
	return i, nil
}
