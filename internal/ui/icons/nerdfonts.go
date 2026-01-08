package icons

import (
	"os"
	"strconv"
	"strings"
)

type ObjectType int

const (
	TypeDatabase ObjectType = iota
	TypeSchema
	TypeTable
	TypeView
	TypeColumn
	TypeIndex
	TypePrimaryKey
	TypeForeignKey
	TypeFunction
	TypeProcedure
	TypeTrigger
	TypeSequence
	TypeType
	TypeExtension
)

const (
	IconDatabase   = "\uf1c0"
	IconSchema     = "\uf07b"
	IconTable      = "\uf0ce"
	IconView       = "\uf06e"
	IconColumn     = "\uf0db"
	IconIndex      = "\uf0e8"
	IconPrimaryKey = "\uf084"
	IconForeignKey = "\uf0c1"
	IconFunction   = "\uf0e7"
	IconProcedure  = "\uf085"
	IconTrigger    = "\uf0e4"
	IconSequence   = "\uf0c9"
	IconType       = "\uf0db"
	IconExtension  = "\uf12e"
)

var nerdFontIcons = map[ObjectType]string{
	TypeDatabase:   IconDatabase,
	TypeSchema:     IconSchema,
	TypeTable:      IconTable,
	TypeView:       IconView,
	TypeColumn:     IconColumn,
	TypeIndex:      IconIndex,
	TypePrimaryKey: IconPrimaryKey,
	TypeForeignKey: IconForeignKey,
	TypeFunction:   IconFunction,
	TypeProcedure:  IconProcedure,
	TypeTrigger:    IconTrigger,
	TypeSequence:   IconSequence,
	TypeType:       IconType,
	TypeExtension:  IconExtension,
}

var useNerdFonts = true

func init() {
	applyEnvConfig()
}

func applyEnvConfig() {
	value := strings.TrimSpace(os.Getenv("LAZYSQL_NERD_FONTS"))
	if value == "" {
		return
	}
	if parsed, err := strconv.ParseBool(value); err == nil {
		useNerdFonts = parsed
	}
}

func SetNerdFonts(enabled bool) {
	useNerdFonts = enabled
}

func Icon(objType ObjectType) string {
	if useNerdFonts {
		if icon, ok := nerdFontIcons[objType]; ok {
			return icon
		}
	}
	if icon, ok := asciiIcons[objType]; ok {
		return icon
	}
	return asciiUnknown
}
