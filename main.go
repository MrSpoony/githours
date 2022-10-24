package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	out, err := exec.Command("git", "log", "--format='%at'").Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(out))
}
