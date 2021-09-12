---
title: "#6 Array and Slice"
author: Brandon Tsai
---

Go has two type of structure to handle list:

- **Array**: Fixed length list of elements.
- **Slice**: flexiable length list that can groe or shrink.

Every element in a Array or Slice must be of same types. The elements are stored sequentially and can be accessed using their index


Array
-------

Declairing an array:

```go
// An array of 5 integers with default value 0.
a := [5]int{}
fmt.Println(a)

// An array with initial value
b := [5]int{2, 4, 6, 8, 10}
fmt.Println(b)

// An array of 5 integers and initial 2 elements
// output = [2 4 0 0 0]
b := [5]int{2, 4}
fmt.Println(b)

//Letting Go compiler infer the length of the array
a := [...]int{1, 3, 5, 7, 9}
fmt.Println(a)
```

Array is pass by value, this means that the value of elements are copied when assigning one array to another. If you make any changes to this copied array, the original one won’t be affected and will remain unchanged. For example:


```go
a := [5]int{2, 4, 6, 8, 10}
b := a

b[0] = 0

fmt.Println(a) // output = [2 , 4, 6, 8, 10]
fmt.Println(b) // output = [0 , 4, 6, 8, 10]
```

Iterating over an array:

```go
a := [5]int{2, 4, 6, 8, 10}

// get value with index
for i := 0; i < len(a); i++ {
    fmt.Println(a[i])
}

// you can use 'range' operator to get index and value of elements.
for index, value := range a {
    fmt.Println(index, value)
}

/*
    Go compiler doesn’t allow creating variables that are never used.
    You can fix this by using an _ (underscore) in place of index
*/
for _, value := range a {
    fmt.Println(value)
}
```


Multidimensional arrays:

```go
a := [2][3]int{
    {1, 3, 5},
    {2, 4, 6}, // This trailing comma is mandatory
}

for i := 0; i < len(a); i++ {
    for j := 0; j < len(a[i]); j++ {
        fmt.Println(a[i][j])
    }
}

// you can use 'range' operator to get index and value of elements.
for _, raw := range a {
    for index, value := range raw {
        fmt.Println(index, value)
    }
}
```


Slice
------

Declair slices:

```go
// Declair an empty slice which length is 0
a := []int{}
fmt.Println(a, len(a))

// Declair a slice of lengh 5 and fill with default value via "make" . function
b := make([]int, 5)
fmt.Println(b, len(b))


// Decalre a slice with initial value
c := []int{2, 4, 6, 8}
fmt.Println(c, len(c))
```

Modify slices:

```go
// Declair an empty slice
a := []int{}
b := []int{2, 4}
fmt.Println(a, len(a))

// add elements via append function
a = append(a, 1)
a = append(a, 2, 3)
fmt.Println(a, len(a))

// apend a slice to another slice.
// the ... lets you pass multiple arguments to a variadic function from a slice
a = append(a, b...)
fmt.Println(a, len(a))

// the is no function for removing elements from slices
// you can use append function to re-slice
index := 3 // the index of element you want to remove.
a = append(a[:index], a[index+1:]...)
fmt.Println(a, len(a))
```

Unlike Array, slice is pass by reference, this means when any changes being made to this copied array, the original one is affected as well. For example:

```go
// Create a slice refer to another slice
a := []int{1, 2, 3, 4}
b := a
b[0] = 2
fmt.Println(a) //output = [2 2 3 4]
fmt.Println(b) //output = [2 2 3 4]

// Create a slice refer to part of another slice
c := a[1:4]
c[0] = 5
fmt.Println(a) //output = [2 5 3 4]
fmt.Println(b) //output = [2 5 3 4]
fmt.Println(c) //output = [5 3 4]
```

If you want to copy the value from one slice to another slice, you should use `copy(dst, src)` function. It copies elements from the source to the destination and returns the number of elements that are copied. The number of elements copied will be the minimum of len(src) and len(dst)

```go
a := []int{1, 2, 3, 4}
b := make([]int, len(a))

copy(b, a)
b[0] = 5

fmt.Println(a, reflect.TypeOf(a).Kind()) // output = [1 2 3 4]
fmt.Println(b, reflect.TypeOf(b).Kind()) // output = [5 2 3 4]

c := []int{}
copy(c, a)                               // nothing copied to c because len(c)= 0
fmt.Println(a, reflect.TypeOf(a).Kind()) // output = [1 2 3 4]
fmt.Println(c, reflect.TypeOf(c).Kind()) // output = []
```

Slice can be used to refer to an array as well.

```go
a := [5]int{0, 2, 4, 6, 8}

// Create an array copy from another array
a1 := a // pass by value
a1[0] = 1
fmt.Println(a, len(a), reflect.TypeOf(a).Kind())    //output = [0 2 4 6 8]
fmt.Println(a1, len(a1), reflect.TypeOf(a1).Kind()) //output = [1 2 4 6 8]

// Create an array copy from part of another array
var a2 [4]int
copy(a2[:], a[0:4])
a2[0] = 2
fmt.Println(a, len(a), reflect.TypeOf(a).Kind())    //output = [0 2 4 6 8]
fmt.Println(a2, len(a2), reflect.TypeOf(a2).Kind()) //output = [2 2 4 6]

// Create a slice refer to an array
s1 := a[:] // pass by reference
s1[0] = 3
fmt.Println(a, len(a), reflect.TypeOf(a).Kind())    //output = [3 2 4 6 8]
fmt.Println(a1, len(a1), reflect.TypeOf(a1).Kind()) //output = [1 2 4 6 8]
fmt.Println(a2, len(a2), reflect.TypeOf(a2).Kind()) //output = [2 2 4 6]
fmt.Println(s1, len(s1), reflect.TypeOf(s1).Kind()) //output = [3 2 4 6 8]

// Create a slice refer to part of an array
s2 := a[0:4]
s2[0] = 4
fmt.Println(a, len(a), reflect.TypeOf(a).Kind())    //output = [4 2 4 6 8]
fmt.Println(a1, len(a1), reflect.TypeOf(a1).Kind()) //output = [1 2 4 6 8]
fmt.Println(a2, len(a2), reflect.TypeOf(a2).Kind()) //output = [2 2 4 6]
fmt.Println(s1, len(s1), reflect.TypeOf(s1).Kind()) //output = [4 2 4 6 8]
fmt.Println(s2, len(s2), reflect.TypeOf(s2).Kind()) //output = [4 2 4 6]
```


Iterating over a slice:

```go
s := []int{2, 4, 6, 8, 10}

// get value with index
for i := 0; i < len(s); i++ {
    fmt.Println(s[i])
}

// you can use 'range' operator as well.
for _, value := range s {
    fmt.Println(value)
}
```


Multidimensional slices:

```go
s := [][]int{
    {2, 4, 6, 8, 10},
    {1, 3, 5},
}

for i := 0; i < len(s); i++ {
    for _, value := range s[i] {
        fmt.Println(value)
    }
}
```
