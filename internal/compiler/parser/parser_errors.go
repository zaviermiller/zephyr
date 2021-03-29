package parser

// ZephyrUnsupportedElementErr is thrown when one of (go-script, template, style)
// are root elements of a Node tree
type ZephyrUnsupportedElementErr struct {
	Element string
}

func (z ZephyrUnsupportedElementErr) Error() string {
	return "Element unsupported: " + z.Element
}
