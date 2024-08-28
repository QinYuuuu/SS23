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
	n     int
	t     int
	p     *big.Int
	coder reedsolomon.RScoder
}

func encrypt(k, m []byte) []byte {
	result := []byte("sa")
	return result
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
	root := make([][]byte, s.n)
	for i := 0; i < s.n; i++ {
		merkleInput := make([][]byte, 0)
		for j := 0; j < s.n; j++ {
			merkleInput[i] = append(k[i][j], chunk[i][j].Data...)
		}
		merkleTree, _ := merkle.NewMerkleTree(merkleInput, hasher.SHA256Hasher)
		root[i] = merkle.Commit(merkleTree)
	}
}
