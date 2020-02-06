# go-errors

[![GoDoc](https://godoc.org/github.com/gildas/go-errors?status.svg)](https://godoc.org/github.com/gildas/go-errors)
This is a library for handling errors in Go language.

|  |   |   |   |
---|---|---|---|
master | [![Build Status](https://dev.azure.com/keltiek/gildas/_apis/build/status/gildas.go-errors?branchName=master)](https://dev.azure.com/keltiek/gildas/_build/latest?definitionId=2&branchName=master) | [![Tests](https://img.shields.io/azure-devops/tests/keltiek/gildas/2/master)](https://dev.azure.com/keltiek/gildas/_build/latest?definitionId=2&branchName=master) | [![coverage](https://img.shields.io/azure-devops/coverage/keltiek/gildas/2/master)](https://dev.azure.com/keltiek/gildas/_build/latest?definitionId=2&branchName=master&view=codecoverage-tab)  
dev | [![Build Status](https://dev.azure.com/keltiek/gildas/_apis/build/status/gildas.go-errors?branchName=dev)](https://dev.azure.com/keltiek/gildas/_build/latest?definitionId=2&branchName=dev) | [![Tests](https://img.shields.io/azure-devops/tests/keltiek/gildas/2/dev)](https://dev.azure.com/keltiek/gildas/_build/latest?definitionId=2&branchName=dev) | [![coverage](https://img.shields.io/azure-devops/coverage/keltiek/gildas/2/dev)](https://dev.azure.com/keltiek/gildas/_build/latest?definitionId=2&branchName=dev&view=codecoverage-tab)  

## Usage

This package combines Go 1.13's [errors package](https://golang.org/pkg/errors) and [github.com/pkg/errors](https://github.com/pkg/errors).

All funcs from both packages are available.

On top of them, I added some sentinels for common errors I need all the time.

All sentinels are from the same type which contains various information. Whenever a sentinel is used, a StackTrace is also generated.

Here is how to use the errors:  
```go
func findme(stuff map[string]string, key string) (string, error) {
    if value, found := stuff[key]; found {
        return value, nil
    }
    return "", errors.NotFound.With(key).WithStack()
}

func main() {
    var allstuff[string]string
    //...
    value, err := findme("key1")
    if errors.Is(err, errors.NotFound) {
        fmt.Fprintf(os.Stderr, "Error: %+v", err)
        // This should print the error its wrapped content and a StackTrace.
    }
}
```

If you plan to do something with the content of the error, you would try that:  
```go
func main() {
    var allstuff[string]string
    //...
    value, err := findme("key1")
    if errors.Is(err, errors.NotFound) {
        var details *errors.Error
        if errors.As(err, &details) {
            fmt.Fprintf(os.Stderr, "Could not find %s", details.What)
        }
    }
}
```

To return an HTTP Status as an error, you could do this:  
```go
func doit() error {
    req, err := http.NewRequest(http.MethodGet, "http://www.acme.org")
    if err != nil {
        return errors.WithStack(err)
    }
    httpclient := http.DefaultClient()
    res, err := httpclient.Do(req)
    if err != nil {
        return errors.WithStack(err)
    }
    if res.StatusCode >= 400 {
        return errors.FromHTTPStatusCode(res.StatusCode)
    }
    return nil
}

func main() {
    err := doit()
    if (errors.Is(err, errors.HTTPBadRequest)) {
        // do something
    }
    // do something else
}
```
