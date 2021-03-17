package zephyr

// Reactive helpers =-=
// The Bind[Type] methods take a pointer to their respective
// type and returns a function used to set the data that will
// also rerender the component. This could change if I end up changing
// the DOM interactions and stuff. I probably dont have it set up right.
// Oh and if this doesnt change, cant wait for generics!

func (c *BaseComponent) BindInt(intPtr *int, args ...interface{}) func(int) {
	// This function could also be implemented by vdom nodes, I wonder if vnodes
	// should all extend their own basecomponent? that would allow to not rerender
	// the whole dom.. yes def

	// check if we need to also set the data
	if len(args) == 1 {
		*intPtr = args[0].(int)
	} else if len(args) > 1 {
		panic("too many args dont abuse me or ill change!!")
	}

	// bind the component listener to the variable's reactive data
	// (really just need the interface impl, can prolly deprecate RD)
	rd := NewRD(*intPtr)

	// though, I do use it here for type safety, but this line could
	// be dumb lmao. i need to figure out when to be seemingly overexplicity
	if _, ok := rd.Data.(int); !ok {
		panic("huh...?")
	}

	rd.Register(&c.Listener)

	// return setter function
	return func(newVal int) {
		// change the ptrs val so that we can just use a regular ol var
		*intPtr = newVal
		// notify the component of change.
		rd.Notify()
	}
}

// repeat for all types until generics
func (c *BaseComponent) BindString(strPtr *string, args ...interface{}) func(string) {
	// This function could also be implemented by vdom nodes, I wonder if vnodes
	// should all extend their own basecomponent? that would allow to not rerender
	// the whole dom.. yes def

	// check if we need to also set the data
	if len(args) == 1 {
		*strPtr = args[0].(string)
	} else if len(args) > 1 {
		panic("too many args dont abuse me or ill change!!")
	}

	// bind the component listener to the variable's reactive data
	// (really just need the interface impl, can prolly deprecate RD)
	rd := NewRD(*strPtr)

	// though, I do use it here for type safety, but this line could
	// be dumb lmao. i need to figure out when to be seemingly overexplicity
	if _, ok := rd.Data.(string); !ok {
		panic("huh...?")
	}

	rd.Register(c.Listener)

	// return setter function
	return func(newVal string) {
		// change the ptrs val so that we can just use a regular ol var
		*strPtr = newVal
		// notify the component of change.
		rd.Notify()
	}
}

func (c *BaseComponent) BindList(arrLocation interface{}) func(interface{}) {
	// This function could also be implemented by vdom nodes, I wonder if vnodes
	// should all extend their own basecomponent? that would allow to not rerender
	// the whole dom.. yes def

	switch arrLocation.(type) {
	case *[]int:
		// bind the component listener to the variable's reactive data
		// (really just need the interface impl, can prolly deprecate RD)
		rd := NewRD(*arrLocation.(*[]int))

		// though, I do use it here for type safety, but this line could
		// be dumb lmao. i need to figure out when to be seemingly overexplicity
		// if _, ok := rd.Data.(string); !ok {
		// 	panic("huh...?")
		// }

		rd.Register(&c.Listener)

		// return setter function
		return func(newArr interface{}) {
			// change the ptrs val so that we can just use a regular ol var
			arrLocation = newArr.(*[]int)
			// fmt.Println(arrLocation)
			// notify the component of change.
			rd.Notify()
		}
	case *[]string:
		// bind the component listener to the variable's reactive data
		// (really just need the interface impl, can prolly deprecate RD)
		rd := NewRD(*arrLocation.(*[]string))

		// though, I do use it here for type safety, but this line could
		// be dumb lmao. i need to figure out when to be seemingly overexplicity
		// if _, ok := rd.Data.(string); !ok {
		// 	panic("huh...?")
		// }

		rd.Register(&c.Listener)

		// return setter function
		return func(newArr interface{}) {
			// change the ptrs val so that we can just use a regular ol var
			arrLocation = newArr.(*[]string)
			// fmt.Println(arrLocation)
			// notify the component of change.
			rd.Notify()
		}
	default:
		return nil

	}
}

// func (c *BaseComponent) CompareData(node *VNode) bool {
// 	for i, data := range c.data {
// 		if !reflect.DeepEqual(data, node.Component.getBase().data[i]) {
// 			return false
// 		}
// 	}
// 	fmt.Println(c, " cached")
// 	return true
// }
