package estructuras

import (
	"fmt"
	"strconv"
)

type PARTITION struct {
	Part_status      [1]byte
	Part_type        [1]byte
	Part_fit         [1]byte
	Part_start       int32
	Part_size        int32
	Part_name        [16]byte
	Part_correlative int32
	Part_id          int32
}

func (p *PARTITION) CrearParticion(Pstart, pSize int, pType, pFit, pName string) {
	p.Part_status[0] = '0' // indica que se creo
	p.Part_start = int32(Pstart)
	p.Part_size = int32(pSize)

	//tipo de particion
	if len(pType) > 0 {
		p.Part_type[0] = pType[0]
	}
	//ajuste de particion
	if len(pFit) > 0 {
		p.Part_fit[0] = pFit[0]
	}
	//asigna el nombre
	copy(p.Part_name[:], pName)
}

func (p *PARTITION) CrearParticionE(pStart, pSize int, pType, pFit, pName string) {
	p.Part_status[0] = '0'       // indica que se creo
	p.Part_start = int32(pStart) //inicio de la particion
	p.Part_size = int32(pSize)   //tamaño de la particion

	//tipo de particion
	if len(pType) > 0 {
		p.Part_type[0] = pType[0]
	}
	//ajuste de particion
	if len(pFit) > 0 {
		p.Part_fit[0] = pFit[0]
	}
	//asigna el nombre
	copy(p.Part_name[:], pName)

}

func (p *PARTITION) PrintPart(ubicacion string) {
	fmt.Printf("Status: %c\n", p.Part_status[0])
	fmt.Printf("Type: %c\n", p.Part_type[0])
	fmt.Printf("Fit: %c\n", p.Part_fit[0])
	fmt.Printf("Start: %d\n", p.Part_start)
	fmt.Printf("Size: %d\n", p.Part_size)
	fmt.Printf("Name: %s\n", string(p.Part_name[:]))
	fmt.Printf("Correlative: %d\n", p.Part_correlative)
	fmt.Printf("ID: %d\n", p.Part_id)
	fmt.Println("-------------------------")

	if p.Part_type[0] == 'E' {
		ebr := &EBR{}
		err := ebr.DeserializeEBR(ubicacion, int(p.Part_start))
		if err != nil {
			fmt.Println("Error al deserializar el EBR:", err)
		} else {
			fmt.Printf("EBR dentro de la partición extendida: %+v\n", ebr)
		}
	}

}

func (particion *PARTITION) MontarParticion(correlativo int, id string) error {
	// se asigna correlativo a la partición
	particion.Part_correlative = int32(correlativo) + 1

	//se asigna el id a la particion
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("error convirtiendo: %v", err)
	}
	particion.Part_id = int32(idInt)

	return nil
}
