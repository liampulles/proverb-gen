package adapter

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/liampulles/proverb-gen/internal/usecase"
)

// --- Gateway front-matter ---

type Gateway interface {
	GenMarkdown(wd string, snippetPaths []string) ([]byte, error)
}

type GatewayImpl struct {
	proverbGenerator usecase.ProverbGenerator
}

var _ Gateway = &GatewayImpl{}

// --- Gateway implementation ---

func NewGatewayImpl(proverbGenerator usecase.ProverbGenerator) *GatewayImpl {
	return &GatewayImpl{
		proverbGenerator: proverbGenerator,
	}
}

func (g *GatewayImpl) GenMarkdown(wd string, snippetPaths []string) ([]byte, error) {
	snippets, err := g.parseFiles(wd, snippetPaths)
	if err != nil {
		return nil, err
	}

	mdBytes, err := g.proverbGenerator.GenMarkdown(snippets)
	if err != nil {
		return nil, fmt.Errorf("usecase error: %w", err)
	}
	return mdBytes, nil
}

func (g *GatewayImpl) parseFiles(wd string, snippetPaths []string) ([]usecase.Snippet, error) {
	snippets := make([]usecase.Snippet, len(snippetPaths))
	for i, snippetPath := range snippetPaths {
		snippet, err := g.parseFile(wd, snippetPath)
		if err != nil {
			return nil, err
		}
		snippets[i] = snippet
	}
	return snippets, nil
}

func (g *GatewayImpl) parseFile(wd string, snippetPath string) (usecase.Snippet, error) {
	htmlBytes, err := os.ReadFile(snippetPath)
	if err != nil {
		return usecase.Snippet{}, fmt.Errorf("could not read %s: %w", snippetPath, err)
	}

	relPath, err := filepath.Rel(wd, snippetPath)
	if err != nil {
		return usecase.Snippet{}, fmt.Errorf("error determining %s relative to %s: %w", snippetPath, wd, err)
	}

	dirs := strings.Split(relPath, string(os.PathSeparator))
	if len(dirs) != 3 {
		return usecase.Snippet{}, fmt.Errorf("unsupported snippet path %s: must be nested within one group folder", relPath)
	}
	group := dirs[1]
	name := dirs[2]

	nameSegs := strings.Split(strings.TrimSuffix(name, ".html"), "|")
	title := nameSegs[0]

	var tags []string
	if len(nameSegs) > 1 {
		tagsStr := nameSegs[1]
		tags = strings.Split(tagsStr, ",")
	}

	return usecase.Snippet{
		Title: title,
		Group: group,
		Tags:  tags,
		HTML:  template.HTML(htmlBytes),
	}, nil
}
