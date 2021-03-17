package parser

import (
	"errors"
	"io"

	"golang.org/x/net/html"
)

// ParseZefr takes a file reader and parses the HTML
// it contains. The function will error if the parsed
// document doesn't contain the required elements, or
// if an unknown element is at the root of the file
func ParseZefr(f io.Reader) error {
	root, err := html.Parse(f)
	if err != nil {
		return err
	}

	// Parse automatically generates the HTML elements, so... dont do that?
	root = root.FirstChild.LastChild

	hasGoScript := false
	hasTemplate := false

	for el := root.FirstChild; el != nil; el = el.NextSibling {
		if el.Type == html.ElementNode {
			switch tag := el.Data; tag {
			case "go-script":
				hasGoScript = true
			case "template":
				hasTemplate = true
			case "style":
			default:
				return ZephyrUnsupportedElementErr{Element: el.Data}
			}

		} else if !(el.Type == html.TextNode || el.Data == "") {
			// not an ElementNode or empty line
			return ZephyrUnsupportedElementErr{Element: el.Data}
		}

	}

	if !hasGoScript {
		return errors.New("Must have <go-script> element")
	}

	if !hasTemplate {
		return errors.New("Must have <template> element")
	}

	return nil
}
