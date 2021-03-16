package cart

import "github.com/zaviermiller/zephyr/examples/ecommerce/src/frontend/items"

type Cart struct {
	Items        []items.Item
	CheckoutStep int
}

func Init() {

}

func NewCart(items []items.Item) Cart {
	return Cart{Items: items, CheckoutStep: -1}
}

func GetCurrentCart() Cart {
	// return cart or create one
	if false {
		//
	} else {

	}
}
