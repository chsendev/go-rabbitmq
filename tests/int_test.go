package tests

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
)

// The Message interface
type Message interface {
	Process()
	Convert([]byte) (Message, error)
}

// The first message type
type Message1 struct {
	Field1 string
	Field2 int
}

func (m Message1) Process() {
	fmt.Printf("Processing Message1: %+v\n", m)
}
func (m Message1) Convert(data []byte) (Message, error) {
	var msg Message1
	err := json.Unmarshal(data, &msg)
	return msg, err
}

// The second message type
type Message2 struct {
	Field3 float64
	Field4 bool
}

func (m Message2) Process() {
	fmt.Printf("Processing Message2: %+v\n", m)
}
func (m Message2) Convert(data []byte) (Message, error) {
	var msg Message2
	err := json.Unmarshal(data, &msg)
	return msg, err
}

// The interface for converting []byte to a specific type
type ByteConverter interface {
	Convert([]byte) (Message, error)
}

// An implementation of ByteConverter for Message1
type Message1Converter struct{}

func (mc Message1Converter) Convert(data []byte) (Message, error) {
	var msg Message1
	err := json.Unmarshal(data, &msg)
	return msg, err
}

// An implementation of ByteConverter for Message2
type Message2Converter struct{}

func (mc Message2Converter) Convert(data []byte) (Message, error) {
	var msg Message2
	err := json.Unmarshal(data, &msg)
	return msg, err
}

// The type of your callback function
type Callback func(msg Message)

func TestB(t *testing.T) {
	// Your callback function
	myCallback := func(msg Message) {
		msg.Process()
	}

	// Your function for receiving messages from RabbitMQ
	receiveMessage := func(converter Message, callback Callback) {
		// Simulate receiving a message from RabbitMQ
		data1 := []byte(`{"Field1":"hello", "Field2":123}`)
		data2 := []byte(`{"Field3":3.14, "Field4":true}`)

		// Convert the []byte data to a Message
		msg1, err := converter.Convert(data1)
		if err != nil {
			fmt.Printf("Failed to convert message: %v\n", err)
			return
		}

		msg2, err := converter.Convert(data2)
		if err != nil {
			fmt.Printf("Failed to convert message: %v\n", err)
			return
		}

		// Call the callback function with the Message
		callback(msg1)
		callback(msg2)
	}

	// Use the functions
	receiveMessage(Message1{}, myCallback)
	receiveMessage(Message2{}, myCallback)
}

func TestAbc(t *testing.T) {
	var poll sync.Pool

	poll.New = func() any {
		return "123"
	}

	a := poll.Get()
	fmt.Println(a)
}
