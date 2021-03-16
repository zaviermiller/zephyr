// +build js,wasm

package main

import (
	"time"

	"github.com/zaviermiller/zephyr/examples/superbasic/src/components/test"
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

var App = zephyr.NewComponent(&AppComponent{})

type AppData struct {
	thickData    []int
	setThickData func(interface{})
}

type AppComponent struct {
	zephyr.BaseComponent
	AppData
}

func (ac *AppComponent) Init() {

	ac.RegisterComponents([]zephyr.Component{
		test.Component,
	})

	ac.setThickData = ac.BindList(&ac.thickData)
	ac.thickData = make([]int, 100)

	ac.setThickData(&ac.thickData)

	go func() {
		for {
			ac.AppendArr()
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (ac *AppComponent) AppendArr() {
	ac.thickData = append(ac.thickData, 0)
	ac.setThickData(&ac.thickData)
}

func (ac *AppComponent) OnesArray() interface{} {
	a := make([]int, len(ac.thickData))
	for i := range a {
		a[i] = i
	}
	return a
}

func (ac *AppComponent) Render() zephyr.VNode {
	return *zephyr.Element("div", nil, []*zephyr.VNode{
		zephyr.Element("h1", nil, []*zephyr.VNode{
			zephyr.DynamicText(ac.OnesArray),
		}),
		zephyr.ChildComponent(test.Component),
	})
}

func main() {
	zefr := zephyr.InitApp(App)
	zefr.Mount("#app")
	// dont terminate
	// done := make(chan struct{}, 1)
	// <-done
}
