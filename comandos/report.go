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
	fmt.Println("Ingresa a parse report")
	com := &REPORT{} //instancia de report
	args := strings.Join(tokens, " ")
	fmt.Println("ER")
	re := regexp.MustCompile(`-name=[^\s]+|-path="[^"]+"|-path=[^\s]+|-id=[^\s]+|-path_file_is="[^"]+"|-path_file_is=[^\s]+`)
	coincidencias := re.FindAllString(args, -1)

	fmt.Println("coincidencias")
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
			nombresValidos := []string{"mbr", "disk", "inode", "bm_inode", "bm_block", "sb", "file", "Ls"}
			if !loContiene(nombresValidos, value) {
				return "", errors.New("el nombre es invalido")
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
	fmt.Println("Ingreso a comandos reportes")
	MBRmontado, err := global.GetMBRPartitionReport(rep.id)
	pathDes := global.ParticionMontada[rep.id]
	//hacer lo de arriba con sbMontado, pathDMontado
	fmt.Println("luego de sacar MBR", "err es: ", err)
	if err != nil {
		return err
	}
	fmt.Println("va a ingresar a switch")
	switch rep.name {
	case "mbr":
		err := reportes.ReporteMBR(MBRmontado, rep.path, pathDes)
		if err != nil {
			fmt.Printf("Error al generar el reporte MBR: %v\n", err)
			return err
		}

	case "disk":
		err := reportes.ReporteDisk(MBRmontado, rep.path, pathDes)
		if err != nil {
			fmt.Printf("Error al generar el reporte Disk: %v\n", err)
			return err
		}
	case "inode":
		fmt.Println("")
	case "bm_inode":
		fmt.Println("")
	}
	return nil
}

func loContiene(list []string, values string) bool {
	for _, valor := range list {
		if valor == values {
			return true
		}
	}
	return false
}
