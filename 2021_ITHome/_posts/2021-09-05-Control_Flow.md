---
title: "#4 Introduction to Control Flow"
author: Brandon Tsai
---


If statement
------------


```go
if (condition1) {

} else if (condition2) {

} else {

}
```

Note that, the parentheses () can be omitted from an if statement. for example:

```go
package main

import "fmt"

func main() {
	var x = 7
	if num%3 == 0 && num%5 == 0 {
		fmt.Printf("%d is divisible by both 3 and 5\n", x)
	} else if x < 0 {
		fmt.Printf("%d is negative\n", x)
	} else {
		fmt.Printf("%d is positive\n", x)
	}
}

```


Switch statement
----------------

Switch evaluates all the cases from top to bottom until a case succeeds. Once a case succeeds, it runs the block of code specified in that case and then stops the evaluation.

```go
package main

import "fmt"

func main() {
	code := "AT"
	switch code {
	case "TW", "ROC":
		fmt.Println("Taiwan")
	case "AU":
		fmt.Println("Australia")
	case "AT":
		{
			fmt.Println("Austria")
			fmt.Println("There is no kangaroos in Austria")
		}
	default:
		fmt.Println("Other Country")
	}
}

```

For loop
--------


You can use `continue` statement to skip running the loop body midway and continue to the next iteration of the loop.
The `break` statement can be used to stop a loop before its normal termination.

```go
package main

import "fmt"

func main() {
	for i := 100; i < 150; i++ {
		if i%2 == 0 {
			continue
		} else if i%7 == 0 {
			fmt.Printf("%d is a the first multiple of 7\n", i)
			break
		}
		fmt.Printf("%d is a the odd number\n", i)
	}
}
```