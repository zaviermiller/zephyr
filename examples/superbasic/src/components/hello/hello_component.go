package hello

import (
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
	"github.com/zaviermiller/zephyr/pkg/core/vdom"
)

type HelloComponent struct {
	*zephyr.BaseComponent
}

func (c *HelloComponent) Init() {
	// greetee := rc.ReactiveString("test")
	// greetee := c.BindString("greetee") // func() *string

	c.DefineProps(map[string]interface{}{
		"recipient": nil,
	})

}

func (c *HelloComponent) Render() vdom.VNode {
	return vdom.BuildElem("h3", nil, []vdom.VNode{vdom.BuildText("Hello " + c.Get("greetee").(string))})

}
