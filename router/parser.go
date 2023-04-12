package router

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"wumpgo.dev/snowflake"
	"wumpgo.dev/wumpgo/objects"
	"wumpgo.dev/wumpgo/objects/permissions"
)

type CommandParser struct {
	currentPath *stack[string]
	handlers    map[string]CommandHandler
}

func (p *CommandParser) Handlers() map[string]CommandHandler {
	return p.handlers
}

func NewParser() *CommandParser {
	return &CommandParser{
		currentPath: ptr(stack[string](make([]string, 0))),
		handlers:    make(map[string]CommandHandler),
	}
}

type GuildCommand interface {
	GuildID() objects.Snowflake
}

type CommandNamer interface {
	CommandName() string
}

type CommandDescriptioner interface {
	Description() string
}

type NameLocalizer interface {
	NameLocalizations() map[string]string
}

type DescriptionLocalizer interface {
	DescriptionLocalizations() map[string]string
}

type DefaultPermissioner interface {
	DefaultPermissions() permissions.PermissionBit
}

type DMPermissioner interface {
	AllowInDMs() bool
}

type AutoCompleter interface {
	AutoComplete(optionName string, value interface{}) []*objects.ApplicationCommandOptionChoice
}

type CommandTyper interface {
	Type() objects.ApplicationCommandType
}

type OptionNameLocalizer interface {
	OptionName(optionName string) map[string]string
}

type OptionDescriptionLocalizer interface {
	OptionDescription(optionName string) map[string]string
}

type ParserError struct {
	Message string
}

func (p *ParserError) Error() string {
	return "command parser: " + p.Message
}

func newParserErrorf(format string, args ...any) *ParserError {
	return &ParserError{Message: fmt.Sprintf(format, args...)}
}

func newFieldErrorf(fieldName string, format string, args ...any) *ParserError {
	args = append([]any{fieldName}, args...)
	return &ParserError{Message: fmt.Sprintf("%s: "+format, args...)}
}

func ptr[T any](in T) *T {
	return &in
}

var (
	userType    = reflect.TypeOf(objects.User{})
	channelType = reflect.TypeOf(objects.Channel{})
	roleType    = reflect.TypeOf(objects.Role{})
	attachType  = reflect.TypeOf(objects.DiscordFile{})
)

func (p *CommandParser) parseFields(v reflect.Value, depth int) ([]objects.ApplicationCommandOption, error) {
	var cmdFlags uint8 = 0

	options := make([]objects.ApplicationCommandOption, 0)
	optionals := make([]objects.ApplicationCommandOption, 0)

	for i := 0; i < v.NumField(); i++ {
		// Ignore unexported fields
		if !v.Type().Field(i).IsExported() {
			continue
		}

		o, err := p.parseOption(v, i, depth)
		if err != nil {
			return nil, err
		}

		if o.OptionType == objects.TypeSubCommand || o.OptionType == objects.TypeSubCommandGroup {
			cmdFlags |= (1 << 0)
		} else {
			cmdFlags |= (1 << 1)
		}

		if cmdFlags == 0x3 {
			return nil, newParserErrorf("cannot have both sub commands an root command options")
		}

		if o.Required {
			options = append(options, *o)
		} else {
			optionals = append(optionals, *o)
		}
	}

	options = append(options, optionals...)

	return options, nil
}

