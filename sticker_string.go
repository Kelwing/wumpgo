// Code generated by "stringer -type StickerType,StickerFormatType -output sticker_string.go"; DO NOT EDIT.

package objects

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[StickerTypeStandard-1]
	_ = x[StickerTypeGuild-2]
}

const _StickerType_name = "StickerTypeStandardStickerTypeGuild"

var _StickerType_index = [...]uint8{0, 19, 35}

func (i StickerType) String() string {
	i -= 1
	if i < 0 || i >= StickerType(len(_StickerType_index)-1) {
		return "StickerType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _StickerType_name[_StickerType_index[i]:_StickerType_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[StickerFormatTypePNG-1]
	_ = x[StickerFormatTypeAPNG-2]
	_ = x[StickerFormatTypeLOTTIE-3]
}

const _StickerFormatType_name = "StickerFormatTypePNGStickerFormatTypeAPNGStickerFormatTypeLOTTIE"

var _StickerFormatType_index = [...]uint8{0, 20, 41, 64}

func (i StickerFormatType) String() string {
	i -= 1
	if i < 0 || i >= StickerFormatType(len(_StickerFormatType_index)-1) {
		return "StickerFormatType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _StickerFormatType_name[_StickerFormatType_index[i]:_StickerFormatType_index[i+1]]
}
