package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/gobwas/glob"
	"github.com/rs/zerolog/log"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/objects/permissions"
)

//go:embed templates/*.gotmpl
var templates embed.FS

var (
	OptionName                       = glob.MustCompile("Name")
	OptionDescription                = glob.MustCompile("Description")
	OptionLocalizedName              = glob.MustCompile("Name.*")
	OptionLocalizedDescription       = glob.MustCompile("Description.*")
	OptionType                       = glob.MustCompile("Type")
	OptionOptionLocalizedName        = glob.MustCompile("Option.*.Name.*")
	OptionOptionLocalizedDescription = glob.MustCompile("Option.*.Description.*")
	OptionDM                         = glob.MustCompile("DM")
	OptionPermissions                = glob.MustCompile("Permissions")
)

var (
	dir = flag.String("dir", ".", "Directory to parse commands from")
)

type ctxKey string

func main() {
	flag.Parse()
	fset := token.NewFileSet()

	log.Info().Str("dir", *dir).Msg("parsing directory")

	pkgs, err := parser.ParseDir(fset, *dir, func(fi fs.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_gen.go")
	}, parser.ParseComments)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse dir")
	}

	ctx := log.Logger.WithContext(context.Background())

	for pkgName, pkg := range pkgs {
		processPackage(ctx, fset, pkgName, pkg)
	}
}

func processPackage(ctx context.Context, fs *token.FileSet, name string, pkg *ast.Package) {
	log.Ctx(ctx).Info().Msgf("processing package %s", name)
	for _, f := range pkg.Files {
		ctx := log.Ctx(ctx).With().Str("package", name).Logger().WithContext(ctx)
		ctx = context.WithValue(ctx, ctxKey("package"), name)
		processFile(ctx, fs, f)
	}
}

func processFile(ctx context.Context, fs *token.FileSet, file *ast.File) {
	log.Ctx(ctx).Info().Msgf("processing file %s", file.Name.Name)

	ctx = log.Ctx(ctx).With().Str("file", file.Name.Name).Logger().WithContext(ctx)
	p, err := doc.NewFromFiles(fs, []*ast.File{file}, *dir)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to parse comments")
		return
	}

	generate(ctx, fmt.Sprintf("%s/%s_cmd_gen.go", *dir, file.Name.Name), p.Types)
}

type tmplStruct struct {
	Name                           string
	CommandName                    string
	CommandDescription             string
	NameLocalizations              map[string]string
	DescriptionLocalizations       map[string]string
	Type                           *objects.ApplicationCommandType
	OptionNameLocalizations        map[string]map[string]string
	OptionDescriptionLocalizations map[string]map[string]string
	DM                             bool
	DefaultPerms                   permissions.PermissionBit
}

type tmplArg struct {
	Package string
	Structs []tmplStruct
	CmdLine string
}

func generate(ctx context.Context, outputFileName string, types []*doc.Type) {
	f, err := os.Create(outputFileName)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("filename", outputFileName).
			Msg("failed to create output file")
	}
	defer f.Close()

	args := tmplArg{
		Structs: make([]tmplStruct, 0),
		Package: ctx.Value(ctxKey("package")).(string),
		CmdLine: strings.Join(os.Args, " "),
	}

	for _, t := range types {
		log.Ctx(ctx).Info().Msgf("parsing comments for type %s", t.Name)
		args.Structs = append(args.Structs, parseComments(ctx, t))
	}

	tmpl, err := template.ParseFS(templates, "templates/*")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load templates")
	}

	if err := tmpl.Execute(f, args); err != nil {
		log.Error().Err(err).Msg("failed to write output")
	}
}

func parseType(arg string) *objects.ApplicationCommandType {
	m := map[string]objects.ApplicationCommandType{
		"ChatInput": objects.CommandTypeChatInput,
		"User":      objects.CommandTypeUser,
		"Message":   objects.CommandTypeMessage,
	}

	if t, ok := m[arg]; ok {
		return &t
	}

	return nil
}

func parseComments(ctx context.Context, t *doc.Type) tmplStruct {
	s := tmplStruct{
		Name:                           t.Name,
		NameLocalizations:              make(map[string]string),
		DescriptionLocalizations:       make(map[string]string),
		OptionNameLocalizations:        make(map[string]map[string]string),
		OptionDescriptionLocalizations: make(map[string]map[string]string),
	}

	commentLines := strings.Split(t.Doc, "\n")
	for _, line := range commentLines {
		if strings.HasPrefix(line, "@") {
			// This is a gen comment
			parts := strings.SplitN(strings.TrimPrefix(line, "@"), " ", 2)
			if len(parts) < 2 {
				log.Ctx(ctx).Warn().Str("line", line).Msg("malformed gen comment")
				continue
			}
			option := parts[0]
			arg := parts[1]

			if OptionName.Match(option) {
				s.CommandName = arg
			} else if OptionDescription.Match(option) {
				s.CommandDescription = arg
			} else if OptionLocalizedName.Match(option) {
				optionParts := strings.Split(option, ".")
				s.NameLocalizations[optionParts[1]] = arg
			} else if OptionLocalizedDescription.Match(option) {
				optionParts := strings.Split(option, ".")
				s.NameLocalizations[optionParts[1]] = arg
			} else if OptionType.Match(option) {
				s.Type = parseType(arg)
			} else if OptionOptionLocalizedName.Match(option) {
				optionParts := strings.Split(option, ".")
				if _, ok := s.OptionNameLocalizations[optionParts[1]]; !ok {
					s.OptionNameLocalizations[optionParts[1]] = make(map[string]string)
				}
				s.OptionNameLocalizations[optionParts[1]][optionParts[3]] = arg
			} else if OptionOptionLocalizedDescription.Match(option) {
				optionParts := strings.Split(option, ".")
				if _, ok := s.OptionDescriptionLocalizations[optionParts[1]]; !ok {
					s.OptionDescriptionLocalizations[optionParts[1]] = make(map[string]string)
				}
				s.OptionDescriptionLocalizations[optionParts[1]][optionParts[3]] = arg
			} else if OptionDM.Match(option) {
				// ignore the error and default to false if the user can't
				// manage to correctly type of the many, MANY valid options
				// for ParseBool
				dm, _ := strconv.ParseBool(arg)
				s.DM = dm
			} else if OptionPermissions.Match(option) {
				perms := strings.Split(arg, ",")
				for _, p := range perms {
					bit, err := permissions.PermissionBitString(strings.TrimSpace(p))
					if err != nil {
						log.Ctx(ctx).Error().Err(err).Msg("")
						continue
					}
					s.DefaultPerms |= bit
				}
			}
		}
	}

	return s
}
