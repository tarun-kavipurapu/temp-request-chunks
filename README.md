# temp-request-chunks
### P2P File Transfer System in Go


## Files and Their Functionalities

### 1. `main.go`

**Functionality:**
- Initializes and starts two peer servers.
- Simulates file requests from one peer to another using goroutines for parallel data transfer.
- Handles the main execution flow and synchronization using waitgroups.

**Key Functions:**
- `main()`: Sets up the peer servers and initiates file requests.

### 2. `encoding.go`

**Functionality:**
- Provides encoding and decoding functionality for messages using Go's `gob` package.
- Defines custom decoders for handling message reading.

**Key Functions:**
- `EncodeToGob(data interface{}) ([]byte, error)`: Encodes data into gob format.
- `DecodeFromGob(gobData []byte, result interface{}) error`: Decodes gob data into the provided interface.
- `DefaultDecoder`: Custom decoder implementation for reading messages.

### 3. `messages.go`

**Functionality:**
- Defines various message and data structures used in the P2P communication.

**Key Types:**
- `FileRegister`: Empty structure for future expansion related to file registration.
- `PeerRegistration`: Empty structure for peer registration.
- `DataMessage`: Generic structure for encapsulating different payloads.
- `RequestChunkData`: Structure for requesting chunk data by file ID.
- `TestSend`: Structure for testing data send, including chunk ID and chunk name.
- `DialOptions`: Structure for additional options when dialing a connection.

### 4. `tcp_transport.go`

**Functionality:**
- Implements the transport layer for handling peer connections and data transmission.
- Manages both client (dialing) and server (listening) functionalities.
- Handles message reading and writing, and synchronization using waitgroups.

**Key Functions:**
- `NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer`: Creates a new TCP peer.
- `Dial(addr string, connType string, opts DialOptions) error`: Dials a connection to a specified address.
- `ListenAndAccept(connType string, opts DialOptions) error`: Listens for incoming connections.
- `handleConn(conn net.Conn, outbound bool, connType string, opts DialOptions)`: Handles a new connection and manages incoming messages.

### 5. `transport.go`

**Functionality:**
- Defines interfaces for peer and transport components.
- Provides abstraction for different types of transport implementations.

**Key Interfaces:**
- `Peer`: Interface representing a peer connection.
- `Transport`: Interface representing transport functionality including dialing, listening, and consuming messages.

## How to Run

1. **Setup and Configuration:**
   - Ensure Go is installed on your system.
   - Clone the repository and navigate to the project directory.

2. **Build and Run:**
   - Use the `go run main.go` command to build and run the project.
   - The main function will initialize and start the peer servers and simulate file requests.

3. **Testing Parallel Requests:**
   - The `main.go` file includes code for testing parallel chunk requests using goroutines.
   - You can adjust the number of chunks or modify the goroutine logic for further testing.

## Additional Information

- The project uses Go's `sync.WaitGroup` for managing concurrency and synchronization.
- Error handling and logging are implemented to aid in debugging and understanding the flow of execution.
- The `DefaultDecoder` and other custom decoders ensure proper handling of incoming messages.

