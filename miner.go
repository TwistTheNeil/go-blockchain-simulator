package main

import (
	"./lib/ds"
	"crypto/sha256"
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"strconv"
)

func OpenConnection(addr string, port int) net.Conn {
	conn, err := net.Dial("tcp", addr+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("ERROR: ")
		fmt.Println(err)
	}
	return conn
}

func GetPortFromPtr(port *int) int {
	x := fmt.Sprintf("%d", *port)
	p, _ := strconv.Atoi(x)
	return p
}

func Mine(payload string, target string) (string, int) {
	for nonce := 0; ; nonce++ {
		h := sha256.New()
		h.Write([]byte(payload + strconv.Itoa(nonce)))
		hash := fmt.Sprintf("%x", h.Sum(nil))
		if hash[:len(target)] == target {
			fmt.Println(target)
			fmt.Println("returning", hash, nonce)
			return hash, nonce
			break
		}
	}
	return "", -1
}

func main() {
	server_address_flag := flag.String("dest", "localhost",
		"Address of the server")
	server_port_flag := flag.Int("port", 9876,
		"Port used to connect to server")

	flag.Parse()

	var server_port int = GetPortFromPtr(server_port_flag)

	server_connection := OpenConnection(*server_address_flag, server_port)
	defer server_connection.Close()

	// Create gob communicating with server
	enc := gob.NewEncoder(server_connection)
	dec := gob.NewDecoder(server_connection)
	var msg ds.Message

	for {
		err := dec.Decode(&msg)
		if err != nil {
			fmt.Println("Event when decoding message from server: ", err)
			if err.Error() == "EOF" {
				break
			}
		}
		fmt.Println(msg)

		hash, nonce := Mine(msg.WorkingBlock.Payload, msg.Target)

		msg.WorkingBlock.Hash = hash
		msg.WorkingBlock.Nonce = nonce

		err = enc.Encode(&msg)
		if err != nil {
			fmt.Println("Event when encoding message for server: ", err)
			if err.Error() == "EOF" {
				break
			}
		}
	}

}
