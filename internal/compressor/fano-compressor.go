package compressor

import (
	"bufio"
	"fano-algorithm/internal/FrequencyCounter"
	"fano-algorithm/pkg/FanoAlgorithm"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type FanoCompressor struct {
	codes []string
}

func NewFanoCompressor() *FanoCompressor {
	return &FanoCompressor{
		codes: make([]string, 256),
	}
}

func (c *FanoCompressor) CreateCodes(fc *FrequencyCounter.FrequencyCounter) {
	encodedLine := FanoAlgorithm.FanoEncoding(fc.Counter)
	for _, code := range encodedLine {
		c.codes[code.Index] = code.Code
	}
}

func (c *FanoCompressor) EncodeFile(file *os.File) error {
	_, err := file.Seek(0, 0)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)
	if err != nil {
		return err
	}
	var sb strings.Builder

	for {
		letter, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		sb.WriteString(c.codes[letter])
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

func (c *FanoCompressor) DecodeFile(file *os.File) error {
	_, err := file.Seek(0, 0)
	if err != nil {
		return err
	}
	var sb strings.Builder
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

			symbol := c.find(code)
			if symbol != -1 {
				sb.WriteByte(byte(symbol))
				break
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

func (c *FanoCompressor) find(code string) int {
	for i := 0; i < len(c.codes); i++ {
		if c.codes[i] == code {
			return i
		}
	}
	return -1
}

func (c *FanoCompressor) SaveCodes() error {
	file, err := os.Create("frequency")
	defer func(fileForMatrix *os.File) {
		err := fileForMatrix.Close()
		if err != nil {
			return
		}
	}(file)

	_, err = fmt.Fprintln(file, c.codes)
	if err != nil {
		return err
	}
	return nil
}

func (c *FanoCompressor) GetCodes() error {
	file, err := os.Open("frequency")
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		firstSymbol := scanner.Text()
		firstSymbol = strings.Trim(firstSymbol, "[]")
		c.codes = strings.Split(firstSymbol, " ")
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func DeleteFCodes() error {
	err := os.Remove("frequency")
	if err != nil {
		return err
	}
	return nil
}
