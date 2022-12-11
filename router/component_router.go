package router

import (
	"fmt"
	"sync"

	"github.com/gobwas/glob"
	"wumpgo.dev/wumpgo/objects"
)

// Component is used to define a Discord component which can be inserted into the router.
type Component interface {
	// Glob is used to define the glob pattern for the component. This should be unique.
	Glob() string
}

type globComponent struct {
	glob glob.Glob
	c    Component
}

// ComponentRouter is used to route Discord components to their respective handlers.
type ComponentRouter struct {
	mu       sync.RWMutex
	routes   map[objects.InteractionType][]globComponent
	setGlobs map[objects.InteractionType]map[string]struct{}
}

// Gets the route. Returns nil if not found.
func (r *ComponentRouter) getRoute(type_ objects.InteractionType, id string) Component {
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
func (r *ComponentRouter) setRoute(type_ objects.InteractionType, pattern string, c Component) error {
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
