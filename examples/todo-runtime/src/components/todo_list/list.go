package todo_list

import (
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

var Component = zephyr.NewComponent(&TodoListComponent{})

type TodoListComponent struct {
	zephyr.BaseComponent

	listItems zephyr.LiveData
}

func (c *TodoListComponent) Init() {
	c.BindProp("items", &c.listItems)

}

func (c *TodoListComponent) Render() *zephyr.VNode {
	return zephyr.Element("div", nil, []*zephyr.VNode{
		// zephyr.DynamicText(zephyr.IndexOfFactory(c.listItems, 0)),
		zephyr.DynamicText(c.listItems),
		// c.RenderFor(c.listItems, func(index int, val interface{}) {
		// 	return c.ChildComponent(todo_list_item.Component, map[string]interface{}{
		// 		"item": val.(LiveStruct)
		// })
		// })
		// c.ChildComponent(todo_list_item.Component, map[string]interface{}{
		// 	"item":
		// })
	})
}
