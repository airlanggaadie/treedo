package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type traversedPath struct {
	FilePath       string
	OccuranceCount uint64
	Occurances     []occurance
	Processed      bool
}

type occurance struct {
	LineOfCode uint64
	Text       string
}

func main() {
	// FIXME: use flag package to parse command line arguments (baseDirectory and word -> default baseDirectory "." and word "TODO:")
	var baseDirectory string
	flag.StringVar(&baseDirectory, "path", ".", "")
	var word string
	flag.StringVar(&word, "word", "TODO,FIXME,REVIEW,HACK,OPTIMIZE", "")

	var words = strings.Split(word, ",")
	var counterFile uint64
	var counterTODOs uint64
	var traversedPaths []traversedPath

	err := filepath.Walk(baseDirectory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && !strings.HasPrefix(path, ".") {
				counterFile++
				traversedPaths = append(traversedPaths, traversedPath{FilePath: path})
			}

			return nil
		})
	if err != nil {
		log.Fatal(err)
	}

	var output strings.Builder
	for i := 0; i < len(traversedPaths); i++ {
		err = processFile(words, &traversedPaths[i])
		if err != nil {
			log.Fatal(err)
		}

		if traversedPaths[i].OccuranceCount == 0 {
			continue
		}

		counterTODOs += traversedPaths[i].OccuranceCount

		output.WriteString("- ")
		output.WriteString(traversedPaths[i].FilePath)
		output.WriteString("\n")

		for _, occ := range traversedPaths[i].Occurances {
			output.WriteString("  :")
			output.WriteString(strconv.FormatUint(occ.LineOfCode, 10))
			output.WriteString(" ")
			output.WriteString(occ.Text)
			output.WriteString("\n")
		}

		output.WriteString("\n\n")
	}

	fmt.Printf("From %d files, you have %d TODOs on this directory.\n\n", counterFile, counterTODOs)
	fmt.Printf(output.String())
}

func processFile(words []string, t *traversedPath) error {
	file, err := os.Open(t.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lineOfCode uint64
	for scanner.Scan() {
		var found bool
		for i := 0; i < len(words); i++ {
			if strings.Contains(scanner.Text(), words[i]) {
				found = true
				break
			}
		}

		lineOfCode++

		if !found {
			continue
		}

		t.Occurances = append(t.Occurances, occurance{
			LineOfCode: lineOfCode,
			Text:       strings.TrimSpace(scanner.Text()),
		})

		t.OccuranceCount++
	}

	t.Processed = true

	return nil
}
