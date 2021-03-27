/*
This file is responsible for implementing the bulk of the so-called
"virtual DOM". In order to keep matching efficient, the VNode struct
pretty much follows the implementation of the standard library HTML
Node struct exactly.
*/
package zephyr

import (
	// TODO: remove

	"math/rand" // TODO: rmeove
	"strconv"
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
	IterativeNode
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

	// Content
	Content interface{}

	// Attrs stores the attributes for the ZNode
	Attrs map[string]interface{}
	// ParsedAttrs is a map[string]string of ready-to-render
	// attributes.
	ParsedAttrs map[string]string

	events    map[string]func(e *DOMEvent)
	listeners map[string]Listener

	// Listener is the nodes listener that
	// triggers a comparison whenever updated.
	// May want to add a way to tell exactly WHAT
	// needs to be updated.
	// Listener *VNodeListener

	RenderChan chan DOMUpdate

	// HTMLNode is the Go representation of the currently rendered
	// HTML tree
	// HTMLNode *html.Node

	// Other node refs
	Parent, FirstChild, LastChild, PrevSibling, NextSibling *VNode

	// Flags
	Static    bool
	Component bool

	// Special fields - may want to interface it up
	ConditionalRenders []ConditionalRender
	CurrentCondition   int
	ConditionUpdated   bool
	Keys               []interface{}
	key                interface{}
	IterRender         func(int, interface{}) *VNode
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

func (node *VNode) GetOrCreateListener(lID string) Listener {
	if node.listeners == nil {
		node.listeners = map[string]Listener{}
	}
	if val, ok := node.listeners[lID]; ok {
		return val
	}
	var newListener Listener
	switch lID {
	case "attr":
		newListener = VNAttrListener{node: node, id: node.DOM_ID + "__attrL"}
	case "content":
		newListener = VNContentListener{node: node, id: node.DOM_ID + "__contentL"}
	case "prop":
		newListener = VNPropListener{node: node, id: node.DOM_ID + "__propL"}
	case "calculator":
		newListener = VNCalculatorListener{node: node, id: node.DOM_ID + "__calculatorL"}
	case "conditional":
		newListener = VNConditionalListener{node: node, id: node.DOM_ID + "__conditionalL"}
	case "iterator":
		newListener = VNIteratorListener{node: node, id: node.DOM_ID + "__iteratorL"}
	default:
		panic("unknown listener type")
	}
	node.listeners[lID] = newListener
	return newListener
}

func GetElID(nodeTag string) string {
	// 0-127
	return "Z" + nodeTag + "-" + strconv.Itoa(int(rand.Uint32()>>25))
}

func (node *VNode) BindEvent(event string, callback func(e *DOMEvent)) *VNode {
	if node.events == nil {
		node.events = map[string]func(*DOMEvent){}
	}
	// add event to vnode
	if node.Component {
		// TODO
		// c := node.Content.(Component).getBase()
		// c.events
		return node
	}

	// domEvents in events.go
	if _, ok := domEvents[event]; !ok {
		panic("that event is not real dawg")
	}
	node.events[event] = callback
	return node
}

// CONDITIONAL RENDERING

func RenderIf(condition interface{}, r *VNode) *VNode {
	// set up listener
	vnode := &VNode{NodeType: ConditionalNode, Component: false, DOM_ID: GetElID("conditional"), CurrentCondition: 0}
	vnode.ConditionalRenders = []ConditionalRender{ConditionalRender{Condition: condition, Render: r}}
	return vnode
}
func (vnode *VNode) RenderElseIf(condition interface{}, r *VNode) *VNode {
	vnode.ConditionalRenders = append(vnode.ConditionalRenders, ConditionalRender{Condition: condition, Render: r})
	return vnode
}

func (vnode *VNode) RenderElse(r *VNode) *VNode {
	vnode.ConditionalRenders = append(vnode.ConditionalRenders, ConditionalRender{Condition: true, Render: r})
	return vnode
}

// ITERATIVE RENDERING

func RenderFor(iterator LiveArray, r func(index int, val interface{}) *VNode) *VNode {
	vnode := &VNode{NodeType: IterativeNode, Component: false, DOM_ID: GetElID("iterator"), Static: false, Content: iterator, IterRender: r}
	keys := vnode.parseIter()
	vnode.Keys = keys
	return vnode
}

func (vnode *VNode) Key(keyVal interface{}) *VNode {
	vnode.key = keyVal
	return vnode
}

func (node *VNode) getNewKeys() (keys []interface{}) {
	if node.NodeType != IterativeNode {
		panic("must be used only on iterativenodes")
	}
	iListener := node.GetOrCreateListener("iterator")
	val := node.Content.(LiveArray).Value(iListener)
	r := node.IterRender
	switch val.(type) {
	case []LiveStruct:
		for i, v := range val.([]LiveStruct) {
			c := r(i, v)
			keys = append(keys, c.key)
		}
	}
	return keys
}

func (node *VNode) parseIter() (keys []interface{}) {
	if node.NodeType != IterativeNode {
		panic("must be used only on iterativenodes")
	}
	iListener := node.GetOrCreateListener("iterator")
	val := node.Content.(LiveArray).Value(iListener)
	r := node.IterRender
	switch val.(type) {
	case []LiveStruct:
		valArr := val.([]LiveStruct)
		var prev *VNode
		for i, v := range valArr {
			c := r(i, v)
			if c.Key == nil {
				panic("iterators must have a key")
			}
			keys = append(keys, c.key)
			if i == 0 {
				node.FirstChild = c
				prev = nil
			}
			if prev != nil {
				prev.NextSibling = c
			}
			c.PrevSibling = prev
			// on the heap, oh well, root elements will be stack
			c.Parent = node
			c.DOM_ID = node.DOM_ID + "-k[" + c.key.(string) + "]"
			prev = c
			node.LastChild = c
		}
	default:
		panic("fuck")
	}
	return keys
}

// Element will create VNodes for the element and all of its children
func Element(tag string, attrs map[string]interface{}, children []*VNode) *VNode {
	vnode := VNode{NodeType: ElementNode, Tag: tag, DOM_ID: GetElID(tag), Component: false}
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
		curr.PrevSibling = prev
		curr.NextSibling = next
		// fmt.Println(curr)
		// on the heap, oh well, root elements will be stack
		curr.Parent = &vnode
		static = static && curr.Static
		prev = curr
	}
	vnode.Attrs = attrs

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
