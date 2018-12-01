package main

import "github.com/mithrandie/readline-csvq"

func main() {
	if err := readline.DialRemote("tcp", ":12344"); err != nil {
		println(err.Error())
	}
}
