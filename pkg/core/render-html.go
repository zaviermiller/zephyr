// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zephyr

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

type writer interface {
	io.Writer
	io.ByteWriter
	WriteString(string) (int, error)
}

// Renders HTML string from VNode tree. This file is
// extended directly from the net/http package; its
// documentation is below.
// TODO: namespaces
//
// Render renders the parse tree n to the given writer.
//
// Rendering is done on a 'best effort' basis: calling Parse on the output of
// Render will always result in something similar to the original tree, but it
// is not necessarily an exact clone unless the original tree was 'well-formed'.
// 'Well-formed' is not easily specified; the HTML5 specification is
// complicated.
//
// Calling Parse on arbitrary input typically results in a 'well-formed' parse
// tree. However, it is possible for Parse to yield a 'badly-formed' parse tree.
// For example, in a 'well-formed' parse tree, no <a> element is a child of
// another <a> element: parsing "<a><a>" results in two sibling elements.
// Similarly, in a 'well-formed' parse tree, no <a> element is a child of a
// <table> element: parsing "<p><table><a>" results in a <p> with two sibling
// children; the <a> is reparented to the <table>'s parent. However, calling
// Parse on "<a><table><a>" does not return an error, but the result has an <a>
// element with an <a> child, and is therefore not 'well-formed'.
//
// Programmatically constructed trees are typically also 'well-formed', but it
// is possible to construct a tree that looks innocuous but, when rendered and
// re-parsed, results in a different tree. A simple example is that a solitary
// text node would become a tree containing <html>, <head> and <body> elements.
// Another example is that the programmatic equivalent of "a<head>b</head>c"
// becomes "<html><head><head/><body>abc</body></html>".
func RenderHTML(w io.Writer, n *VNode) error {
	if x, ok := w.(writer); ok {
		return render(x, n)
	}
	buf := bufio.NewWriter(w)
	if err := render(buf, n); err != nil {
		return err
	}
	return buf.Flush()
}

// plaintextAbort is returned from render1 when a <plaintext> element
// has been rendered. No more end tags should be rendered after that.
var plaintextAbort = errors.New("html: internal error (plaintext abort)")

func render(w writer, n *VNode) error {
	err := render1(w, n)
	if err == plaintextAbort {
		err = nil
	}
	return err
}

func (node *VNode) parseAttrs() (stringAttrs map[string]string, err error) {
	stringAttrs = map[string]string{}
	attrListener := node.GetOrCreateListener("attr")
	for k, v := range node.Attrs {
		switch v.(type) {
		case LiveData:
			stringAttrs[k] = v.(LiveData).string(attrListener)
			// calculated functions get computed and results parsed
		case func(Listener) interface{}:
			eval := v.(func(Listener) interface{})(attrListener)
			// parse result
			switch eval.(type) {
			case LiveData:
				stringAttrs[k] = eval.(LiveData).string(attrListener)
			case string:
				stringAttrs[k] = eval.(string)
			default:
				// make error better
				return nil, errors.New("Must either use LiveData, string, or calculated func that returns one of those.")
			}
		case string:
			stringAttrs[k] = v.(string)
		default:
			return nil, errors.New("Must either use LiveData, string, or calculated func that returns one of those.")
		}
		if k == "class" {
			stringAttrs[k] += " " + node.DOM_ID
		}
	}
	if _, ok := stringAttrs["class"]; !ok {
		stringAttrs["class"] = node.DOM_ID
	}
	return stringAttrs, nil
}

func (node *VNode) parseContent() (parsedContent string, err error) {
	contentListener := node.GetOrCreateListener("content")
	switch node.Content.(type) {
	case LiveData:
		parsedContent = node.Content.(LiveData).string(contentListener)
	// calculated data
	case func(Listener) interface{}:
		evaluated := node.Content.(func(Listener) interface{})(contentListener)
		switch evaluated.(type) {
		case string:
			parsedContent = evaluated.(string)
		default:
			return "", errors.New(node.DOM_ID + " content func return type not supported by render (must be string!): ")
		}
	case string:
		parsedContent = node.Content.(string)
	default:
		// return error lol
		return "", errors.New(node.DOM_ID + " content type not supported by render (must be string!): ")
	}
	return parsedContent, nil
}

