# Simple Blockchain Simulation

This project demonstrates a simple simulation of Blockchain written in Go.

## Prerequisites:
* Go
* Bash (only for the test script)

## Architecture

[TODO: Click here]()

## Project Description

### Quick test run

The script provided will start off $N_MINER miners as background processes and will execute the client to send all the test payload data in the 'test-payload-data' directory.

Run it like:

> bash run-test.sh #default $N_MINER = 2

You may supply an argument to specify the number of miners which will be called.

> bash run-test.sh -m 5

Ctrl+C should kill the server which in turn will send a signal to stop all the miners.

### Server

The server holds the blockchain data structure and will communicate the miners and clients about blocks which need to be mined. It is concurrent and can handle multiple miners and clients.

Run it like

> go run server.go

or

> go run server.go -target 00000000

There are optional arguments which can be passed as well.
```
Usage of server:
  -cport int
    	Port on which server listens on client (default 9875)
  -mport int
    	Port on which server listens on for miner (default 9876)
  -target string
    	Difficulty which should be targeted (default "000")
```

### Miner

The miner will connect to the server and will receive message structs from which it will hash the payload and send back the hash and nonce.

Run it like:

> go run miner.go

or

> go run miner.go -dest server.tld -port 8888

There are optional arguments which can be passed as well.

```
Usage of miner:
  -dest string
    	Address of the server (default "localhost")
  -port int
    	Port used to connect to server (default 9876)
```

### Client

The client will send payload messages to the server in order to be hashed and added to the blockchain.

Run it like:

> go run client.go

or

> go run client.go -payload "$(curl -s https://www.gnu.org/licenses/gpl-3.0.txt)"

There are optional arguments which can be passed as well.

```
Usage of client:
  -addr string
    	Address of the server (default "localhost")
  -payload string
    	Payload to add to blockchain (default "test")
  -port int
    	Port on which the server listens on (default 9875)
```

## My development environment

* Debian sid
* go version go1.9.2 linux/amd64
