


Go modules
----------

https://tutorialedge.net/golang/go-project-structure-best-practices/

https://golang.org/ref/mod

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
