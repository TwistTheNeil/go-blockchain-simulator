package ds

import (
	"crypto/sha256"
	"fmt"
	"net"
	"strconv"
	"time"
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
	Target       string
}

type Miner struct {
	Connection net.Conn
	Mining     bool
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

func (blockchain *Blockchain) BlockchainStatus() {
OUTER:
	for {
		time.Sleep(2 * time.Second)

		if len(blockchain.Blocks) > 1 {
			for i, block := range blockchain.Blocks[:len(blockchain.Blocks)-1] {
				if (block.Hash != blockchain.Blocks[i+1].Prev && block.Hash != "") || (block.Hash == "" && blockchain.Blocks[i+1].Prev == "") {
					blockchain.Complete = false
					if i < blockchain.Last {
						blockchain.Tamper = true
					}
					continue OUTER
				}
			}
		} else if len(blockchain.Blocks) == 1 && blockchain.Blocks[0].Hash == "" {
			continue OUTER
		}
		blockchain.Complete = true
		blockchain.Tamper = false
	}
}
