package main

import (
	"fmt"
	"os"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		return
	}
	destinationFile, err := os.Create(wd + "/test.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	fmt.Println("File created at:", destinationFile.Name())
}
