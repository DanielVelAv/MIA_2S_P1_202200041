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
		err := crearParticionLogica(fdisk, tamanio)
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
		return errors.New("No hay suficiente espacio para crear la particion primaria")
		//fmt.Println("No hay espacio suficiente para crear la particion")
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

	for _, part := range com.Mbr_partitions {
		if part.Part_type[0] == 'E' {
			return errors.New("Ya existe una particion extendida")
		}
	}

	//busca la primera particion disponible
	particionDisponible, inicioParticion, indiceParticion := com.GetFirstA()

	if particionDisponible == nil {
		return errors.New("No hay espacio suficiente para crear la particion")
	}

	particionDisponible.CrearParticionE(inicioParticion, tamanio, fdisk.tipos, fdisk.fit, fdisk.name)

	if particionDisponible != nil {
		com.Mbr_partitions[indiceParticion] = *particionDisponible
	}
	fmt.Println("Particion extendida creada: ")
	err = crearEBRNil(fdisk, inicioParticion, tamanio, fdisk.name)
	if err != nil {
		fmt.Println("Error al crear el EBR:", err)
		return err
	}

	particionDisponible.PrintPart(fdisk.path)

	ebr := &structures.EBR{
		Ebr_part_start: int32(inicioParticion),
		Ebr_part_size:  0,
		Ebr_part_next:  -1,
	}
	err = ebr.DeserializeEBR(fdisk.path, inicioParticion)
	if err != nil {
		fmt.Println("Error al deserializar el EBR bulo:", err)
		return err
	}
	fmt.Printf("EBR encontrado: %+v\n", ebr)

	//se serializa el mbr
	err = com.SerializarMBR(fdisk.path)
	if err != nil {
		fmt.Println("Error al serializar el MBR:", err)
	}

	return nil
}

func crearParticionLogica(fdisk *FDISK, tamanio int) error {
	fmt.Println("Ingresa a crear particion logica")
	mbr, err := obtenerMBR(fdisk.path)
	if err != nil {
		return err
	}

	extendida, err := ExisteExtendida(mbr) //existe extendida puede ser una verificacion solamente
	if err != nil {
		fmt.Println("No se encontro la particion extendida:", err)
		return err
	}
	fmt.Println("Particion extendida encontrada: ", extendida)

	// Verificar el espacio disponible en la partición extendida
	// saco el MBR, verifico el espacio de la extendida usando el ultimo EBR insertado y
	//espacio a insertar de la particion logica
	espacioDisponible, err := EspacioDisponibleExtendida(fdisk, extendida, tamanio)
	if err != nil {
		fmt.Println("Valio v*rga buscando espacio")
	}
	fmt.Println("Espacio disponible en la partición extendida:", espacioDisponible, "tamaño a querer insertar: ", tamanio)
	if espacioDisponible < tamanio {
		return errors.New("No hay suficiente espacio disponible en la partición extendida")
	}
	fmt.Println("Espacio suficiente en la partición extendida")

	ebr := &structures.EBR{}
	err = ebr.DeserializeEBR(fdisk.path, int(extendida.Part_start))

	if err != nil {
		fmt.Println("Error al deserializar el EBR:", err)
		return err
	}

	contadorParticionesLogicas := 0

	//busca el ebr nulo anterior
	for {
		contadorParticionesLogicas++
		if ebr.Ebr_part_next == 0 {
			break
		}
		siguienteEbr := &structures.EBR{}
		err := siguienteEbr.DeserializeEBR(fdisk.path, int(ebr.Ebr_part_next))
		if err != nil {
			return err
		}
		ebr = siguienteEbr
	}

	// Verificar si se excede el límite de 4 particiones lógicas
	/*if contadorParticionesLogicas >= 4 {
		return errors.New("La partición extendida ya contiene el máximo de 4 particiones lógicas")
	}*/

	fmt.Printf("numero de particiones logicas: %d\n", contadorParticionesLogicas)

	iniciologica := int(ebr.Ebr_part_start) + int(ebr.Ebr_part_size)

	fmt.Println("Inicio logico de la particion logica: ", iniciologica)

	nuevaEBR := &structures.EBR{
		Ebr_part_moun:  [1]byte{'0'},
		Ebr_part_fit:   [1]byte{fdisk.fit[0]},
		Ebr_part_start: int32(iniciologica),
		Ebr_part_size:  int32(tamanio),
		Ebr_part_next:  0,
		Ebr_part_name:  [16]byte{},
	}
	copy(nuevaEBR.Ebr_part_name[:], fdisk.name)
	ebr.Ebr_part_next = int32(iniciologica)

	err = ebr.SerializarEBR(fdisk.path, int(ebr.Ebr_part_start))
	if err != nil {
		fmt.Println("Error al actualizar el EBR anterior:", err)
		return err
	}

	err = nuevaEBR.SerializarEBR(fdisk.path, iniciologica)

	if err != nil {
		fmt.Println("Error al crear el nuevo EBR:", err)
		return err
	}

	fmt.Println("Particion logica creada: ")
	//nuevaEBR.PrintPart(fdisk.path)

	//para imprimir todos los ebrs

	ebr = &structures.EBR{}
	err = ebr.DeserializeEBR(fdisk.path, int(extendida.Part_start))
	if err != nil {
		fmt.Println("Error al deserializar el EBR:", err)
		return err
	}

	fmt.Println("Lista de EBRs:")
	for {
		ebr.PrintPart(fdisk.path)
		if ebr.Ebr_part_next == 0 {
			break
		}
		siguienteEbr := &structures.EBR{}
		err := siguienteEbr.DeserializeEBR(fdisk.path, int(ebr.Ebr_part_next))
		if err != nil {
			return err
		}
		ebr = siguienteEbr
	}

	return nil

}

