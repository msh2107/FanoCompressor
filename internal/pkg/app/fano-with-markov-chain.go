package app

import (
	"errors"
	"fano-algorithm/internal/compressor"
	"fano-algorithm/pkg/MarkovChain"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func FanoWithMarkovChain() error {
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

	mc := MarkovChain.NewMarkovChain()
	if fileInfo.IsDir() {
		if mode == "encode" {
			err = mc.AnalyzeDir(name)
			if err != nil {
				return err
			}
		}
		err = workWithDirectoryMC(name, mode, mc)
		if err != nil {
			return err
		}

	} else {
		if mode == "encode" {
			err = mc.AnalyzeFile(name)
			if err != nil {
				return err
			}
		}
		err = workWithFileMC(name, mode, mc)
		if err != nil {
			return err
		}
	}
	if mode == "decode" {
		err = compressor.DeleteMCCodes()
		if err != nil {
			return err
		}
	}

	return nil
}

func workWithDirectoryMC(name, mode string, mc *MarkovChain.MarkovChain) error {

	files, err := os.ReadDir(name)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := filepath.Join(name, file.Name())
		fmt.Println(filePath)
		if file.IsDir() {
			err = workWithDirectoryMC(filePath, mode, mc)
			if err != nil {
				return err
			}
			continue
		}
		err = workWithFileMC(filePath, mode, mc)

		if err != nil {
			return err
		}
	}

	return nil
}

func workWithFileMC(name, mode string, mc *MarkovChain.MarkovChain) error {
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

	ch := compressor.NewFanoWithMarkovChainCompressor()

	if mode == "encode" {
		ch.CreateCodes(mc)
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
