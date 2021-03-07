package main

import (
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
	"github.com/zaviermiller/zephyr/pkg/core/vdom"
	"github.com/zaviermiller/zephyr/pkg/runtime"
)

type RootComponent struct {
	*zephyr.BaseComponent
}

func (rc *RootComponent) Init() {
	rc.DefineData(map[string]interface{}{
		"message": "Hello Zephyr!",
	})

}

// Render is a function that must be implemented by all
// components and is responsible for building the vdom of the
// component.
func (rc *RootComponent) Render() vdom.VNode {
	// localVar := reflect.TypeOf(&hello.HelloComponent{}).String()

	return vdom.BuildElem("h1", nil, []vdom.VNode{vdom.BuildText(rc.GetStr("message"))})
}

// entry-point for Zephyr apps
func main() {
	// ideally we wouldnt need to initialize a variable here, but there is not other way :(
	Root := zephyr.NewComponent(&RootComponent{&zephyr.BaseComponent{}})
	zefr := runtime.InitApp(Root) // initialize plugins here??

	// mount the zephyr app to an element on an HTML doc
	zefr.Mount("#app")

	// DO NOT EDIT/REMOVE - This line prevents the WASM binary from terminating
	done := make(chan struct{}, 0)
	<-done

}
