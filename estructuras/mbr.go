package estructuras

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type MBR struct {
	Mbr_size           int32        // Tamaño del MBR en bytes
	Mbr_creation_date  float32      // Fecha y hora de creación del MBR
	Mbr_disk_signature int32        // Firma del disco
	Mbr_disk_fit       [1]byte      // Tipo de ajuste
	Mbr_partitions     [4]PARTITION // Particiones del MBR
}

func (mbr *MBR) SerializarMBR(ubicacion string) error {

	archivo, err := os.OpenFile(ubicacion, os.O_WRONLY, 0644)

	if err != nil {
		return err
	}
	defer archivo.Close()

	err = binary.Write(archivo, binary.LittleEndian, mbr)

	if err != nil {
		return err
	}

	return nil
}

func (mbr *MBR) Deserialize(path string) error {
	archivo, err := os.Open(path)

	if err != nil {
		return err
	}
	defer archivo.Close()
	//calcula el tamaño
	mbrSize := binary.Size(mbr)

	if mbrSize <= 0 {
		return fmt.Errorf("Tamaño invalido de MBR: %d", mbrSize)
	}
	//lee los bytes
	buffer := make([]byte, mbrSize)
	_, err = archivo.Read(buffer)
	if err != nil {
		return err
	}
	//convierte los bytes a la estructura mbr
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}

	return nil

}

func (mbr *MBR) GetFirstA() (*PARTITION, int, int) {
	offset := binary.Size(mbr) //tamaño de MBR en bytes

	//recorre las particiones
	for i := 0; i < len(mbr.Mbr_partitions); i++ {
		//si la particion esta vacia y el start es -1
		if mbr.Mbr_partitions[i].Part_start == -1 {
			return &mbr.Mbr_partitions[i], offset, i
		} else {
			offset += int(mbr.Mbr_partitions[i].Part_size)
		}
	}
	return nil, -1, -1
}
