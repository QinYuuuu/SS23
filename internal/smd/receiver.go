package smd

import (
	"github.com/QinYuuuu/SS23/crypto/hasher"
	"github.com/QinYuuuu/SS23/crypto/merkle"
	"github.com/vivint/infectious"
	"log"
	"math/big"
	"strconv"
)

type receiver struct {
	n        int // node number
	t        int // threshold
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

func (receiver *receiver) Run() {
	go receiver.handleMsgIn()
	go receiver.handleEcho()
	go receiver.handleVote()
}

func (receiver *receiver) sendMsg(msgs []Message) {
	for _, msg := range msgs {
		if msg.DestID == receiver.id {
			receiver.msgIn <- msg
		} else {
			receiver.msgOut <- msg
		}
	}
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
			//receiver.log.Println("get a vote msg from ", msg.DestID)
			receiver.voteCh <- msg

			break
		case 4:
		default:
			receiver.log.Fatalln("get a wrong type msg")
		}

	}

}

func (receiver *receiver) handleSend() {
	var msg Message
	for {
		select {
		case <-receiver.closeCh:
			return
		case msg = <-receiver.sendCh:
			break
		}
		buf := msg.SendBuf
		var b bool
		for i := 0; i < receiver.n; i++ {
			content := hasher.SHA256Hasher(append([]byte(strconv.Itoa(receiver.id)), buf.s_ij[i].Bytes()...))
			content = append(content, buf.f_ij[i].Data...)
			b, _ = merkle.Verify(buf.root_i[i], buf.witness_ij[i], content, hasher.SHA256Hasher)
			if !b {
				break
			}
		}
		if b {
			// build Merkle tree and send ECHO message
			merkleTree, _ := merkle.NewMerkleTree(buf.root_i, hasher.SHA256Hasher)
			root := merkle.Commit(merkleTree)
			proof := make([]merkle.Witness, receiver.n)
			for i := 0; i < receiver.n; i++ {
				proof[i], _ = merkle.CreateWitness(merkleTree, i)
			}
			msgs := make([]Message, receiver.n)
			for i := 0; i < receiver.n; i++ {
				echobuf := EchoBuf{
					Root:       root,
					Witness_i:  proof[i],
					Root_i:     buf.root_i[i],
					Witness_ij: buf.witness_ij[i],
					S_ij:       buf.s_ij[i],
					F_ij:       buf.f_ij[i],
				}
				msgs[i] = Message{
					Type:    ECHO,
					FromID:  receiver.id,
					DestID:  i,
					EchoBuf: echobuf,
				}

			}
			receiver.sendMsg(msgs)
		}

	}
}

type EchoBuf struct {
	Root       []byte
	Witness_i  merkle.Witness
	Root_i     []byte
	Witness_ij merkle.Witness
	S_ij       *big.Int
	F_ij       infectious.Share
}

func (receiver *receiver) handleEcho() {
	var msg Message
	for {
		select {
		case <-receiver.closeCh:
			return
		case msg = <-receiver.echoCh:
			break
		}
		// get root from echo msg
		id := msg.FromID
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

type voteBuf struct {
	root []byte
}

func (receiver *receiver) handleVote() {
	var msgReceived Message
	for {
		select {
		case <-receiver.closeCh:
			return

		case msgReceived = <-receiver.voteCh:
		}
		// get root from echo msg
		id := msgReceived.FromID
		if !receiver.votes[id] {
			receiver.votes[id] = true
			receiver.voteNum++
		} else {
			// receive duplicate message
			continue
		}
		if receiver.voteNum == receiver.t+1 && !receiver.voted {
			// broadcast vote msg
			msgs := make([]Message, receiver.n)
			for i := 0; i < receiver.n; i++ {
				msgs[i] = Message{
					Type:       VOTE,
					FromID:     receiver.id,
					DestID:     i,
					InstanceID: receiver.id,
					VoteBuf:    msgReceived.VoteBuf,
				}
			}
			receiver.voted = true
			receiver.sendMsg(msgs)
		}
		if receiver.voteNum == receiver.n-receiver.t && receiver.voted && !receiver.done {
			// out put stage
			// reconstruct
		}
	}
}
