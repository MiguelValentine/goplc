package segmentsTypes

type SegmentType uint16

const (
	PORT       SegmentType = 0 << 5
	LOGICAL    SegmentType = 1 << 5
	NETWORK    SegmentType = 2 << 5
	SYMBOLIC   SegmentType = 3 << 5
	DATA       SegmentType = 4 << 5
	DATATYPE_1 SegmentType = 5 << 5
	DATATYPE_2 SegmentType = 6 << 6
)
