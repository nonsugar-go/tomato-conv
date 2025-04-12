package main

import (
	"log"

	"github.com/nonsugar-go/tools/tui"
)

// DevTypeList is a selection list of DevTypes.
func (c ConfInfo) DevTypeList() DevType {
	choices := []string{
		"FortiGate", "PaloAlto",
	}
	selected, err := tui.Select(
		"機器の種類を選択してください", choices)
	if err != nil {
		log.Fatal(err)
	}
	switch selected {
	case choices[0]:
		return DevTypeFortiGate
	case choices[1]:
		return DevTypePaloAlto
	}
	return DevTypeUnknown
}
