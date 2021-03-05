package zephyr

type Listener interface {
	Update()
}

type Subject interface {
	Register(l Listener)
	Notify()
}

// IMPLEMENTATIONS =-=-=

// ReactiveData is the struct that holds the data
// being listened to, and notifies listeners of
// any changes.
type ReactiveData struct {
	// Type is a field to provide artificial type-checking
	Type      string
	Data      interface{}
	Listeners []Listener
}

func newReactiveData(dataType string, data ...interface{}) ReactiveData {
	var rd ReactiveData
	if len(data) > 1 {
		panic("Too much data!")
	}
	if len(data) > 0 {
		rd = ReactiveData{Type: dataType, Data: data[0]}
	} else {
		rd = ReactiveData{Type: dataType, Data: nil}
	}

	return rd
}

func (rd *ReactiveData) Register(l Listener) {
	// make sure listener doesnt exist
	for _, val := range rd.Listeners {
		if l.(ComponentListener).ID == val.(ComponentListener).ID {
			return
		}
	}
	rd.Listeners = append(rd.Listeners, l)
}

func (rd ReactiveData) Notify() {
	for _, l := range rd.Listeners {
		l.Update()
	}
}

type ComponentListener struct {
	ID      string
	Updater func()
}

func (l ComponentListener) Update() {
	// re-render component on update
	l.Updater()
}
