package main

import (
	"project/parallel"
	"os"
    "encoding/binary"
)


func main() {
	// err := os.Remove("metaData.key")
	// if err != nil {
	// 	panic(err)
	// }

	// open meta data file
	file, err := os.Open("metaData.key")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// read key from file
	key := make([]byte, 32) // replace 16 with your key size
	_, err = file.Read(key)
	if err != nil {
		panic(err)
	}

	// read file size from key file
	sizeBytes := make([]byte, 4)
	_, err = file.Read(sizeBytes)
	if err != nil {
		// handle error
	}

	// convert file size to int
	size := int(binary.LittleEndian.Uint32(sizeBytes))

	//Decrypt encyptedP.zip to decryptedP.zip
	err = parallel.DecryptFileP(key, size)
	if err != nil {
		panic(err)
	}
}

