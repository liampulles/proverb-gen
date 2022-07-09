package usecase

import (
	"bytes"
	"fmt"
	"html/template"
)

const proverbsTemplateStr = `---
title: Liam Pulles - Proverbs
description: Proverbs
shareable: true
---
<article>
    <hr>
	<header>
		<h1>Proverbs</h1>
	</header>
	<section>
		<p>Over the years I discover small pieces of cogent advice about programming and being a developer that I find useful.</p> 
		<p>This page tries to collect those thoughts - for my benefit and perhaps others. Email me if you have something good, I'll consider adding it. ;)</p>
	</section>
	{{ range $group, $snippets := . }}
	<section>
		<h2>{{ $group }}</h2>
		{{ range $snippets }}
		<section>
			<h3>{{ .Title }}</h3>
			{{ .HTML }}
			<ul class="proverb-tags">
				{{ range .Tags }}
				<li class="proverb-tag">#{{ . }}</li>
				{{ end }}
			</ul>
		</section>
		{{ end }}
	</section>
	{{ end }}
</article>`

var proverbsTemplate = template.Must(template.New("proverbs").Parse(proverbsTemplateStr))

// --- ProverbGenerator front-matter ---

type Snippet struct {
	Title string
	Group string
	Tags  []string
	HTML  template.HTML
}

type ProverbGenerator interface {
	GenMarkdown(snippets []Snippet) ([]byte, error)
}

type ProverbGeneratorImpl struct{}

var _ ProverbGenerator = &ProverbGeneratorImpl{}

// --- ProverbGenerator implementation ---

func NewProverbGeneratorImpl() *ProverbGeneratorImpl {
	return &ProverbGeneratorImpl{}
}

func (p *ProverbGeneratorImpl) GenMarkdown(snippets []Snippet) ([]byte, error) {
	grouped := p.groupSnippets(snippets)
	return p.templateMarkdown(grouped)
}

func (p *ProverbGeneratorImpl) groupSnippets(snippets []Snippet) map[string][]Snippet {
	grouped := make(map[string][]Snippet)
	for _, snippet := range snippets {
		grouped[snippet.Group] = append(grouped[snippet.Group], snippet)
	}
	return grouped
}

func (p *ProverbGeneratorImpl) templateMarkdown(grouped map[string][]Snippet) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := proverbsTemplate.Execute(buf, grouped); err != nil {
		return nil, fmt.Errorf("template error: %w", err)
	}
	return buf.Bytes(), nil
}
