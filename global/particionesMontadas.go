package global

import (
	structures "MIA_2S_P1_202200041/estructuras"
	"fmt"
)

const FCarnet string = "41"

var (
	ParticionMontada map[string]string = make(map[string]string)
)

// obtiene la particion montada con id
func GetParticionMontadaUID(id string) (*structures.PARTITION, string, error) {
	//path de la particion montada
	path := ParticionMontada[id]

	if path == "" {
		return nil, "", fmt.Errorf("particion %s no montada", id)
	}

	//instancia de MBR
	var mbr structures.MBR

	//deserializa el MBR
	err := mbr.Deserialize(path)
	if err != nil {
		return nil, "", err
	}

	//busca la particion que coincida con el id
	partition, _ := mbr.GetPID(id)
	if partition == nil {
		return nil, "", err
	}

	return partition, path, nil

}

func GetMPartitionReport(id string) (*structures.MBR, *structures.SUBERBLOCK, string, error) {
	path := ParticionMontada[id]
	if path == "" {
		return nil, nil, "", fmt.Errorf("particion %s no montada", id)
	}
	var mbr structures.MBR
	err := mbr.Deserialize(path)
	if err != nil {
		return nil, nil, "", err
	}
	partition, err := mbr.GetPID(id)
	if partition == nil {
		return nil, nil, "", err
	}

	var sb structures.SUBERBLOCK

	return &mbr, &sb, path, nil
}
