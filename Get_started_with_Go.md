# Get started with Go

## Enable dependency tracking for your code

When your code imports packages contained in other modules, you
manage those dependencies through your code's own module. That
module is defined by a `go.mod` file that tracks the modules that
provide those packages. That `go.mod` file stays with your code,
including in your source code repository.

Command to enable dependency tracking:

`go mod init [module_path]`

creates a `go.mod` file. Module path will typically be the repository location where your
source code will be kept. For example, the module path might be `github.com/mymodule`.
If you plan to publish your module for others to use, the module path must be a
location from which Go tools can download your module.

For tutorial, using `example/hello`:

`go mod init example/hello`

## Call code in an external package

1. Visit `pkg.go.dev` and search for external package.

Tutorial uses "quote" package.

2. Locate and click the external package in search results.

For tutorial this is `rsc.io/quote` package.

3. In the **Documentation** section, under **Index**, note the list of functions you can call from your code.

4. At the top of this page, note that package quote is included
in the `rsc.io/quote` module.

5. In your Go code, import and use the external package.

For tutorial:

```go
import "rsc.io/quote"

func main() {
    fmt.Println(quote.Go())
}
```

6. Add new module requirements and sums.

Go will add the `quote` module as a requirement, as well as a `go.sum` file for use in authenticating
the module. Command:

`go mod tidy`

7. Run your code

`go run .`

# Create a Go module

## Start a module that others can use

### Start your module using the `go mod init` command.

Run the go mod init command, giving it your module path - here, use `example.com/greetings`.
If you publish a module, this must be a path from which your module can be downloaded by Go tools.
That would be your code's repository.

```text
go mod init example.com/greetings
```

### Add source code

File `greetings.go`

```go
package greetings

import "fmt"

// Hello returns a greeting for the named person.
func Hello(name string) string {
    message := fmt.Sprintf("Hi, %v. Welcome!", name)
    return message
}
```

## Call your code from another module

### Enable dependency tracking for your code

```text
go mod init example.com/hello
```

### Create file `hello.go`:

```go
package main

import (
    "fmt"
    "example.com/greetings"
)

func main() {
    message := greetings.Hello("Gladys")    // use module
    fmt.Println(message)
}
```

### Replace dependency with local path

For production use, you’d publish the example.com/greetings module from its repository
(with a module path that reflected its published location), where Go tools could find it to
download it. For now, because you haven't published the module yet, you need to adapt the
`example.com/hello` module so it can find the `example.com/greetings` code on your
local file system.

Command:

`go mod edit -replace example.com/greetings=../greetings`

specifies that `example.com/greetings` should be replaced with `../greetings` for the purpose
of locating the dependency.

Command:

`go mod tidy`

synchronize the `example.com/hello` module's dependencies.

To reference a published module, a `go.mod` file would typically omit the replace directive
and use a require directive with a tagged version number at the end.

`require example.com/greetings v1.1.0`

## Add a test

### In the project directory, create a file called `*_test.go`

Ending a file's name with `_test.go` tells the `go test` command that this file contains
test functions.

Test file `greetings_test.go` (example):

```go
package greetings

import (
    "testing"
    "regexp"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {
    name := "Gladys"
    want := regexp.MustCompile(`\b`+name+`\b`)
    msg, err := Hello("Gladys")
    if !want.MatchString(msg) || err != nil {
        t.Fatalf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
    }
}

// TestHelloEmpty calls greetings.Hello with an empty string,
// checking for an error.
func TestHelloEmpty(t *testing.T) {
    msg, err := Hello("")
    if msg != "" || err == nil {
        t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
    }
}
```

### Run the `go test` command to execute the test

## Compile and install the application

`go build` - compiles the packages, along with their dependencies.

`go install` - compiles and installs the packages.

### For install

1. Discover the Go install path, where the go command will install the current package

`go list -f '{{.Target}}'`

2. Add the Go install directory to your system's shell path

Linux or Mac: `export PATH=$PATH:/path/to/your/install/directory`

Windows: `set PATH=%PATH%;C:\path\to\your\install\directory`

Or add GOBIN path (like `$HOME/bin`):

Linux or Mac: `go env -w GOBIN=/path/to/your/bin`

Windows: `go env -w GOBIN=C:\path\to\your\bin`

3. Once you've updated the shell path, run the `go install` command to compile and install
the package.

4. Run the package
