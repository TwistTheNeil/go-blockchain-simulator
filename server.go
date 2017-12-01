package main

import (
	"./lib/ds"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func DisplayStats(miners *map[net.Conn]*ds.Miner, blockchain *ds.Blockchain) {
	for {
		fmt.Println("=========================================================")
		fmt.Println("Blockchain:")
		bc, err := json.MarshalIndent(blockchain, "", "  ")
		if err != nil {
			fmt.Println("DisplayStats() error: ", err)
		}
		fmt.Println(string(bc))

		fmt.Println("Miners:")
		for _, miner := range *miners {
			m, err_m := json.MarshalIndent(miner, "", "  ")
			if err_m != nil {
				fmt.Println("DisplayStats() error: ", err_m)
			}
			fmt.Println(string(m))
		}

		time.Sleep(1 * time.Second)
	}

}

func OpenListener(port int, new_connection chan net.Conn, what chan bool) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("OpenListener() error:", err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("OpenListener() error:", err)
			}
			new_connection <- conn
		}
		what <- true
	}()
}

func Broadcast(miners *map[net.Conn]*ds.Miner, blockchain *ds.Blockchain) {
	for {
		if blockchain.Complete == true {
		} else if blockchain.Last == -1 && len(blockchain.Blocks) == 0 {
		} else if blockchain.Last < len(blockchain.Blocks)-1 {
			msg := ds.Message{blockchain.Blocks[blockchain.Last+1], false, "000"}
			go BroadcastToAllMiners(miners, msg)
		}

		time.Sleep(5 * time.Second)
	}
}

func BroadcastToAllMiners(miners *map[net.Conn]*ds.Miner, msg ds.Message) {
	for conn, miner := range *miners {
		enc := gob.NewEncoder(conn)
		err := enc.Encode(&msg)
		miner.Mining = true
		if err != nil {
			fmt.Println("Broadcast everyone error: ", conn, " -> ", miner, "->", err)
		}
	}
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

	miners := make(map[net.Conn]*ds.Miner)
	clients := make(map[net.Conn]ds.Client)
	blockchain := ds.Blockchain{[]ds.Block{}, -1, false, true, false}

	fmt.Println(len(blockchain.Blocks))

	var miner_port int = GetPortFromPtr(miner_port_flag)
	var client_port int = GetPortFromPtr(client_port_flag)

	new_miner_connection := make(chan net.Conn)
	new_client_connection := make(chan net.Conn)
	remove_miner_connection := make(chan net.Conn)
	remove_client_connection := make(chan net.Conn)
	new_miner_message := make(chan ds.Message)
	miner_chan := make(chan *ds.Miner)
	new_payload := make(chan string)

	// TODO: Get rid of this channel eventually
	what := make(chan bool)

	OpenListener(miner_port, new_miner_connection, what)
	OpenListener(client_port, new_client_connection, what)

	// Start broadcasting
	go Broadcast(&miners, &blockchain)
	go DisplayStats(&miners, &blockchain)
	go blockchain.BlockchainStatus()

	for {
		select {
		case conn := <-new_miner_connection:
			miners[conn] = &ds.Miner{conn, false, 0}
			fmt.Println("Received connection from ", conn)

			go func(conn net.Conn) {
				dec := gob.NewDecoder(conn)
				for {
					var message ds.Message
					err := dec.Decode(&message)
					if err != nil {
						fmt.Println("miner gob decode error: ", err)
						if err.Error() == "EOF" {
							break
						}
					}
					myself := miners[conn]
					myself.Mining = false
					new_miner_message <- message
					miner_chan <- myself
				}
				fmt.Println("Closing miner connection ", conn)
				remove_miner_connection <- conn
			}(conn)

		case msg := <-new_miner_message:
			miner := <-miner_chan
			if msg.WorkingBlock.Index > blockchain.Last {
				if msg.Verify() {
					msg.Mined = true
					if msg.WorkingBlock.Index > 0 {
						msg.WorkingBlock.Prev = blockchain.Blocks[msg.WorkingBlock.Index-1].Hash
					}
					blockchain.Blocks[msg.WorkingBlock.Index] = msg.WorkingBlock
					blockchain.Last += 1
					blockchain.Working = false
					blockchain.Blocks[msg.WorkingBlock.Index].Valid = true

					miners[miner.Connection].Mined++
				}
			}

		case payload := <-new_payload:
			fmt.Println("Got payload: ", payload)
			blockchain_length := len(blockchain.Blocks)
			var prev_hash string
			if blockchain_length == 0 {
				prev_hash = "0000000000000000000000000000000000000000000000000000000000000000"
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
						fmt.Println("client gob decode error: ", err)
						if err.Error() == "EOF" {
							break
						}
					}
					new_payload <- message
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
