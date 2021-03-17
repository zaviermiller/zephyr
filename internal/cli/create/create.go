package create

import "flag"

func Run(prjectName string) {
	runtimeOnly := flag.Bool("runtime", false, "Only use the Zephyr runtime (100% Go)")

	flag.Parse()

	if !*runtimeOnly {

	}
}
