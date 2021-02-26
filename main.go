package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

type Book struct {
	Title      string
	Author     string
	Highlights []string
}

func run() error {
	var in, bookTitle string
	flag.StringVar(&in, "in", "", "path to clippings file")
	flag.StringVar(&bookTitle, "book", "", "filter by book title")
	flag.Parse()

	f, err := os.Open(in)
	if err != nil {
		return err
	}
	defer f.Close()

	highlights := make(map[string]*Book, 0)

	lines := make([]string, 0)

	var parsedHighlights int

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "==========" {
			parsedHighlights++

			title, author := parseTitleAuthor(lines[0])
			if bookTitle != "" && title != bookTitle {
				lines = make([]string, 0)
				continue
			}

			bk, ok := highlights[title]
			if !ok {
				bk = &Book{
					Title:      title,
					Author:     author,
					Highlights: make([]string, 0),
				}
				highlights[title] = bk
			}

			bk.Highlights = append(bk.Highlights, lines[3])
			lines = make([]string, 0)

		} else {
			lines = append(lines, line)
		}
	}

	fmt.Printf("parsed %d highlights\n\n", parsedHighlights)

	if len(highlights) == 0 {
		return errors.New("no highlights found for given book")
	}

	var filename string
	if bookTitle == "" {
		filename = "highlights.txt"
	} else {
		filename = bookTitle + ".txt"
	}
	outFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer outFile.Close()

	for _, bk := range highlights {
		fmt.Printf("+ %s (%d highlights)\n", bk.Title, len(bk.Highlights))

		outFile.Write([]byte(fmt.Sprintf("# %s\n", bk.Title)))
		outFile.Write([]byte(fmt.Sprintf("## %s\n\n", bk.Author)))
		for _, h := range bk.Highlights {
			outFile.Write([]byte(fmt.Sprintf("* %s\n", h)))
		}
	}

	fmt.Println("done")

	return nil
}

func parseTitleAuthor(s string) (book, author string) {
	pattern := regexp.MustCompile(`^(?P<book>.*) \((?P<author>.*)\)`)
	match := pattern.FindStringSubmatch(s)
	return match[1], match[2]
}