func (p *CommandParser) parseCommand(v reflect.Value) (*objects.ApplicationCommand, error) {
	command := &objects.ApplicationCommand{
		AllowUseInDMs: ptr(false),
	}

	t := v.Type()

	if t.Kind() != reflect.Ptr {
		return nil, newParserErrorf("must be a pointer to a command struct")
	}

	v = v.Elem()
	t = v.Type()

	if t.Kind() != reflect.Struct {
		return nil, newParserErrorf("command must be a struct")
	}

	if gc, ok := v.Interface().(GuildCommand); ok {
		command.GuildID = ptr(gc.GuildID())
	}

	if n, ok := v.Interface().(CommandNamer); ok {
		command.Name = n.CommandName()
	} else {
		command.Name = strings.ToLower(t.Name())
	}

	if d, ok := v.Interface().(CommandDescriptioner); ok {
		command.Description = d.Description()
	} else {
		command.Description = command.Name
	}

	if nl, ok := v.Interface().(NameLocalizer); ok {
		command.NameLocalizations = nl.NameLocalizations()
	}

	if dl, ok := v.Interface().(DescriptionLocalizer); ok {
		command.DescriptionLocalizations = dl.DescriptionLocalizations()
	}

	if dp, ok := v.Interface().(DefaultPermissioner); ok {
		command.DefaultPermissions = ptr(dp.DefaultPermissions())
	}

	if dmp, ok := v.Interface().(DMPermissioner); ok {
		command.AllowUseInDMs = ptr(dmp.AllowInDMs())
	}

	if t, ok := v.Interface().(CommandTyper); ok {
		if t.Type() != objects.CommandTypeChatInput && len(command.Options) > 0 {
			return nil, newParserErrorf(
				"options are not allowed for command type %s", t.Type().String(),
			)
		}
		command.Description = ""
		command.Type = ptr(t.Type())
	}

	p.currentPath.Push(command.Name)
	defer p.currentPath.Pop()

	options, err := p.parseFields(v, 0)
	if err != nil {
		return nil, err
	}

	if h, ok := v.Addr().Interface().(CommandHandler); ok {
		p.handlers[command.Name] = h
	}

	command.Options = options

	return command, nil
}

func (p *CommandParser) parseSubCommand(v reflect.Value, depth int) (*objects.ApplicationCommandOption, error) {
	if depth == 2 {
		return nil, newParserErrorf("command nested too deep")
	}
	depth++
	option := &objects.ApplicationCommandOption{
		OptionType: objects.TypeSubCommand,
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, newParserErrorf("sub-command must be a struct")
	}

	if n, ok := v.Interface().(CommandNamer); ok {
		option.Name = n.CommandName()
	} else {
		option.Name = strings.ToLower(v.Type().Name())
	}

	if d, ok := v.Interface().(CommandDescriptioner); ok {
		option.Description = d.Description()
	} else {
		option.Description = option.Name
	}

	if nl, ok := v.Interface().(NameLocalizer); ok {
		option.NameLocalizations = nl.NameLocalizations()
	}

	if dl, ok := v.Interface().(DescriptionLocalizer); ok {
		option.DescriptionLocalizations = dl.DescriptionLocalizations()
	}

	p.currentPath.Push(option.Name)
	defer p.currentPath.Pop()

	options, err := p.parseFields(v, depth)
	if err != nil {
		return nil, err
	}

	if len(options) > 0 && options[0].OptionType == objects.TypeSubCommand {
		option.OptionType = objects.TypeSubCommandGroup
	}

	if h, ok := v.Addr().Interface().(CommandHandler); ok {
		commandPath := strings.Join([]string(*p.currentPath), "/")
		p.handlers[commandPath] = h
	}

	option.Options = options

	return option, nil
}

