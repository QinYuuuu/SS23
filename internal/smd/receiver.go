package smd

type receiver struct {
	id       int
	acquired bool
	voted    bool
	done     bool
}

func (receiver *receiver) handleSend() {
	receiver.acquired = true
}

func (receiver *receiver) handleEcho() {
	receiver.voted = true
}

func (receiver *receiver) handleVote() {
	receiver.voted = false
}
