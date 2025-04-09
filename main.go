package main

import (
	"flag"
	"fmt"
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
		&confInfo.confFilename, "conf", "", "コンフィグ ファイル",
	)
	flag.StringVar(
		&confInfo.outFilename, "out", "", "出力ファイル名",
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

	/*
	   	var style = lipgloss.NewStyle().
	   		BorderStyle(lipgloss.RoundedBorder()).
	   		BorderForeground(lipgloss.Color("228")).
	   		BorderBackground(lipgloss.Color("63")).
	   		BorderTop(true).
	   		BorderLeft(true).
	   		BorderRight(true).
	   		BorderBottom(true).
	   		PaddingLeft(1).
	   		PaddingRight(1).
	   		Width(40)
	   	fmt.Println(style.Render(`🍅 の変換ツール
	   ネットワーク機器の設定ファイルから設定表を作成するツールです`))
	*/

	if confInfo.devType == DevTypeUnknown {
		//confInfo.devType = confInfo.DevTypeList()
	}

	log.Info(fmt.Sprintf("confInfo: %s", confInfo))
}
