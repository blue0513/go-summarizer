package parser

import (
    "strings"

    "golang.org/x/net/html"
)

func ExtractText(n *html.Node) string {
    if n.Type == html.ElementNode && (n.Data == "style" || n.Data == "script") {
        return ""
    }

    if n.Type == html.TextNode {
        return n.Data
    }

    var text strings.Builder
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        text.WriteString(ExtractText(c))
    }
    return text.String()
}
