package main

import (
	// Standard Library
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	// Third party Libraries
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

/*
==============================================================================
Blockchain and Block data structures
==============================================================================
*/

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

/*
==============================================================================
Blockchain and Block Helper functions
==============================================================================
*/
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

// Function to check if a new block is valid
func isBlockValid(oldBlock Block, newBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if calculateBlockHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

// Function for picking a longest chain if there are two competing blocks submitted
func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

/*
==============================================================================
Request Handler functions and supporting data structures
==============================================================================
*/
// Function to handle retreiving the root route
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Blockchain, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

// Message struct for POST requests
type Message struct {
	BPM int
}

//Function to handle POST request to the root route
func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], m.BPM)
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	if isBlockValid(Blockchain[len(Blockchain)-1], newBlock) {
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
		spew.Dump(Blockchain)
	}
	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

// Function to report on the outcome of the POST request
func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

/*
==============================================================================
Orchestration and server functions
==============================================================================
*/

// Function to define HTTP server routes and actions
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

// Function to define and start a web server
func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("PORT")
	log.Println("Listening on ", httpAddr)
	server := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

// The main event
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		currentTime := time.Now()
		// I chose a struct literal here, to make things explicit
		genesisBlock := Block{
			Index:     0,
			Timestamp: currentTime.String(),
			BPM:       0,
			Hash:      "",
			PrevHash:  "",
		}
		genesisBlock.Hash = calculateBlockHash(genesisBlock)
		spew.Dump(genesisBlock)
		Blockchain = append(Blockchain, genesisBlock)
	}()
	log.Fatal(run())
}