func (n *VNode) parseConditional() error {
	conditionalListener := n.GetOrCreateListener("conditional")
	for i, cr := range n.ConditionalRenders {
		condition := cr.Condition
		cr.Render.Parent = n
		switch condition.(type) {
		case bool:
			if condition.(bool) {
				n.FirstChild = cr.Render
				n.ConditionUpdated = n.CurrentCondition != i
				n.CurrentCondition = i
				goto addClass
			}
		case LiveBool:
			cBool := condition.(LiveBool).Value(conditionalListener).(bool)
			if cBool {
				n.FirstChild = cr.Render
				// fmt.Println(n.FirstChild)
				n.ConditionUpdated = n.CurrentCondition != i
				n.CurrentCondition = i
				goto addClass
			}
		case func(l Listener) interface{}:
			eval := condition.(func(l Listener) interface{})(conditionalListener)
			if b, ok := eval.(bool); ok {
				if b {
					n.FirstChild = cr.Render
					n.ConditionUpdated = n.CurrentCondition != i
					n.CurrentCondition = i
					goto addClass
				}
			}
		default:
			return errors.New("conditional must use bool, live bool, or calculated function that returns a bool")
		}
	}
	n.FirstChild = nil
	n.ConditionUpdated = true
	n.CurrentCondition = -1
	n.Tag = ""
	return nil
addClass:
	if _, ok := n.FirstChild.Attrs["class"]; !ok {
		if n.FirstChild.Attrs == nil {
			n.FirstChild.Attrs = map[string]interface{}{
				"class": n.DOM_ID,
			}
		} else {
			n.FirstChild.Attrs["class"] = n.DOM_ID
		}
	}
	// fmt.Println("here right: ", n.FirstChild.Attrs)
	return nil
}

