package zephyr

// The runtime.go file provides the main interface for interacting
// with the Zephyr runtime.
import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"syscall/js"
)

// zephyrJS is a struct representing the js Zephyr var
type zephyrJS struct {
	rootComponent Component
	anchor        string
}

type ZephyrApp struct {
	// Anchor is a JS Value representing an HTMLElement
	// object.
	Anchor         js.Value
	AnchorSelector string

	// UpdateQueue is a channel that receives a DOMUpdate,
	// which holds the id of the DOM element, the op, and data
	UpdateQueue chan DOMUpdate

	// QueueProxy allows updates to be "smushed" together,
	// queuing similar updates as one.
	QueueProxy map[string]DOMUpdate

	// DOMNodes is a map that stores each element by its
	// js.Value, which can be retrieved from the DOMElements
	// map
	DOMNodes map[string]*VNode

	// DOMElements holds each element on the pages js.Value
	// by its id. I think this is faster than calling getElementById
	// in JS
	DOMElements map[string]js.Value

	// might create a prototype on the root element, this will
	// contain its data
	js zephyrJS

	// ComponentInstance is the instance of the root component
	RootComponent Component

	RootNode *VNode
}

// CreateApp creates and returns an instance of a ZephyrApp,
// after doing app-wide initialization (plugins and other stuff maybe)
func CreateApp(rootInstance Component) ZephyrApp {
	rand.Seed(time.Now().Unix())
	app := ZephyrApp{RootComponent: rootInstance, UpdateQueue: make(chan DOMUpdate, 10), QueueProxy: map[string]DOMUpdate{}}

	js.Global().Set("Zephyr", map[string]interface{}{})

	return app
}

