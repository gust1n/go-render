Render
========

A convenience template loader for stdlib's html/template. Takes away the pain of manually having to parse all files for a specific template.

## Installation
```go get github.com/gust1n/go-render/render```

## Usage
    import (
        "github.com/gust1n/go-render/render"
    )

    func main() {
        templates, err := render.Load("templates")
        if err != nil {
            panic(err)
        }
        // Now I have a map[string]*template.Template to use in my handlers
    }

## Examples
Check out the [examples](https://github.com/gust1n/go-render/tree/master/examples) folder for some examples

### Extends
Just use the standard template keyword with a *.html file path

index.html
	
	{{ template "templates/layouts/fullwidth.html" }}

	{{ define "content" }}
	    content of index to be inserted into the fullwidth template
	{{ end }}

This will also work with multi-level support, e.g. 
```index.html ---extends---> layouts/fullwidth.html ---extends---> base.html```

### Include
Automatically parse the right file just by writing the path to it

    {{ define "content" }}
        content of the fullwidth template
        {{ template "includes/widgets/signup.html" . }}
    {{ end }}

### Overwriting define / default value
Any "define" of the same "template" down the extend chain will overwrite the former content
This can be used to define default values for a {{ template }} like so

base.html

    <!DOCTYPE html>
	<html>
	  <head>
	    <title>{{ template "title" }}</title>
	  </head>
	</html>

	{{ define "title" }}Default Title{{ end }}

profile.html

    {{ template "templates/base.html" }}
    {{ define "title" }}Hello World{{ end }}

This would produce panic in std lib parsing but now it works by simply renaming the define's further down the chain not to interrupt the most specific one.

### Custom FuncMap
The Renderer can load a custom FuncMap that is injected into every template. Use it as such:

	var tmplErr error
    templates, tmplErr = render.LoadWithFuncMap("templates", template.FuncMap{
        "greet": func(name string) string {
            return fmt.Sprintf("Hello %s", name)
        },
    })
    if tmplErr != nil {
        panic(tmplErr)
    }

## Credits
Inspired by https://github.com/daemonl/go_sweetpl

## Disclaimer
This is me experimenting and trying to make more use of go templates. I do NOT currently use this in production
