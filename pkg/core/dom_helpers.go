// This file contains all JS-DOM abstractions that the
// vdom will need to communicate with

// These functions are meant to be the only point of interaction between the DOM
// and the WASM binary. These functions should not be run anywhere else. Hopefully
// one day there is a semi/official DOM API for WASM, but for now, we have to use JS :(

package zephyr

import (
	"bytes"
	"fmt"
	"reflect"

	"syscall/js"

	"golang.org/x/net/html"
)

type Document js.Value

type DOMOperation int

const (
	Insert DOMOperation = iota
	Delete
	UpdateAttr
	UpdateContent
	OverwriteInnerHTML
)

type DOMUpdate struct {
	ElementID string
	Data      interface{}
	Operation DOMOperation
}

func GetDocument() Document {
	return Document(js.Global().Get("document"))
}

// QuerySelector is an idiomatic wrapper function around the regular syscall/js
// function
func (d Document) QuerySelector(selector string) js.Value {
	jsDoc := js.Value(d)
	el := jsDoc.Call("querySelector", selector)

	return el
}

func (d Document) GetByID(id string) js.Value {
	el := d.QuerySelector("#" + id)

	return el
}

func domUpdateListener(z *ZephyrApp) {
	for {
		fmt.Println("waiting for update")
		currentUpdate := <-z.UpdateQueue
		fmt.Println("received update: ", currentUpdate, currentUpdate.Data)
		var renderedHTML string
		switch currentUpdate.Data.(type) {
		case *html.Node:
			var bb bytes.Buffer
			html.Render(&bb, currentUpdate.Data.(*html.Node))
			renderedHTML = string(bb.Bytes())
		case string:
			renderedHTML = currentUpdate.Data.(string)
		}

		switch currentUpdate.Operation {
		case Insert:
			parentEl := z.DOMElements[currentUpdate.ElementID]
			fmt.Println("insert ", currentUpdate.Data, "at ", currentUpdate.ElementID)
			parentEl.Call("insertAdjacentHTML", "beforeend", renderedHTML)
		case Delete:
			fmt.Println("delete ", currentUpdate.ElementID)
		case UpdateAttr:
			// UpdateAttr data should be map[string]string
			mapData := currentUpdate.Data.(map[string]string)
			el := Document(z.Anchor).GetByID(currentUpdate.ElementID)
			for key, val := range mapData {
				SetAttribute(el, key, val)
			}
		case UpdateContent:
			fmt.Println("Content updated: ", currentUpdate.ElementID)
			// fmt.Println("updating ", currentUpdate.ElementID)
			el, ok := z.DOMElements[currentUpdate.ElementID]
			if !ok {
				el = Document(z.Anchor).GetByID(currentUpdate.ElementID)
				z.DOMElements[currentUpdate.ElementID] = el
			}
			SetInnerHTML(el, renderedHTML)
			// case OverwriteInnerHTML:
			// 	el := Document(z.Anchor).QuerySelector(currentUpdate.ElementID)
			// 	SetInnerHTML(el, currentUpdate.Data.(string))
		}
	}
}

func SetInnerHTML(el js.Value, content string) {
	js.Global().Get("console").Call("dir", el)
	fmt.Println("set ^ to ", content)
	el.Set("innerHTML", content)
}

func SetAttribute(el js.Value, key, val string) {
	// js.Global().Get("console").Call("dir", el)
	el.Call("setAttribute", key, val)
}

func CompareDOM(z *ZephyrApp) {
	root := z.RootNode
	root.ToHTMLNode()

	var RecurComp func(node VNode)

	RecurComp = func(node VNode) {
		if node.Static && node.FirstChild == nil { /* || node.StaticRoot { */
			return
		}

		if node.NodeType == ElementNode {
			if el, ok := z.DOMNodes[node.DOM_ID]; !ok {
				z.DOMNodes[node.DOM_ID] = *node.HTMLNode
				// insert op
				if _, ok := z.DOMElements[node.DOM_ID]; !ok && node.DOM_ID == root.DOM_ID {
					z.UpdateQueue <- DOMUpdate{Operation: Insert, ElementID: "#app", Data: node.HTMLNode}
					return
				} // else {
				// 	z.UpdateQueue <- DOMUpdate{Operation: Insert, ElementID: node.Parent.DOM_ID, Data: node.HTMLNode}
				// }
				// z.UpdateQueue <- DOMUpdate{Operation: UpdateContent, ElementID: node.Parent.DOM_ID, Data: node}
			} else {
				if !reflect.DeepEqual(el.Attr, node.HTMLNode.Attr) {
					// update attr
					z.UpdateQueue <- DOMUpdate{Operation: UpdateAttr, ElementID: node.DOM_ID, Data: node.HTMLNode.Attr}
				}

				currChild := node.FirstChild
				for currChild != nil {
					RecurComp(*currChild)
					currChild = currChild.NextSibling
				}
			}
		} else if node.NodeType == TextNode {
			otherContent, ok := z.DOMNodes[node.Parent.DOM_ID]
			if !ok {
				panic("error")
			}
			if otherContent.Data != node.HTMLNode.Data {
				// update dom content
				z.UpdateQueue <- DOMUpdate{Operation: UpdateContent, ElementID: node.Parent.DOM_ID, Data: node.HTMLNode.Data}
			}
		}
	}

	RecurComp(root)
}

// RenderHTML will return a string containing the HTML
// for the VNode and all its children
// func RenderHTML(n VNode) string {
// 	htmlString := ""
// 	switch t := n.NodeType; t {
// 	case ElementNode:
// 		htmlString += "<" + n.Tag
// 		if n.Attrs != nil {
// 			for key, val := range n.Attrs {
// 				htmlString += " " + key + "=" + val
// 			}
// 		}
// 		htmlString += ">"
// 		if n.Content != "" {

// 		}
// 		htmlString += renderChildrenHtml(n)
// 		htmlString += "</" + n.Tag + ">"
// 	case TextNode:
// 		switch n.Content.(type) {
// 		case *[]int:
// 			htmlString += html.EscapeString(arrToString(*(n.Content.(*[]int))))
// 		case string:
// 			htmlString += html.EscapeString(n.Content.(string))
// 		}
// 	case CommentNode:
// 		htmlString += "<!--" + n.Content.(string) + "-->"
// 	}
// 	return htmlString
// }

// func arrToString(arr []int) string {
// 	str := "["
// 	for _, item := range arr {
// 		str += strconv.Itoa(item) + " "
// 	}
// 	str += "]"
// 	return str
// }

// // put this into a func since we do it a bunch
// func renderChildrenHtml(n VNode) string {
// 	htmlStr := ""
// 	for _, child := range n.Children {
// 		htmlStr += RenderHTML(child)
// 	}

// 	return htmlStr
// }
