package parallel

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
)

// a struct to hold index and data chunk
type Chunk struct {
    index int
    data  []byte
}

// make global array for order of chunks and chunk size
var order []int

var Order[101]Chunk
var FileSize int
var ChunkSize int

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
    //println("This is before chiper",len(data.data))
    data.data= append(nonce, ciphertext...)
    //println("This is after chiper",(len(data.data)))
    encryptedData <- data
}

func EncryptFileP(key []byte, zipName string) error {
    // Read the contents of the zip file
    zipFile, err := ioutil.ReadFile(zipName)
    if err != nil {
        return err
    }
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
    println("numChunks: ", numChunks)
    FileSize = len(zipFile)
    ChunkSize = FileSize / 100
	println("len(zipFile): ", len(zipFile))
    count := 0
    for i := 0; i < len(zipFile); i += ChunkSize {
        
        wg.Add(1)
        end := i + ChunkSize
        if end > len(zipFile) {
            end = len(zipFile)
        }
        // make encrypted chunk object and send it to channel
        c := Chunk{index: count, data: zipFile[i:end]}
        count++
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

        if err != nil {
            return err
        }
        Order[encryptedChunk.index] = encryptedChunk
    }

    for i := range Order{
        _, err := writer.Write(Order[i].data)
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