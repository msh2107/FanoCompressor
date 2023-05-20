package app

import (
	"errors"
	"fano-algorithm/internal/FrequencyCounter"
	"fano-algorithm/internal/compressor"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func FanoClear() error {
	var name, mode string

	fmt.Print("Введите название файла или директории: ")
	_, err := fmt.Scan(&name)
	if err != nil {
		return err
	}

	fmt.Print("Введите режим работы (encode/decode): ")
	_, err = fmt.Scan(&mode)
	if err != nil {
		return err
	}

	fileInfo, err := os.Stat(name)
	if err != nil {
		return err
	}
	fc := FrequencyCounter.NewFrequencyCounter()
	if fileInfo.IsDir() {
		if mode == "encode" {
			err = fc.CountDirectory(name)
			if err != nil {
				return err
			}
		}
		err = workWithDirectoryF(name, mode, fc)
		if err != nil {
			return err
		}

	} else {
		if mode == "encode" {
			err = fc.CountFile(name)
			if err != nil {
				return err
			}
		}
		err = workWithFileF(name, mode, fc)
		if err != nil {
			return err
		}
	}
	if mode == "decode" {
		err = compressor.DeleteFCodes()
		if err != nil {
			return err
		}
	}

	return nil
}

func workWithDirectoryF(name, mode string, fc *FrequencyCounter.FrequencyCounter) error {

	files, err := os.ReadDir(name)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(name, file.Name())
		fmt.Println(filePath)
		if file.IsDir() {
			err = workWithDirectoryF(filePath, mode, fc)
			if err != nil {
				return err
			}
			continue
		}
		err = workWithFileF(filePath, mode, fc)

		if err != nil {
			return err
		}
	}

	return nil
}

func workWithFileF(name, mode string, fc *FrequencyCounter.FrequencyCounter) error {
	file, err := os.OpenFile(name, os.O_RDWR, 0466)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	ch := compressor.NewFanoCompressor()

	if mode == "encode" {
		ch.CreateCodes(fc)
		err = ch.EncodeFile(file)
		if err != nil {
			return err
		}

		err = ch.SaveCodes()
		if err != nil {
			return err
		}

	} else if mode == "decode" {
		err = ch.GetCodes()
		if err != nil {
			return err
		}

		err = ch.DecodeFile(file)
		if err != nil {
			return err
		}

	} else {
		return errors.New("incorrect mode (only encode/decode)")
	}

	return nil
}
