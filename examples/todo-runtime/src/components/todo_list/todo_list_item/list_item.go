package todo_list_item

import (
	"fmt"

	"github.com/zaviermiller/zephyr/examples/todo-runtime/src/todo"
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

var Component = zephyr.NewComponent(&TodoListItemComponent{})

type TodoListItemComponent struct {
	zephyr.BaseComponent

	todoItem *todo.TodoItem
}

func (c *TodoListItemComponent) Init() {
	c.todoItem = c.BindProp("item").(*todo.TodoItem)
	// fmt.Println(c.todoItem)
	fmt.Println(c.todoItem)
}

func (c *TodoListItemComponent) Render() *zephyr.VNode {
	return zephyr.Element("div", map[string]interface{}{"style": "margin-bottom: 10px;"}, []*zephyr.VNode{
		zephyr.RenderIf(func(l zephyr.Listener) interface{} { return !c.todoItem.IsComplete(l) },
			zephyr.Element("button", map[string]interface{}{
				"style": "margin: 10px 10px;",
			}, []*zephyr.VNode{
				zephyr.StaticText("Complete"),
			}),
		),
		zephyr.RenderIf(c.todoItem.IsComplete,
			zephyr.Element("del", nil, []*zephyr.VNode{
				zephyr.DynamicText(c.todoItem.GetContent),
			}).RenderElse(zephyr.Element("span", nil, []*zephyr.VNode{
				zephyr.DynamicText(c.todoItem.GetContent),
			})),
		),
	})
}
