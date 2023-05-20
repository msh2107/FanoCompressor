package compressor

import (
	"bufio"
	"fano-algorithm/pkg/FanoAlgorithm"
	"fano-algorithm/pkg/MarkovChain"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type FanoWithMarkovChainCompressor struct {
	Matrix      [][]string
	FirstSymbol []string
}

func NewFanoWithMarkovChainCompressor() *FanoWithMarkovChainCompressor {
	matrix := make([][]string, 256)
	for i := range matrix {
		matrix[i] = make([]string, 256)
	}
	return &FanoWithMarkovChainCompressor{
		Matrix:      matrix,
		FirstSymbol: make([]string, 256),
	}
}

func (ch *FanoWithMarkovChainCompressor) CreateCodes(mc *MarkovChain.MarkovChain) {
	for i := 0; i < len(mc.Matrix); i++ {
		encodedLine := FanoAlgorithm.FanoEncoding(mc.Matrix[i])
		for j := 0; j < len(encodedLine); j++ {
			ch.Matrix[i][encodedLine[j].Index] = encodedLine[j].Code
		}
	}

	encodedLine := FanoAlgorithm.FanoEncoding(mc.FirstSymbol)
	for _, c := range encodedLine {
		ch.FirstSymbol[c.Index] = c.Code
	}
}

func (ch *FanoWithMarkovChainCompressor) EncodeFile(file *os.File) error {
	_, err := file.Seek(0, 0)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)
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

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	err = file.Truncate(0)
	if err != nil {
		return err
	}
	codedText := sb.String()
	binaryData := make([]byte, (len(codedText)+7)/8)

	for i := 0; i < len(codedText); i++ {
		if codedText[i] == '1' {
			binaryData[i/8] |= 1 << (7 - i%8)
		}
	}
	_, err = file.Write([]byte{byte(8 - len(codedText)%8)})
	if err != nil {
		return err
	}
	_, err = file.Write(binaryData)
	if err != nil {
		return err
	}

	return nil
}

func (ch *FanoWithMarkovChainCompressor) SaveCodes() error {
	fileForMatrix, err := os.Create("matrix")
	defer func(fileForMatrix *os.File) {
		err := fileForMatrix.Close()
		if err != nil {
			return
		}
	}(fileForMatrix)

	for _, row := range ch.Matrix {
		_, err := fmt.Fprintln(fileForMatrix, row)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	fileForFirstSymbol, err := os.Create("firstSymbol")
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

func (ch *FanoWithMarkovChainCompressor) DecodeFile(file *os.File) error {
	_, err := file.Seek(0, 0)
	if err != nil {
		return err
	}
	firstSymbol := false
	var sb strings.Builder
	prev := 0
	reader := bufio.NewReader(file)
	cut, err := reader.ReadByte()
	cut %= 8
	if err != nil {
		return err
	}
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

	for i := 0; i < len(codedText)-int(cut); {
		j := i + 1
		for ; j <= len(codedText)-int(cut); j++ {
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

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	err = file.Truncate(0)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
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

func DeleteMCCodes() error {
	err := os.Remove("firstSymbol")
	if err != nil {
		return err
	}

	err = os.Remove("matrix")
	if err != nil {
		return err
	}
	return nil
}

func (ch *FanoWithMarkovChainCompressor) GetCodes() error {
	var matrix [][]string
	fileWithMatrix, err := os.Open("matrix")
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(fileWithMatrix)
	for scanner.Scan() {
		rowStr := scanner.Text()
		row := strings.Trim(rowStr, "[]")
		matrix = append(matrix, strings.Split(row, " "))
	}
	ch.Matrix = matrix

	err = fileWithMatrix.Close()
	if err != nil {
		return err
	}

	fileWithFirstSymbol, err := os.Open("firstSymbol")
	if err != nil {
		return err
	}

	scanner = bufio.NewScanner(fileWithFirstSymbol)
	for scanner.Scan() {
		firstSymbol := scanner.Text()
		firstSymbol = strings.Trim(firstSymbol, "[]")
		ch.FirstSymbol = strings.Split(firstSymbol, " ")
	}

	err = fileWithFirstSymbol.Close()
	if err != nil {
		return err
	}

	return nil
}

func (ch *FanoWithMarkovChainCompressor) findInCol(code string, row int) int {
	for i := 0; i < len(ch.Matrix[row]); i++ {
		if ch.Matrix[row][i] == code {
			return i
		}
	}
	return -1
}

func (ch *FanoWithMarkovChainCompressor) findInRow(code string) int {
	for i := 0; i < len(ch.FirstSymbol); i++ {
		if ch.FirstSymbol[i] == code {
			return i
		}
	}
	return -1
}
