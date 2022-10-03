package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

func HandleError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	HandleError(encoder.Encode(i))
	return aBuffer.Bytes()
}

func FromBytes(from []byte, to interface{}) {
	decoder := gob.NewDecoder(bytes.NewReader(from))
	HandleError(decoder.Decode(to))
}
