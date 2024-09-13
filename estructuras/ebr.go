package estructuras

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

func (ebr *EBR) SerializarEBR(ubicacion string, offset int) error {

	archivo, err := os.OpenFile(ubicacion, os.O_WRONLY, 0644)

	if err != nil {
		return err
	}
	defer archivo.Close()

	_, err = archivo.Seek(int64(offset), 0)
	if err != nil {
		return err
	}

	err = binary.Write(archivo, binary.LittleEndian, ebr)

	if err != nil {
		return err
	}

	return nil
}

func (ebr *EBR) DeserializeEBR(ubicacion string, offset int) error {
	archivo, err := os.Open(ubicacion)

	if err != nil {
		return err
	}
	defer archivo.Close()
	_, err = archivo.Seek(int64(offset), 0)
	//calcula el tamaño
	ebrSize := binary.Size(ebr)

	if ebrSize <= 0 {
		return fmt.Errorf("Tamaño invalido de EBR: %d", ebrSize)
	}
	//lee los bytes
	buffer := make([]byte, ebrSize)
	_, err = archivo.Read(buffer)
	if err != nil {
		return err
	}
	//convierte los bytes a la estructura mbr
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, ebr)
	if err != nil {
		return err
	}

	return nil

}

func (ebr *EBR) PrintPart(path string) {
	fmt.Printf("EBR Details:\n")
	fmt.Printf("  Mount: %s\n", string(ebr.Ebr_part_moun[:]))
	fmt.Printf("  Fit: %s\n", string(ebr.Ebr_part_fit[:]))
	fmt.Printf("  Part Start: %d\n", ebr.Ebr_part_start)
	fmt.Printf("  Part Size: %d\n", ebr.Ebr_part_size)
	fmt.Printf("  Part Next: %d\n", ebr.Ebr_part_next)
	fmt.Printf("  Part Name: %s\n", string(ebr.Ebr_part_name[:]))
}
