package runtime

// Only place the runtime should be coupled to the core package?
import (
	"encoding/json"
	"fmt"
	"syscall/js"

	zephyr "github.com/zaviermiller/zephyr/pkg/core"
	"github.com/zaviermiller/zephyr/pkg/core/vdom"
)

// zephyrJS is a struct representing the js Zephyr var
type zephyrJS struct {
	rootComponent zephyr.Component
	anchor        string
}

type ZephyrApp struct {
	// Anchor is a JS Value representing an HTMLElement
	// object.
	Anchor js.Value

	// might create a prototype on the root element, this will
	// contain its data
	js zephyrJS

	// ComponentInstance is the instance of the root component
	RootComponent zephyr.Component

	// VDomRoot is the root VNode for the virtual DOM.
}

// would just pass the struct type into the array but...
func InitApp(rootInstance zephyr.Component) *ZephyrApp {
	app := &ZephyrApp{RootComponent: rootInstance}

	js.Global().Set("Zephyr", map[string]interface{}{})

	// init the app component which kicks off the rest
	app.RootComponent.Init()

	// create listener for component changes
	ListenerFunc := func() {
		vdom := app.RootComponent.Render()
		fmt.Println("update detected! new vdom: " + func() string { b, _ := json.Marshal(vdom); return string(b) }())
		app.CompareAndUpdateDOM(&vdom)
	}

	// move?
	app.RootComponent.CreateListener(zephyr.ComponentListener{ID: "rootListener", Updater: ListenerFunc})

	return app
}

// Mount will mount the ZephyrApp to the given document
// found by the querySelector. It will also begin the
// rendering and patching process
func (z *ZephyrApp) Mount(querySelector string) {

	// Anchor the app to the given element selector
	z.Anchor = vdom.GetDocument().QuerySelector(querySelector)
	js.Global().Get("Zephyr").Set("anchor", z.Anchor)
	// js.Global().Get("Zephyr").Set("rootComponent", js.ValueOf(z.ComponentInstanc))

	newDom := z.RootComponent.Render()

	z.CompareAndUpdateDOM(&newDom)

	// call on mount?
}

func (z *ZephyrApp) CompareAndUpdateDOM(newVDom *vdom.VNode) {
	vdom.SetInnerHTML(z.Anchor, newVDom.RenderHTML())
}
