package main

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"io/ioutil"
	"log"
	"os"
	"project/parallel"
	"time"
	"encoding/binary"
)

// generateKey generates a new random AES key.
func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

// encryptData encrypts the given data with the given key using AES-256 encryption.
func encryptData(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ciphertext := aesgcm.Seal(nil, nonce, data, nil)
	encryptedData := make([]byte, 0, len(nonce)+len(ciphertext))
	encryptedData = append(encryptedData, nonce...)
	encryptedData = append(encryptedData, ciphertext...)
	return encryptedData, nil
}



// decryptData decrypts the given data with the given key using AES-256 encryption.
func decryptData(encryptedData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(encryptedData) < aes.BlockSize {
		return nil, io.ErrUnexpectedEOF
	}
	nonce := encryptedData[:12]
	ciphertext := encryptedData[12:]
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// zip folder to a new zipName.zip
func zipSource(folderPath, zipName string) {
	// Create a new zip file and add the contents of the folder to it
	zipFile, err := os.Create(zipName)
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileContent, err := ioutil.ReadFile(folderPath + "/" + file.Name())
		if err != nil {
			panic(err)
		}
		fileWriter, err := archive.Create(file.Name())
		if err != nil {
			panic(err)
		}
		_, err = fileWriter.Write(fileContent)
		if err != nil {
			panic(err)
		}
	}
}

// Encrypt zipName.zip to encrypted.zip
func encryptFile(key []byte, zipName string) {
	// Read the contents of the zip file and encrypt it
	data, err := ioutil.ReadFile(zipName)
	if err != nil {
		panic(err)
	}
	// fmt.Println("size of the zip folder to be encrypted in bytes: ", len(data))

	encryptedData, err := encryptData(data, key)
	if err != nil {
		panic(err)
	}

	// Write the encrypted data to a new file
	encryptedFile, err := os.Create("encrypted.zip")
	if err != nil {
		panic(err)
	}
	defer encryptedFile.Close()

	_, err = encryptedFile.Write(encryptedData)
	if err != nil {
		panic(err)
	}
}


// Decrypt encrypted.zip to decrypted.zip
func decryptFile(key []byte) {
	// Read the encrypted data from the file and decrypt it
	encryptedData, err := ioutil.ReadFile("encrypted.zip")
	if err != nil {
		panic(err)
	}

	decryptedData, err := decryptData(encryptedData, key)
	if err != nil {
		panic(err)
	}
	
	// Write the decrypted data to a new file
	decryptedFile, err := os.Create("decrypted.zip")
	if err != nil {
		panic(err)
	}
	defer decryptedFile.Close()

	_, err = decryptedFile.Write(decryptedData)
	if err != nil {
		panic(err)
	}
}


func main() {
	//Generate a key for encryption
	key, err := generateKey()
	if err != nil {
		panic(err)
	}

	// key = []byte{164, 251, 124, 160, 192, 76, 167, 99, 197, 62, 83, 81, 98, 17, 42, 11, 209, 147, 211, 24, 160, 49, 82, 62, 163, 137, 218, 18, 113, 122, 57, 100}

	//Set the path of the folder you want to zip and encrypt
	folderPath := "work2"

	//Set the name of the output zip file
	zipName := "folder.zip"

	//Zip folder to 'zipName.zip'
	zipSource(folderPath, zipName)

	// start2 := time.Now()

	// // Encrypt 'zipName.zip' to encypted.zip
	// encryptFile(key, zipName)

	// // Decrypt encypted.zip to decrypted.zip
	// decryptFile(key)
	// log.Printf("main, execution time %s\n", time.Since(start2))	

	start := time.Now()

	//Encrypt 'zipName.zip' to encypted.zip
	parallel.EncryptFileP(key, zipName)
	
	// //Decrypt encypted.zip to decrypted.zip
	// err = parallel.DecryptFileP(key)
	// if err != nil {
	// 	panic(err)
	// }
	log.Printf("main, execution time %s\n", time.Since(start))
	log.Println("Done")
	log.Println(key)
	println("Type of this object is %T\n", key)
	println(parallel.FileSize)

	// create meta data file
	file, err := os.Create("dycryptMain/metaData.key")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// write key to file
	_, err = file.Write(key)
	if err != nil {
		panic(err)
	}

	// convert file size to byte
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(parallel.FileSize))

	// write file size to file
	err = binary.Write(file, binary.LittleEndian, size)
	if err != nil {
		// panic(err)
	}

}
