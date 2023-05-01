package parallel

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"bufio"
	"runtime"
)

// a struct to hold index and data chunk
type Chunk struct {
    index int
    data  []byte
}

// make global array for order of chunks and chunk size
var order []int
var chunkSize int

// encryptDataP encrypts the given data with the given key using AES-256 encryption using goroutines.
func EncryptDataP(data Chunk, key []byte, wg *sync.WaitGroup, encryptedData chan Chunk, errCh chan error) {
    defer wg.Done()

    block, err := aes.NewCipher(key)
    if err != nil {
        errCh <- err
        return
    }

    nonce := make([]byte, 12)
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        errCh <- err
        return
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        errCh <- err
        return
    }

    ciphertext := aesgcm.Seal(nil, nonce, data.data , nil)
    data.data= append(nonce, ciphertext...)
    encryptedData <- data
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

func EncryptFileP(key []byte, zipName string) error {
    // Read the contents of the zip file
    zipFile, err := ioutil.ReadFile(zipName)
    if err != nil {
        return err
    }
	// log.Println(len(zipFile))
	// log.Println(256*256)
    // Create a new encrypted file
    encryptedFile, err := os.Create("encryptedP.zip")
    if err != nil {
        return err
    }
    defer encryptedFile.Close()

    // Create a WaitGroup to wait for all the goroutines to finish
    var wg sync.WaitGroup

    // Create a channel to receive the encrypted data
    encryptedData := make(chan Chunk)

    // Create a channel to receive errors
    errCh := make(chan error)

    // Encrypt the data in chunks using goroutines
    numChunks := runtime.NumCPU()
    print("numChunks: ", numChunks)
    chunkSize = len(zipFile) / numChunks
	print("len(zipFile): ", len(zipFile))
    for i := 0; i < len(zipFile); i += chunkSize {
        wg.Add(1)
        end := i + chunkSize
        if end > len(zipFile) {
            end = len(zipFile)
        }
        // make encrypted chunk object and send it to channel
        c := Chunk{index: i, data: zipFile[i:end]}
        go EncryptDataP(c, key, &wg, encryptedData, errCh)
    }

    // Wait for all the goroutines to finish
    go func() {
        wg.Wait()
        close(encryptedData)
    }()

    // Write the encrypted data to the encrypted file
    writer := bufio.NewWriter(encryptedFile)
    for encryptedChunk := range encryptedData {
        _, err := writer.Write(encryptedChunk.data)
        order = append(order, encryptedChunk.index)
        if err != nil {
            return err
        }
    }
    writer.Flush()
    // Check for any errors during encryption
    select {
    case err := <-errCh:
        return err
    default:
        return nil
    }
}


// decryptFile decrypts the given encrypted file with the given key and writes the decrypted contents to a new file.
func DecryptFileP(key []byte) error {
    encryptedFile, err := os.Open("encryptedP.zip")
    if err != nil {
        return err
    }
    defer encryptedFile.Close()

    encryptedFileInfo, err := encryptedFile.Stat()
    if err != nil {
        return err
    }

    encryptedData := make([]byte, encryptedFileInfo.Size())
    if _, err := encryptedFile.Read(encryptedData); err != nil {
        return err
    }

    decryptData := make([]byte, chunkSize*len(order)) 

    // Make a new array to store encrypted data in order
    encryptDataOrdered := make([]byte, len(encryptedData))
    for i := 0; i < len(order)-1; i++ {
        for j := 0; j < chunkSize; j++ {
            encryptDataOrdered[i*chunkSize+j] = encryptedData[order[i]+j]
        }
    }

    // Decrypt the data in chunks using simple function
    chunkSize = chunkSize+12
  
    for i := 0; i < len(encryptDataOrdered); i += chunkSize {
        end := i + chunkSize
        if end > len(encryptDataOrdered) {
            end = len(encryptDataOrdered)
        }
        data, err:= DecryptDataP(encryptDataOrdered[i:end], key)
        if err != nil {
            return err
        }
        copy(decryptData[i:i+chunkSize], data)
    }

    decryptedFile, err := os.Create("decryptedP.zip")
    if err != nil {
        return err
    }
    defer decryptedFile.Close()

    // Write the decrypted data to the decrypted file
    writer := bufio.NewWriter(decryptedFile)
    if _, err := writer.Write(decryptData); err != nil {
        return err
    }
    writer.Flush()

    return nil
}