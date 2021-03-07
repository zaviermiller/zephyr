// display helpers
package display

var Red = "\033[31m"

var RedHighlight = "\u001b[41;1m"

var Green = "\033[32m"

var GreenHighlight = "\u001b[42m"

var Yellow = "\033[33m"

var Blue = "\033[34m"

var Purple = "\033[35m"

var Bright = "\033[1m"

var Dark = "\033[2m"

var Normal = "\033[0m"

var Clear = "\u001b[2K"

var White = "\u001b[37m"

var Black = "\u001b[30m"

// log helpers
func ZephyrInfo(msg string) string {
	return ("\n\r[" + Purple + "*" + Normal + "] " + msg)
}

func ZephyrError(msg string) string {
	return ("\n\r[" + Red + "*" + Normal + "] " + Red + "ERROR: " + msg + Normal)
}

// func ZenWeirdLog(msg string) {
// 	fmt.Print("\r[" + Purple + "*" + Normal + "] " + msg)
// }
