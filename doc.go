/*
This package combines Go standard errors package (https://pkg.go.dev/errors) and [github.com/pkg/errors](https://github.com/pkg/errors).

All funcs from both packages are available.

On top of them, I added some sentinels for common errors I need all the time.

All sentinels are from the same type which contains various information. Whenever a sentinel is used, a StackTrace is also generated.

Here is how to use the errors:

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

If you plan to do something with the content of the error, you would try that:

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

Note: You should not use "var details errors.Error", use the pointer version "var details *errors.Error" instead.

When several `errors.Error` are chained up, this can be used to extract the ones you want:

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

To return an HTTP Status as an error, you could do this:

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

You can also add more than one _cause_ to an `errors.Error`, turning it into a _multi-error_ container:

	err := errors.Error{}
	err.WithCause(errors.ArgumentInvalid.With("key", "value"))
	err.WithCause(errors.ArgumentMissing.With("key"))
	err.WithCause(fmt.Errorf("some simple string error"))

Finally, errors.Error supports JSON serialization.

	err := errors.InvalidType.With("bogus")
	payload, jerr := json.Marshal(err)
	// ...
*/
package errors
