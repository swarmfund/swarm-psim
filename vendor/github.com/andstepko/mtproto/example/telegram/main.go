package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/andstepko/mtproto"
)

func usage() {
	fmt.Print("Telegram is a simple MTProto tool.\n\nUsage:\n\n")
	fmt.Print("    ./telegram <command> [arguments]\n\n")
	fmt.Print("The commands are:\n\n")
	fmt.Print("    auth  <phone_number>            auth connection by code\n")
	fmt.Print("    msg   <user_id> <msgtext>       send message to user\n")
	fmt.Print("    list                            get contact list\n")
	fmt.Println()
}

func main() {
	var err error

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	commands := map[string]int{"auth": 1, "msg": 2, "list": 0}
	valid := false
	for k, v := range commands {
		if os.Args[1] == k {
			if len(os.Args) < v+2 {
				usage()
				os.Exit(1)
			}
			valid = true
			break
		}
	}

	if !valid {
		usage()
		os.Exit(1)
	}

	f, err := os.OpenFile(os.Getenv("HOME") + "/.telegram_go", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		fmt.Printf("File opening failed: %s\n", err)
		os.Exit(2)
	}

	m, err := mtproto.NewMTProto(f)
	if err != nil {
		fmt.Printf("Create failed: %s\n", err)
		os.Exit(2)
	}

	err = m.Connect()
	if err != nil {
		fmt.Printf("Connect failed: %s\n", err)
		os.Exit(2)
	}

	switch os.Args[1] {
	case "auth":
		err = m.Auth(os.Args[2])
	case "msg":
		userID, _ := strconv.Atoi(os.Args[2])
		err = m.SendMessage(int32(userID), os.Args[3])

	case "list":
		err = m.GetContacts()
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
