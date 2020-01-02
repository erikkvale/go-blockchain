package main

import (
	// Standard Library
	"crypto/sha256"
	"encoding/hex"

	// "encoding/json"
	// "io"
	// "log"
	// "net/http"
	// "os"
	"time"
	// Third party Libraries
	// "github.com/davecgh/go-spew/spew"
	// "github.com/gorilla/mux"
	// "github.com/joho/godotenv"
)

// Block type for a block on my blockchain
type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

// Blockchain slice composed of Block types
var Blockchain []Block

// Function for calculating the hash of a Block
func calculateBlockHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Function to generate a new block
func generateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block
	currentTime := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = currentTime.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateBlockHash(newBlock)
	return newBlock, nil
}
