package comandos

import (
	structures "MIA_2S_P1_202200041/estructuras"
	"MIA_2S_P1_202200041/global"
	"MIA_2S_P1_202200041/utils"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type MOUNT struct {
	path string
	name string
}

func ParserMOUNT(tokens []string) (string, error) {
	fmt.Println("Ingresa a parse mount")
	com := &MOUNT{} //instancia de mount
	args := strings.Join(tokens, " ")
	re := regexp.MustCompile(`-path="[^"]+"|-path=[^\s]+|-name="[^"]+"|-name=[^\s]+`) //ER para los parametros
	coincidencias := re.FindAllString(args, -1)

	for _, match := range coincidencias {
		kv := strings.SplitN(match, "=", 2) //clave, valor
		if len(kv) != 2 {
			return "", fmt.Errorf("formato de parametro inválido: %s", match)
		}
		llave, valor := strings.ToLower(kv[0]), kv[1]

		//excluir comillas
		if strings.HasPrefix(valor, "\"") && strings.HasSuffix(valor, "\"") {
			valor = strings.Trim(valor, "\"")
		}

		switch llave {
		case "-path":
			if valor == "" {
				return "", errors.New("el path no puede estar vacío")
			}
			com.path = valor
		case "-name":
			if valor == "" {
				return "", errors.New("el nombre no puede estar vacío")
			}
			com.name = valor
		default:
			return "", fmt.Errorf("parametro desconocido: %s", llave)
		}
	}

	if com.path == "" {
		return "", errors.New("falta el parametro -path")
	}
	if com.name == "" {
		return "", errors.New("falta el parametro -name")
	}
	fmt.Println("Antes de montar la particion")
	//se monta la particion
	err := MountPartition(com)
	if err != nil {
		return "", err
	}

	return "La Particion ha sido montada correctamente", nil

}

func MountPartition(mount *MOUNT) error {
	var mbr structures.MBR

	//deserializamos el MBR
	err := mbr.Deserialize(mount.path)
	if err != nil {
		fmt.Println("Error al deserializar el MBR:", err)
		return err
	}

	//buscamos la particion que coincida con el nombre
	partition, indiceP := mbr.GetPName(mount.name)
	if partition == nil {
		fmt.Println("Error: No se encontró la partición con el nombre:", mount.name)
		return errors.New("la particion no existe")
	}
	//generacion de id unico
	idParticion, err := GenIdP(mount, indiceP)
	if err != nil {
		fmt.Println("Error al generar el ID de la particion:", err)
		return err
	}

	//se verifica si la particion ya esta montada
	err = verIDM(idParticion)

	if err != nil {
		fmt.Println("Error, particion ya montada:", err)
		return err
	}

	//se guarda la particion montada
	global.ParticionMontada[idParticion] = mount.path

	//modificamos la montacion para indicar que ya esta montada
	partition.MontarParticion(indiceP, idParticion)

	//guarda la particion modificada al mbr
	mbr.Mbr_partitions[indiceP] = *partition

	//serializamos el mbr
	err = mbr.SerializarMBR(mount.path)
	if err != nil {
		fmt.Println("Error al serializar el MBR:", err)
		return err
	}

	fmt.Println("La particion ha sido montada")

	//mostramos la lista de particiones montadas
	fmt.Println("Particiones montadas:")
	for id := range global.ParticionMontada {
		fmt.Println("ID: ", id)
	}

	fmt.Println("Todas las particiones:")
	for i, partition := range mbr.Mbr_partitions {
		fmt.Printf("Partición %d:\n", i+1)
		fmt.Printf("  Nombre: %s\n", partition.Part_name)
		fmt.Printf("  Estado: %s\n", partition.Part_status)
		fmt.Printf("  Tipo: %s\n", partition.Part_type)
		fmt.Printf("  Tamaño: %d\n", partition.Part_size)
		fmt.Printf("  ID de montaje: %d\n", partition.Part_id)
		fmt.Println()
	}

	return nil

}
func GenIdP(mount *MOUNT, indiceP int) (string, error) {
	//asigna una letra
	letter, err := utils.GetLetter(mount.path)
	if err != nil {
		fmt.Println("Error al obtener la letra de la particion:", err)
		return "", err
	}
	idParticion := fmt.Sprintf("%s%d%s", global.FCarnet, indiceP+1, letter)
	return idParticion, nil
}

func verIDM(id string) error {
	err := global.ParticionMontada[id]

	if err != "" {
		return errors.New("la particion " + id + " ya esta montada")
	}
	return nil
}
