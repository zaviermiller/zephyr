package todo_list_item

import (
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

var Component = zephyr.NewComponent(&TodoListItemComponent{})

type TodoListItemComponent struct {
	zephyr.BaseComponent

	todoItem zephyr.LiveStruct
}

func (c *TodoListItemComponent) Init() {
	c.BindProp("item", &c.todoItem)

}

func (c *TodoListItemComponent) Render() *zephyr.VNode {
	return zephyr.Element("div", nil, []*zephyr.VNode{
		// zephyr.DynamicText(zephyr.IndexOfFactory(c.listItems, 0)),
		// zephyr.DynamicText(c.listItems),
		// c.ChildComponent()
	})
}