// Mount will mount the ZephyrApp to the given document
// found by the querySelector. It will also begin the
// rendering and patching process
func (z *ZephyrApp) Mount(querySelector string) {

	// set up context and init component
	// z.RootComponent.getBase().Context = z
	z.RootComponent.Init()
	// fmt.Println(z.RootComponent.getBase().Listener.Updater)
	// Anchor the app to the given element selector
	z.Anchor = GetDocument().QuerySelector(querySelector)
	z.AnchorSelector = querySelector
	js.Global().Get("Zephyr").Set("anchor", z.Anchor)

	// Set up other contexts to make certain things easier
	z.DOMElements = map[string]js.Value{
		querySelector: z.Anchor,
	}
	z.DOMNodes = make(map[string]*VNode)

	// Render the DOM and pass in the update channel.
	z.RootNode = RenderWrapper(z.RootComponent, z.UpdateQueue)

	// initial render
	z.UpdateQueue <- DOMUpdate{Operation: InitialRender, ElementID: z.RootNode.DOM_ID, Data: z.RootNode}

	// recursively find and set events
	// var eventRecur func(*VNode)
	// eventRecur = func(node *VNode) {
	// 	if node.events != nil {
	// 		z.UpdateQueue <- DOMUpdate{Operation: AddEventListeners, ElementID: node.DOM_ID, Data: node.events}
	// 	}
	// 	for c := node.FirstChild; c != nil; c = c.NextSibling {
	// 		eventRecur(c)
	// 	}
	// }
	// eventRecur(z.RootNode)

	// Start listening for DOM updates
	for {
		// fmt.Println("waiting for update")
		currentUpdate := <-z.UpdateQueue
		fmt.Println("received update: ", currentUpdate, currentUpdate.Data)

		// find element in map or on page and insert into map
		el, ok := z.DOMElements[currentUpdate.ElementID]
		// is element alredy on page?
		if !ok {
			el = Document(z.Anchor).GetByID(currentUpdate.ElementID)
			// fmt.Println(currentUpdate.ElementID, el)
			z.DOMElements[currentUpdate.ElementID] = el
		}

		// accept either pre-rendered HTML or html.Node
		var renderedHTML string
		var elId string
		switch currentUpdate.Data.(type) {
		case *VNode:
			var bb bytes.Buffer
			RenderHTML(&bb, currentUpdate.Data.(*VNode))
			elId = currentUpdate.Data.(*VNode).DOM_ID
			renderedHTML = string(bb.Bytes())
		case string:
			renderedHTML = currentUpdate.Data.(string)
		}

		// handle different operations
		switch currentUpdate.Operation {
		case InitialRender:
			z.Anchor.Set("innerHTML", renderedHTML)
		case Insert:
			parentEl := z.DOMElements[currentUpdate.ElementID]
			// fmt.Println("insert ", currentUpdate.Data, "at ", currentUpdate.ElementID)
			parentEl.Call("insertAdjacentHTML", "beforeend", renderedHTML)
		case Delete:
			// fmt.Println("delete ", currentUpdate.ElementID)
		// case UpdateAttr:
		// 	// UpdateAttr data should be html.Attrribute
		// 	newAttr := currentUpdate.Data.(html.Attribute)
		// 	// fmt.Println(newAttr)
		// 	SetAttribute(el, newAttr.Key, newAttr.Val)
		// delete(z.QueueProxy, strconv.Itoa(int(currentUpdate.Operation))+"."+currentUpdate.ElementID)
		case UpdateAttrs:
		case UpdateConditional:
			fmt.Println("received conditional render update: ", renderedHTML)
			// node := currentUpdate.Data.(*VNode)
			if _, ok := z.DOMNodes[currentUpdate.ElementID]; !ok {
				// insert node at parent
				// el := z.Anchor.Call("getElementById", node.Parent.DOM_ID)
				el := GetFirstElemWithClass(z.Anchor, currentUpdate.ElementID)
				el.Call("insertAdjacentHTML", "beforeend", renderedHTML)
			}
			// replace node
			el := GetFirstElemWithClass(z.Anchor, currentUpdate.ElementID)
			ReplaceElement(el, renderedHTML)
		case SetAttrs:
			// SetAttrs data should be map[string]string
			mapData := currentUpdate.Data.(map[string]string)
			for key, val := range mapData {
				SetAttribute(el, key, val)
			}
		case Replace:
			el := z.DOMElements[currentUpdate.ElementID]
			ReplaceElement(el, renderedHTML)
		case AddEventListeners:
			// fmt.Println("test")
			// el := z.DOMElements[currentUpdate.ElementID]
			el := GetFirstElemWithClass(z.Anchor, currentUpdate.ElementID)
			for ev, cb := range currentUpdate.Data.(map[string]func(*DOMEvent)) {
				AddEventListener(el, ev, cb)
			}
		case RemoveAttr:
			// TODO
			// RemoveAttr(el, currentUpdate.Data)
		case UpdateContent:
			// fmt.Println("Content updated: ", currentUpdate.ElementID)
			el := GetFirstElemWithClass(z.Anchor, currentUpdate.ElementID)
			fmt.Println(el, currentUpdate.ElementID)
			// fmt.Println("updating ", currentUpdate.ElementID)
			SetInnerHTML(el, renderedHTML)
			// case OverwriteInnerHTML:
			// 	el := Document(z.Anchor).QuerySelector(currentUpdate.ElementID)
			// 	SetInnerHTML(el, currentUpdate.Data.(string))
		}
		if currentUpdate.ElementID == z.AnchorSelector {
			// el = Document(z.Anchor).QuerySelector(currentUpdate.ElementID)
			// z.DOMElements[currentUpdate.ElementID] = el
			// when parent is passed
		} else {
			el = Document(z.Anchor).GetByID(currentUpdate.ElementID)
			z.DOMElements[currentUpdate.ElementID] = el
		}
		if elId != "" {
			z.DOMElements[elId] = Document(z.Anchor).GetByID(elId)
			fmt.Println("test: ", elId, z.DOMElements)
		}
	}
}
