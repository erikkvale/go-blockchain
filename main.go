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

type Block struct {
	Index int
	Timestamp string
	BPM int
	Hash string
	PrevHash string
}

var Blockchain []Block