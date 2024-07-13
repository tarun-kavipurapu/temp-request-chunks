package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
)

type Decoder interface {
	Reader(io.Reader, *msg) error
}

func EncodeToGob(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to encode data: %w", err)
	}
	return buf.Bytes(), nil
}

// DecodeFromGob decodes a gob byte slice into the provided interface
func DecodeFromGob(gobData []byte, result interface{}) error {
	buf := bytes.NewBuffer(gobData)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(result); err != nil {
		return fmt.Errorf("failed to decode data: %w", err)
	}
	return nil
}

type GOBDecoder struct{}

func (dec GOBDecoder) Decode(r io.Reader, msg *msg) error {
	return gob.NewDecoder(r).Decode(msg)
}

type DefaultDecoder struct {
}

func (dec DefaultDecoder) Reader(r io.Reader, msg *msg) error {
	bufReader := bufio.NewReader(r)

	peekBuf, err := bufReader.Peek(1)
	if err != nil {
		return err
	}

	if peekBuf[0] == IncomingStream {
		msg.Stream = true
		// Advance the reader past the peeked byte
		if _, err := bufReader.ReadByte(); err != nil {
			return err
		}
		return nil
	}

	buf := make([]byte, 1028)
	n, err := bufReader.Read(buf)
	if err != nil {
		return err
	}
	msg.Payload = buf[:n]

	return nil
}
