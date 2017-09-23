package main

import (
	"flag"
	"fmt"
	"os"

	"contnet.org/lib/cnm-go"
	"contnet.org/lib/cnm-go/cnmfmt"
)

func main() {
	var filename string
	// future flags here
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("ERROR missing CNM filepath to load")
		os.Exit(2)
	} else {
		filename = flag.Args()[0]
	}

	// open file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR opening file:", err)
		return
	}
	defer file.Close()

	// parse document
	doc, err := cnm.ParseDocument(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR parsing file:", err)
	}

	// display metadata and contents
	fmt.Printf("Title:\t%s\n", doc.Title)
	//TODO how is it marked that there is no sitemap?
	if doc.Site.Path != "" {
		fmt.Println("Sitemap:", doc.Site)
		doc.Site.WriteIndent(os.Stdout, 1)
	}
	//TODO is Links nil if no links?
	if doc.Links != nil && len(doc.Links) > 0 {
		fmt.Println("Links:")
		for _, link := range doc.Links {
			fmt.Printf("\tURL=%s Name=%s Description=%s\n", link.URL, link.Name, link.Description)
		}
	}
	pushContent(doc.Content)
}

func pushContent(content cnm.Block) {
	// content blocks
	switch block := content.(type) {
	// these are all content blocks
	case *cnm.ContentBlock:
		fmt.Printf("Content: Name=%s Args=%v\n", block.Name(), block.Args())
		// is container block
	case *cnm.EmbedBlock:
		fmt.Printf("Embed: Name=%s Args=%v URL=%s Description=%s Type=%s\n", block.Name(), block.Args(), block.URL, block.Description, block.Type)
	case *cnm.HeaderBlock: // table header row
		fmt.Printf("Header: Name=%s Args=%v\n", block.Name(), block.Args())
		// is container block
	case *cnm.ListBlock:
		fmt.Printf("List: Name=%s Args=%v Ordered=%t\n", block.Name(), block.Args(), block.Ordered())
		// is container block
	case *cnm.RawBlock:
		fmt.Printf("Raw text: Name=%s Args=%s Syntax=%s Text=%s\n", block.Name(), block.Args(), block.Syntax, block.Contents.Text)
	case *cnm.RowBlock:
		fmt.Printf("Row: Name=%s Args=%v\n", block.Name(), block.Args())
		// is container block
	case *cnm.SectionBlock:
		fmt.Printf("Section: Name=%s Args=%v, Title=%s\n", block.Name(), block.Args(), block.Title())
		// is container block
	case *cnm.TableBlock:
		fmt.Printf("Table: Name=%s Args=%v\n", block.Name(), block.Args())
		// is also a container block
	case *cnm.TextBlock:
		fmt.Printf("Text: Name=%s Args=%v, Format=%s:\n", block.Name(), block.Args(), block.Format)
		switch textContent := block.Contents.(type) {
		case cnmfmt.TextFmtContents:
			fmt.Printf("\tFmt: ")
			textContent.WriteIndent(os.Stdout, 0)
		case *cnm.TextPlainContents:
			for _, paragraph := range textContent.Paragraphs {
				fmt.Printf("\tPlain paragraph: %s\n", paragraph)
			}
		case *cnm.TextPreContents:
			fmt.Printf("\tPreformatted: %s\n", textContent.Text)
		default:
			panic(fmt.Sprintf("Unknown text block type %T", textContent))
		}
	default:
		panic(fmt.Sprintf("Unknown content block type %T", block))
	}

	// is it also a container block = does it have children?
	if block, isContainerBlock := content.(cnm.ContainerBlock); isContainerBlock {
		//TODO these checks necessary?
		if block.Children() != nil && len(block.Children()) > 0 {
			for _, child := range block.Children() {
				pushContent(child)
			}
		}
		return
	}
}
