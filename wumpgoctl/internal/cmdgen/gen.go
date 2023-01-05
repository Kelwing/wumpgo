package cmdgen

import (
	"context"
	"embed"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/objects/permissions"
)

//go:embed templates/*.gotmpl
var templates embed.FS

var (
	OptionNameRe                       = regexp.MustCompile(`@Name (?P<name>[-_\p{L}\p{N}]{1,32})`)
	OptionDescriptionRe                = regexp.MustCompile(`@Description (?P<description>[-_ \p{L}\p{N}]{1,100})`)
	OptionLocalizedNameRe              = regexp.MustCompile(`@Name.(?P<locale>[-A-Za-z]{2,5}) (?P<name>[-_\p{L}\p{N}]{1,32})`)
	OptionLocalizedDescriptionRe       = regexp.MustCompile(`@Description.(?P<locale>[-A-Za-z]{2,5}) (?P<description>[-_ \p{L}\p{N}]{1,100})`)
	OptionTypeRe                       = regexp.MustCompile(`@Type (?P<type>[A-Za-z]+)`)
	OptionOptionLocalizedNameRe        = regexp.MustCompile(`@Option.(?P<option>[-_\p{L}\p{N}]{1,32}).Name.(?P<locale>[-A-Za-z]{2,5}) (?P<name>[-_\p{L}\p{N}]{1,32})`)
	OptionOptionLocalizedDescriptionRe = regexp.MustCompile(`@Option.(?P<option>[-_\p{L}\p{N}]{1,32}).Description.(?P<locale>[-A-Za-z]{2,5}) (?P<description>[-_ \p{L}\p{N}]{1,100})`)
	OptionDMRe                         = regexp.MustCompile(`(?i)@DM (?P<value>true|false|yes|no|1|0)`)
	OptionPermissionsRe                = regexp.MustCompile(`@Permissions (?P<value>([A-Za-z]+,? ?)+|[0-9]+)`)
)

func Gen(ctx context.Context, dir string) {
	log.Ctx(ctx).Info().Str("dir", dir).Msg("parsing directory")

	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, dir, func(fi fs.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_gen.go")
	}, parser.ParseComments)
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("failed to parse dir")
	}

	for pkgName, pkg := range pkgs {
		processPackage(ctx, fset, pkgName, pkg, dir)
	}
}

type ctxKey string

func processPackage(ctx context.Context, fs *token.FileSet, name string, pkg *ast.Package, dir string) {
	log.Ctx(ctx).Info().Msgf("processing package %s", name)
	for fileName, f := range pkg.Files {
		ctx := log.Ctx(ctx).With().Str("package", name).Logger().WithContext(ctx)
		ctx = context.WithValue(ctx, ctxKey("package"), name)
		processFile(ctx, fs, fileName, f, dir)
	}
}

func processFile(ctx context.Context, fs *token.FileSet, name string, file *ast.File, dir string) {
	log.Ctx(ctx).Info().Msgf("processing file %s", filepath.Base(name))

	ctx = log.Ctx(ctx).With().Str("file", file.Name.Name).Logger().WithContext(ctx)
	p, err := doc.NewFromFiles(fs, []*ast.File{file}, dir)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to parse comments")
		return
	}

	baseName := strings.Split(filepath.Base(name), ".")[0]
	outputName := fmt.Sprintf("%s/%s_cmd_gen.go", dir, baseName)
	generate(ctx, outputName, p.Types)
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
	Imports []string
}

func generate(ctx context.Context, outputFileName string, types []*doc.Type) {
	args := tmplArg{
		Structs: make([]tmplStruct, 0),
		Package: ctx.Value(ctxKey("package")).(string),
		CmdLine: strings.Join(os.Args, " "),
	}

	for _, t := range types {
		log.Ctx(ctx).Info().Msgf("parsing comments for type %s", t.Name)
		s, ok := parseComments(ctx, &args, t)
		if ok {
			args.Structs = append(args.Structs, s)
		}
	}

	tmpl, err := template.ParseFS(templates, "templates/*")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load templates")
	}

	if len(args.Structs) > 0 {
		f, err := os.Create(outputFileName)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Str("filename", outputFileName).
				Msg("failed to create output file")
		}
		defer f.Close()
		if err := tmpl.Execute(f, args); err != nil {
			log.Error().Err(err).Msg("failed to write output")
		}
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

