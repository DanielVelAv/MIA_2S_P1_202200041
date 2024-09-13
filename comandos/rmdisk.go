package comandos

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type RMDISK struct {
	path string
}

func ParserRMDISK(tokens []string) (string, error) {
	fmt.Println("Ingresa a parse rmdisk")
	com := &RMDISK{} //instancia

	args := strings.Join(tokens, " ")

	re := regexp.MustCompile(`-path="[^"]+"|-path=[^\s]+`)

	coincidencias := re.FindAllString(args, -1)

	for _, match := range coincidencias {

		div := strings.SplitN(match, "=", 2)
		if len(div) != 2 {
			return "", fmt.Errorf("Formato de parametro inválido: %v", match)
		}
		key, value := strings.ToLower(div[0]), div[1]

		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		switch key {

		case "-path":

			if value == "" {
				return "", errors.New("el path se encuentra vacio")
			}
			com.path = value
		default:
			return "", fmt.Errorf("parametro desconocido: %v", key)
		}

	}

	if com.path == "" {
		return "", errors.New("faltan parametro : -path")
	}

	err := comandosRmdisk(com)

	if err != nil {
		fmt.Println("Error en la ejecucion del comando:", err)
	}

	return "Disco eliminado correctamente", nil //devuelve el comando creado

}

func comandosRmdisk(com *RMDISK) error {
	fmt.Printf("ingreso a comands rmdisk")
	err := eliminarDisco(com) //se debe modificar, primero elimina el disco y luego el archivo
	if err != nil {
		fmt.Println("Error al eliminar el disco", err)
		return err
	}
	return nil //no hay errores en la ejecucion del comando
}

func eliminarDisco(com *RMDISK) error {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Esta seguro que desea eliminar el disco? (s/n)")
	eleccion, _ := reader.ReadString('\n')
	eleccion = strings.TrimSpace(eleccion)

	if eleccion == "s" {

		err := os.Remove(com.path)

		if err != nil {
			fmt.Println("Error al eliminar el disco", err)
			return err
		}
		fmt.Println("Disco eliminado exitosamente")

	} else if eleccion == "n" {
		fmt.Println("Operacion cancelada")
		return nil
	} else {
		fmt.Println("Respuesta invalida. Operacion cancelada")
		return errors.New("respuesta invalida")
	}

	return nil
}
