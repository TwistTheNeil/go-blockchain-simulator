package ds

import "net"

type Block struct {
	Nonce   int
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

type Miner struct {
	Connection net.Conn
	Mined      int
}

type Client struct {
	Connection net.Conn
}
