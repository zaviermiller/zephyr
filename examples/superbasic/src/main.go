package main

import (
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
	"github.com/zaviermiller/zephyr/pkg/runtime"
)

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
