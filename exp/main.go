package main

import (
	"fmt"

	"github.com/shipsaw/scenario/rand"
)

func main() {
	fmt.Println(rand.String(10))
	fmt.Println(rand.RememberToken())
}