func render1(w writer, n *VNode) error {
	// if !n.Component {
	attrs, err := n.parseAttrs()
	if err != nil {
		return err
	}
	n.ParsedAttrs = attrs
	// if n.events != nil {
	// 	n.RenderChan <- DOMUpdate{Operation: AddEventListeners, ElementID: n.DOM_ID, Data: n.events}
	// }
	// } else {
	// 	n.ParsedAttrs = map[string]string{
	// 		"class": "" + n.DOM_ID,
	// 	}
	// }
	// Render non-element nodes; these are the easy cases.
	switch n.NodeType {
	case ErrorNode:
		return errors.New("html: cannot render an ErrorNode node")
	case TextNode:
		if n.Static {
			return escape(w, n.Content.(string))
		}
		parsed, err := n.parseContent()
		if err != nil {
			return err
		}
		n.Tag = parsed
		return escape(w, parsed)
	case DocumentNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := render1(w, c); err != nil {
				return err
			}
		}
		return nil
	case ElementNode:
		// No-op.

	case CommentNode:
		if _, err := w.WriteString("<!--"); err != nil {
			return err
		}
		if _, err := w.WriteString(n.Content.(string)); err != nil {
			return err
		}
		if _, err := w.WriteString("-->"); err != nil {
			return err
		}
		return nil
	case DoctypeNode:
		if _, err := w.WriteString("<!DOCTYPE "); err != nil {
			return err
		}
		if _, err := w.WriteString(n.Tag); err != nil {
			return err
		}
		if n.ParsedAttrs != nil {
			var p, s string
			for k, v := range n.ParsedAttrs {
				switch k {
				case "public":
					p = v
				case "system":
					s = v
				}
			}
			if p != "" {
				if _, err := w.WriteString(" PUBLIC "); err != nil {
					return err
				}
				if err := writeQuoted(w, p); err != nil {
					return err
				}
				if s != "" {
					if err := w.WriteByte(' '); err != nil {
						return err
					}
					if err := writeQuoted(w, s); err != nil {
						return err
					}
				}
			} else if s != "" {
				if _, err := w.WriteString(" SYSTEM "); err != nil {
					return err
				}
				if err := writeQuoted(w, s); err != nil {
					return err
				}
			}
		}
		return w.WriteByte('>')
	case RawNode:
		//??
		_, err := w.WriteString(n.Tag)
		return err
	case ConditionalNode:
		n.parseConditional()
		n.Tag = n.FirstChild.Tag
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := render1(w, c); err != nil {
				return err
			}
		}
		return nil
	case IterativeNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			// c.Attrs = map[string]interface{}{"class": n.DOM_ID}
			if err := render1(w, c); err != nil {
				return err
			}
		}
		return nil
	default:
		return errors.New("html: unknown node type")
	}

	// Render the <xxx> opening tag.
	if err := w.WriteByte('<'); err != nil {
		return err
	}
	// write the tag
	if _, err := w.WriteString(n.Tag); err != nil {
		return err
	}
	for k, v := range n.ParsedAttrs {
		if err := w.WriteByte(' '); err != nil {
			return err
		}
		// namespaces not supported atm
		// if a.Namespace != "" {
		// 	if _, err := w.WriteString(a.Namespace); err != nil {
		// 		return err
		// 	}
		// 	if err := w.WriteByte(':'); err != nil {
		// 		return err
		// 	}
		// }
		if _, err := w.WriteString(k); err != nil {
			return err
		}
		if _, err := w.WriteString(`="`); err != nil {
			return err
		}
		if err := escape(w, v); err != nil {
			return err
		}
		if err := w.WriteByte('"'); err != nil {
			return err
		}
	}
	if voidElements[n.Tag] {
		if n.FirstChild != nil {
			return fmt.Errorf("html: void element <%s> has child nodes", n.Tag)
		}
		_, err := w.WriteString("/>")
		return err
	}
	if err := w.WriteByte('>'); err != nil {
		return err
	}

	// Add initial newline where there is danger of a newline beging ignored.
	if c := n.FirstChild; c != nil && c.NodeType == TextNode && strings.HasPrefix(c.Tag, "\n") {
		switch n.Tag {
		case "pre", "listing", "textarea":
			if err := w.WriteByte('\n'); err != nil {
				return err
			}
		}
	}

	// Render any child nodes.
	switch n.Tag {
	case "iframe", "noembed", "noframes", "noscript", "plaintext", "script", "style", "xmp":
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.NodeType == TextNode {
				if _, err := w.WriteString(c.Tag); err != nil {
					return err
				}
			} else {
				if err := render1(w, c); err != nil {
					return err
				}
			}
		}
		if n.Tag == "plaintext" {
			// Don't render anything else. <plaintext> must be the
			// last element in the file, with no closing tag.
			return plaintextAbort
		}
	default:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := render1(w, c); err != nil {
				return err
			}
		}
	}

	// Render the </xxx> closing tag.
	if _, err := w.WriteString("</"); err != nil {
		return err
	}
	if _, err := w.WriteString(n.Tag); err != nil {
		return err
	}
	return w.WriteByte('>')
}

// writeQuoted writes s to w surrounded by quotes. Normally it will use double
// quotes, but if s contains a double quote, it will use single quotes.
// It is used for writing the identifiers in a doctype declaration.
// In valid HTML, they can't contain both types of quotes.
func writeQuoted(w writer, s string) error {
	var q byte = '"'
	if strings.Contains(s, `"`) {
		q = '\''
	}
	if err := w.WriteByte(q); err != nil {
		return err
	}
	if _, err := w.WriteString(s); err != nil {
		return err
	}
	if err := w.WriteByte(q); err != nil {
		return err
	}
	return nil
}

