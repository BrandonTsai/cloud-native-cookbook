---
title: "Packages and modules"
author: Brandon Tsai
---

Flat file structure
-------------------

For most small project, a flat folder structure is good enough and  recommended. Starting with flat structure can keep you focus on the delivering the service quickly without the influence of the complicated structure.

```
# example of flat structure.
helloworld/
- main.go
- main_test.go
- tools.go
- tools_test.go
```


Within a small application, we can just use the main package and functionalities imported from Go’s core library packages.
However, The code base could become bigger and more complex as time goes on. We might import massive third party packages from internet. Or, we might want to divide the complicate code base into semantic groups of functionality and unify them into a customized shared package within our project.

At this point, we can use [Go Moduless](https://go.dev/blog/using-go-modules) to manages these dependencies and the customize packages.


Go modules
----------


### go mod

module maintenance
`go mod init example.com/myproject` command is used to generated a go.mod file in the current directory, which will be viewed as the root directory of a module called

`go mod vendor`
Allows you to build the vendor folder using your go modules file

`go mod tidy`
Adds missing and remove unused modules
the go mod tidy command is used to add missing module dependencies into and remove unused module dependencies from the go.mod file by analyzing all the source code of the current project.

`go mod download <package_name>`
It downloads your module to the local cache. This helps when you are building your project with large size external packages. It’s better to cache them which certainly helps increasing build speed.
