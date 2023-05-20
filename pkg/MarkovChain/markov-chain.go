package MarkovChain

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

type MarkovChain struct {
	Matrix      [][]int
	FirstSymbol []int
}

func NewMarkovChain() *MarkovChain {
	matrix := make([][]int, 256)
	for i := range matrix {
		matrix[i] = make([]int, 256)
	}
	return &MarkovChain{
		Matrix:      matrix,
		FirstSymbol: make([]int, 256),
	}
}

func (mc *MarkovChain) AnalyzeDir(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())
		if file.IsDir() {

			err = mc.AnalyzeDir(filePath)
			if err != nil {
				return err
			}
			continue
		}
		err = mc.AnalyzeFile(filePath)

	}
	return nil
}

func (mc *MarkovChain) AnalyzeFile(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0466)
	if err != nil {
		return err
	}
	var letter byte
	reader := bufio.NewReader(file)
	prev, err := reader.ReadByte()
	mc.FirstSymbol[prev]++
	if err != nil {
		return err
	}
	for {
		letter, err = reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		mc.Matrix[prev][letter]++
		prev = letter
		mc.FirstSymbol[prev]++
	}

	return nil
}
