# go-errors

![GoVersion](https://img.shields.io/github/go-mod/go-version/gildas/go-errors)
[![GoDoc](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/gildas/go-errors) 
[![License](https://img.shields.io/github/license/gildas/go-errors)](https://github.com/gildas/go-errors/blob/master/LICENSE) 
[![Report](https://goreportcard.com/badge/github.com/gildas/go-errors)](https://goreportcard.com/report/github.com/gildas/go-errors)  

![master](https://img.shields.io/badge/branch-master-informational)
[![Test](https://github.com/gildas/go-errors/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/gildas/go-errors/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/gildas/go-errors/branch/master/graph/badge.svg?token=gFCzS9b7Mu)](https://codecov.io/gh/gildas/go-errors/branch/master)

![dev](https://img.shields.io/badge/branch-dev-informational)
[![Test](https://github.com/gildas/go-errors/actions/workflows/test.yml/badge.svg?branch=dev)](https://github.com/gildas/go-errors/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/gildas/go-errors/branch/dev/graph/badge.svg?token=gFCzS9b7Mu)](https://codecov.io/gh/gildas/go-errors/branch/dev)

This is a library for handling errors in Go language.

## Usage

This package combines Go standard [errors package](https://pkg.go.dev/errors) and [github.com/pkg/errors](https://github.com/pkg/errors).

All funcs from both packages are available.

On top of them, I added some sentinels for common errors I need all the time.

All sentinels are from the same type which contains various information. Whenever a sentinel is used, a StackTrace is also generated.

Here is how to use the errors:  
```go
func findme(stuff map[string]string, key string) (string, error) {
    if value, found := stuff[key]; found {
        return value, nil
    }
    return "", errors.NotFound.With("key", key)
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

Note: You should not use `var details errors.Error`, use the pointer version `var details *errors.Error` instead.

When several `errors.Error` are chained up, this can be used to extract the ones you want:
```go
func main() {
    var allstuff[string]string
    //...
    value, err := findme("key1")
    if errors.Is(err, errors.NotFound) {
        details := errors.NotFound.Clone()
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

You can chain multiple errors together (See [Example](https://pkg.go.dev/github.com/gildas/go-errors#example-WrapErrors)):

```go
func doit() error {
    // ... some processing
    var urlError *url.Error // errors from the net/url package

    return errors.WrapErrors(err, errors.NotFound.With("key", "key1"), urlError)
}
```

You can also use a `errors.MultiError` to collect multiple errors together (See [Example](https://pkg.go.dev/github.com/gildas/go-errors#example-MultiError)):

```go
me := errors.MultiError{}
me.Append(errors.NotFound.With("key", "key1"))
me.Append(errors.NotFound.With("key", "key2"))
fmt.Println(me.AsError())
```

Finally, `errors.Error` supports JSON serialization:
```go
err := errors.InvalidType.With("bogus")
payload, jerr := json.Marshal(err)
// ...
```