// Section 12.1.2, "Elements", gives this list of void elements. Void elements
// are those that can't have any contents.
var voidElements = map[string]bool{
	"area":   true,
	"base":   true,
	"br":     true,
	"col":    true,
	"embed":  true,
	"hr":     true,
	"img":    true,
	"input":  true,
	"keygen": true, // "keygen" has been removed from the spec, but are kept here for backwards compatibility.
	"link":   true,
	"meta":   true,
	"param":  true,
	"source": true,
	"track":  true,
	"wbr":    true,
}

// escape.go
// These replacements permit compatibility with old numeric entities that
// assumed Windows-1252 encoding.
// https://html.spec.whatwg.org/multipage/syntax.html#consume-a-character-reference
var replacementTable = [...]rune{
	'\u20AC', // First entry is what 0x80 should be replaced with.
	'\u0081',
	'\u201A',
	'\u0192',
	'\u201E',
	'\u2026',
	'\u2020',
	'\u2021',
	'\u02C6',
	'\u2030',
	'\u0160',
	'\u2039',
	'\u0152',
	'\u008D',
	'\u017D',
	'\u008F',
	'\u0090',
	'\u2018',
	'\u2019',
	'\u201C',
	'\u201D',
	'\u2022',
	'\u2013',
	'\u2014',
	'\u02DC',
	'\u2122',
	'\u0161',
	'\u203A',
	'\u0153',
	'\u009D',
	'\u017E',
	'\u0178', // Last entry is 0x9F.
	// 0x00->'\uFFFD' is handled programmatically.
	// 0x0D->'\u000D' is a no-op.
}

// unescapeEntity reads an entity like "&lt;" from b[src:] and writes the
// corresponding "<" to b[dst:], returning the incremented dst and src cursors.
// Precondition: b[src] == '&' && dst <= src.
// attribute should be true if parsing an attribute value.
func unescapeEntity(b []byte, dst, src int, attribute bool) (dst1, src1 int) {
	// https://html.spec.whatwg.org/multipage/syntax.html#consume-a-character-reference

	// i starts at 1 because we already know that s[0] == '&'.
	i, s := 1, b[src:]

	if len(s) <= 1 {
		b[dst] = b[src]
		return dst + 1, src + 1
	}

	if s[i] == '#' {
		if len(s) <= 3 { // We need to have at least "&#.".
			b[dst] = b[src]
			return dst + 1, src + 1
		}
		i++
		c := s[i]
		hex := false
		if c == 'x' || c == 'X' {
			hex = true
			i++
		}

		x := '\x00'
		for i < len(s) {
			c = s[i]
			i++
			if hex {
				if '0' <= c && c <= '9' {
					x = 16*x + rune(c) - '0'
					continue
				} else if 'a' <= c && c <= 'f' {
					x = 16*x + rune(c) - 'a' + 10
					continue
				} else if 'A' <= c && c <= 'F' {
					x = 16*x + rune(c) - 'A' + 10
					continue
				}
			} else if '0' <= c && c <= '9' {
				x = 10*x + rune(c) - '0'
				continue
			}
			if c != ';' {
				i--
			}
			break
		}

		if i <= 3 { // No characters matched.
			b[dst] = b[src]
			return dst + 1, src + 1
		}

		if 0x80 <= x && x <= 0x9F {
			// Replace characters from Windows-1252 with UTF-8 equivalents.
			x = replacementTable[x-0x80]
		} else if x == 0 || (0xD800 <= x && x <= 0xDFFF) || x > 0x10FFFF {
			// Replace invalid characters with the replacement character.
			x = '\uFFFD'
		}

		return dst + utf8.EncodeRune(b[dst:], x), src + i
	}

	// Consume the maximum number of characters possible, with the
	// consumed characters matching one of the named references.

	for i < len(s) {
		c := s[i]
		i++
		// Lower-cased characters are more common in entities, so we check for them first.
		if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
			continue
		}
		if c != ';' {
			i--
		}
		break
	}

	entityName := string(s[1:i])
	if entityName == "" {
		// No-op.
	} else if attribute && entityName[len(entityName)-1] != ';' && len(s) > i && s[i] == '=' {
		// No-op.
	} else if x := entity[entityName]; x != 0 {
		return dst + utf8.EncodeRune(b[dst:], x), src + i
	} else if x := entity2[entityName]; x[0] != 0 {
		dst1 := dst + utf8.EncodeRune(b[dst:], x[0])
		return dst1 + utf8.EncodeRune(b[dst1:], x[1]), src + i
	} else if !attribute {
		maxLen := len(entityName) - 1
		if maxLen > longestEntityWithoutSemicolon {
			maxLen = longestEntityWithoutSemicolon
		}
		for j := maxLen; j > 1; j-- {
			if x := entity[entityName[:j]]; x != 0 {
				return dst + utf8.EncodeRune(b[dst:], x), src + j + 1
			}
		}
	}

	dst1, src1 = dst+i, src+i
	copy(b[dst:dst1], b[src:src1])
	return dst1, src1
}

