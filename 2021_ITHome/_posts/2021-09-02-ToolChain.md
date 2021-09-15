---
title: "#2 Installation and basic toolchain introduction"
author: Brandon Tsai
---

Installation
-----------

1. Download the package from https://golang.org/doc/install
2. Open the package file you downloaded and follow the prompts to install Go.
3. Open Terminal and verify that Go is installed by command `go version`



Go support in VSCode
--------------------

Install the go extension for VSCode to use feaures like code navigation, IntelliSense, symbol search, formatting and testing.

[](02-go-extension-vscode.png)


Refer https://code.visualstudio.com/docs/languages/go for more Hotkey usage of this extension.


The Go Toolchain
------------

We will use the example repo to help us understand the CLI commands for Go.

```bash
$ git clone https://go.googlesource.com/example
$ cd example
```

### go run

`go run` command is used for compiling and running a simple Go program.
It is not recommended to compile and run large Go projects.
Instead, you should use the commands `go build` for large project to build and executable the binary files.

```bash
$ cd hello
$ go run .
Hello, Go examples!
```

### go build

`go build` compiles the packages with their dependencies and writes the resulting executable to in current folder, but it does not install the results.

```bash
$ cd hello
$ go build
$ ls
hello		hello.go
$ ./hello 
Hello, Go examples!

```

**Cross Compile**

By default, `go build` command generate project binary for the current system’s platform.
You can use `GOOS` and `GOARCH` environment variables to specify the target system platform and architecture for different types of machine.
Example:

```bash
$ export GOARCH="amd64"
$ export GOOS="linux"
$ go build
```


**X Flag**

`go build -x` allows you to see how Go is performing the build process. This command helps to debug the process when you are doing a conditional build.

```bash
$ go build -x .
WORK=/var/folders/7z/k_5hdgqx3vq1qrtxrk3619rw0000gn/T/go-build2107802638
mkdir -p $WORK/b001/
cat >$WORK/b001/importcfg.link << 'EOF' # internal
packagefile golang.org/x/example/hello=/Users/brandon/Library/Caches/go-build/7c/7c302081fd9a21baf0f4d985062baa1ed447ad7a1e7a01e03c2c72871664a6fa-d
packagefile fmt=/usr/local/go/pkg/darwin_amd64/fmt.a
packagefile golang.org/x/example/stringutil=/Users/brandon/Library/Caches/go-build/d3/d35bf9fd2fb681cfcc6b3d61a8ea6df43aa584e3a39cfbcedcb8b8499b43c26e-d
packagefile runtime=/usr/local/go/pkg/darwin_amd64/runtime.a
packagefile errors=/usr/local/go/pkg/darwin_amd64/errors.a
packagefile internal/fmtsort=/usr/local/go/pkg/darwin_amd64/internal/fmtsort.a
packagefile io=/usr/local/go/pkg/darwin_amd64/io.a
packagefile math=/usr/local/go/pkg/darwin_amd64/math.a
packagefile os=/usr/local/go/pkg/darwin_amd64/os.a
packagefile reflect=/usr/local/go/pkg/darwin_amd64/reflect.a
packagefile strconv=/usr/local/go/pkg/darwin_amd64/strconv.a
packagefile sync=/usr/local/go/pkg/darwin_amd64/sync.a
packagefile unicode/utf8=/usr/local/go/pkg/darwin_amd64/unicode/utf8.a
packagefile internal/abi=/usr/local/go/pkg/darwin_amd64/internal/abi.a
packagefile internal/bytealg=/usr/local/go/pkg/darwin_amd64/internal/bytealg.a
packagefile internal/cpu=/usr/local/go/pkg/darwin_amd64/internal/cpu.a
packagefile internal/goexperiment=/usr/local/go/pkg/darwin_amd64/internal/goexperiment.a
packagefile runtime/internal/atomic=/usr/local/go/pkg/darwin_amd64/runtime/internal/atomic.a
packagefile runtime/internal/math=/usr/local/go/pkg/darwin_amd64/runtime/internal/math.a
packagefile runtime/internal/sys=/usr/local/go/pkg/darwin_amd64/runtime/internal/sys.a
packagefile internal/reflectlite=/usr/local/go/pkg/darwin_amd64/internal/reflectlite.a
packagefile sort=/usr/local/go/pkg/darwin_amd64/sort.a
packagefile math/bits=/usr/local/go/pkg/darwin_amd64/math/bits.a
packagefile internal/itoa=/usr/local/go/pkg/darwin_amd64/internal/itoa.a
packagefile internal/oserror=/usr/local/go/pkg/darwin_amd64/internal/oserror.a
packagefile internal/poll=/usr/local/go/pkg/darwin_amd64/internal/poll.a
packagefile internal/syscall/execenv=/usr/local/go/pkg/darwin_amd64/internal/syscall/execenv.a
packagefile internal/syscall/unix=/usr/local/go/pkg/darwin_amd64/internal/syscall/unix.a
packagefile internal/testlog=/usr/local/go/pkg/darwin_amd64/internal/testlog.a
packagefile internal/unsafeheader=/usr/local/go/pkg/darwin_amd64/internal/unsafeheader.a
packagefile io/fs=/usr/local/go/pkg/darwin_amd64/io/fs.a
packagefile sync/atomic=/usr/local/go/pkg/darwin_amd64/sync/atomic.a
packagefile syscall=/usr/local/go/pkg/darwin_amd64/syscall.a
packagefile time=/usr/local/go/pkg/darwin_amd64/time.a
packagefile unicode=/usr/local/go/pkg/darwin_amd64/unicode.a
packagefile internal/race=/usr/local/go/pkg/darwin_amd64/internal/race.a
packagefile path=/usr/local/go/pkg/darwin_amd64/path.a
EOF
mkdir -p $WORK/b001/exe/
cd .
/usr/local/go/pkg/tool/darwin_amd64/link -o $WORK/b001/exe/a.out -importcfg $WORK/b001/importcfg.link -buildmode=exe -buildid=Ti7_47Jz2to15h_7uJ7o/DGcSnpuydnHt86uC2z3J/oTJijITJicVoW5OA9_wI/Ti7_47Jz2to15h_7uJ7o -extld=clang /Users/brandon/Library/Caches/go-build/7c/7c302081fd9a21baf0f4d985062baa1ed447ad7a1e7a01e03c2c72871664a6fa-d
/usr/local/go/pkg/tool/darwin_amd64/buildid -w $WORK/b001/exe/a.out # internal
mv $WORK/b001/exe/a.out hello
rm -r $WORK/b001/

$ ls
hello		hello.go
```

