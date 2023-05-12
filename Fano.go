package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	var dir, name, mode string

	fmt.Print("Выберите вид объекта (dir/file): ")
	_, err := fmt.Scan(&dir)
	if err != nil {
		log.Fatal(err)
	}
	switch dir {
	case "dir":
		fmt.Print("Введите название директории: ")
		_, err := fmt.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print("Выберите режим работы (encode/decode): ")
		_, err = fmt.Scan(&mode)
		if err != nil {
			log.Fatal(err)
		}

		err = workWithDirectory(name, mode)
		if err != nil {
			log.Fatal(err)
		}

	case "file":
		fmt.Print("Введите название файла: ")
		_, err := fmt.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}

		err = workWithFile(name, mode)
		if err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatal("Incorrect type of file")
	}
}

func workWithDirectory(name, mode string) error {

	files, err := os.ReadDir(name)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), "Matrix") || strings.HasSuffix(file.Name(), "FirstSymbol") {
			continue
		}
		filePath := filepath.Join(name, file.Name())
		fmt.Println(filePath)
		if file.IsDir() {
			err = workWithDirectory(filePath, mode)
			if err != nil {
				return err
			}
			continue
		}
		err = workWithFile(filePath, mode)

		if err != nil {
			return err
		}
	}

	return nil
}

func workWithFile(name, mode string) error {
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

	t := time.Now()
	ch := NewCodesHelpers(file)

	if mode == "encode" {
		err = ch.EncodeFile()
		if err != nil {
			return err
		}
	} else if mode == "decode" {
		err = ch.DecodeFile()
		if err != nil {
			return err
		}
	} else {
		return errors.New("incorrect mode")
	}
	fmt.Println("Время выполнения: ", time.Since(t))
	return nil
}