func (p *CommandParser) parseOption(v reflect.Value, i, depth int) (*objects.ApplicationCommandOption, error) {
	t := v.Type().Field(i)
	if t.Anonymous {
		return p.parseSubCommand(v.Field(i), depth)
	}

	option := &objects.ApplicationCommandOption{
		Required: true,
	}

	kind := t.Type.Kind()
	fieldType := t.Type

	if kind == reflect.Ptr {
		kind = fieldType.Elem().Kind()
		fieldType = fieldType.Elem()
	}

	switch kind {
	case reflect.String:
		option.OptionType = objects.TypeString
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		option.OptionType = objects.TypeInteger
	case reflect.Float64:
		option.OptionType = objects.TypeNumber
	case reflect.Bool:
		option.OptionType = objects.TypeBoolean
	case reflect.Uint64: // Should be a snowflake
		option.OptionType = objects.TypeMentionable
	case reflect.Struct:
		if fieldType.AssignableTo(userType) {
			option.OptionType = objects.TypeUser
		} else if fieldType.AssignableTo(channelType) {
			option.OptionType = objects.TypeChannel
		} else if fieldType.AssignableTo(roleType) {
			option.OptionType = objects.TypeRole
		} else if fieldType.AssignableTo(attachType) {
			option.OptionType = objects.TypeAttachment
		} else {
			return nil, newParserErrorf(
				"unknown struct type %s for field %s", fieldType.Name(), t.Name)
		}
	default:
		return nil, newParserErrorf(
			"unknown type %s for field %s", fieldType.Name(), t.Name,
		)
	}

	tagData := t.Tag.Get("discord")
	tagParts := strings.Split(tagData, ",")

	if len(tagParts) == 0 || tagParts[0] == "" {
		option.Name = strings.ToLower(t.Name)
	} else {
		option.Name = tagParts[0]
	}

	option.Description = option.Name

	for _, p := range tagParts[1:] {
		args := strings.Split(p, ":")
		switch args[0] {
		case "optional":
			option.Required = false
		case "autocomplete":
			option.Autocomplete = true
		case "description":
			if len(args) != 2 {
				return nil, newFieldErrorf(t.Name, "no argument for description tag")
			}

			option.Description = args[1]
		case "minLength":
			if option.OptionType != objects.TypeString {
				return nil, newFieldErrorf(t.Name, "minLength only supported for strings")
			}
			if len(args) != 2 {
				return nil, newFieldErrorf(t.Name, "no argument for minLength tag")
			}
			min, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return nil, newFieldErrorf(t.Name, "argument for minLength tag is not int")
			}
			option.MinLength = min
		case "maxLength":
			if option.OptionType != objects.TypeString {
				return nil, newFieldErrorf(t.Name, "maxLength only supported for strings")
			}
			if len(args) != 2 {
				return nil, newFieldErrorf(t.Name, "no argument for maxLength tag")
			}
			max, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return nil, newFieldErrorf(t.Name, "argument for maxLength tag is not int")
			}
			option.MaxLength = max
		case "minValue":
			if option.OptionType != objects.TypeInteger && option.OptionType != objects.TypeNumber {
				return nil, newFieldErrorf(t.Name, "minValue only supported for ints or floats")
			}
			if len(args) != 2 {
				return nil, newFieldErrorf(t.Name, "no argument for minValue tag")
			}
			option.MinValue = json.Number(args[1])
		case "maxValue":
			if option.OptionType != objects.TypeInteger && option.OptionType != objects.TypeNumber {
				return nil, newFieldErrorf(t.Name, "maxValue only supported for ints or floats")
			}
			if len(args) != 2 {
				return nil, newFieldErrorf(t.Name, "no argument for maxValue tag")
			}
			option.MaxValue = json.Number(args[1])
		case "channelTypes":
			if option.OptionType != objects.TypeChannel {
				return nil, newFieldErrorf(t.Name, "channelTypes only supported for channels")
			}
			if len(args) != 2 {
				return nil, newFieldErrorf(t.Name, "missing types for channelTypes tag")
			}
			types := strings.Split(args[1], ";")
			option.ChannelTypes = make([]objects.ChannelType, len(types))
			for i, ct := range types {
				typeID, err := strconv.ParseInt(ct, 10, 64)
				if err != nil {
					return nil, newFieldErrorf(t.Name, "invalid channel type format: %s", ct)
				}
				realct := objects.ChannelType(typeID)
				if strings.HasPrefix(realct.String(), "ChannelType(") {
					return nil, newFieldErrorf(t.Name, "invalid channel type: %d", realct)
				}

				option.ChannelTypes[i] = realct
			}
		}
	}

	choicesTagData, ok := t.Tag.Lookup("choices")
	if ok {
		choicesTagParts := strings.Split(choicesTagData, ",")
		option.Choices = make([]objects.ApplicationCommandOptionChoice, len(choicesTagParts))
		for i, p := range choicesTagParts {
			keyValue := strings.Split(p, ":")
			if len(keyValue) != 2 {
				return nil, newFieldErrorf(t.Name, "choices tag must be a comma separate list of name:value pairs")
			}
			name := keyValue[0]
			value := keyValue[1]

			option.Choices[i] = objects.ApplicationCommandOptionChoice{
				Name:  name,
				Value: value,
			}
		}
	}

	if onl, ok := v.Interface().(OptionNameLocalizer); ok {
		option.NameLocalizations = onl.OptionName(t.Name)
	}

	if odl, ok := v.Interface().(OptionDescriptionLocalizer); ok {
		option.DescriptionLocalizations = odl.OptionDescription(t.Name)
	}

	return option, nil
}

