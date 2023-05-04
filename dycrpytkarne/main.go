package main

import (
	"log"
	"project/parallel"
	"time"
)


func main() {
	// Generate a key for encryption
	// key, err := generateKey()
	// if err != nil {
	// 	panic(err)
	// }

	// Set the path of the folder you want to zip and encrypt
	// folderPath := "work2"

	// // Set the name of the output zip file
	// zipName := "folder.zip"

	// // Zip folder to 'zipName.zip'
	// zipSource(folderPath, zipName)

	//key := [100]byte{35, 152, 210 106 244 164 222 91 49 31 112 46 91 210 31 132 218 7 168 58 59 124 137 68 46 204 227 189 159 142 105 229]
	
	//key := []byte{230, 129, 203, 150, 85, 121, 133, 196, 150, 77, 147, 107, 255, 23, 176, 246, 4, 204, 4, 29, 99, 89, 137, 14, 247, 107, 103, 250, 35, 15, 104, 144}
	key := []byte{164, 251, 124, 160, 192, 76, 167, 99, 197, 62, 83, 81, 98, 17, 42, 11, 209, 147, 211, 24, 160, 49, 82, 62, 163, 137, 218, 18, 113, 122, 57, 100}

	start := time.Now()

	//Decrypt encypted.zip to decrypted.zip
	err := parallel.DecryptFileP(key)
	if err != nil {
		panic(err)
	}
	log.Printf("main, execution time %s\n", time.Since(start))
	log.Println("Done")
	log.Println(key)
}

