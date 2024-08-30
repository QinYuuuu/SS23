package smd

import "log"

type receiver struct {
	id       int
	acquired bool
	voted    bool
	done     bool
	msgIn    chan Message // in SMD protocol, msg received
	msgOut   chan Message // in SMD protocol, msg sent

	root []byte

	echos   int          // count echo msg
	votes   int          // count vote msg
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

func (receiver *receiver) handleSend() {
	receiver.acquired = true
}

func (receiver *receiver) handleEcho() {
	receiver.voted = true
}

func (receiver *receiver) handleVote() {
	receiver.voted = false
}
