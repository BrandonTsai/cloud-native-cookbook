---
title: "Call Linux cmmand"
author: Brandon Tsai
---


https://zetcode.com/golang/exec-command/


Basic Linux command example
---------------------------

Example:

```go
package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("ls", "-al", "/Users/brandon/Documents/projects/Go-Exercises")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	fmt.Println(string(stdout))
}
```

Call Linux Command with Pipe
----------------------------