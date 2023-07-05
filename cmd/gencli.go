package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alviezhang/gensub/internal"
)

func main() {
	subtype := flag.String("s", "", "Subscription type of the generated file")
	filename := flag.String("f", "", "Subscription type of the generated file")

	flag.Parse()

	result, err := internal.Generate(*filename, *subtype)
	if err != nil {
		fmt.Println("Usage: gencli -s <subtype> -f <filename>")
		os.Exit(1)
	}

	fmt.Println(result)
}
