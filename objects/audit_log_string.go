// Code generated by "stringer -type=AuditLogEvent -trimprefix=AuditLogEvent -output audit_log_string.go"; DO NOT EDIT.

package objects

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AuditLogEventGuildUpdate-1]
	_ = x[AuditLogEventChannelCreate-10]
	_ = x[AuditLogEventChannelUpdate-11]
	_ = x[AuditLogEventChannelDelete-12]
	_ = x[AuditLogEventOverwriteCreate-13]
	_ = x[AuditLogEventOverwriteUpdate-14]
	_ = x[AuditLogEventOverwriteDelete-15]
	_ = x[AuditLogEventMemberKick-20]
	_ = x[AuditLogEventMemberPrune-21]
	_ = x[AuditLogEventMemberBanAdd-22]
	_ = x[AuditLogEventMemberBanRemove-23]
	_ = x[AuditLogEventMemberUpdate-24]
	_ = x[AuditLogEventMemberRoleUpdate-25]
	_ = x[AuditLogEventMemberMove-26]
	_ = x[AuditLogEventMemberDisconnect-27]
	_ = x[AuditLogEventBotAdd-28]
	_ = x[AuditLogEventRoleCreate-30]
	_ = x[AuditLogEventRoleUpdate-31]
	_ = x[AuditLogEventRoleDelete-32]
	_ = x[AuditLogEventInviteCreate-40]
	_ = x[AuditLogEventInviteUpdate-41]
	_ = x[AuditLogEventInviteDelete-42]
	_ = x[AuditLogEventWebhookCreate-50]
	_ = x[AuditLogEventWebhookUpdate-51]
	_ = x[AuditLogEventWebhookDelete-52]
	_ = x[AuditLogEventEmojiCreate-60]
	_ = x[AuditLogEventEmojiUpdate-61]
	_ = x[AuditLogEventEmojiDelete-62]
	_ = x[AuditLogEventMessageDelete-72]
	_ = x[AuditLogEventMessageBulkDelete-73]
	_ = x[AuditLogEventMessagePin-74]
	_ = x[AuditLogEventMessageUnpin-75]
	_ = x[AuditLogEventIntegrationCreate-80]
	_ = x[AuditLogEventIntegrationUpdate-81]
	_ = x[AuditLogEventIntegrationDelete-82]
	_ = x[AuditLogEventStageInstanceCreate-83]
	_ = x[AuditLogEventStageInstanceUpdate-84]
	_ = x[AuditLogEventStageInstanceDelete-85]
	_ = x[AuditLogEventStickerCreate-90]
	_ = x[AuditLogEventStickerUpdate-91]
	_ = x[AuditLogEventStickerDelete-92]
	_ = x[AuditLogEventGuildScheduledEventCreate-100]
	_ = x[AuditLogEventGuildScheduledEventUpdate-101]
	_ = x[AuditLogEventGuildScheduledEventDelete-102]
	_ = x[AuditLogEventThreadCreate-110]
	_ = x[AuditLogEventThreadUpdate-111]
	_ = x[AuditLogEventThreadDelete-112]
	_ = x[AuditLogApplicationCommandPermissionUpdate-121]
	_ = x[AuditLogAutoModerationRuleCreate-140]
	_ = x[AuditLogAutoModerationRuleUpdate-141]
	_ = x[AuditLogAutoModerationRuleDelete-142]
	_ = x[AuditLogAutoModerationBlockMessage-143]
	_ = x[AuditLogAutoModerationFlagToChannel-144]
	_ = x[AuditLogAutoModerationUserCommunicationDisabled-145]
}

