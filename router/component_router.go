package router

import (
	"fmt"
	"sync"

	"github.com/gobwas/glob"
	"wumpgo.dev/wumpgo/objects"
)

// ComponentRouterContext is used to define the context for a component response.
type ComponentRouterContext struct {
}

// ComponentResponder is used to respond to a component interaction.
type ComponentResponder interface {
}

// ComponentHandler is used to define a handler for a component result.
type ComponentHandler interface {
	// Component is used ot handle routing to the component type.
	Component(ctx *ComponentRouterContext, c ComponentResponder) error
}

// ModalSubmitHandler is used to define a handler for a modal submit.
type ModalSubmitHandler interface {
	// ModalSubmit is used to handle routing to the modal submit type.
	ModalSubmit(ctx *ComponentRouterContext, c ComponentResponder) error
}

type globComponent struct {
	glob glob.Glob
	c    any
}

type routeManager struct {
	mu       sync.RWMutex
	routes   map[objects.InteractionType][]globComponent
	setGlobs map[objects.InteractionType]map[string]struct{}
}

// Gets the route. Returns nil if not found.
func (r *routeManager) getRoute(type_ objects.InteractionType, id string) any {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.routes == nil {
		return nil
	}

	if routes, ok := r.routes[type_]; ok {
		for _, route := range routes {
			if route.glob.Match(id) {
				return route.c
			}
		}
	}

	return nil
}

// Sets the route. Returns an error if something goes wrong or if the glob is not unique.
func (r *routeManager) setRoute(type_ objects.InteractionType, pattern string, c any) error {
	// Check if the component is nil.
	if c == nil {
		return fmt.Errorf("component is nil")
	}

	// Lock the routers mutex.
	r.mu.Lock()
	defer r.mu.Unlock()

	// Make sure the root maps we will need to add to/verify against exist.
	if r.routes == nil {
		r.routes = map[objects.InteractionType][]globComponent{}
	}
	if r.setGlobs == nil {
		r.setGlobs = map[objects.InteractionType]map[string]struct{}{}
	}

	// Compile the glob to make sure it is valid.
	g, err := glob.Compile(pattern)
	if err != nil {
		return err
	}

	// Check if the glob was used before.
	x, ok := r.setGlobs[type_]
	if ok {
		// There is a map for the type, check if the glob was used before.
		if _, ok = x[pattern]; ok {
			return fmt.Errorf("glob %q already used", pattern)
		}
	} else {
		// There is no map for the type, create one and we also know thie glob was not used before.
		x = map[string]struct{}{}
		r.setGlobs[type_] = x
	}

	// Add the glob to the map as used.
	x[pattern] = struct{}{}

	// Add the route to the routes map.
	r.routes[type_] = append(r.routes[type_], globComponent{glob: g, c: c})

	// Return no errors!
	return nil
}

// ComponentRouter is used to route Discord components to their respective handlers.
type ComponentRouter struct {
	routeManager
}

// Register is used to register a component handler. The component handler must also implement one or many of
// ComponentHandler and ModalSubmitHandler.
func (r *ComponentRouter) Register(pattern string, c any) error {
	implements := false

	if x, ok := c.(ComponentHandler); ok {
		if err := r.setRoute(objects.InteractionComponent, pattern, x); err != nil {
			return err
		}
		implements = true
	}

	if x, ok := c.(ModalSubmitHandler); ok {
		if err := r.setRoute(objects.InteractionModalSubmit, pattern, x); err != nil {
			return err
		}
		implements = true
	}

	if !implements {
		return fmt.Errorf("component does not implement any of ComponentHandler or ModalSubmitHandler")
	}
	return nil
}

// MustRegister is used to register a component handler. The component handler must also implement one or many of
// ComponentHandler and ModalSubmitHandler. This function panics if an error occurs.
func (r *ComponentRouter) MustRegister(pattern string, c any) {
	if err := r.Register(pattern, c); err != nil {
		panic(err)
	}
}

// Route is used to route a component interaction to its respective handler.
func (r *ComponentRouter) Route() {

}
