

Flag
-----

```
package main

import (
	"flag"
	"fmt"
)

func main() {

	debugPtr := flag.Bool("debug", false, "show debug log")
	pathPtr := flag.String("path", "~/", "the folder path")
	flag.Parse()

	if *debugPtr {
		fmt.Println("Debug mode")
	}
	fmt.Println(*pathPtr)

}

```


Args
----

