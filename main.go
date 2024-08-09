package main

import (
	"flag"
	"fmt"
)

func main() {
	configFlag := flag.String("c", CONFIG_FILE, "Path to Configuration File")
	debugFlag := flag.Bool("d", false, "Show debug logs")
	flag.Parse()

	err := Setup(*configFlag, *debugFlag)
	if err != nil {
		fmt.Println(err)
		return
	}
}
