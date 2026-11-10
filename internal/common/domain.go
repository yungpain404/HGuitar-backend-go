package common

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// Language mã ngôn ngữ nội dung.
type Language string

const (
	LanguageEN Language = "en"
	LanguageVI Language = "vi"
)

func (l Language) String() string { return string(l) }

// FromCode chuyển chuỗi mã sang Language, mặc định en nếu không hợp lệ.
func FromCode(code string) Language {
	switch strings.ToLower(code) {
	case "vi":
		return LanguageVI
	default:
		return LanguageEN
	}
}

// PublicationStatus trạng thái công khai.
type PublicationStatus string

const (
	PublicationStatusPublished  PublicationStatus = "PUBLISHED"
	PublicationStatusDraft      PublicationStatus = "DRAFT"
	PublicationStatusArchived   PublicationStatus = "ARCHIVED"
)

// GuitarType loại đàn.
type GuitarType string

const (
	GuitarTypeAcoustic GuitarType = "ACOUSTIC"
	GuitarTypeElectric GuitarType = "ELECTRIC"
	GuitarTypeClassical GuitarType = "CLASSICAL"
	GuitarTypeBass     GuitarType = "BASS"
)

// ChordCategory thể loại hợp âm.
type ChordCategory string

const (
	ChordCategoryMajor      ChordCategory = "MAJOR"
	ChordCategoryMinor      ChordCategory = "MINOR"
	ChordCategorySeventh    ChordCategory = "SEVENTH"
	ChordCategoryMinorSeventh ChordCategory = "MINOR_SEVENTH"
	ChordCategoryPower      ChordCategory = "POWER"
)

// ContentBlockType loại khối nội dung bài học.
type ContentBlockType string

const (
	ContentBlockTypeText      ContentBlockType = "TEXT"
	ContentBlockTypeImage     ContentBlockType = "IMAGE"
	ContentBlockTypeVideo     ContentBlockType = "VIDEO"
	ContentBlockTypeAudio     ContentBlockType = "AUDIO"
	ContentBlockTypeChordDiagram ContentBlockType = "CHORD_DIAGRAM"
	ContentBlockTypeStrumming ContentBlockType = "STRUMMING"
)

// Role vai trò người dùng.
type Role string

const (
	RoleStudent       Role = "STUDENT"
	RoleContentEditor Role = "CONTENT_EDITOR"
	RoleAdmin         Role = "ADMIN"
)

// ScanValue hỗ trợ scan enum từ DB.
func ScanEnum[T ~string](src any, target *T) error {
	switch v := src.(type) {
	case []byte:
		*target = T(string(v))
	case string:
		*target = T(v)
	default:
		return fmt.Errorf("cannot scan %T into enum", src)
	}
	return nil
}

// ValueEnum hỗ trợ ghi enum vào DB.
func ValueEnum[T ~string](v T) (driver.Value, error) {
	return string(v), nil
}

// RolesToStrings chuyển danh sách Role thành chuỗi.
func RolesToStrings(roles []Role) []string {
	out := make([]string, len(roles))
	for i, r := range roles {
		out[i] = string(r)
	}
	return out
}

// StringsToRoles chuyển danh sách chuỗi thành Role.
func StringsToRoles(strs []string) []Role {
	out := make([]Role, len(strs))
	for i, s := range strs {
		out[i] = Role(s)
	}
	return out
}
