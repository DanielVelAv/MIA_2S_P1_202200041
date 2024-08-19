package estructuras

import (
	"encoding/binary"
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
