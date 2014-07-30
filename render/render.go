package render

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var re_extends *regexp.Regexp = regexp.MustCompile("{{ extends [\"']?([^'\"}']*)[\"']? }}")
var re_defineTag *regexp.Regexp = regexp.MustCompile("{{ ?define \"([^\"]*)\" ?\"?([a-zA-Z0-9]*)?\"? ?}}")
var re_templateTag *regexp.Regexp = regexp.MustCompile("{{ ?template \"([^\"]*)\" ?([^ ]*)? ?}}")
var re_includeTag *regexp.Regexp = regexp.MustCompile("{{ ?include \"([^\"]*)\" ?([^ ]*)? ?}}")

var ErrTmplNotFound = errors.New("template not found")
var ErrTmplEmpty = errors.New("template is empty")

type Renderer struct {
	templates map[string]*template.Template
	FuncMap   map[string]interface{} //template.FuncMap
}

type namedTemplate struct {
	Name string
	Src  string
}

func file_content(path string) (string, error) {
	// Read the file content of the template
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	s := string(b)

	if len(s) < 1 {
		return "", ErrTmplEmpty
	}

	return s, nil
}

func (r *Renderer) add(stack *[]*namedTemplate, path string) error {
	tplSrc, err := file_content(path)
	if err != nil {
		return err
	}

	// Check if template contains "extend" keyword
	extendsMatches := re_extends.FindStringSubmatch(tplSrc)
	if len(extendsMatches) == 2 {
		// Perform recursion until no more extend's are found
		err := r.add(stack, extendsMatches[1])
		if err != nil {
			return err
		}
		// Remove the extend code
		tplSrc = re_extends.ReplaceAllString(tplSrc, "")
	}

	// Add included files
	tplSrc = addIncluded(tplSrc)

	// Add template to the stack
	*stack = append((*stack), &namedTemplate{
		Name: path,
		Src:  tplSrc,
	})
	return nil
}

// addIncluded recursively checks for include blocks and simply includes the file content
func addIncluded(src string) string {
	includedMatches := re_includeTag.FindStringSubmatch(src)
	if len(includedMatches) < 2 {
		return src
	}

	// Check if template contains "include" block
	src = re_includeTag.ReplaceAllStringFunc(src, func(raw string) string {
		parsed := re_includeTag.FindStringSubmatch(raw)
		includePath := parsed[1]

		content, err := file_content(includePath)
		if err != nil {
			panic(err)
		}

		return content
	})

	return addIncluded(src)
}

func (r *Renderer) assemble(path string) (*template.Template, error) {
	// The stack holds our template extend stack
	stack := []*namedTemplate{}

	err := r.add(&stack, path)
	if err != nil {
		return nil, err
	}

	// The rootTemplate holds our stack of parsed files
	var rootTemplate *template.Template

	// Replace 'define' blocks with UIDs to support overwriting the same block with the most specific template.
	// The 'defines' Map should contain the most specific definition with the block identifier as key (given that the stack was
	// properly ordered general -> specific)
	// This has to be separate loop to get all 'define' blocks before starting to replace definitions
	defines := map[string]string{}
	defineIdx := 0
	for _, namedTemplate := range stack {
		// If has a 'define' block
		namedTemplate.Src = re_defineTag.ReplaceAllStringFunc(namedTemplate.Src, func(raw string) string {
			parsed := re_defineTag.FindStringSubmatch(raw)
			blockName := fmt.Sprintf("BLOCK_%d", defineIdx)

			// Keep track of which definition belongs to which define statement
			defines[parsed[1]] = blockName
			defineIdx += 1

			return "{{ define \"" + blockName + "\" }}"
		})
	}

	for i, namedTemplate := range stack {
		// Replace all 'template' statements with the UID from above.
		namedTemplate.Src = re_templateTag.ReplaceAllStringFunc(namedTemplate.Src, func(raw string) string {
			parsed := re_templateTag.FindStringSubmatch(raw)
			origName := parsed[1]
			replacedName, ok := defines[origName]

			// Default the import var to . if not set
			dot := "."
			if len(parsed) == 3 && len(parsed[2]) > 0 {
				dot = parsed[2]
			}
			if ok {
				return fmt.Sprintf(`{{ template "%s" %s }}`, replacedName, dot)
			} else {
				return ""
			}
		})

		// Holds template we're currently dealing with
		var currentTmpl *template.Template

		// If first iteration, this should be the root template
		if i == 0 {
			currentTmpl = template.New(namedTemplate.Name)
			rootTemplate = currentTmpl
		} else { // Otherwise "inherit" from the root template
			currentTmpl = rootTemplate.New(namedTemplate.Name)
		}

		// Add our custom FuncMap (must be added before parsing)
		currentTmpl.Funcs(r.FuncMap)

		_, err := currentTmpl.Parse(namedTemplate.Src)
		if err != nil {
			return nil, err
		}
	}

	return rootTemplate, nil
}

func generateTemplateName(base, path string) string {
	return filepath.ToSlash(path[len(base)+1:])
}

// LoadTemplates loads and parses all *.html templates in specified directory.
// It also handles the recursive scan up the "extend"-chain
func (r *Renderer) LoadTemplates(dirPath string) error {
	if r.templates == nil {
		r.templates = make(map[string]*template.Template)
	}

	// Traverse the passed dir and parse all *.html templates
	filepath.Walk(dirPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if fi.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		tmpl, err := r.assemble(path)
		if err != nil {
			panic(err)
		}

		r.templates[generateTemplateName(dirPath, path)] = tmpl
		return nil
	})

	return nil
}

// ExecuteTemplate tries to inject the passed data into passed template name and write to the ResponseWriter
func (r *Renderer) ExecuteTemplate(rw http.ResponseWriter, name string, data interface{}) error {
	if _, found := r.templates[name]; !found {
		return ErrTmplNotFound
	}

	// Always (try to) write to buffer first to properly catch errors
	buf := new(bytes.Buffer)
	if err := r.templates[name].Execute(buf, data); err != nil {
		return err
	}
	_, err := buf.WriteTo(rw)
	return err
}

func New() *Renderer {
	return &Renderer{
		templates: make(map[string]*template.Template),
	}
}
