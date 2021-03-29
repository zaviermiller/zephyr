package todo

import (
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

type TodoItem struct {
	zephyr.LiveStructImpl
	Completed bool
	Content   string
	Position  int
}

func NewTodoItem(content string, pos int) TodoItem {
	item := TodoItem{Completed: false, Content: content, Position: pos}

	return item
}

func (item *TodoItem) ToggleComplete() {
	item.Completed = !item.Completed
	item.Notify()
}

func (item *TodoItem) IsComplete(l zephyr.Listener) interface{} {
	item.Register(l)
	return item.Completed
}

func (item *TodoItem) GetContent(l zephyr.Listener) interface{} {
	item.Register(l)
	return item.Content
}

// func (item *TodoItem) UpdatePosition() interface{} {
// 	// item.Position =
// }
