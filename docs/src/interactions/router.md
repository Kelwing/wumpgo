# Router

The router package contains a fairly opinionated and high level application command and component router that builds upon the `App` from the interactions package.  It allows you to declare your slash commands in Go using either comments (through `wumpgoctl gen`) or interface implementations.  The former generates the latter, allowing you to write and maintain less code.

It also allows you to route component interactions using an http-style router.  It allows you to specify your component handlers as routes (using an http-style path for the custom ID), and even supports path parameters.