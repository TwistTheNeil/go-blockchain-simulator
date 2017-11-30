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
		fmt.Println("Error while opening connection to server: ", err)
	}
	return conn
}

func GetPortFromPtr(port *int) int {
	x := fmt.Sprintf("%d", *port)
	p, _ := strconv.Atoi(x)
	return p
}

func Mine(payload string, target string, mined_message chan ds.Block, messenger chan bool) {
	fmt.Println("mine start")
	for nonce := 0; ; nonce++ {
		fmt.Println("mine 1")
		select {
		case cancel := <-messenger:
			if cancel == true {
				return
			}
		default:
			h := sha256.New()
			h.Write([]byte(payload + strconv.Itoa(nonce)))
			hash := fmt.Sprintf("%x", h.Sum(nil))
			if hash[:len(target)] == target {
				fmt.Println("mine done!")
				mined_message <- ds.Block{nonce, 0, payload, "", hash, false}
				return
			}
		}
	}
	fmt.Println("mine cancelled")
	mined_message <- ds.Block{-1, 0, payload, "", "", false}
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

	mine_cancel := make(chan bool)
	new_message := make(chan ds.Message)
	mined_block := make(chan ds.Block)

	// Receive message from server
	go func() {
		for {
			err := dec.Decode(&msg)
			if err != nil {
				fmt.Println("Event when decoding message from server: ", err)
				if err.Error() == "EOF" {
					break
				}
			}
			fmt.Println("decoded msg")
			if msg.Mined == true {
				fmt.Println("cancel chan sent")
				//				mine_cancel <- true
			} else {
				fmt.Println("send new msg from select{}")
				new_message <- msg
			}
		}
	}()

	for {
		select {
		case new_msg := <-new_message:
			fmt.Println("go mine")
			go Mine(new_msg.WorkingBlock.Payload, new_msg.Target, mined_block, mine_cancel)

		case new_block := <-mined_block:
			fmt.Println("Do i ever get here?")
			if new_block.Nonce > -1 {
				fmt.Println("nonce > 1")
				msg.WorkingBlock.Hash = new_block.Hash
				msg.WorkingBlock.Nonce = new_block.Nonce

				// Send message to server
				go func(msg ds.Message) {
					for {
						err := enc.Encode(&msg)
						if err != nil {
							fmt.Println("Event when encoding message for server: ", err)
							if err.Error() == "EOF" {
								break
							}
						}
					}
				}(msg)
			}
		}
	}
}
