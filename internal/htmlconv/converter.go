package htmlconv

import (
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// Convert transforms HTML (including Confluence storage format) into GitHub-flavored Markdown.
func Convert(htmlStr string) string {
	if htmlStr == "" {
		return ""
	}

	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return htmlStr
	}

	var b strings.Builder
	walkNode(&b, doc)
	return strings.TrimSpace(b.String())
}

func walkNode(b *strings.Builder, n *html.Node) {
	switch n.Type {
	case html.TextNode:
		b.WriteString(n.Data)
		return
	case html.ElementNode:
		renderElement(b, n)
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkNode(b, c)
	}
}

func renderElement(b *strings.Builder, n *html.Node) {
	tag := n.Data

	// Handle Confluence storage format macros
	if tag == "ac:structured-macro" {
		renderMacro(b, n)
		return
	}
	if tag == "ac:link" {
		renderACLink(b, n)
		return
	}
	if tag == "ac:emoticon" {
		return // skip emoticons
	}
	if tag == "ac:plain-text-body" || tag == "ac:rich-text-body" {
		walkChildren(b, n)
		return
	}
	if strings.HasPrefix(tag, "ac:") || strings.HasPrefix(tag, "ri:") {
		walkChildren(b, n)
		return
	}

	switch tag {
	case "h1":
		b.WriteString("\n# ")
		walkChildren(b, n)
		b.WriteString("\n\n")
	case "h2":
		b.WriteString("\n## ")
		walkChildren(b, n)
		b.WriteString("\n\n")
	case "h3":
		b.WriteString("\n### ")
		walkChildren(b, n)
		b.WriteString("\n\n")
	case "h4":
		b.WriteString("\n#### ")
		walkChildren(b, n)
		b.WriteString("\n\n")
	case "h5":
		b.WriteString("\n##### ")
		walkChildren(b, n)
		b.WriteString("\n\n")
	case "h6":
		b.WriteString("\n###### ")
		walkChildren(b, n)
		b.WriteString("\n\n")
	case "p":
		walkChildren(b, n)
		b.WriteString("\n\n")
	case "br":
		b.WriteString("\n")
	case "strong", "b":
		b.WriteString("**")
		walkChildren(b, n)
		b.WriteString("**")
	case "em", "i":
		b.WriteString("*")
		walkChildren(b, n)
		b.WriteString("*")
	case "code":
		b.WriteString("`")
		walkChildren(b, n)
		b.WriteString("`")
	case "pre":
		b.WriteString("\n```\n")
		walkChildren(b, n)
		b.WriteString("\n```\n\n")
	case "a":
		href := getAttr(n, "href")
		b.WriteString("[")
		walkChildren(b, n)
		b.WriteString("](")
		b.WriteString(href)
		b.WriteString(")")
	case "ul":
		b.WriteString("\n")
		renderList(b, n, false)
		b.WriteString("\n")
	case "ol":
		b.WriteString("\n")
		renderList(b, n, true)
		b.WriteString("\n")
	case "li":
		walkChildren(b, n)
	case "table":
		renderTable(b, n)
	case "blockquote":
		var inner strings.Builder
		walkChildren(&inner, n)
		for _, line := range strings.Split(strings.TrimSpace(inner.String()), "\n") {
			b.WriteString("> ")
			b.WriteString(line)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	case "hr":
		b.WriteString("\n---\n\n")
	case "img":
		alt := getAttr(n, "alt")
		src := getAttr(n, "src")
		b.WriteString("![")
		b.WriteString(alt)
		b.WriteString("](")
		b.WriteString(src)
		b.WriteString(")")
	default:
		walkChildren(b, n)
	}
}

func renderMacro(b *strings.Builder, n *html.Node) {
	macroName := getAttr(n, "ac:name")

	switch macroName {
	case "code":
		lang := getMacroParam(n, "language")
		b.WriteString("\n```")
		b.WriteString(lang)
		b.WriteString("\n")
		body := getMacroBody(n)
		b.WriteString(body)
		if !strings.HasSuffix(body, "\n") {
			b.WriteString("\n")
		}
		b.WriteString("```\n\n")

	case "info", "note", "warning", "tip":
		label := strings.ToUpper(macroName[:1]) + macroName[1:]
		var inner strings.Builder
		walkChildren(&inner, n)
		text := strings.TrimSpace(inner.String())
		b.WriteString("> **")
		b.WriteString(label)
		b.WriteString(":** ")
		b.WriteString(text)
		b.WriteString("\n\n")

	case "toc":
		// Remove table of contents macro

	default:
		walkChildren(b, n)
	}
}

func renderACLink(b *strings.Builder, n *html.Node) {
	pageTitle := ""
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "ri:page" {
			pageTitle = getAttr(c, "ri:content-title")
		}
	}
	if pageTitle != "" {
		b.WriteString("*")
		b.WriteString(pageTitle)
		b.WriteString("*")
	} else {
		walkChildren(b, n)
	}
}

func renderList(b *strings.Builder, n *html.Node, ordered bool) {
	idx := 1
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "li" {
			if ordered {
				b.WriteString(strconv.Itoa(idx))
				b.WriteString(". ")
				idx++
			} else {
				b.WriteString("- ")
			}
			var inner strings.Builder
			walkChildren(&inner, c)
			b.WriteString(strings.TrimSpace(inner.String()))
			b.WriteString("\n")
		}
	}
}

func renderTable(b *strings.Builder, n *html.Node) {
	var rows [][]string
	collectRows(n, &rows)
	if len(rows) == 0 {
		return
	}

	// Compute column widths
	cols := 0
	for _, row := range rows {
		if len(row) > cols {
			cols = len(row)
		}
	}

	b.WriteString("\n")
	// Header row
	b.WriteString("|")
	for i := 0; i < cols; i++ {
		cell := ""
		if i < len(rows[0]) {
			cell = rows[0][i]
		}
		b.WriteString(" ")
		b.WriteString(cell)
		b.WriteString(" |")
	}
	b.WriteString("\n")

	// Separator
	b.WriteString("|")
	for i := 0; i < cols; i++ {
		b.WriteString("---|")
	}
	b.WriteString("\n")

	// Data rows
	for _, row := range rows[1:] {
		b.WriteString("|")
		for i := 0; i < cols; i++ {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			b.WriteString(" ")
			b.WriteString(cell)
			b.WriteString(" |")
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")
}

func collectRows(n *html.Node, rows *[][]string) {
	if n.Type == html.ElementNode && n.Data == "tr" {
		var cells []string
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && (c.Data == "td" || c.Data == "th") {
				var cellBuf strings.Builder
				walkChildren(&cellBuf, c)
				cells = append(cells, strings.TrimSpace(cellBuf.String()))
			}
		}
		*rows = append(*rows, cells)
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectRows(c, rows)
	}
}

func walkChildren(b *strings.Builder, n *html.Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkNode(b, c)
	}
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		attrKey := a.Key
		if a.Namespace != "" {
			attrKey = a.Namespace + ":" + a.Key
		}
		if attrKey == key {
			return a.Val
		}
	}
	return ""
}

func getMacroParam(n *html.Node, paramName string) string {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "ac:parameter" {
			name := getAttr(c, "ac:name")
			if name == paramName && c.FirstChild != nil {
				return c.FirstChild.Data
			}
		}
	}
	return ""
}

func getMacroBody(n *html.Node) string {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "ac:plain-text-body" {
			var b strings.Builder
			walkChildren(&b, c)
			return b.String()
		}
	}
	return ""
}

