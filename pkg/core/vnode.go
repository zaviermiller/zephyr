/*
This file is responsible for implementing the bulk of the so-called
"virtual DOM". In order to keep matching efficient, the VNode struct
pretty much follows the implementation of the standard library HTML
Node struct exactly.
*/
package zephyr

import (
	"fmt"       // TODO: remove
	"math/rand" // TODO: rmeove
	"reflect"
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
	Namespace, Key string
	Value          interface{}
}

// VNode struct is a simple intermediary between the stdlib html.Node
// and a Zephyr Component instance. There are also a few extra fields for
// optimizing patching. Its all on the heap :( I feel like this is bad, but
// I don't know how else to have ref variables inited in the func. FP maybe?
// Will investigate.
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
	Attrs map[string]interface{}

	// Listener is the nodes listener that
	// triggers a comparison whenever updated.
	// May want to add a way to tell exactly WHAT
	// needs to be updated.
	Listener *VNodeListener

	// HTMLNode is the Go representation of the currently rendered
	// HTML tree
	HTMLNode *html.Node

	// Other node refs
	Parent, FirstChild, LastChild, PrevSibing, NextSibling *VNode

	// Flags
	Static    bool
	Comment   bool
	Component bool
}

// the js part of this is very very temporary, in fact the whole function is
// going to just try and build it up.
// func (node *VNode) BuildAttrs(attrs map[string]interface{}) {
// 	// rand.Seed(time.Now().Unix())
// 	zAttrs := []ZephyrAttr{}
// 	createdFuncs := map[string]string{}

// 	for key, val := range attrs {
// 		// dont redefine the funcs, silly zephyr!
// 		if _, ok := createdFuncs[node.Tag]; ok {
// 			// zAttrs[key] = jsFunc
// 			continue
// 		}
// 		zAttrs = append(zAttrs, ZephyrAttr{Key: key, Value: val})
// 	}
// 	node.Attrs = zAttrs
// }

func arrToString(arr []int) string {
	str := "[ "
	for _, item := range arr {
		str += strconv.Itoa(item) + " "
	}
	str += "]"
	return str
}

func (node *VNode) ToHTMLNode() *html.Node {
	htmlNode := &html.Node{Type: html.NodeType(node.NodeType)}
	if node.NodeType == TextNode {
		switch node.Content.(type) {
		case ZephyrString:
			htmlNode.Data = node.Content.(ZephyrString).string(node.Listener)
		// computed MUST IMPLEMENT - custom type??
		case func(*VNodeListener) interface{}:
			//  := node.Content.(func(*VNodeListener) interface{})
			// TODO
			evaluated := node.Content.(func(*VNodeListener) interface{})(node.Listener)

			switch evaluated.(type) {
			case string:
				htmlNode = &html.Node{Data: evaluated.(string), Type: html.NodeType(node.NodeType)}
			case int, int8, int16, int32, int64, uint:
				htmlNode = &html.Node{Data: strconv.Itoa(evaluated.(int)), Type: html.NodeType(node.NodeType)}
			case []int:
				htmlNode = &html.Node{Data: arrToString(evaluated.([]int)), Type: html.NodeType(node.NodeType)}
			default:
				fmt.Println(node.DOM_ID+" func return type not supported by ToHTMLNode: ", reflect.TypeOf(evaluated).String())
			}
		case string:
			htmlNode.Data = node.Content.(string)
		default:
			fmt.Println(node.DOM_ID+" type not supported by ToHTMLNode: ", reflect.TypeOf(node.Content).String())
		}
	} else if node.NodeType == ElementNode {
		// htmlNode = &html.Node{Data: node.Tag, Type: html.NodeType(node.NodeType)}
		htmlNode.Data = node.Tag
	} else {
		fmt.Println(node.NodeType)
	}

	attrs := []html.Attribute{}
	for key, val := range node.Attrs {
		switch val.(type) {
		// other zdata handler
		case string:
			attrs = append(attrs, html.Attribute{Namespace: "", Key: key, Val: val.(string)})
		case func(*VNodeListener) interface{}:
			// computed attrs can only be strings
			// fmt.Println(attr.Value)
			attrs = append(attrs, html.Attribute{Namespace: "", Key: key, Val: val.(func(*VNodeListener) interface{})(node.Listener).(string)})
		}
	}
	attrs = append(attrs, html.Attribute{Key: "id", Val: node.DOM_ID})
	htmlNode.Attr = attrs
	node.HTMLNode = htmlNode
	return htmlNode
}

// ToHTMLTree builds the HTMl node and its children
// used for first runs.
func (node *VNode) ToHTMLTree() *html.Node {
	htmlNode := node.ToHTMLNode()
	currChild := node.FirstChild
	var htmlChild *html.Node
	for currChild != nil {
		if htmlChild != nil {
			htmlChild.NextSibling = currChild.ToHTMLTree()
			prev := htmlChild
			htmlChild = htmlChild.NextSibling
			htmlChild.PrevSibling = prev
		} else {
			htmlChild = currChild.ToHTMLTree()
		}
		if currChild.PrevSibing == nil {
			htmlNode.FirstChild = htmlChild
		} else if currChild.NextSibling == nil {
			htmlNode.LastChild = htmlChild
		}
		currChild = currChild.NextSibling
	}
	node.HTMLNode = htmlNode
	// fmt.Println(htmlNode)
	return htmlNode
}

func GetElID(nodeTag string) string {
	// 0-127
	return "Z" + nodeTag + "-" + strconv.Itoa(int(rand.Uint32()>>25))
}

// ChildComponent calls the render func of a child component
func (parent *BaseComponent) ChildComponent(c Component, props map[string]interface{}) *VNode {
	// set context based on parent
	parentBase := parent.getBase()
	base := c.getBase()
	base.Context = parentBase.Context

	// parse and pass props
	base.props = props
	if base.props == nil {
		base.props = make(map[string]interface{})
	}
	// fmt.Println(c, base.props)
	// initialize component
	InitWrapper(c)

	// render component
	node := RenderWrapper(c)
	return node
}

// Element will create VNodes for the element and all of its children
func Element(tag string, attrs map[string]interface{}, children []*VNode) *VNode {
	vnode := VNode{NodeType: ElementNode, Tag: tag, DOM_ID: GetElID(tag), Component: false, Listener: &VNodeListener{}}
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
		// on the heap, oh well, root elements will be stack
		curr.Parent = &vnode
		static = static && curr.Static
		prev = curr
	}
	vnode.Attrs = make(map[string]interface{})
	for key, attr := range attrs {
		switch attr.(type) {
		// handle by type
		case ZephyrString:
			vnode.Attrs[key] = attr.(ZephyrData).Value(vnode.Listener)

		case func() interface{}:
			fmt.Println("TEst")
		default:
			vnode.Attrs[key] = attr
		}
	}
	vnode.Static = static
	return &vnode
}

func StaticText(content string) *VNode {
	vnode := VNode{NodeType: TextNode, Content: content, Static: true, Component: false}

	return &vnode
}

// works with computed props!
func DynamicText(dynamicData interface{}) *VNode {
	// dynamicVal := evalFunc()
	var vnode *VNode
	vnode = &VNode{NodeType: TextNode, Component: false, Static: false, DOM_ID: GetElID("dynamicText")}
	vnode.Listener = &VNodeListener{id: vnode.DOM_ID} // something idk
	vnode.Content = dynamicData
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
