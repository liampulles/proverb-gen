package cli

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/liampulles/proverb-gen/internal/adapter"
)

const templateFileName = "_proverb_template.md"

// --- Engine front-matter ---

type Engine interface {
	Run(wd string, args []string) error
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

func (e *EngineImpl) Run(wd string, args []string) error {
	tmplPath := path.Join(wd, templateFileName)

	options, err := parseFlags(args)
	if err != nil {
		return err
	}

	mdBytes, err := e.gateway.InitMarkdown(tmplPath, options.Title, options.Tags)
	if err != nil {
		return fmt.Errorf("adapter error: %w", err)
	}

	if _, err := os.Stdout.Write(mdBytes); err != nil {
		return fmt.Errorf("could not write to stdout: %w", err)
	}
	return nil
}

type options struct {
	Title string
	Tags  []string
}

func parseFlags(args []string) (options, error) {
	// Define and run the flag set.
	flagSet := flag.NewFlagSet("proverb-gen", flag.ContinueOnError)
	titlePtr := flagSet.String("title", "Some Title", "Title for the proverb")
	tagsPtr := flagSet.String("tags", "general,code-design", "Tags associated with the proverb, comma separated.")
	if err := flagSet.Parse(args[1:]); err != nil {
		return options{}, err
	}

	tags := strings.Split(*tagsPtr, ",")

	return options{
		Title: *titlePtr,
		Tags:  tags,
	}, nil
}
