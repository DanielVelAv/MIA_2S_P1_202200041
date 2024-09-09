package estructuras

import (
	"encoding/binary"
	"os"
)

type EBR struct {
	Ebr_part_moun  [1]byte
	Ebr_part_fit   [1]byte
	Ebr_part_start int32
	Ebr_part_size  int32
	Ebr_part_next  int32
	Ebr_part_name  [16]byte
}

func (ebr *EBR) SerializarEBR(ubicacion string) error {

	archivo, err := os.OpenFile(ubicacion, os.O_WRONLY, 0644)

	if err != nil {
		return err
	}
	defer archivo.Close()

	err = binary.Write(archivo, binary.LittleEndian, ebr)

	if err != nil {
		return err
	}

	return nil
}