func unmarshalOptions(dst any, choices []*objects.ApplicationCommandDataOption, resolvable *objects.ResolvedData, logger ...zerolog.Logger) error {
	log := zerolog.Nop()
	if len(logger) > 0 {
		log = logger[0]
		log.Debug().Msg("logging unmarshalling")
	}
	val := reflect.ValueOf(dst)

	if !val.IsValid() {
		log.Warn().Msg("destination is not valid")
		return newParserErrorf("dst is not valid")
	}

	if val.Type().Kind() == reflect.Interface {
		val = val.Elem()
	}

	if val.Type().Kind() != reflect.Ptr {
		log.Warn().Msg("destination is not a pointer")
		return newParserErrorf("dst must be a pointer")
	}

	if val.Type().Elem().Kind() != reflect.Struct {
		log.Warn().Msg("destination must be a pointer to a command definition struct")
		return newParserErrorf("dst must be a pointer to a command definition struct")
	}

	choiceMap := make(map[string]*objects.ApplicationCommandDataOption)
	for _, c := range choices {
		choiceMap[c.Name] = c
	}

	for i := 0; i < val.Elem().Type().NumField(); i++ {
		f := val.Elem().Type().Field(i)
		fv := val.Elem().Field(i)
		tag := f.Tag

		l := log.With().Str("field_name", fv.Type().Name()).Logger()

		l.Debug().Msg("processing field")

		tagVal, ok := tag.Lookup("discord")
		v := strings.Split(tagVal, ",")[0]
		if !ok || v == "" {
			// infer option name from field name
			v = strings.ToLower(f.Name)
		}

		c, ok := choiceMap[v]
		if !ok {
			l.Debug().Msg("tag not found in users choices")
			continue
		}

		vv := reflect.ValueOf(c.Value)
		if fv.Type().Kind() == reflect.Ptr {
			l.Debug().Msg("field value is a pointer")
			fv = fv.Elem()
		}

		if fv.Type().Kind() == reflect.Struct {
			l.Debug().Msg("field is a struct")
			if vv.Type().Kind() == reflect.String {
				// Value should be a Snowflake
				sn, err := snowflake.SnowflakeFromString(vv.String())
				l.Debug().Str("snowflake", vv.String()).Msg("resolving snowflake")
				if err == nil {
					if fv.Type().AssignableTo(userType) {
						log.Debug().Msg("is assignable to user")
						u, ok := resolvable.Users[sn]
						if ok {
							if fv.Kind() == reflect.Ptr {
								fv.Set(reflect.ValueOf(&u))
							} else {
								fv.Set(reflect.ValueOf(u))
							}
						}
					} else if fv.Type().AssignableTo(channelType) {
						c, ok := resolvable.Channels[sn]
						if ok {
							fv.Set(reflect.ValueOf(c))
						}
					} else if fv.Type().AssignableTo(roleType) {
						r, ok := resolvable.Roles[sn]
						if ok {
							fv.Set(reflect.ValueOf(r))
						}
					} else if fv.Type().AssignableTo(attachType) {
						a, ok := resolvable.Attachments[sn]
						if ok {
							fv.Set(reflect.ValueOf(a))
						}
					}
				}
			} else {
				l.Debug().Str("field_type", vv.Type().Kind().String()).Msg("field is not a string")
			}
		}

		if fv.CanSet() && fv.Type().AssignableTo(vv.Type()) {
			fv.Set(vv)
		} else if vv.CanConvert(fv.Type()) {
			nv := vv.Convert(fv.Type())
			fv.Set(nv)
		} else {
			l.Warn().Str("value_type", vv.Type().Name()).Msg("value not assignable or convertable to field")
		}
	}

	return nil
}
