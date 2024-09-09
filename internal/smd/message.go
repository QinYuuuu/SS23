package smd

import "fmt"

const (
	SEND int = 0
	ECHO int = 1
	VOTE int = 2
)

type Message struct {
	FromID     int
	DestID     int
	InstanceID int
	Type       int
	SendBuf    sendBuf
	EchoBuf    echoBuf
	VoteBuf    voteBuf
}

// String formats the Message for debug output.
func (m Message) String() string {
	t := ""
	if m.Type == 0 {
		t += "Send"
	}
	if m.Type == 1 {
		t += "Echo"
	}
	if m.Type == 2 {
		t += "Vote"
	}
	return fmt.Sprintf("%v from node %d", t, m.FromID)
}