// unescape unescapes b's entities in-place, so that "a&lt;b" becomes "a<b".
// attribute should be true if parsing an attribute value.
func unescape(b []byte, attribute bool) []byte {
	for i, c := range b {
		if c == '&' {
			dst, src := unescapeEntity(b, i, i, attribute)
			for src < len(b) {
				c := b[src]
				if c == '&' {
					dst, src = unescapeEntity(b, dst, src, attribute)
				} else {
					b[dst] = c
					dst, src = dst+1, src+1
				}
			}
			return b[0:dst]
		}
	}
	return b
}

// lower lower-cases the A-Z bytes in b in-place, so that "aBc" becomes "abc".
func lower(b []byte) []byte {
	for i, c := range b {
		if 'A' <= c && c <= 'Z' {
			b[i] = c + 'a' - 'A'
		}
	}
	return b
}

const escapedChars = "&'<>\"\r"

func escape(w writer, s string) error {
	i := strings.IndexAny(s, escapedChars)
	for i != -1 {
		if _, err := w.WriteString(s[:i]); err != nil {
			return err
		}
		var esc string
		switch s[i] {
		case '&':
			esc = "&amp;"
		case '\'':
			// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
			esc = "&#39;"
		case '<':
			esc = "&lt;"
		case '>':
			esc = "&gt;"
		case '"':
			// "&#34;" is shorter than "&quot;".
			esc = "&#34;"
		case '\r':
			esc = "&#13;"
		default:
			panic("unrecognized escape character")
		}
		s = s[i+1:]
		if _, err := w.WriteString(esc); err != nil {
			return err
		}
		i = strings.IndexAny(s, escapedChars)
	}
	_, err := w.WriteString(s)
	return err
}

// EscapeString escapes special characters like "<" to become "&lt;". It
// escapes only five such characters: <, >, &, ' and ".
// UnescapeString(EscapeString(s)) == s always holds, but the converse isn't
// always true.
func EscapeString(s string) string {
	if strings.IndexAny(s, escapedChars) == -1 {
		return s
	}
	var buf bytes.Buffer
	escape(&buf, s)
	return buf.String()
}

// UnescapeString unescapes entities like "&lt;" to become "<". It unescapes a
// larger range of entities than EscapeString escapes. For example, "&aacute;"
// unescapes to "á", as does "&#225;" and "&xE1;".
// UnescapeString(EscapeString(s)) == s always holds, but the converse isn't
// always true.
func UnescapeString(s string) string {
	for _, c := range s {
		if c == '&' {
			return string(unescape([]byte(s), false))
		}
	}
	return s
}
