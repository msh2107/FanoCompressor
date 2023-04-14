package main

import (
	"bufio"
	"os"
)

type MarkovChain struct {
	Matrix      [][]int
	FirstSymbol []int
}

func NewMarkovChain() *MarkovChain {
	mtrx := make([][]int, 256)
	for i := range mtrx {
		mtrx[i] = make([]int, 256)
	}
	return &MarkovChain{
		Matrix:      mtrx,
		FirstSymbol: make([]int, 256),
	}
}

func (mc *MarkovChain) AnalyzeText(file *os.File) {
	var letter byte
	reader := bufio.NewReader(file)
	prev, err := reader.ReadByte()
	mc.FirstSymbol[prev]++
	if err != nil {
		return
	}
	for {
		letter, err = reader.ReadByte()
		if err != nil {
			return
		}
		mc.Matrix[prev][letter]++
		prev = letter
		mc.FirstSymbol[prev]++
	}
}
