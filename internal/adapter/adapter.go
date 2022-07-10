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
	GenMarkdown(wd string, snippetPaths []string, imagePaths []string) ([]byte, error)
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

func (g *GatewayImpl) GenMarkdown(wd string, snippetPaths []string, imagePaths []string) ([]byte, error) {
	imageMap, err := g.parseImageFiles(wd, imagePaths)
	if err != nil {
		return nil, err
	}

	snippets, err := g.parseSnippetFiles(wd, snippetPaths, imageMap)
	if err != nil {
		return nil, err
	}

	mdBytes, err := g.proverbGenerator.GenMarkdown(snippets)
	if err != nil {
		return nil, fmt.Errorf("usecase error: %w", err)
	}
	return mdBytes, nil
}

func (g *GatewayImpl) parseSnippetFiles(wd string, snippetPaths []string, imageMap map[string]imageInfo) ([]usecase.Snippet, error) {
	snippets := make([]usecase.Snippet, len(snippetPaths))
	for i, snippetPath := range snippetPaths {
		snippet, err := g.parseSnippetFile(wd, snippetPath, imageMap)
		if err != nil {
			return nil, err
		}
		snippets[i] = snippet
	}
	return snippets, nil
}

func (g *GatewayImpl) parseSnippetFile(wd string, snippetPath string, imageMap map[string]imageInfo) (usecase.Snippet, error) {
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

	// If it doesn't exist, it will be empty - which works.
	imageInfo := imageMap[title]

	return usecase.Snippet{
		Title:        title,
		Group:        group,
		Tags:         tags,
		HTML:         template.HTML(htmlBytes),
		ImageRelPath: template.URL(imageInfo.RelPath),
		ImageText:    imageInfo.Text,
	}, nil
}

type imageInfo struct {
	RelPath string
	Text    string
}

func (g *GatewayImpl) parseImageFiles(wd string, imagePaths []string) (map[string]imageInfo, error) {
	imageMap := make(map[string]imageInfo)
	for _, imagePath := range imagePaths {
		name, info, err := g.parseImageFile(wd, imagePath)
		if err != nil {
			return nil, err
		}
		imageMap[name] = info
	}
	return imageMap, nil
}

func (g *GatewayImpl) parseImageFile(wd string, imagePath string) (string, imageInfo, error) {
	relPath, err := filepath.Rel(wd, imagePath)
	if err != nil {
		return "", imageInfo{}, fmt.Errorf("error determining %s relative to %s: %w", imagePath, wd, err)
	}

	base := filepath.Base(imagePath)
	fullName := strings.TrimSuffix(base, filepath.Ext(base))

	segs := strings.Split(fullName, "|")
	name := segs[0]
	var text string
	if len(segs) > 1 {
		text = segs[1]
	}

	return name, imageInfo{
		RelPath: relPath,
		Text:    text,
	}, nil
}
