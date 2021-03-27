package todo

import (
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

type TodoItem struct {
	zephyr.LiveStructImpl
	Completed bool
	Content   string
}

func NewTodoItem(content string) TodoItem {
	item := TodoItem{Completed: false, Content: content}

	return item
}

func (item *TodoItem) Complete() {
	item.Completed = true
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
