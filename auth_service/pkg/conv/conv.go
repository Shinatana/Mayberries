package conv

import "unsafe"

func StrToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func BytesToStr(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
