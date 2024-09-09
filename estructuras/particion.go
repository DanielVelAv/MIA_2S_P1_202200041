package estructuras

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
