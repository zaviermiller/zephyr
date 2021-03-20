// +build js,wasm

package main

import (
	"fmt"
	"strconv"
	"time"

	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

var App = zephyr.NewComponent(&AppComponent{})

type AppData struct {
	stringTest zephyr.ZephyrData
}

type AppComponent struct {
	zephyr.BaseComponent
	AppData
}

func (ac *AppComponent) Init() {

	ac.stringTest = ac.NewLiveString("initial value")

	fmt.Println(ac.StrLength())

	go func() {
		for {
			ac.stringTest.Set(ac.stringTest.Value(nil).(string) + " z")
			time.Sleep(1000 * time.Millisecond)
		}
	}()
}

func (ac *AppComponent) StrLength() interface{} {
	str, _ := ac.stringTest.Value(ac.Listener).(string)
	fmt.Println(len(str))
	return strconv.Itoa(len(str))
}

func (ac *AppComponent) Render() *zephyr.VNode {
	return zephyr.Element("div", nil, []*zephyr.VNode{
		zephyr.Element("h1", nil, []*zephyr.VNode{
			zephyr.DynamicText(ac.stringTest),
		}),
		// ac.ChildComponent(test.Component, map[string]interface{}{"propArr": ac.stringTest, "propComputed": ac.StrLength}),
		// test.Component.RenderWithProps()
	})
}

func main() {
	zefr := zephyr.CreateApp(App)
	zefr.Mount("#app")
}
