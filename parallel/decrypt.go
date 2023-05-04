package parallel

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"

	// "crypto/rand"
	"io"
	// "io/ioutil"
	"os"
	// "runtime"
	// "sync"
)

// decryptFile decrypts the given encrypted file with the given key and writes the decrypted contents to a new file.
func DecryptFileP(key []byte) error {
 
    dat, err := os.ReadFile("..\\encryptedP.zip")
    if err != nil {
        return err
    }
    ChunkSize = 487+28

    decryptedFile, err := os.Create("decryptedP.zip")
    if err != nil {
        return err
    }
    defer decryptedFile.Close()

    // Write the decrypted data to the decrypted file
    writer := bufio.NewWriter(decryptedFile)

    for i := 0; i < len(dat); i += ChunkSize {
        end := i + ChunkSize
        if end > len(dat) {
            end = len(dat)
        }
        data, err:= DecryptDataP(dat[i:end], key)
        if err != nil {
            return err
        }
        _, err = writer.Write(data)
    }
  
    writer.Flush()

    return nil
}


// decryptDataP decrypts the given data with the given key using AES-256 encryption using goroutines.
func DecryptDataP(encryptedData []byte, key []byte) ([]byte, error) {

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    if len(encryptedData) < aes.BlockSize {
        return  nil, io.ErrUnexpectedEOF
    }

    nonce := encryptedData[:12]
    ciphertext := encryptedData[12:]

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return  nil, err
    }

    return plaintext, nil
}
