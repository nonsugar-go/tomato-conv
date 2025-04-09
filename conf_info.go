//go:generate stringer -type DevType -trimprefix DevType
//go:generate stringer -type OutType -trimprefix OutType
package main

import "fmt"

// ConfInfo is the configuration information.
type ConfInfo struct {
	devType      DevType // device type
	confFilename string  // config filename
	outType      OutType // output type
	outFilename  string  // output filename
}

// String implements fmt.Stringer interface
func (c ConfInfo) String() string {
	return fmt.Sprintf(
		"{devType: %s, confFilename: %q, outType: %s, outFilename: %q}",
		c.devType, c.confFilename, c.outType, c.outFilename,
	)
}

// DevType indicates the device type.
type DevType int

const (
	DevTypeUnknown DevType = iota
	DevTypeFortiGate
	DevTypePaloAlto
)

// OutType indicates the type of output.
type OutType int

const (
	OutTypeUnknown OutType = iota
	OutTypeExcel
	OutTypeJson
	OutTypeCsv
)
