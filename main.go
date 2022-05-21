package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// FIXME: use flag package to parse command line arguments (path and word -> default path "." and word "TODO:")
	path := "."
	word := strings.ToLower("TODO:")

	// todo:

	counterFile := 0
	counterTODOs := 0
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && !strings.Contains(path, ".git") {
				counterFile++
				fmt.Println(counterFile, path, info.Size())
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				scanner.Split(bufio.ScanLines)

				counterLine := 0
				for scanner.Scan() {
					counterLine++
					if strings.Contains(strings.ToLower(scanner.Text()), word) { // make it upper-case for handle the case insensitive
						counterTODOs++
						fmt.Println("--", counterLine, strings.TrimSpace(scanner.Text())) // trim space for better output (no extra space after line)
					}
				}
			}
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("from %d files, you have %d TODOs on this directory. \n", counterFile, counterTODOs)
}
