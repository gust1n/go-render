Render
========

A thin layer on top of go(lang)s std library html/template.
Adds "extends" and overwriting "define" with the most specific content.

Loosely based on https://github.com/daemonl/go_sweetpl

## Examples
Check out the [examples](https://github.com/gust1n/go-render/tree/master/examples) folder for a more thorough examples

### Extends
Just use the new "extends" keyword
	
	{{ extends "templates/layouts/fullwidth.html" }}

	{{ define "content" }}
	    content of the fullwidth template
	{{ end }}

This will also work with multi-level support, e.g. 
```index.html ---extends---> layouts/fullwidth.html ---extends---> base.html```

### Include
Simple include functionality, the package simply replaces the include block with the content of the file

    {{ define "content" }}
        content of the fullwidth template
        {{ include "includes/widgets/signup.html" }}
    {{ end }}

### Overwriting define / default value
Any "define" of the same "template" down the extend chain will overwrite the former content
This can be used to pass down default values for a {{ template }} like so

base.html

    <!DOCTYPE html>
	<html>
	  <head>
	    <title>{{ template "title" }}</title>
	  </head>
	</html>

	{{ define "title" }}Default Title{{ end }}

profile.html

    {{ extends "templates/base.html" }}
    {{ define "title" }}Hello World{{ end }}

This would produce panic in std lib parsing but now it works.

## Installation
```go get github.com/gust1n/go-render```

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

### Custom FuncMap
The Renderer type has a custom FuncMap that is injected into every template. Use it as such:

	var tmplErr error
    templates, tmplErr = render.LoadWithFuncMap("templates", template.FuncMap{
        "greet": func(name string) string {
            return fmt.Sprintf("Hello %s", name)
        },
    })
    if tmplErr != nil {
        panic(tmplErr)
    }

## Disclaimer
This is me experimenting and trying to make more use of go templates. I do NOT currently use this in production