const _AuditLogEvent_name = "GuildUpdateChannelCreateChannelUpdateChannelDeleteOverwriteCreateOverwriteUpdateOverwriteDeleteMemberKickMemberPruneMemberBanAddMemberBanRemoveMemberUpdateMemberRoleUpdateMemberMoveMemberDisconnectBotAddRoleCreateRoleUpdateRoleDeleteInviteCreateInviteUpdateInviteDeleteWebhookCreateWebhookUpdateWebhookDeleteEmojiCreateEmojiUpdateEmojiDeleteMessageDeleteMessageBulkDeleteMessagePinMessageUnpinIntegrationCreateIntegrationUpdateIntegrationDeleteStageInstanceCreateStageInstanceUpdateStageInstanceDeleteStickerCreateStickerUpdateStickerDeleteGuildScheduledEventCreateGuildScheduledEventUpdateGuildScheduledEventDeleteThreadCreateThreadUpdateThreadDeleteAuditLogApplicationCommandPermissionUpdateAuditLogAutoModerationRuleCreateAuditLogAutoModerationRuleUpdateAuditLogAutoModerationRuleDeleteAuditLogAutoModerationBlockMessageAuditLogAutoModerationFlagToChannelAuditLogAutoModerationUserCommunicationDisabled"

var _AuditLogEvent_map = map[AuditLogEvent]string{
	1:   _AuditLogEvent_name[0:11],
	10:  _AuditLogEvent_name[11:24],
	11:  _AuditLogEvent_name[24:37],
	12:  _AuditLogEvent_name[37:50],
	13:  _AuditLogEvent_name[50:65],
	14:  _AuditLogEvent_name[65:80],
	15:  _AuditLogEvent_name[80:95],
	20:  _AuditLogEvent_name[95:105],
	21:  _AuditLogEvent_name[105:116],
	22:  _AuditLogEvent_name[116:128],
	23:  _AuditLogEvent_name[128:143],
	24:  _AuditLogEvent_name[143:155],
	25:  _AuditLogEvent_name[155:171],
	26:  _AuditLogEvent_name[171:181],
	27:  _AuditLogEvent_name[181:197],
	28:  _AuditLogEvent_name[197:203],
	30:  _AuditLogEvent_name[203:213],
	31:  _AuditLogEvent_name[213:223],
	32:  _AuditLogEvent_name[223:233],
	40:  _AuditLogEvent_name[233:245],
	41:  _AuditLogEvent_name[245:257],
	42:  _AuditLogEvent_name[257:269],
	50:  _AuditLogEvent_name[269:282],
	51:  _AuditLogEvent_name[282:295],
	52:  _AuditLogEvent_name[295:308],
	60:  _AuditLogEvent_name[308:319],
	61:  _AuditLogEvent_name[319:330],
	62:  _AuditLogEvent_name[330:341],
	72:  _AuditLogEvent_name[341:354],
	73:  _AuditLogEvent_name[354:371],
	74:  _AuditLogEvent_name[371:381],
	75:  _AuditLogEvent_name[381:393],
	80:  _AuditLogEvent_name[393:410],
	81:  _AuditLogEvent_name[410:427],
	82:  _AuditLogEvent_name[427:444],
	83:  _AuditLogEvent_name[444:463],
	84:  _AuditLogEvent_name[463:482],
	85:  _AuditLogEvent_name[482:501],
	90:  _AuditLogEvent_name[501:514],
	91:  _AuditLogEvent_name[514:527],
	92:  _AuditLogEvent_name[527:540],
	100: _AuditLogEvent_name[540:565],
	101: _AuditLogEvent_name[565:590],
	102: _AuditLogEvent_name[590:615],
	110: _AuditLogEvent_name[615:627],
	111: _AuditLogEvent_name[627:639],
	112: _AuditLogEvent_name[639:651],
	121: _AuditLogEvent_name[651:693],
	140: _AuditLogEvent_name[693:725],
	141: _AuditLogEvent_name[725:757],
	142: _AuditLogEvent_name[757:789],
	143: _AuditLogEvent_name[789:823],
	144: _AuditLogEvent_name[823:858],
	145: _AuditLogEvent_name[858:905],
}

func (i AuditLogEvent) String() string {
	if str, ok := _AuditLogEvent_map[i]; ok {
		return str
	}
	return "AuditLogEvent(" + strconv.FormatInt(int64(i), 10) + ")"
}
