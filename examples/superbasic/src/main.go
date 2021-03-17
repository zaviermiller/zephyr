// +build js,wasm

package main

import (
	"strconv"
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

func (ac *AppComponent) ArrLength() interface{} {
	return strconv.Itoa(len(ac.thickData))
}

func (ac *AppComponent) Render() *zephyr.VNode {
	return zephyr.Element("div", nil, []*zephyr.VNode{
		zephyr.Element("h1", nil, []*zephyr.VNode{
			zephyr.DynamicText(ac.OnesArray),
		}),
		ac.ChildComponent(test.Component, map[string]interface{}{"propArr": &ac.thickData, "propComputed": ac.ArrLength}),
	})
}

func main() {
	zefr := zephyr.CreateApp(App)
	zefr.Mount("#app")
}
