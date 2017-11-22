package ds

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
