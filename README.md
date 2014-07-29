Render
========

A thin layer on top of go(lang)s std library html/template.
Adds "extends" and overwriting "define" with the most specific content.

Loosely based on https://github.com/daemonl/go_sweetpl

## Examples

### Extends
Just use the new "extends" keyword
	{{ extends "templates/layouts/fullwidth.html" }}

	{{ define "content" }}
	    content of the fullwidth template
	{{ end }}

This will also work with multi-level support, e.g. 
index.html ---extends---> layouts/fullwidth.html ---extends---> base.html

Check out the examples folder for a more thorough example of this

### Overwriting define / default value
Any "define" of the same "template" down the extend chain will overwrite the former content
This can be used to pass down default values for a {{ template }} like so

base.html

    <!DOCTYPE html>
	<html>
	  <head>
	    <title>{{ template "title" }}</title>
	    {{ template "style" }}
	  </head>
	</html>

	{{ define "title" }}Default Title{{ end }}

profile.html

    {{ extends "templates/base.html" }}
    {{ define "title" }}Hello World{{ end }}

This would produce panic in std lib parsing but now it works.

## Installation
```go get ```

## Usage
    import "github.com/gust1n/go-render/render"

    var rnd *render.Renderer

    func main() {
    	// Now I can
    	rnd.ExecuteTemplate()
    	// probably into a http handler
    }

    func init() {
    	rnd = render.New()

		rnd.LoadTemplates("my/template/path")
    }

### Custom FuncMap
The Renderer type has a custom FuncMap that is injected into every template. Use it as such:

	rnd = render.New()
	rnd.FuncMap = template.FuncMap{
        "greet": func(name string) string {
            return fmt.Sprintf("Hello %s", name)
        },
    },