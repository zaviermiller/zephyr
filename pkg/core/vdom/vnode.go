package vdom

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"syscall/js"
	// "syscall/js"
)

type VNodeType int

// ZNode types, follows https://www.w3schools.com/jsref/prop_node_nodetype.asp
const (
	ElementNode VNodeType = iota
	TextNode
	CommentNode
	// DocumentNode
	// DocTypeNode
)

// ZephyrAttrs represents the HTML attributes for
// a VNode in a key/value map.
// e.g. <input type="text" /> -> "type": "text"
type ZephyrAttrs map[string]string

// ZNode is the struct containing information
// about each node on the virtual DOM
type VNode struct {
	// NodeType is the virtual nodes DOM node type
	NodeType VNodeType

	// Tag is the HTML tag - only used for ElementNodes
	Tag string

	// Content - only used for Text/CommentNodes
	Content string

	// Attrs stores the attributes for the ZNode
	Attrs ZephyrAttrs

	// Component

	// Children stores the ZNode children
	Children []VNode

	Parent *VNode

	Static  bool
	Comment bool
}

func (node *VNode) BuildAttrs(attrs map[string]interface{}) {
	zAttrs := ZephyrAttrs{}
	for key, val := range attrs {
		switch val.(type) {

		case string:
			zAttrs[key] = val.(string)

		case func():
			funcName := node.Tag + strconv.Itoa(rand.Int()%32)
			js.Global().Set(funcName, js.FuncOf(
				func(this js.Value, args []js.Value) interface{} {
					val, _ := val.(func())
					val()
					return nil
				}))
			zAttrs[key] = funcName + "()"

			// Most likely an event function
		case func(js.Value):
			funcName := node.Tag + strconv.Itoa(rand.Int()%32)
			js.Global().Set(funcName, js.FuncOf(
				func(this js.Value, args []js.Value) interface{} {
					// todo
					val.(func(js.Value))(args[0])
					return nil
				}))
			zAttrs[key] = funcName + "(this)"
		default:
			fmt.Println(reflect.TypeOf(val))
		}
	}
	node.Attrs = zAttrs
}

func BuildElem(tag string, attrs map[string]interface{}, children []VNode) VNode {
	vnode := VNode{NodeType: ElementNode, Tag: tag, Children: children}
	vnode.BuildAttrs(attrs)

	return vnode
}

func BuildText(content string) VNode {
	vnode := VNode{NodeType: TextNode, Content: content}

	return vnode
}

func BuildComment(commentMsg string) VNode {
	vnode := VNode{NodeType: CommentNode, Content: commentMsg}

	return vnode
}

// RenderHTML will return a string containing the HTML
// for the VNode and all its children
func (n VNode) RenderHTML() string {
	htmlString := ""
	switch t := n.NodeType; t {
	case ElementNode:
		htmlString += "<" + n.Tag
		if n.Attrs != nil {
			for key, val := range n.Attrs {
				htmlString += " " + key + "=" + val
			}
		}
		htmlString += ">"
		htmlString += n.renderChildrenHtml()
		htmlString += "</" + n.Tag + ">"
	case TextNode:
		htmlString += n.Content
	case CommentNode:
		htmlString += "<!--" + n.Content + "-->"
	}
	return htmlString
}

// put this into a func since we do it a bunch
func (n VNode) renderChildrenHtml() string {
	htmlStr := ""
	for _, child := range n.Children {
		htmlStr += child.RenderHTML()
	}

	return htmlStr
}
