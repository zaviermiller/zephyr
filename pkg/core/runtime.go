package zephyr

// The runtime.go file provides the main interface for interacting
// with the Zephyr runtime.
import (
	"bytes"
	"math/rand"
	"time"

	"syscall/js"

	"golang.org/x/net/html"
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

	// DOMNodes is a map that stores each element by its
	// js.Value, which can be retrieved from the DOMElements
	// map
	DOMNodes map[string]html.Node

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
	app := ZephyrApp{RootComponent: rootInstance, UpdateQueue: make(chan DOMUpdate, 1)}

	js.Global().Set("Zephyr", map[string]interface{}{})

	return app
}

// Mount will mount the ZephyrApp to the given document
// found by the querySelector. It will also begin the
// rendering and patching process
func (z *ZephyrApp) Mount(querySelector string) {

	// set up context and init component
	z.RootComponent.getBase().Context = z
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
	z.DOMNodes = make(map[string]html.Node)

	// Render the DOM
	z.RootNode = RenderWrapper(z.RootComponent)

	go z.CompareDOM(z.RootNode)

	// Start listening for DOM updates
	for {
		// fmt.Println("waiting for update")
		currentUpdate := <-z.UpdateQueue
		// fmt.Println("received update: ", currentUpdate, currentUpdate.Data)

		// find element in map or on page and insert into map
		el, ok := z.DOMElements[currentUpdate.ElementID]
		// is element alredy on page?
		if !ok {
			el = Document(z.Anchor).GetByID(currentUpdate.ElementID)
			z.DOMElements[currentUpdate.ElementID] = el
		}

		// accept either pre-rendered HTML or html.Node
		var renderedHTML string
		switch currentUpdate.Data.(type) {
		case *html.Node:
			var bb bytes.Buffer
			html.Render(&bb, currentUpdate.Data.(*html.Node))
			renderedHTML = string(bb.Bytes())
		case string:
			renderedHTML = currentUpdate.Data.(string)
		}

		// handle different operations
		switch currentUpdate.Operation {
		case Insert:
			parentEl := z.DOMElements[currentUpdate.ElementID]
			// fmt.Println("insert ", currentUpdate.Data, "at ", currentUpdate.ElementID)
			parentEl.Call("insertAdjacentHTML", "beforeend", renderedHTML)
		case Delete:
			// fmt.Println("delete ", currentUpdate.ElementID)
		case UpdateAttr:
			// UpdateAttr data should be html.Attrribute
			newAttr := currentUpdate.Data.(html.Attribute)
			SetAttribute(el, newAttr.Key, newAttr.Val)
		case SetAttrs:
			// SetAttrs data should be map[string]string
			mapData := currentUpdate.Data.(map[string]string)
			for key, val := range mapData {
				SetAttribute(el, key, val)
			}
		case RemoveAttr:
			// TODO
			// RemoveAttr(el, currentUpdate.Data)
		case UpdateContent:
			// fmt.Println("Content updated: ", currentUpdate.ElementID)
			// fmt.Println("updating ", currentUpdate.ElementID)
			SetInnerHTML(el, renderedHTML)
			// case OverwriteInnerHTML:
			// 	el := Document(z.Anchor).QuerySelector(currentUpdate.ElementID)
			// 	SetInnerHTML(el, currentUpdate.Data.(string))
		}
	}
}
