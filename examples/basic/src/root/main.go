package root

import (
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

// import core/reactivity

type RootComponent struct {
	// Extend Component struct
	*zephyr.BaseComponent
}

// ac.Set("counter", 0)

// func (rc RootComponent) SomeComputedProp() string {
// 	return "Current count: " + strconv.Itoa(rc.Get("counter").(int))
// }
