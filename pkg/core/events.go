package zephyr

import "syscall/js"

type Emitter interface {
	Call()
}

type DOMEvent struct {
	Target js.Value
}

func (e DOMEvent) Call() {

}
