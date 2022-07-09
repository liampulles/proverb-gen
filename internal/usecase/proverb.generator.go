package usecase

import (
	"bytes"
	"fmt"
	"html/template"
)

// --- ProverbGenerator front-matter ---

type ProverbGenerator interface {
	InitMarkdown(tmpl *template.Template, title string, tags []string) ([]byte, error)
}

type ProverbGeneratorImpl struct{}

var _ ProverbGenerator = &ProverbGeneratorImpl{}

// --- ProverbGenerator implementation ---

func NewProverbGeneratorImpl() *ProverbGeneratorImpl {
	return &ProverbGeneratorImpl{}
}

func (p *ProverbGeneratorImpl) InitMarkdown(tmpl *template.Template, title string, tags []string) ([]byte, error) {
	data := templateData{
		Title: title,
		Tags:  tags,
	}
	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, data); err != nil {
		return nil, fmt.Errorf("template error: %w", err)
	}

	return buf.Bytes(), nil
}

type templateData struct {
	Title string
	Tags  []string
}
