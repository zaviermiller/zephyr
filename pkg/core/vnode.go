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

	// The following are types that don't follow the Go HTML package. They are
	// used for conditional and iterative rendering
	ConditionalNode
)

// ZephyrAttrs represents the HTML attributes for
// a VNode in a key/value map.
// e.g. <input type="text" /> -> "type": "text"
type ZephyrAttr struct {
	// Namespace is currently unused
	Namespace, Key string
	// Value should be a regular Go
	// data type.
	Value interface{}
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

	// Special fields - may want to interface it up
	ConditionalRenders []ConditionalRender
	CurrentCondition   int
	ConditionUpdated   bool
}

type ConditionalRender struct {
	Condition interface{}
	Render    *VNode
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
func (node *VNode) ToHTMLNode() *html.Node {
	htmlNode := &html.Node{Type: html.NodeType(node.NodeType)}
	switch node.NodeType {
	case TextNode:
		switch node.Content.(type) {
		case LiveData:
			htmlNode.Data = node.Content.(LiveData).string(node.Listener)
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
				htmlNode = &html.Node{Data: arrToString(evaluated.([]int), node.Listener), Type: html.NodeType(node.NodeType)}
			default:
				fmt.Println(node.DOM_ID+" func return type not supported by ToHTMLNode: ", reflect.TypeOf(evaluated).String())
			}
		case string:
			htmlNode.Data = node.Content.(string)
		case int:
			htmlNode.Data = strconv.Itoa(node.Content.(int))
		default:
			fmt.Println(node.DOM_ID+" type not supported by ToHTMLNode: ", reflect.TypeOf(node.Content).String())
		}
	case ElementNode:
		// htmlNode = &html.Node{Data: node.Tag, Type: html.NodeType(node.NodeType)}
		htmlNode.Data = node.Tag
	case ConditionalNode:
		for i, cr := range node.ConditionalRenders {
			condition := cr.Condition
			switch condition.(type) {
			case bool:
				if condition.(bool) {
					node.FirstChild = cr.Render
					val, ok := node.FirstChild.Attrs["id"]
					if !ok {
						node.FirstChild.Attrs["id"] = node.DOM_ID
					} else {
						node.FirstChild.Attrs["id"] = val.(string) + " " + node.DOM_ID
					}
					node.HTMLNode = node.FirstChild.ToHTMLTree()
					node.ConditionUpdated = false
					if node.CurrentCondition != i {
						node.ConditionUpdated = true
					}
					node.HTMLNode.Attr = append(node.HTMLNode.Attr, html.Attribute{Key: "id", Val: node.DOM_ID})
					node.CurrentCondition = i
					return node.HTMLNode
				}
			case LiveBool:
				cBool := condition.(LiveBool).Value(node.Listener).(bool)
				if cBool {
					node.FirstChild = cr.Render
					val, ok := node.FirstChild.Attrs["id"]
					if !ok {
						node.FirstChild.Attrs["id"] = node.DOM_ID
					} else {
						node.FirstChild.Attrs["id"] = val.(string) + " " + node.DOM_ID
					}
					node.HTMLNode = node.FirstChild.ToHTMLTree()
					node.ConditionUpdated = false
					if node.CurrentCondition != i {
						node.ConditionUpdated = true
					}
					node.HTMLNode.Attr = append(node.HTMLNode.Attr, html.Attribute{Key: "id", Val: node.DOM_ID})
					node.CurrentCondition = i
					return node.HTMLNode
				}
			case func() interface{}:
				eval := condition.(func() interface{})()
				if _, ok := eval.(bool); ok {
					if eval.(bool) {
						node.FirstChild = cr.Render
						val, ok := node.FirstChild.Attrs["id"]
						if !ok {
							node.FirstChild.Attrs["id"] = node.DOM_ID
						} else {
							node.FirstChild.Attrs["id"] = val.(string) + " " + node.DOM_ID
						}
						node.HTMLNode = node.FirstChild.ToHTMLTree()
						node.ConditionUpdated = false
						if node.CurrentCondition != i {
							node.ConditionUpdated = true
						}
						node.CurrentCondition = i
						node.HTMLNode.Attr = append(node.HTMLNode.Attr, html.Attribute{Key: "id", Val: node.DOM_ID})
						return node.HTMLNode
					}
				}
			}
		}
		return nil
	default:
		fmt.Println(node.NodeType)
	}
	attrs := []html.Attribute{}
	// attrs are parsed and any zephyr data is calculated
	for key, val := range node.Attrs {
		switch val.(type) {
		case LiveData:
			parsedAttr := html.Attribute{Key: key, Val: val.(LiveData).string(node.Listener)}
			attrs = append(attrs, parsedAttr)
		// calculated functions get computed and results parsed
		case func(*VNodeListener) interface{}:
			eval := val.(func(*VNodeListener) interface{})(node.Listener)
			var parsedAttr html.Attribute
			// parse result
			switch eval.(type) {
			case LiveData:
				parsedAttr = html.Attribute{Key: key, Val: eval.(LiveData).string(node.Listener)}
			case string:
				parsedAttr = html.Attribute{Key: key, Val: eval.(string)}
			case int, int8, int16, int32, int64:
				parsedAttr = html.Attribute{Key: key, Val: strconv.Itoa(eval.(int))}
			default:
				panic("please use a string")
			}
			attrs = append(attrs, parsedAttr)
		case string:
			parsedAttr := html.Attribute{Key: key, Val: val.(string)}
			attrs = append(attrs, parsedAttr)
		default:
			panic("Please use string or LiveData or CalculatorFunc")
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
	if node.NodeType == ConditionalNode {
		return htmlNode
	}
	currChild := node.FirstChild
	var htmlChild *html.Node
	for currChild != nil {
		if htmlChild != nil {
			htmlChild.NextSibling = currChild.ToHTMLTree()
			if htmlChild.NextSibling != nil {
				prev := htmlChild
				htmlChild = htmlChild.NextSibling
				htmlChild.PrevSibling = prev
			}
		} else {
			htmlChild = currChild.ToHTMLTree()
			fmt.Println(htmlChild)
		}
		if currChild.PrevSibing == nil {
			htmlNode.FirstChild = htmlChild
		} else if currChild.NextSibling == nil {
			htmlNode.LastChild = htmlChild
		}
		currChild = currChild.NextSibling
	}
	node.HTMLNode = htmlNode
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

func RenderIf(condition interface{}, r func(*VNodeListener) *VNode) *VNode {
	vnode := &VNode{NodeType: ConditionalNode, Component: false, Listener: &VNodeListener{}, DOM_ID: GetElID("conditional")}
	vnode.ConditionalRenders = []ConditionalRender{ConditionalRender{Condition: condition, Render: r(vnode.Listener)}}
	return vnode
}

func (vnode *VNode) RenderElseIf(condition interface{}, r func(*VNodeListener) *VNode) *VNode {
	vnode.ConditionalRenders = append(vnode.ConditionalRenders, ConditionalRender{Condition: condition, Render: r(vnode.Listener)})
	return vnode
}

func (vnode *VNode) RenderElse(r func(*VNodeListener) *VNode) *VNode {
	vnode.ConditionalRenders = append(vnode.ConditionalRenders, ConditionalRender{Condition: true, Render: r(vnode.Listener)})
	return vnode
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
		// fmt.Println(curr)
		// on the heap, oh well, root elements will be stack
		curr.Parent = &vnode
		static = static && curr.Static
		prev = curr
	}
	vnode.Attrs = make(map[string]interface{})
	for key, attr := range attrs {
		// switch attr.(type) {
		// // handle by type
		// case LiveData:
		// 	vnode.Attrs[key] = attr.(LiveData)

		// case func(*VNodeListener) interface{}:
		// 	vnode.Attrs[key] = attr.(func(*VNodeListener) interface{})
		// default:
		// 	// vnode.Attrs[key] = attr
		// 	panic("type not supported")
		// }
		vnode.Attrs[key] = attr
	}
	vnode.Static = static
	// fmt.Println(vnode.Attrs)
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
	// fmt.Println(vnode)
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
