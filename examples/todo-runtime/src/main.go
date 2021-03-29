// +build js,wasm

package main

import (
	"time"

	"github.com/zaviermiller/zephyr/examples/todo-runtime/src/components/todo_list"
	"github.com/zaviermiller/zephyr/examples/todo-runtime/src/todo"
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

type AppComponent struct {
	zephyr.BaseComponent

	// Data
	todoItems  zephyr.LiveArray
	showField  zephyr.LiveBool
	fieldInput zephyr.LiveString
}

func (ac *AppComponent) Init() {

	item1 := todo.NewTodoItem("clean up")

	ac.todoItems = zephyr.NewLiveArray([]zephyr.LiveStruct{
		&item1,
	})
	ac.showField = zephyr.NewLiveBool(false)
	ac.fieldInput = zephyr.NewLiveString("make a mess")

	go func() {
		time.Sleep(1 * time.Second)
		// ac.showField.Set(true)
		// item2 := todo.NewTodoItem("make a mess")
		// ac.todoItems.Set(append(ac.todoItems.Value(nil).([]zephyr.LiveStruct), &item2))
		ac.AddTodoItem()
		// time.Sleep(2 * time.Second)
		// item1.Complete()
	}()
}

func (ac *AppComponent) AddTodoItem() {
	itemContent := ac.fieldInput.Value(nil).(string)
	newItem := todo.NewTodoItem(itemContent)
	// fmt.Println()
	ac.todoItems.Append((&newItem))
	ac.fieldInput.Set("")
	ac.showField.Set(false)
}

// func (ac *AppComponent)

func (ac *AppComponent) Render() *zephyr.VNode {
	return zephyr.Element("div", nil, []*zephyr.VNode{
		ac.ComponentWithProps(&todo_list.TodoListComponent{}, map[string]interface{}{"items": ac.todoItems}),
		zephyr.RenderIf(func(l zephyr.Listener) interface{} { return !ac.showField.Value(l).(bool) },
			zephyr.Element("button", nil, []*zephyr.VNode{
				zephyr.StaticText("New item"),
			}).BindEvent("click", func(e *zephyr.DOMEvent) { ac.showField.Set(true) }),
		).RenderElse(zephyr.Element("div", nil, []*zephyr.VNode{
			zephyr.Element("input", map[string]interface{}{"type": "text", "value": ac.fieldInput}, nil).BindEvent("input", func(e *zephyr.DOMEvent) {
				ac.fieldInput.Set(e.Target.Get("value").String())
			}),
			zephyr.Element("button", nil, []*zephyr.VNode{
				zephyr.StaticText("New item"),
			}).BindEvent("click", func(e *zephyr.DOMEvent) { ac.AddTodoItem() }),
		}),
		// test.Component.RenderWithProps()
		),
		zephyr.Element("p", nil, []*zephyr.VNode{
			zephyr.DynamicText(ac.todoItems),
		}),
	})
}

func main() {
	zefr := zephyr.CreateApp(&AppComponent{})
	zefr.Mount("#app")
}
