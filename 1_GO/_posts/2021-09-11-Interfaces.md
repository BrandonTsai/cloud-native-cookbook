---
title: "#11 Interface & any"
author: Brandon Tsai
---

Refer https://gobyexample.com/interfaces


Go does not have Polymorphism like Object-Oriented language.
But it provide `Interface` type which defines a set of methods.
A variable of interface type can hold any value that implements those methods.


```go
// _Interfaces_ are named collections of method
// signatures.

package main

import (
	"fmt"
	"math"
)

// Here's a basic interface for geometric shapes.
type geometry interface {
	area() float64
}

// For our example we'll implement this interface on
// `rect` and `circle` types.
type rect struct {
	width, height float64
}
type circle struct {
	radius float64
}

// To implement an interface in Go, we just need to
// implement all the methods in the interface. Here we
// implement `geometry` on `rect`s.
func (r rect) area() float64 {
	return r.width * r.height
}


// The implementation for `circle`s.
func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}


// If a variable has an interface type, then we can call
// methods that are in the named interface. Here's a
// generic `measure` function taking advantage of this
// to work on any `geometry`.
func measure(g geometry) {
	fmt.Println(g)
	fmt.Println(g.area())
}

func main() {
	r1 := rect{width: 3, height: 4}
	r2 := rect{width: 10, height: 5}
	c1 := circle{radius: 5}
	geometryList := [3]geometry{r1, r2, c1}

	for _, g := range geometryList {
		measure(g)
	}
}
```


Array of interface objects/values?
-------

can we do following?
```go
a := [5]geometry{}

for _, value := range a {
    measure(value)
}
```


Exercise
--------

> 1. Implement an order system, use different structure for different type of dishes, and use interface to count the total price.


```go
package main

import (
	"fmt"
)

type dish interface {
	getPrice() float64
}

type coffee struct {
	name  string
	sugar int
	size  string
	price float64
}

type cake struct {
	name  string
	price float64
}

func (c coffee) getPrice() float64 {
	return c.price
}

func (c cake) getPrice() float64 {
	return c.price
}

func main() {
	largeLatte := coffee{name: "Latte", sugar: 1, size: "large", price: 4.5}
	smallFlatWhite := coffee{name: "Flat White", sugar: 0, size: "small", price: 3.5}
	cheeseCake := cake{name: "Cheese Cake", price: 6.5}
	tiramisu := cake{name: "Tiramisu", price: 8}
	order1 := [5]dish{largeLatte, smallFlatWhite, smallFlatWhite, cheeseCake, tiramisu}

	totalPrice := 0.0
	for _, i := range order1 {
		totalPrice += i.getPrice()
		fmt.Println(i)
	}
	fmt.Println(totalPrice)
}
```