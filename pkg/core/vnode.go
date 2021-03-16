/*
This file is responsible for implementing the bulk of the so-called
"virtual DOM". In order to keep matching efficient, the VNode struct
pretty much follows the implementation of the standard library HTML
Node struct exactly.
*/
package zephyr

import (
	// unneeded

	"fmt" // unneeded
	"math/rand"
	"reflect" // unneeded
	"strconv"

	"golang.org/x/net/html"
)

type VNodeType int

// (follows html package exactly)
const (
	ErrorNode VNodeType = iota
	TextNode
	DocumentNode
	ElementNode
	CommentNode
	DoctypeNode
	// RawNode nodes are not returned by the parser, but can be part of the
	// Node tree passed to func Render to insert raw HTML (without escaping).
	// If so, this package makes no guarantee that the rendered HTML is secure
	// (from e.g. Cross Site Scripting attacks) or well-formed.
	RawNode
)

// ZephyrAttrs represents the HTML attributes for
// a VNode in a key/value map.
// e.g. <input type="text" /> -> "type": "text"
type ZephyrAttr struct {
	// Namespace is currently unused
	Namespace, Key, Value string
}

// The VNode struct is a simple intermediary between the stdlib html.Node
// and a Zephyr Component instance. There are also a few extra fields for
// optimizing patching
type VNode struct {
	// NodeType is the virtual nodes DOM node type
	NodeType VNodeType

	// Tag is the HTML tag - only used for ElementNodes
	Tag string

	// DOM_ID is the auto-generated ID that is used to update the node
	DOM_ID string

	// Content - only used for Text/CommentNodes
	Content interface{}

	// Attrs stores the attributes for the ZNode
	Attrs []ZephyrAttr

	// Component responsible for this vnode
	Component Component

	// HTMLNode is the Go representation of the currently rendered
	// HTML tree
	HTMLNode *html.Node
	// Listener is alerted when data the node cares
	// about is updated
	// Listener *zephyr.VDomListener

	// Other node refs
	Parent, FirstChild, LastChild, PrevSibing, NextSibling *VNode

	Static  bool
	Comment bool
}

// the js part of this is very very temporary, in fact the whole function is
// going to just try and build it up.
func (node *VNode) BuildAttrs(attrs map[string]interface{}) {
	// rand.Seed(time.Now().Unix())
	zAttrs := []ZephyrAttr{}
	createdFuncs := map[string]string{}

	for key, val := range attrs {
		// dont redefine the funcs, silly zephyr!
		if _, ok := createdFuncs[node.Tag]; ok {
			// zAttrs[key] = jsFunc
			continue
		}
		switch val.(type) {
		case string:
			zAttrs = append(zAttrs, ZephyrAttr{Key: key, Value: val.(string)})
		// this is some kind of method that returns a value
		default:
			fmt.Println(reflect.TypeOf(val).String())
		}
		node.Attrs = zAttrs
	}
}

func arrToString(arr []int) string {
	str := "[ "
	for _, item := range arr {
		str += strconv.Itoa(item) + " "
	}
	str += "]"
	return str
}

// ToHTMLNode will build the HTML node that can then be
// compared to the VNodes N
func (node *VNode) ToHTMLNode() *html.Node {
	htmlNode := &html.Node{}
	if node.NodeType == TextNode {
		switch node.Content.(type) {
		case string:
			htmlNode = &html.Node{Data: node.Content.(string), Type: html.NodeType(node.NodeType)}

		case *string:
			htmlNode = &html.Node{Data: *node.Content.(*string), Type: html.NodeType(node.NodeType)}

		case *[]int:
			htmlNode = &html.Node{Data: arrToString(*node.Content.(*[]int)), Type: html.NodeType(node.NodeType)}
		case []int:
			// htmlNode = &html.Node{Data: arrToString(node.Content.([]int)), Type: html.NodeType(node.NodeType)}

		// some computed prop
		case func() interface{}:
			evaluated := node.Content.(func() interface{})()

			switch evaluated.(type) {
			case *string:
				htmlNode = &html.Node{Data: evaluated.(string), Type: html.NodeType(node.NodeType)}
			case *[]int:
				htmlNode = &html.Node{Data: arrToString(*evaluated.(*[]int)), Type: html.NodeType(node.NodeType)}
			case []int:
				htmlNode = &html.Node{Data: arrToString(evaluated.([]int)), Type: html.NodeType(node.NodeType)}
			default:
				fmt.Println(reflect.TypeOf(evaluated).String())
			}
		default:
			fmt.Println(reflect.TypeOf(node.Content).String())
		}
	} else if node.NodeType == ElementNode {
		htmlNode = &html.Node{Data: node.Tag, Type: html.NodeType(node.NodeType)}
	} else {
		fmt.Println(node.NodeType)
	}
	attrs := []html.Attribute{}
	for _, attr := range node.Attrs {
		attrs = append(attrs, html.Attribute{Namespace: attr.Namespace, Key: attr.Key, Val: attr.Value})
	}
	attrs = append(attrs, html.Attribute{Key: "id", Val: node.DOM_ID})
	currChild := node.FirstChild
	var htmlChild *html.Node
	for currChild != nil {
		if htmlChild != nil {
			htmlChild.NextSibling = currChild.ToHTMLNode()
			prev := htmlChild
			htmlChild = htmlChild.NextSibling
			htmlChild.PrevSibling = prev
		} else {
			htmlChild = currChild.ToHTMLNode()
		}
		if currChild.PrevSibing == nil {
			htmlNode.FirstChild = htmlChild
		} else if currChild.NextSibling == nil {
			htmlNode.LastChild = htmlChild
		}
		currChild = currChild.NextSibling
	}
	htmlNode.Attr = attrs
	node.HTMLNode = htmlNode
	// fmt.Println(htmlNode)
	return htmlNode
}

func GetElID(nodeTag string) string {
	return "z-" + nodeTag + strconv.Itoa(rand.Int()%128)
}

// ChildComponent calls the render func of a child component
func ChildComponent(c Component) *VNode {
	node := RenderWrapper(c)
	return &node
}

func (node *VNode) Props(map[string]interface{}) VNode {

}

// Element will create VNodes for the element and all of its children
func Element(tag string, attrs map[string]interface{}, children []*VNode) VNode {
	vnode := VNode{NodeType: ElementNode, Tag: tag, DOM_ID: GetElID(tag)}
	var prev *VNode = nil
	var next *VNode = nil
	static := true
	// linked list of children
	for i, curr := range children {
		if i < len(children)-1 {
			next = children[i+1]
		} else {
			next = nil
			vnode.LastChild = curr
		}

		if i == 0 {
			vnode.FirstChild = curr
		}
		curr.PrevSibing = prev
		curr.NextSibling = next
		curr.Parent = &vnode
		static = static && curr.Static
		prev = curr
	}
	vnode.BuildAttrs(attrs)
	vnode.Static = static
	return vnode
}

func StaticText(content string) VNode {
	vnode := VNode{NodeType: TextNode, Content: content, Static: true}

	return vnode
}

func DynamicText(dynamicData interface{}) VNode {
	// dynamicVal := evalFunc()

	// type stuff
	vnode := VNode{NodeType: TextNode, Content: dynamicData}
	return vnode
	// switch dynamicVal.(type) {
	// case *[]int:
	// default:
	// 	return &VNode{NodeType: TextNode}
	// }
}

func Comment(commentMsg string) VNode {
	vnode := VNode{NodeType: CommentNode, Content: commentMsg}

	return vnode
}

// func Child() VNode {

// }
