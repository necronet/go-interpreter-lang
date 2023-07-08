package main

import("fmt"
	"os"
	"os/user"
	"necronet.info/interpreter/repl"
)

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is a random interpreter language!\n", user.Username)
	fmt.Printf("Type whatever command you feel like\n")
	repl.Start(os.Stdin, os.Stdout)

}
