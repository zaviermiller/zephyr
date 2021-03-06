package runtime

// Only place the runtime should be coupled to the core package?
import (
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

	js zephyrJS

	// ComponentInstance is the instance of the root component
	ComponentInstance zephyr.Component

	// VDomRoot is the root VNode for the virtual DOM.
	VDomRoot *vdom.VNode
}

// would just pass the struct type into the array but...
func InitApp(rootInstance zephyr.Component) ZephyrApp {
	app := ZephyrApp{ComponentInstance: rootInstance, VDomRoot: nil}

	js.Global().Set("Zephyr", map[string]interface{}{})

	// init the app component which kicks off the rest
	app.ComponentInstance.Init()

	// instantiate component and its child components

	return app
}

// Mount will mount the ZephyrApp to the given document
// found by the querySelector. It will also begin the
// rendering and patching process
func (z *ZephyrApp) Mount(querySelector string) {

	// Anchor the app to the given element selector
	z.Anchor = GetDocument().QuerySelector(querySelector)
	js.Global().Get("Zephyr").Set("anchor", z.Anchor)
	// js.Global().Get("Zephyr").Set("rootComponent", js.ValueOf(z.ComponentInstanc))

	// create listener for component changes
	ListenerFunc := func() {
		vdom := z.ComponentInstance.Render()
		z.UpdateDOM(&vdom)
	}

	z.ComponentInstance.CreateListener(zephyr.ComponentListener{ID: "rootListener", Updater: ListenerFunc})

	newDom := z.ComponentInstance.Render()

	z.UpdateDOM(&newDom)

	// call on mount?
}

func (z *ZephyrApp) UpdateDOM(newVDom *vdom.VNode) {
	// if z.VDomRoot == nil {
	z.VDomRoot = newVDom
	// }
	SetInnerHTML(z.Anchor, z.VDomRoot.RenderHTML())
}
