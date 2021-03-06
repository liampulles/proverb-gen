package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/liampulles/proverb-gen/internal/adapter"
)

const proverbsRelDir = "_proverbs"

var proverbImagesRelDir = path.Join("images", "proverbs")

// --- Engine front-matter ---

type Engine interface {
	Run(wd string) error
}

type EngineImpl struct {
	gateway adapter.Gateway
}

var _ Engine = &EngineImpl{}

// --- Engine implementation ---

func NewEngineImpl(gateway adapter.Gateway) *EngineImpl {
	return &EngineImpl{
		gateway: gateway,
	}
}

func (e *EngineImpl) Run(wd string) error {
	imagePaths, err := e.readImagePaths(wd)
	if err != nil {
		return err
	}

	snippetPaths, err := e.readSnippetPaths(wd)
	if err != nil {
		return err
	}

	mdBytes, err := e.gateway.GenMarkdown(wd, snippetPaths, imagePaths)
	if err != nil {
		return fmt.Errorf("adapter error: %w", err)
	}

	if _, err := os.Stdout.Write(mdBytes); err != nil {
		return fmt.Errorf("could not write to stdout: %w", err)
	}
	return nil
}

func (e *EngineImpl) readSnippetPaths(wd string) ([]string, error) {
	proverbsDir := filepath.Join(wd, proverbsRelDir)

	var snippetPaths []string
	if err := filepath.WalkDir(proverbsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".html" {
			snippetPaths = append(snippetPaths, path)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("walk error in %s: %w", wd, err)
	}
	return snippetPaths, nil
}

func (e *EngineImpl) readImagePaths(wd string) ([]string, error) {
	proverbImagesDir := filepath.Join(wd, proverbImagesRelDir)

	var imagePaths []string
	if err := filepath.WalkDir(proverbImagesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		imagePaths = append(imagePaths, path)

		return nil
	}); err != nil {
		return nil, fmt.Errorf("walk error in %s: %w", wd, err)
	}
	return imagePaths, nil
}
