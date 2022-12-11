package scaffolding

import (
	"embed"
	"io/fs"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
)

const (
	templatesDir = "templates"
)

type LoadedTemplates map[string]*template.Template

var (
	//go:embed templates/*
	files     embed.FS
	templates LoadedTemplates
)

func init() {
	err := loadTemplates()
	if err != nil {
		log.Error().Err(err).Msg("failed to load scaffolding templates")
	}
}

func loadTemplates() error {
	if templates == nil {
		templates = make(LoadedTemplates)
	}

	return fs.WalkDir(files, templatesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		pt, err := template.ParseFS(files, path)
		if err != nil {
			return err
		}

		templateName := strings.ReplaceAll(
			strings.TrimPrefix(path, templatesDir+"/"),
			"tmpl",
			"",
		)

		templates[templateName] = pt

		return nil
	})
}

func Templates() LoadedTemplates {
	return templates
}
