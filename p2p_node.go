package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var peers = make(map[string]bool)
var mu sync.Mutex

func sendMessageToPeers(message string, senderAddr string) {
	mu.Lock()
	defer mu.Unlock()
	for peerAddr := range peers {
		if peerAddr != senderAddr {
			go func(peer string) {
				conn, err := net.Dial("tcp", peer)
				if err != nil {
					log.Println("Failed to connect to peer:", peer, err)
					return
				}
				defer conn.Close()
				_, err = conn.Write([]byte(message))
				if err != nil {
					log.Println("Failed to send message to peer:", peer, err)
				}
			}(peerAddr)
		}
	}
}
func StartServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start server on port %s: %v", port, err)
	}
	defer listener.Close()

	fmt.Println("ServerStarted on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to acept connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	defer conn.Close()

	peerAddr := conn.RemoteAddr().String()
	mu.Lock()
	peers[peerAddr] = true
	mu.Unlock()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}
		fmt.Printf("Received message from: %s: %s", peerAddr, message)

		sendMessageToPeers(message, peerAddr)
	}
}

func ConncetToPeer(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Failed to connect to peer at %s: %v", address, err)
	}
	defer conn.Close()

	fmt.Printf("Connected to peer:", address)

	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Error reading peer message:", err)
				return
			}
			fmt.Println("Message from peer:", message)
		}
	}()

	userInput := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		message, _ := userInput.ReadString('\n')

		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Println("Failed to send message to peer:", err)
			return
		}
	}
}
func main() {
	go StartServer("8080")

	fmt.Print("Enter the address of peer to connect to (or leave empty): ")

	var peerAddress string
	fmt.Scanln(&peerAddress)

	if peerAddress != "" {
		go ConncetToPeer(peerAddress)
	}

	select {}
}
