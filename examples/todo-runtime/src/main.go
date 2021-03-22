// +build js,wasm

package main

import (
	"time"

	"github.com/zaviermiller/zephyr/examples/todo-runtime/src/components/todo_list"
	"github.com/zaviermiller/zephyr/examples/todo-runtime/src/todo"
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

var App = zephyr.NewComponent(&AppComponent{})

type AppComponent struct {
	zephyr.BaseComponent

	// Data
	todoItems zephyr.LiveArray
	showField zephyr.LiveBool
}

func (ac *AppComponent) Init() {

	item1 := todo.NewTodoItem("clean up")

	ac.todoItems = ac.NewLiveArray([]zephyr.LiveStructIface{
		&item1,
	})
	ac.showField = ac.NewLiveBool(false)

	go func() {
		time.Sleep(1 * time.Second)
		item2 := todo.NewTodoItem("make a mess")
		ac.todoItems.Set(append(ac.todoItems.Value(nil).([]zephyr.LiveStructIface), &item2))
		time.Sleep(2 * time.Second)
		item1.Complete()
		ac.showField.Set(true)
		time.Sleep(2 * time.Second)
		ac.showField.Set(false)
	}()
}

// func (ac *AppComponent)

func (ac *AppComponent) Render() *zephyr.VNode {
	return zephyr.Element("div", nil, []*zephyr.VNode{
		ac.ChildComponent(todo_list.Component, map[string]interface{}{"items": ac.todoItems}),
		zephyr.Element("button", nil, []*zephyr.VNode{
			zephyr.StaticText("New item"),
		}),
		zephyr.Element("br", nil, nil),
		zephyr.RenderIf(ac.showField, func(l *zephyr.VNodeListener) *zephyr.VNode {
			return zephyr.Element("input", map[string]interface{}{"type": "text", "placeholder": "eg. Pick up the daughter"}, nil)
		}).RenderElse(func(l *zephyr.VNodeListener) *zephyr.VNode {
			return zephyr.Element("p", nil, []*zephyr.VNode{zephyr.StaticText("\nlick my balls")})
		}),
		// test.Component.RenderWithProps()
	})
}

func main() {
	zefr := zephyr.CreateApp(App)
	zefr.Mount("#app")
}
