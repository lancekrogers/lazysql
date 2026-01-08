package icons

const asciiUnknown = "?"

const (
	ASCIIDatabase   = "DB"
	ASCIISchema     = "SC"
	ASCIITable      = "TB"
	ASCIIView       = "VW"
	ASCIIColumn     = "CL"
	ASCIIIndex      = "IX"
	ASCIIPrimaryKey = "PK"
	ASCIIForeignKey = "FK"
	ASCIIFunction   = "FN"
	ASCIIProcedure  = "PR"
	ASCIITrigger    = "TR"
	ASCIISequence   = "SQ"
	ASCIIType       = "TY"
	ASCIIExtension  = "EX"
)

var asciiIcons = map[ObjectType]string{
	TypeDatabase:   ASCIIDatabase,
	TypeSchema:     ASCIISchema,
	TypeTable:      ASCIITable,
	TypeView:       ASCIIView,
	TypeColumn:     ASCIIColumn,
	TypeIndex:      ASCIIIndex,
	TypePrimaryKey: ASCIIPrimaryKey,
	TypeForeignKey: ASCIIForeignKey,
	TypeFunction:   ASCIIFunction,
	TypeProcedure:  ASCIIProcedure,
	TypeTrigger:    ASCIITrigger,
	TypeSequence:   ASCIISequence,
	TypeType:       ASCIIType,
	TypeExtension:  ASCIIExtension,
}
