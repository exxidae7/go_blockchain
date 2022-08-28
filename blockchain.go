package main

import (
	"fmt"
	"log"
	"strings"
	"time"
	"encoding/json"
	"crypto/sha256"
)

type Block struct {
	nonce int 
	previousHash [32]byte
	timestamp int64
	transactions []*Transactions	
}
func NewBlock(nonce int , previousHash [32]byte , transactions []*Transactions) *Block {
	b := new(Block)
	b.timestamp = time.Now().UnixNano()
	b.nonce = nonce
	b.previousHash = previousHash
	b.transactions = transactions
	return b
}


func (b *Block) Print() { 
	fmt.Printf("timestamp:      %d\n", b.timestamp)
	fmt.Printf("nonce:          %d\n", b.nonce)
	fmt.Printf("previous hash:  %x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
}
func (b *Block) Hash() [32]byte {
	m , _ := json.Marshal(b)
	fmt.Println(string(m))
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct{
		Timestamp int64 `json:"timestamp"`
		Nonce int `json:"nonce"`
		PreviousHash [32]byte `json:"previous_hash"`
		Transactions []*Transactions `json:"transactions"`
	}{
		Timestamp: b.timestamp,
		Nonce : b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})

}
type Blockchain struct { 
	transactionPool []*Transactions
	chain           []*Block
}

type Transactions struct {
	senderBlockchainAddress string
	recipientBlockchainAddress string
	value	 	 float32
}

func NewTransactions(sender string, recipient string , value float32) *Transactions {
	return &Transactions{sender , recipient, value}
}

func (t *Transactions) Print(){
	fmt.Printf("%s\n", strings.Repeat("_", 25))
	fmt.Printf("Sender Blockchain Address:           %s\n", t.senderBlockchainAddress)
	fmt.Printf("Recipient Blockchain Address:        %s\n", t.recipientBlockchainAddress)
	fmt.Printf("Value: 				    %.1f\n" , t.value)
}

func (t *Transactions) MarshalJSON() ([]byte, error){
	return json.Marshal(struct {
		Sender string `json:"sender_blockchain_address"`
		Recipient string `json:"recipient_blockchain_address"`
		Value  float32 `json:"value"`
	}{
		Sender: t.recipientBlockchainAddress,
		Recipient: t.recipientBlockchainAddress,
		Value: t.value,		
	})
}

func NewBlockchain() *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0 , b.Hash())
	return bc
}
func (bc *Blockchain) Print() {
	for i , block := range bc.chain { 
		fmt.Printf("%s chain %d %s\n" , strings.Repeat("=" ,25) , i,
			strings.Repeat("=",25))
		block.Print()
	}
	fmt.Printf("%s\n" , strings.Repeat("*" , 55))
}
func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain) - 1]
}

func (bc *Blockchain) AddTransaction(sender string , recipient string , value float32){
	t := NewTransactions(sender ,recipient , value)
	bc.transactionPool = append(bc.transactionPool, t)
}

func (bc *Blockchain) CreateBlock(nonce int , previousHash [32]byte) *Block {
	b := NewBlock(nonce , previousHash , bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transactions{}
	return b
}

func init(){
	log.SetPrefix("Blockchain: ")
}

func main(){
	// block := &Block{nonce: 1}
	// fmt.Printf("%x\n", block.Hash())
	
	blockChain := NewBlockchain()
	blockChain.Print()


	blockChain.AddTransaction("A" , "B" , 1.0)
	previousHash := blockChain.LastBlock().Hash()
	blockChain.CreateBlock(5, previousHash)
	blockChain.Print()

	blockChain.AddTransaction("C" , "D" , 4.0)
	blockChain.AddTransaction("X" , "Y" , 2.0)

	previousHash = blockChain.LastBlock().Hash()
	blockChain.CreateBlock(2 , previousHash)
	blockChain.Print()



	

}