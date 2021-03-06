package hello

import (
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
	"github.com/zaviermiller/zephyr/pkg/core/vdom"
)

type Props struct {
	message string
}

type HelloComponent struct {
	*zephyr.BaseComponent

	//props????
	// message string
}

func (c *HelloComponent) Init() {
	c.DefineData(map[string]interface{}{
		"greetee": "Zephyr",
	})
}

func (c *HelloComponent) Render() vdom.VNode {
	return vdom.BuildElem("h3", nil, []vdom.VNode{vdom.BuildText("Hello " + c.Get("greetee").(string))})

}
