package todo

import (
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

type TodoItem struct {
	zephyr.LiveStruct
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
