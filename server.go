package main

import (
	"./lib/ds"
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
)

func OpenListener(port int, new_connection chan net.Conn, what chan bool) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("OMG 1: Neil do something!")
		fmt.Println(err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("OMG 2: Neil do something!")
			}
			new_connection <- conn
		}
		what <- true
	}()

}

func Broadcast(miners map[net.Conn]ds.Miner, msg ds.Message) {
	for conn, miner := range miners {
		enc := gob.NewEncoder(conn)
		err := enc.Encode(&msg)
		if err != nil {
			fmt.Println("Broadcast error: ", conn, " -> ", miner, "->", err)
		}
	}
}

func UpdateBlockchain() {
	// broadcast to all miners if nothing is going on
}

func GetPortFromPtr(port *int) int {
	x := fmt.Sprintf("%d", *port)
	p, _ := strconv.Atoi(x)
	return p
}

func main() {
	miner_port_flag := flag.Int("mport", 9876,
		"Port on which server listens on for miner")
	client_port_flag := flag.Int("cport", 9875,
		"Port on which server listens on client")

	flag.Parse()

	miners := make(map[net.Conn]ds.Miner)
	clients := make(map[net.Conn]ds.Client)
	blockchain := ds.Blockchain{[]ds.Block{}, -1, false, true}

	fmt.Println(len(blockchain.Blocks))

	var miner_port int = GetPortFromPtr(miner_port_flag)
	var client_port int = GetPortFromPtr(client_port_flag)

	new_miner_connection := make(chan net.Conn)
	new_client_connection := make(chan net.Conn)
	remove_miner_connection := make(chan net.Conn)
	remove_client_connection := make(chan net.Conn)
	new_miner_message := make(chan ds.Message)
	new_payload := make(chan string)

	// TODO: Get rid of this channel eventually
	what := make(chan bool)

	OpenListener(miner_port, new_miner_connection, what)
	OpenListener(client_port, new_client_connection, what)

	for {
		select {
		case conn := <-new_miner_connection:
			miners[conn] = ds.Miner{conn, 0}
			fmt.Println("Received connection from ", conn)

			go func(conn net.Conn) {
				dec := gob.NewDecoder(conn)
				for {
					var message ds.Message
					err := dec.Decode(&message)
					if err != nil {
						fmt.Println("GOB DECODE ERROR: ", err)
						break
					}
					fmt.Println("GOT HERE")
					fmt.Println(message)
					new_miner_message <- message
					fmt.Println("FIN HERE")
				}
				fmt.Println("Attempting to close connection ", conn)
				remove_miner_connection <- conn
			}(conn)

		case msg := <-new_miner_message:
			if msg.WorkingBlock.Index > blockchain.Last {
				if msg.Verify() {
					msg.Mined = true
					blockchain.Blocks[msg.WorkingBlock.Index] = msg.WorkingBlock
					blockchain.Last += 1
					Broadcast(miners, msg)
				}
			}

		case payload := <-new_payload:
			fmt.Println("Got payload: ", payload)
			blockchain_length := len(blockchain.Blocks)
			var prev_hash string
			if blockchain_length == 0 {
				prev_hash = "0000000000000000000000000000000000000000000000000000000000000000"
			} else {
				prev_hash = blockchain.Blocks[len(blockchain.Blocks)-1].Hash
			}
			payload_block := ds.Block{0, len(blockchain.Blocks), payload, prev_hash, "", false}
			blockchain.Blocks = append(blockchain.Blocks, payload_block)
			blockchain.Complete = false
			fmt.Println(len(blockchain.Blocks))
			fmt.Println(blockchain)

		case conn := <-new_client_connection:
			clients[conn] = ds.Client{conn}
			fmt.Println("Received client connection: ", conn)

			go func(conn net.Conn) {
				dec := gob.NewDecoder(conn)
				for {
					var message string
					err := dec.Decode(&message)
					if err != nil {
						fmt.Println("GOB DECODE ERROR: ", err)
						break
					}
					fmt.Println("GOT HERE")
					fmt.Println(message)
					new_payload <- message
					fmt.Println("FIN HERE")
				}
				fmt.Println("Closing client connection ", conn)
				remove_client_connection <- conn
			}(conn)

		case conn := <-remove_miner_connection:
			delete(miners, conn)

		case conn := <-remove_client_connection:
			delete(clients, conn)
		}
	}

	// TODO: remove stuff beyond this point
	<-what
	<-what
	fmt.Println("Debug: Success")
	os.Exit(0)
}
