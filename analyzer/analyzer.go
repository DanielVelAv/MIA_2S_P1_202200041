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
		return "", nil
	case "rmdisk":
		fmt.Println(indTokens[1:])
		return comandos.ParserRMDISK(indTokens[1:])
		return "", nil
	case "fdisk":
		fmt.Println(indTokens[1:])
		return comandos.ParserFDISK(indTokens[1:])
		return "", nil
	default:
		return nil, fmt.Errorf("command not found: %s", indTokens[0])
	}
}
