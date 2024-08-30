package smd

import (
	"github.com/QinYuuuu/SS23/crypto/merkle"
	"log"
)

type receiver struct {
	n        int
	t        int
	id       int
	acquired bool
	voted    bool
	done     bool
	msgIn    chan Message // in SMD protocol, msg received
	msgOut   chan Message // in SMD protocol, msg sent

	root []byte

	echos   []bool       // store nodes have sent echo
	votes   []bool       // store nodes have sent vote
	echoNum int          // count echo msg
	voteNum int          // count vote msg
	sendCh  chan Message // store send msg
	echoCh  chan Message // store echo msg
	voteCh  chan Message // store vote msg
	closeCh chan bool
	log     log.Logger
}

// messages router
func (receiver *receiver) handleMsgIn() {
	receiver.log.Println("inside handleMsgIn")
	var msg Message
	for !receiver.done {
		msg = <-receiver.msgIn
		//map messages by type 1: Ready; 2: Echo; 3: CallHelp; 4: Help
		switch msg.Type {
		case 0:
			//get a send msg
			//receiver.log.Println("get a send msg from ", msg.DestID)
			if !receiver.acquired {

				receiver.sendCh <- msg
				receiver.acquired = true
			}
			break
		case 1:
			//get a echo msg
			//receiver.log.Println("get a echo msg from ", msg.DestID)
			if !receiver.voted {
				receiver.echoCh <- msg
			}
			break
		case 2:
			//get a vote msg
			if !receiver.voted {
				receiver.voteCh <- msg
			}
			break
		case 4:
		default:
			receiver.log.Fatalln("get a wrong type msg")
		}

	}

}

type echobuf struct {
	root       []byte
	witness_i  merkle.Witness
	root_i     []byte
	witness_ij merkle.Witness
	s_ij       []byte
	f_ij       []byte
}

func (receiver *receiver) handleEcho() {
	var msg Message
	for {
		select {
		case <-receiver.closeCh:
			return

		case msg = <-receiver.echoCh:
		}
		// get root from echo msg
		id := msg.Type
		if !receiver.echos[id] {
			receiver.echos[id] = true
			receiver.echoNum++
		} else {
			// receive duplicate message
			continue
		}
		if receiver.echoNum == receiver.n-receiver.t && !receiver.voted {
			// send vote msg
			receiver.voted = true
		}
	}

}

func (receiver *receiver) handleVote() {
	var msg Message
	for {
		select {
		case <-receiver.closeCh:
			return

		case msg = <-receiver.voteCh:
		}
		// get root from echo msg
		id := msg.Type
		if !receiver.votes[id] {
			receiver.votes[id] = true
			receiver.voteNum++
		} else {
			// receive duplicate message
			continue
		}
		if receiver.voteNum == receiver.t+1 && !receiver.voted {
			// send vote msg
			votemsg := Message{
				Type: VOTE,
			}
			receiver.voted = true
		}
		if receiver.voteNum == receiver.n-receiver.t && !receiver.done {
			// out put stage

		}
	}
}
