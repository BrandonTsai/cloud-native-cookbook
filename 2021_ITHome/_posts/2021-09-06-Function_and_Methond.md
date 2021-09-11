Call funstion from resuable packages.
Refer: https://www.callicoder.com/golang-packages/

pass arguments
return variables

Introduction of Function
-------------------------

```go
package main

import (
	"fmt"
)

func main() {

	message := getHelloWorld("Taiwan")
	fmt.Printf(message)
}

func getHelloWorld(country string) string {
	return "Hello " + country
}
```
