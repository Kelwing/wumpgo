package scaffolding

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"regexp"
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
	funcMap   = template.FuncMap{
		"ToLower": strings.ToLower,
		"ToUpper": strings.ToUpper,
		"Bashify": Bashify,
	}
)

var (
	nameExtRe   = regexp.MustCompile(`(?P<Name>[A-Za-z]+)\.go(?P<Ext>[a-z]+)`)
	nameNoExtRe = regexp.MustCompile(`(?P<Name>[A-Za-z]+)\.gotemplate`)
)

func resolveFilename(s string) (string, error) {
	matches := nameNoExtRe.FindStringSubmatch(s)
	if matches != nil {
		return matches[nameNoExtRe.SubexpIndex("Name")], nil
	}

	matches = nameExtRe.FindStringSubmatch(s)
	if matches != nil {
		return matches[nameExtRe.SubexpIndex("Name")] + "." + matches[nameExtRe.SubexpIndex("Ext")], nil
	}

	return "", fmt.Errorf("invalid filename: %s", s)
}

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

	return fs.WalkDir(files, templatesDir, func(fp string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		templateName, err := resolveFilename(strings.TrimPrefix(fp, templatesDir+"/"))
		if err != nil {
			return err
		}

		pt, err := template.New(path.Base(fp)).Funcs(funcMap).ParseFS(files, fp)
		if err != nil {
			return err
		}

		rootPath, _ := path.Split(fp)
		relRoot := strings.SplitN(rootPath, "/", 2)[1]

		templates[path.Join(relRoot, templateName)] = pt

		return nil
	})
}

func Templates() LoadedTemplates {
	return templates
}
