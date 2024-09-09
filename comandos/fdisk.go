package comandos

import (
	structures "MIA_2S_P1_202200041/estructuras"
	utils "MIA_2S_P1_202200041/utils"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type FDISK struct {
	size  int
	unit  string
	path  string
	tipos string
	fit   string
	name  string
}

func ParserFDISK(tokens []string) (*FDISK, error) {
	com := &FDISK{} //instancias

	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`-size=\d+|-unit=[kKmM]|-path="[^"]+"|-path=[^\s]+|-type=[pPeEIL]|-fit=[bBfFwW]{2}|-name="[^"]+"|-name=[^\s]+`)
	matches := re.FindAllString(args, -1)

	for _, match := range matches {
		kv := strings.SplitN(match, "=", 2) //clave, valor
		if len(kv) != 2 {
			return nil, fmt.Errorf("formato de parametro inválido: %s", match)
		}
		llave, valor := strings.ToLower(kv[0]), kv[1]

		//excluir comillas
		if strings.HasPrefix(valor, "\"") && strings.HasSuffix(valor, "\"") {
			valor = strings.Trim(valor, "\"")
		}

		switch llave {
		case "-size":
			size, err := strconv.Atoi(valor)
			if err != nil || size <= 0 {
				return nil, errors.New("el tamaño debe ser un entero positivo")
			}
			com.size = size
		case "-unit":
			valor = strings.ToUpper(valor)
			if valor != "K" && valor != "M" && valor != "B" {
				return nil, errors.New("la unidad debe ser K, M O B0")
			}
			com.unit = valor
		case "-path":
			if valor == "" {
				return nil, errors.New("el path no puede estar vacío")
			}
			com.path = valor
		case "-type":
			valor = strings.ToUpper(valor)
			if valor != "P" && valor != "E" && valor != "L" {
				return nil, errors.New("el tipo debe ser P, E o L")
			}
			com.tipos = valor
		case "-fit":
			valor = strings.ToUpper(valor)
			if valor != "BF" && valor != "FF" && valor != "WF" {
				return nil, errors.New("el ajuste debe ser BF, FF o WF")
			}
			com.fit = valor
		case "-name":
			if valor == "" {
				return nil, errors.New("el nombre no puede estar vacío")
			}
			com.name = valor
		default:
			return nil, fmt.Errorf("parametro desconocido: %v", llave)
		}
	}
	//verifica que se hayan ingresado los de tipo obligatorio, size path name

	if com.size == 0 {
		return nil, errors.New("faltan parámetros requeridos: -size")
	}
	if com.name == "" {
		return nil, errors.New("faltan parámetros requeridos: -name")
	}
	if com.path == "" {
		return nil, errors.New("faltan parámetros requeridos: -path")
	}

	if com.unit == "" {
		com.unit = "K"
	}
	if com.tipos == "" {
		com.tipos = "P"
	}
	if com.fit == "" {
		com.fit = "WF"
	}

	err := comandosFdisk(com)
	if err != nil {
		fmt.Println("Error en la ejecucion del comando:", err)
	}

	return com, nil //devuelve el comand creado

}

func comandosFdisk(fdisk *FDISK) error {
	tamanio, err := utils.ConvertirTamanioBy(fdisk.size, fdisk.unit)
	if err != nil {
		fmt.Println("Error al convertir el tamaño:", err)
		return err
	}
	if fdisk.tipos == "P" {
		err := crearParticionPrimaria(fdisk, tamanio)
		if err != nil {
			fmt.Println("Error al crear la particion primaria:", err)
			return err
		}
	} else if fdisk.tipos == "E" {
		err := crearParticionExtendida(fdisk, tamanio)
		if err != nil {
			fmt.Println("Error al crear la particion extendida:", err)
			return err
		}
	} else if fdisk.tipos == "L" {
		//err := crearParticionLogica(fdisk, tamanio)
		if err != nil {
			fmt.Println("Error al crear la particion logica:", err)
			return err
		}
	}
	return nil
}

func crearParticionPrimaria(fdisk *FDISK, tamanio int) error {
	var com structures.MBR //instancia de mbr

	err := com.Deserialize(fdisk.path) // se deserializa el mbr para obtener la informacion
	if err != nil {                    //en caso que de eror
		fmt.Println("Error al deserializar el MBR:", err)
		return err
	}

	//busca la primera particion disponible
	particionDisponible, inicioParticion, indiceParticion := com.GetFirstA()

	if particionDisponible == nil {
		fmt.Println("No hay espacio suficiente para crear la particion")
	}
	particionDisponible.CrearParticion(inicioParticion, tamanio, fdisk.tipos, fdisk.fit, fdisk.name)

	if particionDisponible != nil {
		com.Mbr_partitions[indiceParticion] = *particionDisponible
	}

	//se serializa el mbr
	err = com.SerializarMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error al serializar el MBR:", err)
	}

	return nil
}

func crearParticionExtendida(fdisk *FDISK, tamanio int) error {
	var com structures.MBR //instancia de mbr

	err := com.Deserialize(fdisk.path) // se deserializa el mbr para obtener la informacion
	if err != nil {                    //en caso que de eror
		fmt.Println("Error al deserializar el MBR:", err)
		return err
	}

	//busca la primera particion disponible
	particionDisponible, inicioParticion, indiceParticion := com.GetFirstA()

	if particionDisponible == nil {
		fmt.Println("No hay espacio suficiente para crear la particion")
	}
	particionDisponible.CrearParticionE(inicioParticion, tamanio, fdisk.tipos, fdisk.fit, fdisk.name)

	if particionDisponible != nil {
		com.Mbr_partitions[indiceParticion] = *particionDisponible
	}
	err = crearEBR(fdisk, inicioParticion, tamanio, fdisk.name)
	if err != nil {
		fmt.Println("Error al crear el EBR:", err)
		return err
	}
	//se serializa el mbr
	err = com.SerializarMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error al serializar el MBR:", err)
	}

	return nil
}

//func crearParticionLogica{

//}

func crearEBR(com *FDISK, pStart int, pSize int, pName string) error {
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
	var pNext int
	pNext = pStart + pSize
	actualN := pName + "EBR"

	// Crear el EBR con los valores proporcionados
	ebr := &structures.EBR{
		Ebr_part_moun:  [1]byte{0},
		Ebr_part_fit:   [1]byte{fitBy},
		Ebr_part_start: int32(pStart),
		Ebr_part_size:  int32(pSize),
		Ebr_part_next:  int32(pNext),
		Ebr_part_name:  [16]byte{},
	}
	copy(ebr.Ebr_part_name[:], actualN)

	// Serializar el MBR en el archivo
	err := ebr.SerializarEBR(com.path)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}