func crearEBRNil(com *FDISK, pStart int, pSize int, pName string) error {

	// Crear el EBR con los valores proporcionados
	ebr := &structures.EBR{
		Ebr_part_moun:  [1]byte{0},
		Ebr_part_fit:   [1]byte{'W'},
		Ebr_part_start: int32(pStart),
		Ebr_part_size:  int32(0),
		Ebr_part_next:  int32(0),
		Ebr_part_name:  [16]byte{},
	}
	copy(ebr.Ebr_part_name[:], "")

	// Serializar el MBR en el archivo
	err := ebr.SerializarEBR(com.path, pStart)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return nil
}

func ExisteExtendida(mbr *structures.MBR) (*structures.PARTITION, error) {
	for _, part := range mbr.Mbr_partitions {
		if part.Part_type[0] == 'E' {
			return &part, nil
		}
	}
	return nil, errors.New("No hay particion extendida")
}

func EspacioDisponibleExtendida(fdisk *FDISK, extendida *structures.PARTITION, tamanio int) (int, error) {

	ebr := &structures.EBR{}
	err := ebr.DeserializeEBR(fdisk.path, int(extendida.Part_start))
	if err != nil {
		return 0, err
	}
	//fmt.Println("Todo bien deserializando", ebr)
	//inicialmente tiene todo
	espacioDisponible := int(extendida.Part_size)
	//fmt.Println("espacio disponible antes de verificar ebrs: ", espacioDisponible)

	//restar lo ocupado por el primer ebr
	espacioDisponible -= int(ebr.Ebr_part_size)
	//fmt.Println("espacio disponible despues de verificar el primer ebr: ", espacioDisponible)

	for ebr.Ebr_part_next != 0 {
		siguienteEBR := &structures.EBR{}
		err := siguienteEBR.DeserializeEBR(fdisk.path, int(ebr.Ebr_part_next))
		if err != nil {
			return 0, err
		}
		espacioDisponible -= int(siguienteEBR.Ebr_part_size)
		ebr = siguienteEBR
	}

	// Restar el espacio ocupado por el último EBR nulo
	espacioDisponible -= int(ebr.Ebr_part_size)
	//fmt.Println("espacio disponible despues de verificar todos los ebrs: ", espacioDisponible)
	return espacioDisponible, nil

}

func obtenerMBR(path string) (*structures.MBR, error) {
	var mbr structures.MBR
	err := mbr.Deserialize(path)
	if err != nil {
		return nil, err
	}
	return &mbr, nil
}
