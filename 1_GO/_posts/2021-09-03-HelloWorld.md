---
title: "#3 Deep into Hello World!"
author: Brandon Tsai
---

Let's start by understanding the hello.go example

```go
package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	fmt.Println(stringutil.Reverse("!selpmaxe oG ,olleH"))
}
```

hello.go include three parts: `packages`, `import` and `func main()`.

### packages

The first statement of every go source file must be a package declaration.
Package is a way for collecting related Go code together.
A packages can have many files inside of it. For example:

main.go

```go
// main.go
package main
import (
	"fmt"
)

func main() {
	fmt.Println("Hello MAIN!")
}
```

help.go

```go
// help.go
package main
import (
	"fmt"
)

func help() {
	fmt.Println("Help!")
}
```


There are two type of packages

- Executable: Generate a file that we can run.
- Reuseable: Code dependencies or libraries that we can reuse.

Use package name `main` to specify this package can be compiled and then executed. Inside the main package, it must has a func call 'main'

Other package name defines a package that cab be used as a dependency.
We will discuss how to use reuseable packages later.



### import

Use to import code from other packages.
Although there are some standard libraries such as math, fmt ,debug ... etc.
We still need to use import to link the library to our package.
you can find out more standard libraries on [https://pkg.go.dev/std](https://pkg.go.dev/std)

Besides, you can import third party packages from internet as well.
for example `import "golang.org/x/example/stringutil"`
This required you to install the package via command `go get <package>` before building your own package.


### func main()

The entry of our execuable code. This function is required for main package.


GOPATH
-------

GOPATH is a variable that defines the root of your workspace. By default, the workspace directory is a directory that is named go within your user home directory (~/go for Linux and MacOS, %USERPROFILE%/go for Windows). GOPATH stores your code base and all the files that are necessary for your development. You can use another directory as your workspace by configuring GOPATH for different scopes. GOPATH is the root of your workspace and contains the following folders:

src/: location of Go source code (for example, .go, .c, .g, .s).

pkg/: location of compiled package code (for example, .a).

bin/: location of compiled executable programs built by Go.

