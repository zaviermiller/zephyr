// +build js,wasm

package main

import (
	"time"

	"github.com/zaviermiller/zephyr/examples/superbasic/src/components/test"
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

var App = zephyr.NewComponent(&AppComponent{})

type AppData struct {
	stringTest zephyr.LiveData
}

type AppComponent struct {
	zephyr.BaseComponent
	AppData
}

func (ac *AppComponent) Init() {

	ac.stringTest = ac.NewLiveString("initial value")

	go func() {
		for {
			ac.UpdateStr(ac.stringTest.Value(nil).(string) + "z")
			time.Sleep(5 * time.Millisecond)
		}
	}()
}

func (ac *AppComponent) StrLength(l *zephyr.VNodeListener) interface{} {
	str, _ := ac.stringTest.Value(l).(string)
	// fmt.Println(len(str))
	return len(str)
}

func (ac *AppComponent) UpdateStr(newStr string) {
	ac.stringTest.Set(newStr)
}

func (ac *AppComponent) Render() *zephyr.VNode {
	return zephyr.Element("div", nil, []*zephyr.VNode{
		zephyr.Element("h1", nil, []*zephyr.VNode{
			zephyr.DynamicText(ac.StrLength),
		}),
		ac.ChildComponent(test.Component, map[string]interface{}{"prop": ac.stringTest}),
		// test.Component.RenderWithProps()
	})
}

func main() {
	zefr := zephyr.CreateApp(App)
	zefr.Mount("#app")
}
