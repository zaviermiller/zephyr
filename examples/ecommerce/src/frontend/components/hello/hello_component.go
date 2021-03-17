package hello

import (
	zephyr "github.com/zaviermiller/zephyr/pkg/core"
)

type Props struct {
	// zephyr.BaseProps

	greetee string
}
type HelloComponent struct {
	zephyr.BaseComponent
}

func (c *HelloComponent) Init() {
	// greetee := rc.ReactiveString("test")
	// greetee := c.BindString("greetee") // func() *string

	// c.setGreetee = c.BindString(&c.greetee)

	// c.DefineProps(map[string]interface{}{
	// 	"recipient": nil,
	// })
	// c.DefineProps(Props{recipient: "default value"})

}

func (c *HelloComponent) Render() zephyr.VNode {
	return zephyr.BuildElem("h3", nil, []zephyr.VNode{zephyr.BuildText("Hello ")})

}
