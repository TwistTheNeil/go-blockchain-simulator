package ds

import (
	"crypto/sha256"
	"fmt"
	"net"
	"strconv"
)

type Block struct {
	Nonce   int
	Index   int
	Payload string
	Prev    string
	Hash    string
	Valid   bool
}

type Blockchain struct {
	Blocks   []Block
	Last     int
	Tamper   bool
	Complete bool
}

type Message struct {
	WorkingBlock Block
	Mined        bool
}

type Miner struct {
	Connection net.Conn
	Mined      int
}

type Client struct {
	Connection net.Conn
}

func (m *Message) Verify() bool {
	h := sha256.New()
	h.Write([]byte(m.WorkingBlock.Payload + strconv.Itoa(m.WorkingBlock.Nonce)))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	if hash == m.WorkingBlock.Hash {
		return true
	}
	return false
}
