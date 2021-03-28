package todo_list_item

import (
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
}

func (c *TodoListItemComponent) Render() *zephyr.VNode {
	return zephyr.Element("div", map[string]interface{}{"style": "margin-bottom: 10px;"}, []*zephyr.VNode{
		zephyr.RenderIf(func(l zephyr.Listener) interface{} {
			// fmt.Println("todo: ", c.todoItem, c.todoItem.IsComplete(l))
			return !c.todoItem.IsComplete(l).(bool)
		},
			zephyr.Element("button", map[string]interface{}{
				"style": "margin: 10px 10px;",
			}, []*zephyr.VNode{
				zephyr.StaticText("Complete"),
			}).BindEvent("click", func(e *zephyr.DOMEvent) { c.todoItem.Complete() }),
		),
		zephyr.RenderIf(c.todoItem.IsComplete,
			zephyr.Element("del", nil, []*zephyr.VNode{
				zephyr.DynamicText(c.todoItem.GetContent),
			}),
		).RenderElse(zephyr.Element("span", nil, []*zephyr.VNode{
			zephyr.DynamicText(c.todoItem.GetContent),
		})),
	})
}
