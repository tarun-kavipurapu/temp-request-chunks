package main

type FileRegister struct {
}

// it is empty because i dont need to send any additional data for peer registration to central server if i need i will add it later
type PeerRegistration struct {
}
type DataMessage struct {
	Payload any
}

type RequestChunkData struct {
	FileId string
}
type TestSend struct {
	ChunkId   int
	ChunkName string
}
type DialOptions struct {
	IntValue  int
	ExtraInfo string
}
