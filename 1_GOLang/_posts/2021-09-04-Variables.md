---
title: "#4 Introduction to Variables"
author: Brandon Tsai
---

Decaliring Variables
---------------------

Let's update the hello.go example with a variable as following.

```go
package main

import "fmt"

func main() {
    var message string = "Hello world!"
    fmt.Println(message)
}
```

A variable declairation is composed of 3 components:

| var                                                   | hello                        | string                    |
| ----------------------------------------------------- | ---------------------------- | ------------------------- |
| Inform Go compier that we are creating a new variable | the name of the new variable | The associated data type. |


Go is Static Tye language. It cares the type of value that is going to be assigned to a variable. You can not assign a different type of value to a variable. for example:

```go
var message string = "Hello world"
message = 100 // ERROR! You can not assign a integer value to a string variable
```


Default Value
----------

Any variable declared without an initial value will have a default value assigned.

| Type              | Default Value |
| ----------------- | ------------- |
| bool              | false         |
| string            | ""            |
| int, int32, int64 | 0             |
| flost32, float64  | 0.0           |


```go
package main

import "fmt"

func main() {
	var message string
	fmt.Println(message)

	message = "Hello world!"
	fmt.Println(message)
}
```

Type Inference
--------------

Although Go is a Statically typed language, you do not need to specify the type of every variable you declare.
You can use `:=` to create a new varibale.
Go compiler will analysis the type of value and define the varible type as same as the value.


```go
package main

import "fmt"

func main() {
	// var hello string = "Hello world"
	message := "Hello World"
	message = "Hello Taiwan!"
	fmt.Println(hello)
}

```

Name Convention
---------------

The convention in Go is to use `MixedCaps` or `mixedCaps` (simply camelCase) rather than underscores to write multi-word names.

| Convention | Usage                                                        |
| ---------- | ------------------------------------------------------------ |
| MixedCaps  | If an identifier needs to be visible outside the package     |
| mixedCaps  | If you don't have the intention to use it in another package |


Type Conversion
---------------

**Numeric types conversion**

Golang doesn’t allow you to mix numeric types in an expression. For example:

```go
a := 2021
b := 3.1415

result := a + b // Error (invalid operation: mismatched types int and float64)

result := float64(a) + b // Works

```


**String and other basic data types conversion:**

Use the [strconv](https://pkg.go.dev/strconv) package to conversions to/from string representations of basic data types

```go
package main

import (
	"fmt"
	"reflect"
	"strconv"
)

func main() {
	// Convert String to Bool, Float, Integer
	b, err := strconv.ParseBool("true")
	fmt.Println(b, err, reflect.TypeOf(b))

	f, err := strconv.ParseFloat("3.1415", 64)
	fmt.Println(f, err, reflect.TypeOf(f))

	i, err := strconv.ParseInt("-42", 10, 64)
	fmt.Println(i, err, reflect.TypeOf(i))

	// Convert Bool, Float, Integer to String.
	sb := strconv.FormatBool(b)
	sf := strconv.FormatFloat(f, 'G', -1, 64)
	si := strconv.FormatInt(i, 10)

	fmt.Printf(sb + "\t" + sf + "\t" + si + "\n")
}
```

Constants Variables
-------------

We can declair `constant` variables to represent fixed values

for example

```go
const capital string = "Taipei"
const population2020 int32 = 23568378
const (
	a = 5 + population2020 // Valid
	b = population2020 / 5 // Valid
)

func main() {
	fmt.Printf("Hello %s !\n", capital)
	fmt.Printf("Taiwan's Population: %d!\n", population2020)
}
```


Access Environment Variables
----------------------

https://www.callicoder.com/go-read-write-environment-variables/

Sometime, we need to access the system environment variable at runtine so that we can make the same application work in different environments like Dev, UAT, and Production.

The `os` packages provide functions to work with environment variables, such as:

| Funtion            | Usage                                                                                                                                                   |
| ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Setenv(key, value) | Set an environment variable.                                                                                                                            |
| Getenv(key)        | Get an environment variable.                                                                                                                            |
| Unsetenv(key)      | Unset an environment variable.                                                                                                                          |
| LookupEnv(key)     | Get the value of an environment variable and a boolean indicating whether the environment variable is present or not. It returns a string and a boolean |

for example:

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	// Set Environment Variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "admin")
	os.Setenv("DB_PASS", "test123")

	// Get Environment Variables
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	fmt.Printf("Host: %s, Port: %s\n", user, pass)

	// Unset Environment Variable
	os.Unsetenv("DB_HOST")
	fmt.Printf("Try to get host: %s\n", os.Getenv("DB_HOST"))


	// Get the value of an environment variable and a boolean indicating whether the environment variable is set or not.
	database, ok := os.LookupEnv("DB_NAME")
	if !ok {
		fmt.Println("DB_NAME is not present")
	} else {
		fmt.Printf("Database Name: %s\n", database)
	}
}

```
