package router

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/Postcord/interactions"
	"github.com/Postcord/objects"
	"github.com/Postcord/rest"
)

// ModalGenerationCtx is used to generate a modal.
type ModalGenerationCtx struct {
	// Defines the interaction which started this.
	*objects.Interaction

	// Path is used to define the path of the modal.
	Path string `json:"path"`
}

// ModalRouterCtx is used to define the context for the modal event.
type ModalRouterCtx struct {
	// Defines the error handler.
	errorHandler ErrorHandler

	// Defines the global allowed mentions configuration.
	globalAllowedMentions *objects.AllowedMentions

	// Defines the void ID generator.
	voidGenerator

	// Defines the response builder.
	responseBuilder

	// Context is a context.Context passed from the HTTP handler.
	Context context.Context

	// Defines the interaction which started this.
	*objects.Interaction

	// Params is used to define any URL parameters.
	Params map[string]string `json:"params"`

	// ModalItems is used to define the modal items.
	ModalItems map[string]string `json:"modal_items"`

	// RESTClient is used to define the REST client.
	RESTClient rest.RESTClient `json:"rest_client"`
}

// ModalContentItem is used to define a item in a modal.
type ModalContentItem struct {
	// Short defines if the content of the modal is short text.
	Short bool `json:"short"`

	// Label is used to define the label of the item.
	Label string `json:"label"`

	// Key is used to define the key in the map.
	Key string `json:"key"`

	// Placeholder is used to define the placeholder for the input.
	// If this is blank, there will be no placeholder.
	Placeholder string `json:"placeholder"`

	// Value defines what the field should be pre-filled with.
	Value string `json:"value"`

	// Required defines if the content is required.
	Required bool `json:"required"`

	// MinLength defines the minimum length of the content.
	MinLength uint `json:"min_length"`

	// MaxLength defines the maximum length of the content.
	MaxLength uint `json:"max_length"`
}

// ModalContent defines the content of the modal.
type ModalContent struct {
	// Path is used to define the path to the modal.
	Path string `json:"path"`

	// Contents is used to define the contents of the modal.
	Contents func(*ModalGenerationCtx) (name string, contents []ModalContentItem) `json:"-"`

	// Function is the function that will be called when the modal is executed.
	Function func(*ModalRouterCtx) error `json:"-"`
}

// ModalRouter is used to route modals.
type ModalRouter struct {
	routes map[string]*ModalContent
	tree   *node
}

// ResponseDataBuilder is used to
type ResponseDataBuilder interface {
	ResponseData() *objects.InteractionApplicationCommandCallbackData
}

// Used to prepare the router.
func (f *ModalRouter) prep() {
	if f.routes == nil {
		f.routes = map[string]*ModalContent{}
	}
}

// AddModal is used to add a modal to the router.
func (f *ModalRouter) AddModal(modal *ModalContent) {
	f.prep()
	f.routes[modal.Path] = modal
}

// ModalPathNotFound is thrown when the modal path is not found.
var ModalPathNotFound = errors.New("modal path not found")

