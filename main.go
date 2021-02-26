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

type Highlight struct {
	Book   string
	Author string
	Text   string
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

	highlights := make([]Highlight, 0)

	lines := make([]string, 0)

	var parsedHighlights int

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "==========" {
			parsedHighlights++

			book, author := parseBookAuthor(lines[0])
			if bookTitle != "" && book != bookTitle {
				lines = make([]string, 0)
				continue
			}

			h := Highlight{
				Book:   book,
				Author: author,
				Text:   lines[3],
			}

			highlights = append(highlights, h)
			lines = make([]string, 0)
		} else {
			lines = append(lines, line)
		}
	}

	fmt.Printf("parsed %d highlights\n", parsedHighlights)

	if len(highlights) == 0 {
		return errors.New("no highlights found for given book")
	}

	outFile, err := os.OpenFile(bookTitle+".txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer outFile.Close()

	fmt.Printf("writing %d highlights\n", len(highlights))

	outFile.Write([]byte(fmt.Sprintf("# %s\n", highlights[0].Book)))
	outFile.Write([]byte(fmt.Sprintf("## %s\n\n", highlights[0].Author)))

	for _, h := range highlights {
		outFile.Write([]byte(fmt.Sprintf("* %s\n", h.Text)))
	}

	fmt.Println("done")

	return nil
}

func parseBookAuthor(s string) (book, author string) {
	pattern := regexp.MustCompile(`^(?P<book>.*) \((?P<author>.*)\)`)
	match := pattern.FindStringSubmatch(s)
	return match[1], match[2]
}
