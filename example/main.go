package main

import (
	"cuslog"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Printf("level-1: %s\n", cuslog.CurLevel())

	cuslog.Debug("This is Debug info aaa")

	cuslog.SetOptions(cuslog.WithLevel(cuslog.InfoLevel))
	fmt.Printf("level-2: %s\n", cuslog.CurLevel())

	cuslog.Debug("Already show DEBUG message?")
	cuslog.Info("This is INFO level")

	fd, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("create file error")
	}
	defer fd.Close()

	l := cuslog.New(cuslog.WithLevel(cuslog.InfoLevel),
		cuslog.WithOutput(fd),
		cuslog.WithFormatter(&cuslog.TextFormatter{IgnoreBasicFields: false}),
	)

	l.Info("custom log ")

}
