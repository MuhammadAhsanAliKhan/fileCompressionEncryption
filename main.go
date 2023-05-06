package main

import (
	"archive/zip"
	//"crypto/aes"
	//"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"io"

	//"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"project/parallel"
	"time"
)

// generateKey generates a new random AES key.
func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}


func zipSource2(folderToZip, outputFilePath string) {
    
    // Create a new zip archive.
    zipFile, err := os.Create(outputFilePath)
    if err != nil {
        panic(err)
    }
    defer zipFile.Close()
    zipWriter := zip.NewWriter(zipFile)
    defer zipWriter.Close()

    // Walk through all the directories and files in the folder to zip.
    err = filepath.Walk(folderToZip, func(filePath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        // Skip the folder to zip itself.
        if filePath == folderToZip {
            return nil
        }
        
        // Create a new zip header for the current file or directory.
        zipHeader, err := zip.FileInfoHeader(info)
        if err != nil {
            return err
        }
        zipHeader.Name = filePath[len(folderToZip)+1:]

        // If the current path is a directory, add a trailing slash to the name.
        if info.IsDir() {
            zipHeader.Name += "/"
        }

        // Create a new zip file entry and write the contents to it.
        zipFileEntry, err := zipWriter.CreateHeader(zipHeader)
        if err != nil {
            return err
        }
        if !info.IsDir() {
            fileToZip, err := os.Open(filePath)
            if err != nil {
                return err
            }
            defer fileToZip.Close()
            _, err = io.Copy(zipFileEntry, fileToZip)
            if err != nil {
                return err
            }
        }

        return nil
    })

    if err != nil {
        panic(err)
    }
    
	println("Successfully zipped all folders and files!")

}



func main() {
	//Generate a key for encryption
	key, err := generateKey()
	if err != nil {
		panic(err)
	}
	//Set the path of the folder you want to zip and encrypt
	folderPath := "work"

	//Set the name of the output zip file
	zipName := "folder.zip"
	zipSource2(folderPath, zipName)
	start := time.Now()

	//Encrypt 'zipName.zip' to encypted.zip
	parallel.EncryptFileP(key, zipName)
	
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
		panic(err)
	}

}
