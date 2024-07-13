package main

import (
	"encoding/gob"
	"log"
	"os"
	"sync"
	"time"
)

func init() {
	gob.Register(TestSend{})

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// Initialize Transport Options for both peer servers
	opts1 := TransportOpts{
		ListenAddr: "localhost:5002",
		Decoder:    &DefaultDecoder{}, // Assume GobDecoder is defined elsewhere
	}
	opts2 := TransportOpts{
		ListenAddr: "localhost:5003",
		Decoder:    &DefaultDecoder{},
	}

	// Create two PeerServer instances
	peerServer1 := NewPeerServer(opts1)
	peerServer2 := NewPeerServer(opts2)

	dialOptsPeer1 := DialOptions{
		IntValue:  0,
		ExtraInfo: "Peer1  listener",
	}
	dialOptsPeer2 := DialOptions{
		IntValue:  0,
		ExtraInfo: "Peer2  listener",
	}

	// Start both PeerServers
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := peerServer1.Start(dialOptsPeer1); err != nil {
			log.Fatalf("Failed to start peer server 1: %v\n", err)
		}

	}()
	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := peerServer2.Start(dialOptsPeer2); err != nil {
			log.Fatalf("Failed to start peer server 2: %v\n", err)
		}

	}()

	time.Sleep(2 * time.Second)

	// Wait for the signal that both servers are ready<-ready1

	// Simulate a file request from peerServer1 to peerServer2
	var wg1 sync.WaitGroup

	//only for loop without go  routines runs the code perfectly as intended
	for i := 1; i <= 10; i++ {
		wg1.Add(1)
		go func(i int) {
			defer wg1.Done()
			if err := peerServer1.TestFileRequest(i, "127.0.0.1:5003"); err != nil {
				log.Printf("Failed to request file from peerServer2: %v\n", err)
			}

		}(i)
	}

	wg1.Wait()
	wg.Wait()
	log.Println("Servers shut down successfully.")
}
