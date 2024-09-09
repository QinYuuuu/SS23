package smd

import (
	"github.com/QinYuuuu/SS23/crypto/hasher"
	"github.com/QinYuuuu/SS23/crypto/merkle"
	"github.com/QinYuuuu/SS23/reedsolomon"
	"github.com/QinYuuuu/SS23/utils/polynomial"
	"github.com/vivint/infectious"
	"math/big"
	"strconv"
)

type Sender struct {
	id          int
	n           int
	t           int
	p           *big.Int
	coder       reedsolomon.RScoder
	sendChannel chan Message
}

func encrypt(k, m []byte) []byte {
	result := []byte("sa")
	return result
}

type sendBuf struct {
	root_i     [][]byte
	witness_ij []merkle.Witness
	s_ij       []*big.Int
	f_ij       []infectious.Share
}

func (s *Sender) Input(message [][]byte) {
	if len(message) != s.n {
		return
	}
	c := make([][]byte, s.n)
	k := make([][][]byte, s.n)
	for i := 0; i < s.n; i++ {
		f, _ := polynomial.NewRand(s.t, s.p)
		k[i] = make([][]byte, s.n)
		for j := 0; j < s.n; j++ {
			sij := f.EvalMod(new(big.Int).SetInt64(int64(j)), s.p)
			k[i][j] = hasher.SHA256Hasher(append([]byte(strconv.Itoa(j)), sij.Bytes()...))
		}
		c[i] = k[i][0]
	}
	chunk := make([][]infectious.Share, s.n)
	for i := 0; i < s.n; i++ {
		chunk[i] = s.coder.Encode(c[i])
	}
	// compute merkle tree with root r_i
	root := make([][]byte, s.n)
	proof := make([][]merkle.Witness, s.n)
	for i := 0; i < s.n; i++ {
		merkleInput := make([][]byte, 0)
		for j := 0; j < s.n; j++ {
			merkleInput[i] = append(k[i][j], chunk[i][j].Data...)
		}
		merkleTree, _ := merkle.NewMerkleTree(merkleInput, hasher.SHA256Hasher)
		root[i] = merkle.Commit(merkleTree)
		proof[i] = make([]merkle.Witness, s.n)
		for j := 0; j < s.n; j++ {
			proof[i][j], _ = merkle.CreateWitness(merkleTree, j)
		}
	}
	msgs := make([]Message, s.n)
	for i := 0; i < s.n; i++ {
		msgs[i] = Message{
			FromID: s.id,
			DestID: i,
			Type:   SEND,
		}
	}
}
