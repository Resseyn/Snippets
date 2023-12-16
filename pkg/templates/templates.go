package templates

import (
	"SnippetsTESTBYGUIDE/internal/loggers"
	"SnippetsTESTBYGUIDE/pkg/models"
	"html/template"
	"path/filepath"
)

var TemplateCache map[string]*template.Template

type TemplateData struct {
	User     *models.User
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func NewTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl.html"))
	if err != nil {
		loggers.Logger.Println(err)
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.ParseFiles(page)
		if err != nil {
			loggers.Logger.Println(err)
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl.html"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
