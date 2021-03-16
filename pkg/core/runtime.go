package zephyr

// The runtime.go file provides the main interface for interacting
// with the Zephyr runtime.
import (
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
	Anchor js.Value

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

	RootNode VNode
}

// would just pass the struct type into the array but...
func InitApp(rootInstance Component) *ZephyrApp {
	app := &ZephyrApp{RootComponent: rootInstance, UpdateQueue: make(chan DOMUpdate, 1)}

	rand.Seed(time.Now().Unix())

	js.Global().Set("Zephyr", map[string]interface{}{})

	// init the app component which kicks off the rest
	app.RootComponent.Init()

	return app
}

// Mount will mount the ZephyrApp to the given document
// found by the querySelector. It will also begin the
// rendering and patching process
func (z *ZephyrApp) Mount(querySelector string) {
	// fmt.Println(z.RootComponent.getBase().Listener.Updater)
	// create listener for root changes
	ListenerFunc := func() {
		// vdom := app.RootComponent.Render()
		// fmt.Println("update detected! new vdom: " + func() string { b, _ := json.Marshal(app.RootNode); return string(b) }())
		go CompareDOM(z)
	}
	// move?
	z.RootComponent.getBase().SetListenerUpdater(ListenerFunc)
	// Anchor the app to the given element selector
	z.Anchor = GetDocument().QuerySelector(querySelector)
	js.Global().Get("Zephyr").Set("anchor", z.Anchor)
	z.DOMElements = map[string]js.Value{
		querySelector: z.Anchor,
	}
	z.DOMNodes = make(map[string]html.Node)
	z.UpdateQueue = make(chan DOMUpdate, 1)

	z.RootNode = z.RootComponent.Render()
	domUpdateListener(z)

	// call on mount?
}
