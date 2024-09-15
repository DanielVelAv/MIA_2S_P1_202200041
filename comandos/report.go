package comandos

import (
	"MIA_2S_P1_202200041/global"
	"MIA_2S_P1_202200041/reportes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type REPORT struct {
	name         string
	path         string
	id           string
	path_file_Is string
}

func ParserREP(tokens []string) (string, error) {
	com := &REPORT{} //instancia de report

	args := strings.Join(tokens, " ")

	re := regexp.MustCompile(`-name=[^\s]+|-path="[^"]+"|-path=[^\s]+|-id=[^\s]+|-path_file_is="[^"]+"|-path_file_is=[^\s]+`)
	coincidencias := re.FindAllString(args, -1)

	if len(coincidencias) != len(tokens) {
		for _, token := range tokens {
			if !re.MatchString(token) {
				return "", fmt.Errorf("formato de parametro inválido: %s", token)
			}
		}
	}

	for _, match := range coincidencias {
		clave := strings.SplitN(match, "=", 2)
		if len(clave) != 2 {
			return "", fmt.Errorf("formato de parametro inválido: %s", match)
		}
		key, value := strings.ToLower(clave[0]), clave[1]

		//quita comillas
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		switch key {
		case "-id":
			if value == "" {
				return "", errors.New("el id no puede venir vacio")
			}
			com.id = value
		case "-name":
			if value != "mbr" || value != "disk" || value != "inode" || value != "block" || value != "bm_inode" || value != "bm_block" || value != "sb" || value != "file" || value != "Ls" {
				return "", errors.New("el nombre no es valido, debe ser mbr, disk, inode, block, bm_inode, bm_block, sb, file o Ls")
			}
			com.name = value
		case "-path":
			if value == "" {
				return "", errors.New("el path no puede venir vacio")
			}
			com.path = value
		case "-path_file_Is":
			com.path_file_Is = value
		default:
			return "", fmt.Errorf("parametro desconocido: %s", key)
		}

	}

	if com.id == "" || com.name == "" || com.path == "" {
		return "", errors.New("faltan parametro : -id, -name, -path")
	}

	err := comandReportes(com)
	if err != nil {
		return "", err
	}

	return "Reporte: El reporte se genero correctamente", nil

}

func comandReportes(rep *REPORT) error {
	MBRmontado, sbMontado, pathDMontado, err := global.GetMPartitionReport(rep.id)
	if err != nil {
		return err
	}
	switch rep.name {
	case "mbr":
		err := reportes.ReporteMBR(MBRmontado, rep.path)
		if err != nil {
			fmt.Printf("Error al generar el reporte MBR: %v\n", err)
			return err
		}
	case "disk":
	}
	return nil
}
