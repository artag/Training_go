# Basics of Go

[Get started with Go](Get_started_with_Go.md)

[A Tour of Go](A_Tour_of_Go.md)

## Tools

The main Go distribution includes tools for building, testing, and analyzing code:

* `go build`, which builds Go binaries using only information in the source files themselves, no separate makefiles
* `go test`, for unit testing and microbenchmarks
* `go fmt`, for formatting code
* `go install`, for retrieving and installing remote packages
* `go vet`, a static analyzer looking for potential errors in code
* `go run`, a shortcut for building and executing code
* `godoc`, for displaying documentation or serving it via HTTP
* `gorename`, for renaming variables, functions, and so on in a type-safe way
* `go generate`, a standard way to invoke code generators
* `go mod`, for creating a new module, adding dependencies, upgrading dependencies, etc.

It also includes *profiling* and *debugging* support, *runtime* instrumentation (for example,
to track *garbage collection* pauses), and a *race condition* tester.

Third-party tools (adds to the standard distribution):

* `gocode`, enables code autocompletion in many text editors
* `goimports`, automatically adds/removes package imports as needed
* `errcheck`, detects code that might unintentionally ignore errors.
