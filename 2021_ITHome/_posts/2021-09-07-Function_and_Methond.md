---
title: "#7 Function and Custom Type Declarations"
author: Brandon Tsai
---
Introduction of Function
-------------------------

We can use the `func` keyword to declair a function as following: 

![](images/07_go_function_declair.png)


The input parameters and return type(s) are optional, such as `func main()`. Besides, a function can return multiple values as well. For example:

```go
package main

import (
	"errors"
	"fmt"
)

func main() {
	hello()

	rate, _ := getRate(53, 100)
	fmt.Println(rate, "%")
}

func hello() {
	fmt.Println("Hello world")
}

func getRate(numerator int, denominator int) (float64, error) {
	if denominator <= 0 {
		err := errors.New("denominator cannot be zero or negative")
		return 0.0, err
	}

	rate := float64(numerator) / float64(denominator) * 100
	return rate, nil
}
```


On above example, we return an `error` type value to indicate an abnormal situation and use blank identifier `_` to ignore some of the results from a function that returns multiple values.


We can put functions to another go file as well. For example:

utils.go

```go
package main

import "fmt"

func hello() {
	fmt.Println("Hello world")
}

```

main.go

```go
package main

func main() {
	hello()
}

```

Then you can run command "go run *.go" to test the result

```bash
$ go run *.go
Hello world

```



Custom Type Declarations and methods
--------

GO apporach instead of Object Orient approach

Why Custom type and reverse function?