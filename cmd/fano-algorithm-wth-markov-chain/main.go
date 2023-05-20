package main

import (
	"fano-algorithm/internal/pkg/app"
	"fmt"
	"log"
)

func main() {
	var mode int
	fmt.Println("Выберите режим кодирования:\n1. Чистый Фано\n2. Фано на цепи Маркова")
	fmt.Print("Ваш выбор: ")
	_, err := fmt.Scan(&mode)
	if err != nil {
		log.Fatal(err)
	}
	switch mode {
	case 1:
		err = app.FanoClear()
		if err != nil {
			log.Fatal(err)
		}
	case 2:

		err = app.FanoWithMarkovChain()
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Неверный режим")
	}

}
