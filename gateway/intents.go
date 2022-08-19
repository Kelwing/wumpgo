package gateway

import (
	"fmt"
	"strings"
)

var intentMap = map[string]Intent{
	"guilds":                   IntentsGuilds,
	"guild_members":            IntentsGuildMembers,
	"guild_bans":               IntentsGuildBans,
	"guild_emojis":             IntentsGuildEmojis,
	"guild_integrations":       IntentsGuildIntegrations,
	"guild_webhooks":           IntentsGuildWebhooks,
	"guild_invites":            IntentsGuildInvites,
	"guild_voice_states":       IntentsGuildVoiceStates,
	"guild_presences":          IntentsGuildPresences,
	"guild_messages":           IntentsGuildMessages,
	"guild_message_reactions":  IntentsGuildMessageReactions,
	"guild_message_typing":     IntentsGuildMessageTyping,
	"direct_messages":          IntentsDirectMessages,
	"direct_message_reactions": IntentsDirectMessageReactions,
	"direct_message_typing":    IntentsDirectMessageTyping,
	"all_without_priviledged":  IntentsAllWithoutPrivileged,
	"all":                      IntentsAll,
	"none":                     IntentsNone,
}

func ParseIntents(s []string) (Intent, error) {
	var intents Intent
	for _, intent := range s {
		intent = strings.TrimSpace(intent)
		if intent == "" {
			continue
		}
		if i, ok := intentMap[intent]; !ok {
			return intents, fmt.Errorf("unknown intent: %s", intent)
		} else {
			intents |= i
		}
	}
	return intents, nil
}
