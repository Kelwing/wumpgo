// Code generated by "stringer -type=ChannelType,ChannelFlag,ChannelForumSortOrder -trimprefix=Channel -output channel_string.go"; DO NOT EDIT.

package objects

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ChannelTypeGuildText-0]
	_ = x[ChannelTypeDM-1]
	_ = x[ChannelTypeGuildVoice-2]
	_ = x[ChannelTypeGroupDM-3]
	_ = x[ChannelTypeGuildCategory-4]
	_ = x[ChannelTypeGuildAnnouncement-5]
	_ = x[ChannelTypeAnnouncementThread-10]
	_ = x[ChannelTypePublicThread-11]
	_ = x[ChannelTypePrivateThread-12]
	_ = x[ChannelTypeGuildStageVoice-13]
	_ = x[ChannelTypeGuildDirectory-14]
	_ = x[ChannelTypeGuildForum-15]
}

const (
	_ChannelType_name_0 = "TypeGuildTextTypeDMTypeGuildVoiceTypeGroupDMTypeGuildCategoryTypeGuildAnnouncement"
	_ChannelType_name_1 = "TypeAnnouncementThreadTypePublicThreadTypePrivateThreadTypeGuildStageVoiceTypeGuildDirectoryTypeGuildForum"
)

var (
	_ChannelType_index_0 = [...]uint8{0, 13, 19, 33, 44, 61, 82}
	_ChannelType_index_1 = [...]uint8{0, 22, 38, 55, 74, 92, 106}
)

func (i ChannelType) String() string {
	switch {
	case i <= 5:
		return _ChannelType_name_0[_ChannelType_index_0[i]:_ChannelType_index_0[i+1]]
	case 10 <= i && i <= 15:
		i -= 10
		return _ChannelType_name_1[_ChannelType_index_1[i]:_ChannelType_index_1[i+1]]
	default:
		return "ChannelType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ChannelFlagPinned-2]
	_ = x[ChannelFlagRequireTag-16]
}

const (
	_ChannelFlag_name_0 = "FlagPinned"
	_ChannelFlag_name_1 = "FlagRequireTag"
)

func (i ChannelFlag) String() string {
	switch {
	case i == 2:
		return _ChannelFlag_name_0
	case i == 16:
		return _ChannelFlag_name_1
	default:
		return "ChannelFlag(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ChannelForumSortOrderLatestActivity-0]
	_ = x[ChannelForumSortOrderCreationDate-1]
}

const _ChannelForumSortOrder_name = "ForumSortOrderLatestActivityForumSortOrderCreationDate"

var _ChannelForumSortOrder_index = [...]uint8{0, 28, 54}

func (i ChannelForumSortOrder) String() string {
	if i >= ChannelForumSortOrder(len(_ChannelForumSortOrder_index)-1) {
		return "ChannelForumSortOrder(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ChannelForumSortOrder_name[_ChannelForumSortOrder_index[i]:_ChannelForumSortOrder_index[i+1]]
}
