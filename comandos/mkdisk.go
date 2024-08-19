package comandos

import (
	structures "MIA_2S_P1_202200041/estructuras"
	utils "MIA_2S_P1_202200041/utils"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MKDISK struct { //estructura con los atributos de mkdisk
	size int
	unit string
	fit  string
	path string
}

func ParserMkDisk(tokens []string) (*MKDISK, error) {
	com := &MKDISK{} //instancia

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mkdisk
	re := regexp.MustCompile(`-size=\d+|-unit=[kKmM]|-fit=[bBfFwW]{2}|-path="[^"]+"|-path=[^\s]+`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar diferentes parámetros
		switch key {
		case "-size":
			// Convierte el valor del tamaño a un entero
			size, err := strconv.Atoi(value)
			if err != nil || size <= 0 {
				return nil, errors.New("el tamaño debe ser un número entero positivo")
			}
			com.size = size
		case "-unit":
			// Verifica que la unidad sea "K" o "M"
			if value != "K" && value != "M" {
				return nil, errors.New("la unidad debe ser K o M")
			}
			com.unit = strings.ToUpper(value)
		case "-fit":
			// Verifica que el ajuste sea "BF", "FF" o "WF"
			value = strings.ToUpper(value)
			if value != "BF" && value != "FF" && value != "WF" {
				return nil, errors.New("el ajuste debe ser BF, FF o WF")
			}
			com.fit = value
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, errors.New("el path no puede estar vacío")
			}
			com.path = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -size y -path hayan sido proporcionados
	if com.size == 0 {
		return nil, errors.New("faltan parámetros requeridos: -size")
	}
	if com.path == "" {
		return nil, errors.New("faltan parámetros requeridos: -path")
	}

	// Si no se proporcionó la unidad, se establece por defecto a "M"
	if com.unit == "" {
		com.unit = "M"
	}

	// Si no se proporcionó el ajuste, se establece por defecto a "FF"
	if com.fit == "" {
		com.fit = "FF"
	}

	// Crear el disco con los parámetros proporcionados
	err := comandosMkdisk(com)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return com, nil // Devuelve el comando MKDISK creado

}

func comandosMkdisk(com *MKDISK) error {
	sizeB, err := utils.ConvertirTamanioBy(com.size, com.unit)

	if err != nil {
		fmt.Println("Error converting size:", err)
		return err
	}

	// Crear el disco con el tamaño proporcionado
	err = crearDisco(com, sizeB)
	if err != nil {
		fmt.Println("Error creating disk:", err)
		return err
	}

	// Crear el MBR con el tamaño proporcionado
	err = crearMBR(com, sizeB)
	if err != nil {
		fmt.Println("Error creating MBR:", err)
		return err
	}

	return nil
}

func crearDisco(com *MKDISK, sizeB int) error {
	// Crear las carpetas necesarias
	err := os.MkdirAll(filepath.Dir(com.path), os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return err
	}

	// Crear el archivo binario
	file, err := os.Create(com.path)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	// Escribir en el archivo usando un buffer de 1 MB
	buffer := make([]byte, 1024*1024) // Crea un buffer de 1 MB
	for sizeB > 0 {
		writeSize := len(buffer)
		if sizeB < writeSize {
			writeSize = sizeB // Ajusta el tamaño de escritura si es menor que el buffer
		}
		if _, err := file.Write(buffer[:writeSize]); err != nil {
			return err // Devuelve un error si la escritura falla
		}
		sizeB -= writeSize // Resta el tamaño escrito del tamaño total
	}
	return nil
}

func crearMBR(com *MKDISK, sizeB int) error {
	// Seleccionar el tipo de ajuste
	var fitBy byte
	switch com.fit {
	case "FF":
		fitBy = 'F'
	case "BF":
		fitBy = 'B'
	case "WF":
		fitBy = 'W'
	default:
		fmt.Println("Invalid fit type")
		return nil
	}

	// Crear el MBR con los valores proporcionados
	mbr := &structures.MBR{
		Mbr_size:           int32(sizeB),
		Mbr_creation_date:  float32(time.Now().Unix()),
		Mbr_disk_signature: rand.Int31(),
		Mbr_disk_fit:       [1]byte{fitBy},
		Mbr_partitions: [4]structures.PARTITION{
			{Part_status: [1]byte{'N'}, Part_type: [1]byte{'N'}, Part_fit: [1]byte{'N'}, Part_start: -1, Part_size: -1, Part_name: [16]byte{'P'}, Part_correlative: 1, Part_id: -1},
			{Part_status: [1]byte{'N'}, Part_type: [1]byte{'N'}, Part_fit: [1]byte{'N'}, Part_start: -1, Part_size: -1, Part_name: [16]byte{'P'}, Part_correlative: 2, Part_id: -1},
			{Part_status: [1]byte{'N'}, Part_type: [1]byte{'N'}, Part_fit: [1]byte{'N'}, Part_start: -1, Part_size: -1, Part_name: [16]byte{'P'}, Part_correlative: 3, Part_id: -1},
			{Part_status: [1]byte{'N'}, Part_type: [1]byte{'N'}, Part_fit: [1]byte{'N'}, Part_start: -1, Part_size: -1, Part_name: [16]byte{'N'}, Part_correlative: 4, Part_id: -1},
		},
	}

	// Serializar el MBR en el archivo
	err := mbr.SerializarMBR(com.path)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}
