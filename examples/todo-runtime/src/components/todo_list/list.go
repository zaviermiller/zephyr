package todo_list

import (
	"github.com/zaviermiller/zephyr/examples/todo-runtime/src/components/todo_list/todo_list_item"
	"github.com/zaviermiller/zephyr/examples/todo-runtime/src/todo"
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

type TodoListComponent struct {
	zephyr.BaseComponent

	listItems zephyr.LiveArray
}

func (c *TodoListComponent) Init() {
	c.listItems = c.BindProp("items").(zephyr.LiveArray)

}

func (c *TodoListComponent) Render() *zephyr.VNode {
	return zephyr.Element("div", nil, []*zephyr.VNode{
		// zephyr.DynamicText(zephyr.IndexOfFactory(c.listItems, 0)),
		zephyr.RenderFor(c.listItems, func(index int, val interface{}) *zephyr.VNode {
			return c.ComponentWithProps(&todo_list_item.TodoListItemComponent{}, map[string]interface{}{
				"item": val,
			}).Key(val.(*todo.TodoItem).Content)
		}),
		// c.ChildComponent(todo_list_item.Component, map[string]interface{}{
		// 	"item":
		// })
	})
}