### go install

Compile and install packages and dependencies to $GOPATH/bin folder.
Default GOPATH is:

 - $HOME/go on Unix-like systems
 - %USERPROFILE%\go on Windows

for example:

```
$ go install
$ ls ~/go/bin/
hello
```

### go get

Add/Upgrade/Downgrade/Remove a single dependency package to $GOPATH/pkg/.

Install package:

```bash
$ go get github.com/gobuffalo/flect
go: downloading github.com/gobuffalo/flect v0.2.3
go get: added github.com/gobuffalo/flect v0.2.3
```

Downgrade/Upgrade to particular version

```bash
$ go get github.com/gobuffalo/flect@v0.2.2
go: downloading github.com/gobuffalo/flect v0.2.2
go get: downgraded github.com/gobuffalo/flect v0.2.3 => v0.2.2
```

Uninstall package

```bash
$ go get github.com/gobuffalo/flect@none
go get: removed github.com/gobuffalo/flect v0.2.2
```

### go fmt

To format Go source code with a consistent coding style.

### go test

Executes test functions (Usually in test files). You can add the -v flag to lists all of the tests and their results.

```bash
$ cd stringutil
$ go test
PASS
ok  	golang.org/x/example/stringutil	0.004s
```

### go tool

This is the most powerful tool that Go provides natively. We all owe to Go’s engineers for coming up with this utility. Right from diagnosing the latency problems to finding bottlenecks with heap and threads in the code, Go profiling tools is something a production-grade system wants!

`go tool trace`

trace lets you collect your program traces and helps to visualize it.

Given a trace file produced by 'go test':

```
$ go test -trace=trace.out
PASS
ok  	golang.org/x/example/stringutil	0.005s
```

Open a web browser displaying trace:

```
$ go tool trace trace.out
2021/09/08 01:08:28 Parsing trace...
2021/09/08 01:08:28 Splitting trace...
2021/09/08 01:08:28 Opening browser. Trace viewer is listening on http://127.0.0.1:53082
```


`go tool pprof`

pprof lets you collect CPU profiles, traces, and heap profiles for your Go programs. 
Please refer https://pkg.go.dev/net/http/pprof for more details.