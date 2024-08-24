package smd

import (
	"github.com/QinYuuuu/SS23/crypto/hasher"
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

func encrypt(k, m []byte) []byte{

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
			k[i][j] = hasher.SHA256Hasher(append([]byte(strconv.Itoa(j)),sij.Bytes()...))
		}
		c[i] = k[i][0] ^ k[i][0]
	}
	chunk := make([][]infectious.Share, s.n)
	for i := 0; i < s.n; i++ {
		chunk[i] = s.coder.Encode(c[i])
	}
	merkleInput := make([][]byte, s.n)
	for i := 0; i < s.n; i++ {
		merkleInput[i] := make()
		for j := 0; j < s.n; j++ {

		}
		 := f[i][]
	}
}
