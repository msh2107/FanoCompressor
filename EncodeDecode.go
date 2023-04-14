package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type CodesHelpers struct {
	Matrix      [][]string
	FirstSymbol []string
	file        *os.File
}

func NewCodesHelpers(file *os.File) *CodesHelpers {
	mtrx := make([][]string, 256)
	for i := range mtrx {
		mtrx[i] = make([]string, 256)
	}
	return &CodesHelpers{
		Matrix:      mtrx,
		FirstSymbol: make([]string, 256),
		file:        file,
	}
}

func (ch *CodesHelpers) CreateCodes() {
	mc := NewMarkovChain()
	mc.AnalyzeText(ch.file)
	for i := 0; i < len(mc.Matrix); i++ {
		tmp := FanoEncoding(mc.Matrix[i])
		for j := 0; j < len(tmp); j++ {
			ch.Matrix[i][tmp[j].Index] = tmp[j].Code
		}
	}

	tmp := FanoEncoding(mc.FirstSymbol)
	for _, c := range tmp {
		ch.FirstSymbol[c.Index] = c.Code
	}
}

func (ch *CodesHelpers) EncodeFile() error {
	ch.CreateCodes()
	_, err := ch.file.Seek(0, 0)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(ch.file)
	prev, err := reader.ReadByte()
	if err != nil {
		return err
	}
	var sb strings.Builder
	sb.WriteString(ch.FirstSymbol[prev])

	for {
		letter, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		sb.WriteString(ch.Matrix[prev][letter])
		prev = letter
	}
	fmt.Println()

	_, err = ch.file.Seek(0, 0)
	if err != nil {
		return err
	}
	err = ch.file.Truncate(0)
	if err != nil {
		return err
	}
	codedText := sb.String()
	binBytes := make([]byte, (len(codedText)+7)/8)
	for i := 0; i < len(codedText); i++ {
		if codedText[i] == '1' {
			binBytes[i/8] |= 1 << uint(7-i%8)
		}
	}
	_, err = ch.file.Write(binBytes)
	if err != nil {
		return err
	}

	fileForMatrix, err := os.Create("matrix.txt")
	defer func(fileForMatrix *os.File) {
		err := fileForMatrix.Close()
		if err != nil {
			return
		}
	}(fileForMatrix)

	for _, row := range ch.Matrix {
		_, err := fmt.Fprintln(fileForMatrix, row)
		if err != nil {
			return err
		}
	}

	fileForFirstSymbol, err := os.Create("firstSymbol.txt")
	defer func(fileForMatrix *os.File) {
		err := fileForMatrix.Close()
		if err != nil {
			return
		}
	}(fileForMatrix)

	_, err = fmt.Fprintln(fileForFirstSymbol, ch.FirstSymbol)
	if err != nil {
		return err
	}

	return nil
}

func (ch *CodesHelpers) DecodeFile() error {

	err := ch.scanMatrix()
	if err != nil {
		return err
	}

	err = ch.scanFirstSymbol()
	if err != nil {
		return err
	}

	_, err = ch.file.Seek(0, 0)
	if err != nil {
		return err
	}
	firstSymbol := false
	var sb strings.Builder
	prev := 0
	reader := bufio.NewReader(ch.file)
	var buffer strings.Builder

	for {
		bits, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		for i := 7; i >= 0; i-- {
			mask := uint8(1) << uint8(i)
			bit := (bits & mask) >> uint8(i)
			buffer.WriteString(strconv.Itoa(int(bit)))

		}
	}

	codedText := buffer.String()

	for i := 0; i < len(codedText); {
		j := i + 1
		for ; j <= len(codedText); j++ {
			code := codedText[i:j]
			if !firstSymbol {
				symbol := ch.findInRow(code)
				if symbol != -1 {
					sb.WriteByte(byte(symbol))
					prev = symbol
					firstSymbol = true
					break
				}
			} else {
				symbol := ch.findInCol(code, prev)
				if symbol != -1 {
					sb.WriteByte(byte(symbol))
					prev = symbol
					break
				}
			}
		}

		i = j
	}

	_, err = ch.file.Seek(0, 0)
	if err != nil {
		return err
	}
	err = ch.file.Truncate(0)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(ch.file)
	_, err = writer.WriteString(sb.String())
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (ch *CodesHelpers) scanMatrix() error {
	fileWithMatrix, err := os.Open("matrix.txt")
	defer func(fileWithMatrix *os.File) {
		err := fileWithMatrix.Close()
		if err != nil {
			return
		}
	}(fileWithMatrix)
	if err != nil {
		return err
	}

	var matrix [][]string

	scanner := bufio.NewScanner(fileWithMatrix)
	for scanner.Scan() {
		rowStr := scanner.Text()
		row := strings.Trim(rowStr, "[]")
		matrix = append(matrix, strings.Split(row, " "))
	}
	ch.Matrix = matrix
	return nil
}

func (ch *CodesHelpers) scanFirstSymbol() error {
	fileWithFirstSymbol, err := os.Open("firstSymbol.txt")
	defer func(fileWithFirstSymbol *os.File) {
		err := fileWithFirstSymbol.Close()
		if err != nil {
			return
		}
	}(fileWithFirstSymbol)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(fileWithFirstSymbol)
	for scanner.Scan() {
		firstSymbol := scanner.Text()
		firstSymbol = strings.Trim(firstSymbol, "[]")
		ch.FirstSymbol = strings.Split(firstSymbol, " ")
	}

	return nil
}

func (ch *CodesHelpers) findInCol(code string, row int) int {
	for i := 0; i < len(ch.Matrix[row]); i++ {
		if ch.Matrix[row][i] == code {
			return i
		}
	}
	return -1
}

func (ch *CodesHelpers) findInRow(code string) int {
	for i := 0; i < len(ch.FirstSymbol); i++ {
		if ch.FirstSymbol[i] == code {
			return i
		}
	}
	return -1
}
