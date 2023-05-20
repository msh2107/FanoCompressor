package FrequencyCounter

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

type FrequencyCounter struct {
	Counter []int
}

func NewFrequencyCounter() *FrequencyCounter {
	return &FrequencyCounter{
		Counter: make([]int, 256),
	}
}

func (fc *FrequencyCounter) CountFile(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY, 0466)
	if err != nil {
		return err
	}
	var b byte
	reader := bufio.NewReader(file)
	for {
		b, err = reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		fc.Counter[b]++
	}
	return nil
}

func (fc *FrequencyCounter) CountDirectory(name string) error {
	files, err := os.ReadDir(name)
	if err != nil {
		return err
	}
	for _, file := range files {
		filePath := filepath.Join(name, file.Name())
		if file.IsDir() {

			err = fc.CountDirectory(filePath)
			if err != nil {
				return err
			}
			continue
		}
		err = fc.CountFile(filePath)

	}
	return nil
}
