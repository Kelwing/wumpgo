# Quick Start

wumpgo is modular, but that leads to it also having many moving parts.  Getting started quickly by hand can a complicated process.  To speed that up and reduce the time it takes to get started to near-zero we provide a tool called `wumpgoctl`.  This tool provides both scaffolding for getting started quickly, and a codegen tool for generating interface implementations for slash commands.  As such, installing it is highly recommended.

## Installing wumpgoctl

wumpgoctl is a quick go install away.

```bash
go install wumpgo.dev/wumpgo/wumpgoctl@latest
```

## Scaffolding

Scaffolding a new project with wumpgoctl is quick and easy.
We recommend taking a look at the help command to fully understand what you are doing `wumpgoctl init --help`

If you want to get started quickly with a common configuration, the below command should be enough

```bash
export BOTNAME=MyCoolBot
export BOTDIRECTORY=$(echo "${BOTNAME}" | tr '[:upper:]' '[:lower:]')
mkdir ${BOTDIRECTORY}
wumpgoctl init --name "${BOTNAME}" --gateway --local --codegen --root ${BOTDIRECTORY}
```

## Running your new bot

You can test the newly scaffolded bot with
```bash
go generate ./...
DISCORD_PUBLIC_KEY=xxx DISCORD_TOKEN=xxx go run . gateway
```