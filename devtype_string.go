// Code generated by "stringer -type DevType -trimprefix DevType"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DevTypeUnknown-0]
	_ = x[DevTypeFortiGate-1]
	_ = x[DevTypePaloAlto-2]
}

const _DevType_name = "UnknownFortiGatePaloAlto"

var _DevType_index = [...]uint8{0, 7, 16, 24}

func (i DevType) String() string {
	if i < 0 || i >= DevType(len(_DevType_index)-1) {
		return "DevType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _DevType_name[_DevType_index[i]:_DevType_index[i+1]]
}
