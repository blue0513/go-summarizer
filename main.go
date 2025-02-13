package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/blue0513/go-summarizer/openai"
	"github.com/blue0513/go-summarizer/parser"
	"github.com/blue0513/go-summarizer/request"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

func main() {
	var url, lang string
	flag.StringVar(&url, "url", "", "URL to summarize")
	flag.StringVar(&lang, "lang", "English", "Language for summarization")
	flag.Parse()

	if url == "" {
		fmt.Println("Error: URL is required")
		return
	}

	html, err := request.Fetch(url)
	if err != nil {
		fmt.Println("Error: parsing HTML:", err)
		return
	}

	page := parser.Extract(html)
	ctx := context.Background()
	s := spinner.New(spinner.CharSets[36], 100*time.Millisecond)

	s.Start()
	res, err := openai.Summarize(ctx, page, lang)
	s.Stop()

	if err != nil {
		fmt.Println("Error: summarize text", err)
		return
	}

	color.Green("---- Summary ----")
	fmt.Println(res)

	if err = openai.Ask(ctx, page); err != nil {
		log.Fatal(err)
	}
}
