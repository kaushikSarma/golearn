package main

import (
	"flag"
	"fmt"
)

func getDependencyListInFile() []string {
	return []string{}
}

func getFilesInDirectory() []string {
	return []string{}
}

func getDirectories() []string {
	return []string{}
}

func main() {
	sourceDir := flag.String("source", "", "")
	flag.Parse()
	if string(*sourceDir) == "" {
		fmt.Print("Please provide root directory of your JS project")
		return
	}
	fmt.Printf("JS Source: %v\n", string(*sourceDir))
}
