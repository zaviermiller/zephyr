// This file contains all JS-DOM abstractions that the
// vdom will need to communicate with

// These functions are meant to be the only point of interaction between the DOM
// and the WASM binary. These functions should not be run anywhere else. Hopefully
// one day there is a semi/official DOM API for WASM, but for now, we have to use JS :(

package zephyr

import (
	"fmt"
	"strconv"
	"syscall/js"
)

type Document js.Value

type DOMOperation int

const (
	InitialRender DOMOperation = iota
	Insert
	Delete
	UpdateAttr
	SetAttrs
	RemoveAttr
	UpdateContent
	Replace
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
	// teehee
	return el
}

func SetInnerHTML(el js.Value, content string) {
	// fmt.Println("set ^ to ", content)
	el.Set("innerHTML", content)
}

func ReplaceElement(el js.Value, newEl string) {
	// parent := el.Get("parentNode")
	// tmpEl := js.Global().Get("document").Call("createElement", "div")
	el.Set("outerHTML", newEl)
	// js.Global().Get("console").Call("dir", tmpEl.Get("firstChild"))
	// parent.Call("replaceChild", el, tmpEl.Get("firstChild"))
}

func SetAttribute(el js.Value, key, val string) {
	// js.Global().Get("console").Call("dir", el)
	el.Call("setAttribute", key, val)
}

func RemoveAttribute(el js.Value, key string) {
	// js.Global().Get("console").Call("dir", el)
	el.Call("removeAttribute", key)
}

// CompareNode will compare the currently rendered component subtree
// and a newly generated one. It sends over any updates through
// the UpdateQueue, where they will then be processed. The
// currently rendered DOM is stored in the DOMNodes map, which
// allows for quick reads for comparisons.
func (z *ZephyrApp) CompareNode(root *VNode) {
	// fmt.Println("node: ", root.DOM_ID)
	if root.DOM_ID == z.RootNode.DOM_ID {
		root.ToHTMLTree()
		z.UpdateQueue <- DOMUpdate{Operation: InitialRender, ElementID: z.AnchorSelector, Data: root.HTMLNode}
		fmt.Println("\nInitial render...\n\n")
		return
	} else {
		root.ToHTMLNode()
	}

	var RecurComp func(node VNode)

	RecurComp = func(node VNode) {
		if node.Static && node.FirstChild == nil { /* || node.StaticRoot { */
			return
		}

		// switch
		if node.NodeType == ElementNode {
			// check if node exists already
			el, ok := z.DOMNodes[node.DOM_ID]
			if !ok {
				// create node in map if it doesn't exist
				z.DOMNodes[node.DOM_ID] = *node.HTMLNode
				el = *node.HTMLNode
				// initial render
				if _, ok := z.DOMElements[node.DOM_ID]; !ok {
					domElem := GetDocument().GetByID(node.DOM_ID)
					if domElem.Equal(js.Null()) {
						if node.HTMLNode == nil {
							return
						}
						// insert at parent if it doesnt exist on dom but is ready
						z.UpdateQueue <- DOMUpdate{Operation: Insert, ElementID: node.Parent.DOM_ID, Data: node.HTMLNode}
						return
					}
					z.DOMElements[node.DOM_ID] = domElem
				}
				// z.UpdateQueue <- DOMUpdate{Operation: UpdateContent, ElementID: node.Parent.DOM_ID, Data: node}
			}
			// check attrs
			for _, val := range el.Attr {
				_, ok := node.Attrs[val.Key]
				if !ok {
					// remove attr
					z.UpdateQueue <- DOMUpdate{Operation: RemoveAttr, ElementID: node.DOM_ID, Data: val}
					continue
				}
				// set arr
				for _, newVal := range node.HTMLNode.Attr {
					if newVal.Key == val.Key {
						if newVal.Val == val.Val {
							break
						} else {
							// fmt.Println("mismatched attr, sending update: ", newVal.Val, val.Val)
							// z.QueueUpdate(UpdateAttr, node.DOM_ID, newVal)
							z.UpdateQueue <- DOMUpdate{Operation: UpdateAttr, ElementID: node.DOM_ID, Data: newVal}
							break
						}
					}
				}
			}
			// currChild := node.FirstChild
			// for currChild != nil {
			// 	RecurComp(*currChild)
			// 	currChild = currChild.NextSibling
			// }
		} else if node.NodeType == TextNode {
			otherContent, ok := z.DOMNodes[node.Parent.DOM_ID]
			if !ok || otherContent.Data != node.HTMLNode.Data {
				// update dom content
				z.UpdateQueue <- DOMUpdate{Operation: UpdateContent, ElementID: node.Parent.DOM_ID, Data: node.HTMLNode.Data}
			}
		} else if node.NodeType == ConditionalNode {
			el, ok := z.DOMElements[node.DOM_ID]
			if !ok {
				js.Global().Get("console").Call("dir", el)
				el = GetDocument().GetByID(node.DOM_ID)
				if el.Equal(js.Null()) {
					z.UpdateQueue <- DOMUpdate{Operation: Insert, ElementID: node.Parent.DOM_ID, Data: node.HTMLNode}
					return
				}
			}
			if node.ConditionUpdated {
				z.UpdateQueue <- DOMUpdate{Operation: Replace, ElementID: node.DOM_ID, Data: node.HTMLNode}
			}
			newEl := GetDocument().GetByID(node.DOM_ID)
			z.DOMElements[node.DOM_ID] = newEl
			js.Global().Get("console").Call("dir", newEl)

		}

		z.DOMNodes[node.DOM_ID] = *node.HTMLNode
	}

	RecurComp(*root)
}

func (z *ZephyrApp) QueueUpdate(op DOMOperation, id string, data interface{}) {
	updateID := strconv.Itoa(int(op)) + "." + id
	val, ok := z.QueueProxy[updateID]
	if !ok {
		z.QueueProxy[updateID] = DOMUpdate{Operation: op, ElementID: id, Data: data}
		return
	}
	val.Data = data
}
