package main

import "fmt"

func Debug(msg string) {
	fmt.Println("[DEBUG] => " + msg)
}

func Message(msg string) {
	fmt.Println("[MESSAGE] => " + msg)
}

func Warning(msg string) {
	fmt.Println("[WARNING] => " + msg)
}

func Error(msg error) {
	fmt.Println("[ERROR] => " + msg.Error())
}
