package adapter

import (
	"fmt"
	"html/template"
	"os"

	"github.com/liampulles/proverb-gen/internal/usecase"
)

// --- Gateway front-matter ---

type Gateway interface {
	InitMarkdown(tmplPath string, title string, tags []string) ([]byte, error)
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

func (g *GatewayImpl) InitMarkdown(tmplPath string, title string, tags []string) ([]byte, error) {
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %w", tmplPath, err)
	}

	tmpl, err := template.New("proverb").Parse(string(tmplBytes))
	if err != nil {
		return nil, fmt.Errorf("could not parse template from %s: %w", tmplPath, err)
	}

	mdBytes, err := g.proverbGenerator.InitMarkdown(tmpl, title, tags)
	if err != nil {
		return nil, fmt.Errorf("usecase error: %w", err)
	}
	return mdBytes, nil
}
