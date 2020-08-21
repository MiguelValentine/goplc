package tag

type DataType uint16

const (
	NULL            DataType = 0x00
	BOOL            DataType = 0xc1
	SINT            DataType = 0xc2
	INT             DataType = 0xc3
	DINT            DataType = 0xc4
	LINT            DataType = 0xc5
	USINT           DataType = 0xc6
	UINT            DataType = 0xc7
	UDINT           DataType = 0xc8
	REAL            DataType = 0xca
	LREAL           DataType = 0xcb
	STIME           DataType = 0xcc
	DATE            DataType = 0xcd
	TIME_AND_DAY    DataType = 0xce
	DATE_AND_STRING DataType = 0xcf
	STRING          DataType = 0xd0
	WORD            DataType = 0xd1
	DWORD           DataType = 0xd2
	BIT_STRING      DataType = 0xd3
	LWORD           DataType = 0xd4
	STRING2         DataType = 0xd5
	FTIME           DataType = 0xd6
	LTIME           DataType = 0xd7
	ITIME           DataType = 0xd8
	STRINGN         DataType = 0xd9
	SHORT_STRING    DataType = 0xda
	TIME            DataType = 0xdb
	EPATH           DataType = 0xdc
	ENGUNIT         DataType = 0xdd
	STRINGI         DataType = 0xde
	STRUCT          DataType = 0xa002
)
