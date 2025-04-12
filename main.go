//go:generate go-winres make
package main

import (
	"flag"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"

	"github.com/nonsugar-go/tools/tui"
)

func main() {
	confInfo := ConfInfo{
		devType: DevTypeUnknown,
		outType: OutTypeExcel,
	}
	var devTypeStr string
	flag.StringVar(
		&devTypeStr, "dev", "", "機器の種類 {fgt|pa}",
	)
	flag.StringVar(
		&confInfo.confFilename, "conf", "", "設定ファイル",
	)
	flag.StringVar(
		&confInfo.outFilename, "out", "", "出力ファイル",
	)
	flag.Parse()
	switch {
	case strings.EqualFold(devTypeStr, "fortigate") ||
		strings.EqualFold(devTypeStr, "fgt") ||
		strings.EqualFold(devTypeStr, "fg"):
		confInfo.devType = DevTypeFortiGate
	case strings.EqualFold(devTypeStr, "paloalto") ||
		strings.EqualFold(devTypeStr, "pa"):
		confInfo.devType = DevTypePaloAlto
	}

	tui.Title("トマト🍅 の変換ツール")
	tui.MsgBox("ネットワーク機器の設定ファイルから設定表を作成するツールです")

	if confInfo.devType == DevTypeUnknown {
		confInfo.devType = confInfo.DevTypeList()
	}
	if confInfo.confFilename == "" {
		ext := map[DevType][]string{
			DevTypeFortiGate: {".conf"},
			DevTypePaloAlto:  {".xml", ".tgz"},
		}
		var err error
		for {
			if confInfo.confFilename, err = tui.FilePicker(
				"", ext[confInfo.devType]); err != nil {
				log.Fatal("設定ファイルが選択できていません")
			}
			break
		}
	}
	if confInfo.outFilename == "" {
		ext := filepath.Ext(confInfo.confFilename)
		confInfo.outFilename = strings.TrimRight(
			confInfo.confFilename, ext) + ".xlsx"
	}

	tui.PrintTable([]string{"項目", "選択した値"},
		[][]string{
			{"機器の種類", confInfo.devType.String()},
			{"設定ファイル", confInfo.confFilename},
			{"出力の種類", confInfo.outType.String()},
			{"出力ファイル", confInfo.outFilename},
		})
	tui.PressAnyKey()
}
