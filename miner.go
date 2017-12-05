package main

import (
	"./lib/ds"
	"crypto/sha256"
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func OpenConnection(addr string, port int) (net.Conn, error) {
	var err error
	var conn net.Conn

	for i := 0; i < 5; i++ {
		conn, err = net.Dial("tcp", addr+":"+strconv.Itoa(port))
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	return conn, err
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

	fmt.Println("My PID:", os.Getpid())
	var server_port int = GetPortFromPtr(server_port_flag)

	server_connection, err := OpenConnection(*server_address_flag, server_port)
	if err != nil {
		fmt.Println("Error while connecting to server:", err)
		os.Exit(2)
	}
	defer server_connection.Close()

	// Create gob communicating with server
	enc := gob.NewEncoder(server_connection)
	dec := gob.NewDecoder(server_connection)

	for {
		var msg ds.Message

		err := dec.Decode(&msg)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Info: Server closed connection")
				break
			} else {
				fmt.Println("Event when decoding message from server: ", err)
			}
		}

		hash, nonce := Mine(msg.WorkingBlock.Payload, msg.Target)

		msg.WorkingBlock.Hash = hash
		msg.WorkingBlock.Nonce = nonce

		err = enc.Encode(&msg)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Info: Server closed connection")
				break
			} else {
				fmt.Println("Event when encoding message for server: ", err)
			}
		}
	}

}
