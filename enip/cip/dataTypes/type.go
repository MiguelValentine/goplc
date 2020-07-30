package dataTypes

type TypeString string
type Code uint8

// Vol1 Appendix C
const (
	// Logical Boolean with values TRUE and FALSE
	BOOL TypeString = "BOOL"
	// Signed 8–bit integer value
	SINT TypeString = "SINT"
	// Signed 16–bit integer value
	INT TypeString = "INT"
	// Signed 32–bit integer value
	DINT TypeString = "DINT"
	// Signed 64–bit integer value
	LINT TypeString = "LINT"
	// Unsigned 8–bit integer value
	USINT TypeString = "USINT"
	// Unsigned 16–bit integer value
	UINT TypeString = "UINT"
	// Unsigned 32–bit integer value
	UDINT TypeString = "UDINT"
	// Unsigned 64–bit integer value
	ULINT TypeString = "ULINT"
	// 32–bit floating point value
	REAL TypeString = "REAL"
	// 64–bit floating point value
	LREAL TypeString = "LREAL"
	// Synchronous time information
	STIME TypeString = "STIME"
	// Date information
	DATE TypeString = "DATE"
	// Time of day
	TIME_OF_DAY TypeString = "TIME_OF_DAY"
	// Date and time of day
	DATE_AND_TIME TypeString = "DATE_AND_TIME"
	// character string (1 byte per character)
	STRING TypeString = "STRING"
	// bit string - 8-bits
	BYTE TypeString = "BYTE"
	// bit string - 16-bits
	WORD TypeString = "WORD"
	// bit string - 32-bits
	DWORD TypeString = "DWORD"
	// bit string - 64-bits
	LWORD TypeString = "LWORD"
	// character string (2 bytes per character)
	STRING2 TypeString = "STRING2"
	// Duration (high resolution)
	FTIME TypeString = "FTIME"
	// Duration (long)
	LTIME TypeString = "LTIME"
	// Duration (short)
	ITIME TypeString = "ITIME"
	// character string (N bytes per character)
	STRINGN TypeString = "STRINGN"
	// character sting (1 byte per character, 1 byte length indicator)
	SHORT_STRING TypeString = "SHORT_STRING"
	// Duration (milliseconds)
	TIME TypeString = "TIME"
	// CIP path segments
	EPATH TypeString = "EPATH"
	// Engineering Units
	ENGUNIT TypeString = "ENGUNIT"
	// International Character String
	STRINGI TypeString = "STRINGI"
)

var TypeStringMap map[TypeString]Code
var TypeCodeMap map[Code]TypeString

func reg(code Code, typeString TypeString) {
	TypeStringMap[typeString] = code
	TypeCodeMap[code] = typeString
}

func init() {
	TypeStringMap = make(map[TypeString]Code)
	TypeCodeMap = make(map[Code]TypeString)

	reg(0xc1, BOOL)
	reg(0xc2, SINT)
	reg(0xc3, INT)
	reg(0xc4, DINT)
	reg(0xc5, LINT)
	reg(0xc6, USINT)
	reg(0xc7, UINT)
	reg(0xc8, UDINT)
	reg(0xc9, ULINT)
	reg(0xca, REAL)
	reg(0xcb, LREAL)
	reg(0xcc, STIME)
	reg(0xcd, DATE)
	reg(0xce, TIME_OF_DAY)
	reg(0xcf, DATE_AND_TIME)
	reg(0xd0, STRING)
	reg(0xd1, BYTE)
	reg(0xd2, WORD)
	reg(0xd3, DWORD)
	reg(0xd4, LWORD)
	reg(0xd5, STRING2)
	reg(0xd6, FTIME)
	reg(0xd7, LTIME)
	reg(0xd8, ITIME)
	reg(0xd9, STRINGN)
	reg(0xda, SHORT_STRING)
	reg(0xdb, TIME)
	reg(0xdc, EPATH)
	reg(0xdd, ENGUNIT)
	reg(0xde, STRINGI)
}

func IsValidTypeCode(code uint8) bool {
	_, ok := TypeCodeMap[Code(code)]
	return ok
}

func GetTypeCodeString(code uint8) string {
	v, _ := TypeCodeMap[Code(code)]
	return string(v)
}
