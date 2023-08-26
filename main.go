package main

import (
	"crypto/md5"
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
	NextHash     string
	Position     string
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

//var Blockchain *Blockchain

func (bc *Blockchain) AddBlock() {

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/new", newBook).Methods("POST")
	r.HandleFunc("/block", writeBlock).Methods("POST")
	r.HandleFunc("/", getBlockChain).Methods("GET")
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
		w.Write([]byte("error decoding the book checkout "))
		return
	}
}
func getBlockChain(w http.ResponseWriter, r *http.Request) {

}