// Builds the router.
func (f *ModalRouter) build(loader loaderPassthrough) interactions.HandlerFunc {
	// Build the tree.
	tree := node{}
	for route, form := range f.routes {
		tree.addRoute(route, &routeContext{
			i: form,
			r: route,
		})
	}
	f.tree = &tree

	// Return the handler.
	return func(reqCtx context.Context, ctx *objects.Interaction) (resp *objects.InteractionResponse) {
		// Get the value from the tree.
		var data objects.ApplicationModalInteractionData
		if err := json.Unmarshal(ctx.Data, &data); err != nil {
			return loader.errHandler(err)
		}
		params := map[string]string{}
		val := f.tree.getValue(data.CustomID, params)
		if val == nil {
			return loader.errHandler(ModalPathNotFound)
		}

		// Create the rest tape if this is wanted.
		r := loader.rest
		tape := tape{}
		var returnedErr string
		errHandler := loader.errHandler
		if loader.generateFrames {
			r = &restTape{
				tape: &tape,
				rest: r,
			}
			errHandler = func(err error) *objects.InteractionResponse {
				returnedErr = err.Error()
				return loader.errHandler(err)
			}
		}

		// Handle if the function panics.
		defer func() {
			if errGeneric := recover(); errGeneric != nil {
				// Shouldn't try and return from defer.
				loader.errHandler(ungenericError(errGeneric))
			}
		}()

		// Handle test frames.
		defer func() {
			if loader.generateFrames {
				// Now we have all the data, we can generate the frame.
				fr := frame{ctx, tape, returnedErr, resp}
				go fr.write("testframes", "modals", strings.ReplaceAll(val.r, "/", "_"))
			}
		}()

		// Call the context.
		modalItems := map[string]string{}
		for _, row := range data.Components {
			for _, x := range row.Components {
				if x.CustomID != "" {
					modalItems[x.CustomID] = x.Value
				}
			}
			if row.CustomID != "" {
				modalItems[row.CustomID] = row.Value
			}
		}
		rctx := &ModalRouterCtx{
			errorHandler:          loader.errHandler,
			globalAllowedMentions: loader.globalAllowedMentions,
			Interaction:           ctx,
			Context:               reqCtx,
			Params:                params,
			ModalItems:            modalItems,
			RESTClient:            r,
		}
		if err := val.i.(*ModalContent).Function(rctx); err != nil {
			resp = errHandler(err)
			return
		}
		resp = rctx.buildResponse(false, loader.errHandler, loader.globalAllowedMentions)
		return
	}
}

// MultipleModalResponses is thrown when a modal response is triggered within another.
var MultipleModalResponses = errors.New("multiple modal responses")

// UnknownContextType is thrown when a context type is not from Postcord.
var UnknownContextType = errors.New("unknown context type")

func uint2IntPtr(x uint) *int {
	if x == 0 {
		return nil
	}
	y := int(x)
	return &y
}

// SendModalResponse is used to send the modal response with the given context.
// The passed through context is expected to be one of a valid Postcord type.
// The router will need to be built before you can use this function.
func (f *ModalRouter) SendModalResponse(ctx ResponseDataBuilder, path string) error {
	// Get the value from the tree.
	m := map[string]string{}
	val := f.tree.getValue(path, m)
	if val == nil {
		return ModalPathNotFound
	}

	// Mark the response type as a modal.
	var interaction *objects.Interaction
	switch x := ctx.(type) {
	case *ModalRouterCtx:
		return MultipleModalResponses
	case *CommandRouterCtx:
		x.respType = objects.ResponseModal
		interaction = x.Interaction
	case *ComponentRouterCtx:
		x.respType = objects.ResponseModal
		interaction = x.Interaction
	default:
		return UnknownContextType
	}

	// Cast the form content from the data.
	formContent := val.i.(*ModalContent)

	// Build the response.
	data := ctx.ResponseData()
	data.CustomID = path
	formName, formContents := formContent.Contents(&ModalGenerationCtx{
		Interaction: interaction,
		Path:        path,
	})
	data.Title = formName
	components := make([]*objects.Component, len(formContents))
	for i, modalContent := range formContents {
		style := objects.TextStyleParagraph
		if modalContent.Short {
			style = objects.TextStyleShort
		}
		components[i] = &objects.Component{
			Type: objects.ComponentTypeActionRow,
			Components: []*objects.Component{
				{
					Type:        objects.ComponentTypeInputText,
					Label:       modalContent.Label,
					Style:       objects.ButtonStyle(style),
					CustomID:    modalContent.Key,
					Placeholder: modalContent.Placeholder,
					Value:       modalContent.Value,
					Required:    modalContent.Required,
					MinLength:   uint2IntPtr(modalContent.MinLength),
					MaxLength:   uint2IntPtr(modalContent.MaxLength),
				},
			},
		}
	}
	data.Components = components

	// Return no errors.
	return nil
}
