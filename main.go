package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var name string
	var mode string

	fmt.Print("Введите имя файла: ")
	_, err := fmt.Scan(&name)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(name, os.O_RDWR, 0466)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	if err != nil {
		log.Fatal(err)
	}

	ch := NewCodesHelpers(file)

	flag.StringVar(&mode, "mode", "", "Выберите режим работы (encode/decode)")
	flag.Parse()

	if mode == "encode" {
		err = ch.EncodeFile()
		if err != nil {
			log.Fatal(err)
		}
	} else if mode == "decode" {
		err = ch.DecodeFile()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Incorrect mode")
	}
}
