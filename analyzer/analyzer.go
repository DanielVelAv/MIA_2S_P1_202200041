package analyzer

import (
	comandos "MIA_2S_P1_202200041/comandos"
	"errors"
	"fmt"
	"strings"
)

func Analyzer(input string) (interface{}, error) {
	fmt.Println("ingresa a analyzer")
	indTokens := strings.Fields(input)

	if len(indTokens) == 0 {
		return nil, errors.New("no command found")
	}
	fmt.Println(indTokens, input)
	fmt.Println(indTokens[0])
	switch indTokens[0] {
	case "mkdisk":
		fmt.Println(indTokens[1:])
		return comandos.ParserMkDisk(indTokens[1:])

	default:
		return nil, fmt.Errorf("command not found: %s", indTokens[0])
	}
}
