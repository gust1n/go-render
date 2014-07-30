package render

import "html/template"

// Load prepares and parses all templates from the passed basePath
func Load(basePath string) (map[string]*template.Template, error) {
	rnd := new(basePath)
	if err := rnd.loadTemplates(); err != nil {
		return nil, err
	}

	return rnd.templates, nil
}

// LoadWithFuncMap prepares and parses all templates from the passed basePath and injects
// a custom template.FuncMap into each template
func LoadWithFuncMap(basePath string, funcMap template.FuncMap) (map[string]*template.Template, error) {
	rnd := new(basePath)
	rnd.funcMap = funcMap

	if err := rnd.loadTemplates(); err != nil {
		return nil, err
	}

	return rnd.templates, nil
}