func parseComments(ctx context.Context, args *tmplArg, t *doc.Type) (tmplStruct, bool) {
	s := tmplStruct{
		Name:                           t.Name,
		NameLocalizations:              make(map[string]string),
		DescriptionLocalizations:       make(map[string]string),
		OptionNameLocalizations:        make(map[string]map[string]string),
		OptionDescriptionLocalizations: make(map[string]map[string]string),
	}

	hasOpts := false

	commentLines := strings.Split(t.Doc, "\n")
	for _, line := range commentLines {
		if strings.HasPrefix(line, "@") {
			hasOpts = true
			if matches := OptionNameRe.FindStringSubmatch(line); len(matches) > 0 {
				s.CommandName = matches[OptionNameRe.SubexpIndex("name")]
			} else if matches := OptionDescriptionRe.FindStringSubmatch(line); len(matches) > 0 {
				s.CommandDescription = matches[OptionDescriptionRe.SubexpIndex("description")]
			} else if matches := OptionLocalizedNameRe.FindStringSubmatch(line); len(matches) > 0 {
				s.NameLocalizations[matches[OptionLocalizedNameRe.SubexpIndex("locale")]] = matches[OptionLocalizedNameRe.SubexpIndex("name")]
			} else if matches := OptionLocalizedDescriptionRe.FindStringSubmatch(line); len(matches) > 0 {
				s.DescriptionLocalizations[matches[OptionLocalizedDescriptionRe.SubexpIndex("locale")]] = matches[OptionLocalizedDescriptionRe.SubexpIndex("name")]
			} else if matches := OptionTypeRe.FindStringSubmatch(line); len(matches) > 0 {
				args.Imports = append(args.Imports, "wumpgo.dev/wumpgo/objects")
				s.Type = parseType(matches[OptionTypeRe.SubexpIndex("type")])
			} else if matches := OptionOptionLocalizedNameRe.FindStringSubmatch(line); len(matches) > 0 {
				option := matches[OptionOptionLocalizedNameRe.SubexpIndex("option")]
				locale := matches[OptionOptionLocalizedNameRe.SubexpIndex("locale")]
				name := matches[OptionOptionLocalizedNameRe.SubexpIndex("name")]
				if _, ok := s.OptionNameLocalizations[option]; !ok {
					s.OptionNameLocalizations[option] = make(map[string]string)
				}
				s.OptionNameLocalizations[option][locale] = name
			} else if matches := OptionOptionLocalizedDescriptionRe.FindStringSubmatch(line); len(matches) > 0 {
				option := matches[OptionOptionLocalizedNameRe.SubexpIndex("option")]
				locale := matches[OptionOptionLocalizedNameRe.SubexpIndex("locale")]
				name := matches[OptionOptionLocalizedNameRe.SubexpIndex("name")]
				if _, ok := s.OptionDescriptionLocalizations[option]; !ok {
					s.OptionDescriptionLocalizations[option] = make(map[string]string)
				}
				s.OptionDescriptionLocalizations[option][locale] = name
			} else if matches := OptionDMRe.FindStringSubmatch(line); len(matches) > 0 {
				// ignore the error and default to false if the user can't
				// manage to correctly type any of the many, MANY valid options
				// for ParseBool
				dm, _ := strconv.ParseBool(matches[OptionDMRe.SubexpIndex("value")])
				s.DM = dm
			} else if matches := OptionPermissionsRe.FindStringSubmatch(line); len(matches) > 0 {
				args.Imports = append(args.Imports, "wumpgo.dev/wumpgo/objects/permissions")

				value := matches[OptionPermissionsRe.SubexpIndex("value")]

				if p, err := strconv.ParseUint(value, 10, 64); err == nil {
					// Raw integer permissions
					s.DefaultPerms = permissions.PermissionBit(p)
					log.Ctx(ctx).Debug().Uint64("perms", uint64(s.DefaultPerms)).Msg("raw permissions found")
				} else {
					perms := strings.Split(value, ",")
					for _, p := range perms {
						log.Ctx(ctx).Debug().Str("perm", p).Msg("string perms found")
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
	}

	return s, hasOpts
}
