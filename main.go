package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/blue0513/go-summarizer/openai"
	"github.com/blue0513/go-summarizer/parser"
	"github.com/blue0513/go-summarizer/request"
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

	text := parser.Extract(html)
	res, err := openai.Summarize(context.Background(), text, lang)
	if err != nil {
		fmt.Println("Error: summarize text", err)
		return
	}

	fmt.Println("Summarized......:\n\n", res)
}
