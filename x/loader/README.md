# Loader

Loader is an experimental plugin loader for wumpgo.  Under the hood it uses the standard plugin package, but defines an interface for creating and loading plugins containing gateway handlers and slash commands, and registering them with a gateway receiver and command router respectively.