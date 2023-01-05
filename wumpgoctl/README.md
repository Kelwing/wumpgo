# wumpgoctl

wumpgoctl is a command line application that provides a collection of tools to make bot creation with wumpgo easier.

## Installation

```bash
go install wumpgo.dev/wumpgo/wumpgoctl@latest
```

## Example usage

You can scaffold a new project
```bash
mkdir mybot && cd mybot
wumpgoctl init --pkg github.com/USERNAME/mybot --name MyBot --gateway --http
```

See `wumpgoctl init --help` for more information.

## Tools

### init

The init subcommand scaffolds new projects for you, making the time from 0 to bot extremely short.  You can get up and running with a very basic bot in under a minute.

###  gen

The gen subcommand provides a comment parser and code generation tool for implementing the router interfaces to customize your slash commands.

You can checkout the [gen](./docs/cmdgen.md) documentation for more information.