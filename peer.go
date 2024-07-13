package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type PeerServer struct {
	peerLock    sync.RWMutex
	peers       map[string]Peer //conn
	Transport   Transport
	cServerAddr string
}

func NewPeerServer(opts TransportOpts) *PeerServer {
	transport := NewTCPTransport(opts)

	peerServer := &PeerServer{
		peers:     make(map[string]Peer),
		peerLock:  sync.RWMutex{},
		Transport: transport,
	}

	transport.OnPeer = peerServer.OnPeer

	return peerServer
}

func (p *PeerServer) OnPeer(peer Peer, connType string) error {
	p.peerLock.Lock()
	defer p.peerLock.Unlock()

	if connType == "peer-server" {
		p.peers[peer.RemoteAddr().String()] = peer
	} else {
		log.Printf("peer connected as a listener so do nothing")
	}

	log.Printf("connected with remote %s", peer.RemoteAddr())
	return nil
}

func (p *PeerServer) Start(info DialOptions) error {
	fmt.Printf("[%s] starting Peer ...\n", p.Transport.Addr())

	err := p.Transport.ListenAndAccept("peer-server", info)
	if err != nil {
		return err
	}

	p.loop()
	return nil
}

func (p *PeerServer) loop() {
	defer func() {
		log.Println("Central Server has stopped due to error or quit action")
		p.Transport.Close()
	}()
	for msg := range p.Transport.Consume() {
		// log.Println((msg.Payload))
		dataMsg := new(DataMessage)
		buf := bytes.NewReader(msg.Payload)
		// log.Println(buf)
		decoder := gob.NewDecoder(buf)
		if err := decoder.Decode(dataMsg); err != nil {
			log.Println("failed to decode data: %w", err)
		}

		if err := p.handleMessage(msg.From, dataMsg); err != nil {
			log.Println("Error handling Mesage:", err)
		}
	}
}

func (p *PeerServer) handleMessage(from string, dataMsg *DataMessage) error {
	switch v := dataMsg.Payload.(type) {
	case TestSend:
		return p.TestFileSend(from, v)
	}

	return nil
}

func (p *PeerServer) TestFileSend(from string, v TestSend) error {
	filePath := fmt.Sprintf("D:\\Devlopement\\p2p-file-transfer\\chunks\\017934bd678d02194f615f9e27cab2a72839fc28653daaf81ac7933f2945467d\\%s", v.ChunkName)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	fileInfo, err := file.Stat()

	p.peerLock.RLock()
	peer := p.peers[from]
	p.peerLock.RUnlock()

	peer.Send([]byte{IncomingStream})

	binary.Write(peer, binary.LittleEndian, fileInfo.Size())
	n, err := io.Copy(peer, file)
	if err != nil {
		return err
	}
	fmt.Printf("[%s] written (%d) bytes over the network to %s with chunk id (%d)\n", p.Transport.Addr(), n, from, v.ChunkId)

	return nil
}

func (p *PeerServer) TestFileRequest(i int, ipAddr string) error {

	dialOptsPeer := DialOptions{
		IntValue:  i,
		ExtraInfo: "Dialer",
	}

	// var dataMessage pkg.DataMessage
	err := p.Transport.Dial(ipAddr, "peer-server", dialOptsPeer)

	if err != nil {
		return err
	}

	// for index, element := range p.peers {
	// 	log.Println(index, element)
	// }

	dataMessage := DataMessage{
		Payload: TestSend{
			ChunkId:   i,
			ChunkName: fmt.Sprintf("chunk_%d.chunk", i),
		},
	}
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(dataMessage); err != nil {
		fmt.Printf("Failed to encode data: %v\n", err)

	}
	p.peerLock.RLock()
	peer, exists := p.peers[ipAddr]
	p.peerLock.RUnlock()

	if !exists {
		return fmt.Errorf("peer with address %s not found", ipAddr)
	}

	// Send the request

	if err := peer.Send(buf.Bytes()); err != nil {
		fmt.Printf("Failed to send data: %v\n", err)

	}

	// Wait for response stream
	time.Sleep(time.Millisecond * 500)

	var fileSize int64

	if err := binary.Read(peer, binary.LittleEndian, &fileSize); err != nil {
		fmt.Printf("Failed to read file size: %v\n", err)

	}

	file, err := os.Create(fmt.Sprintf("chunk_%d.chunk", i))
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)

	}
	defer file.Close()

	if _, err := io.CopyN(file, peer, fileSize); err != nil {
		fmt.Printf("Failed to copy file: %v\n", err)

	}

	fmt.Printf("Successfully received and saved chunk_%d.chunk\n", i)
	// p.peerLock.Lock()
	peer.CloseStream(i)
	// p.peerLock.Unlock()

	return nil

}
