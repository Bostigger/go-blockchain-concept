package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"time"
)

type Book struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedDate string `json:"published_date"`
	ISBN          string `json:"isbn"`
}
type Block struct {
	Data         BookCheckout
	PreviousHash string
	Hash         string
	Position     int
	TimeStamp    string
}
type Blockchain struct {
	blocks []*Block
}

type BookCheckout struct {
	BookId       string    `json:"book_id"`
	User         string    `json:"user"`
	CheckOutDate time.Time `json:"check_out_date"`
	IsGenesis    bool      `json:"is_genesis"`
}

var BlockChain *Blockchain

func (b *Block) generateHash() {
	bytes, _ := json.Marshal(b.Data)
	data := string(b.Position) + b.TimeStamp + string(bytes) + b.PreviousHash
	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))
}

func (bc *Blockchain) AddBlock(data BookCheckout) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	block := CreateBlock(prevBlock, data)
	if validBlock(block, prevBlock) {
		bc.blocks = append(bc.blocks, block)
	}
}

func CreateBlock(preBlock *Block, data BookCheckout) *Block {
	block := &Block{}
	block.Data = data
	block.Position = preBlock.Position + 1
	block.PreviousHash = preBlock.Hash
	block.TimeStamp = time.Now().String()
	block.generateHash()
	return block
}

func validBlock(block, prevBlock *Block) bool {
	if block.PreviousHash != prevBlock.Hash {
		return false
	}
	if !block.validateHash(block.Hash) {
		return false
	}
	if prevBlock.Position+1 != block.Position {
		return false
	}
	return true
}

func (block *Block) validateHash(hash string) bool {
	block.generateHash()
	if block.Hash != hash {
		return false
	}
	return true
}

func NewBlockChain() *Blockchain {
	return &Blockchain{
		[]*Block{GenesisBlock()},
	}
}

func GenesisBlock() *Block {
	return CreateBlock(&Block{}, BookCheckout{IsGenesis: true})
}
func main() {
	BlockChain = NewBlockChain()
	r := mux.NewRouter()
	r.HandleFunc("/new", newBook).Methods("POST")
	r.HandleFunc("/block", writeBlock).Methods("POST")
	r.HandleFunc("/", getBlockChain).Methods("GET")

	go func() {
		for _, block := range BlockChain.blocks {
			fmt.Printf("Prev.Hash : %x\n", block.PreviousHash)
			bytes, _ := json.MarshalIndent(block.Data, "", "")
			fmt.Printf("Data %v\n", string(bytes))
			fmt.Printf("Hash %v\n", block.Hash)
			fmt.Println()

		}
	}()

	log.Println("server listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}

func newBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not create: %v", err)
		w.Write([]byte("error creating a new book"))
		return
	}
	h := md5.New()
	_, err := io.WriteString(h, book.ID+book.PublishedDate)
	if err != nil {
		return
	}
	book.ID = fmt.Sprintf("%x", h.Sum(nil))
	res, err := json.MarshalIndent(book, "", "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not marshal payload %v", err)
		w.Write([]byte("coundnt save the book"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func writeBlock(w http.ResponseWriter, r *http.Request) {
	var checkOutItem BookCheckout
	if err := json.NewDecoder(r.Body).Decode(&checkOutItem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		w.Write([]byte("error writing block "))
		return
	}

	BlockChain.AddBlock(checkOutItem)
}
func getBlockChain(w http.ResponseWriter, r *http.Request) {
	bcbytes, err := json.MarshalIndent(BlockChain.blocks, "", "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	io.WriteString(w, string(bcbytes))
}
