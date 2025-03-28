package analyzer

import (
	"MIA_2S_P1_202200041/comandos"
	"errors"
	"fmt"
	"strings"
)

func Analyzer(input string) (interface{}, error) {
	fmt.Println("ingresa a analyzer")
	indTokens := strings.Fields(input)

	if len(indTokens) == 0 {
		return "", errors.New("no command found")
	}

	switch indTokens[0] {
	case "mkdisk":
		fmt.Println(indTokens[1:])
		return comandos.ParserMkDisk(indTokens[1:])
	case "rmdisk":
		fmt.Println(indTokens[1:])
		return comandos.ParserRMDISK(indTokens[1:])
	case "fdisk":
		fmt.Println(indTokens[1:])
		return comandos.ParserFDISK(indTokens[1:])
	case "mount":
		fmt.Println(indTokens[1:])
		return comandos.ParserMOUNT(indTokens[1:])
	case "rep":
		fmt.Println(indTokens[1:])
		return comandos.ParserREP(indTokens[1:])
	default:
		return nil, fmt.Errorf("command not found: %s", indTokens[0])
	}
}
