package reportes

import (
	structures "MIA_2S_P1_202200041/estructuras"
	"MIA_2S_P1_202200041/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ReporteDisk(mbr *structures.MBR, path string, pathDes string) error {
	err := utils.CrearDirectorios(path)
	if err != nil {
		return err
	}

	nombreArchivo, ImagenGenerada := utils.ObtenerNombreArchivo(path)
	//comprueba si hay una particion extendidad
	hayE, tExtendida := hayExtendida(mbr)
	fmt.Println("Espacio total Disco: ", mbr.Mbr_size)
	fmt.Println("Hay extendida: ", hayE, "Tamaño extendida: ", tExtendida)
	nL, espacioTL, EspacioDisponibleE := numeroLogicas(mbr, pathDes, hayE)
	println("Numero de logicas: ", nL, "espacio total logicas", espacioTL, nombreArchivo, ImagenGenerada)

	nPrimarias, espTP := numeroPrimarias(mbr, tExtendida)

	fmt.Println("Numero de primarias: ", nPrimarias, "Espacio total primarias: ", espTP)

	//? variables con espacio total primarias y extendida si hay
	//espacioDisponibleMBR := espacioDMBR()

	//cuenta cuantas primarias hay y si hay extendida,primarias no cambia estructura, y en caso de que no haya extendida
	// el espacio de extendida queda como espacio vacio
	pathFinal := strings.Split(pathDes, "/")
	ultimo := pathFinal[len(pathFinal)-1]
	NombreDisco := strings.TrimSuffix(ultimo, ".mia")

	contenido := fmt.Sprintf(`digraph G {
  		fontname="Helvetica,Arial,sans-serif";
  		node [fontname="Helvetica,Arial,sans-serif", shape=plaintext, style=filled, fillcolor=lightgrey];
		graph [label = "%s" labelloc="t"]
		a0 [label=<
    		<TABLE border="1" cellspacing="0" cellpadding="10" style="filled" gradientangle="315">
			<TR>
				<TD rowspan="2" border="1">MBR</TD>`, NombreDisco)
	contenido += espTP
	if hayE == true {
		numcolumnas := nL * 2
		if EspacioDisponibleE == 0 {

		} else {
			numcolumnas += 1
		}
		contenido += fmt.Sprintf(`
		<TD colspan="%d">Extendida</TD>
	    
		</TR>

         `, numcolumnas)
		contenido += espacioTL
	} else {
		contenido += "</TR>"
	}

	contenido += "</TABLE>>];}"

	file, err := os.Create(nombreArchivo)
	if err != nil {
		return fmt.Errorf("error al crear el archivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(contenido)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo: %v", err)
	}

	com := exec.Command("dot", "-Tpng", nombreArchivo, "-o", ImagenGenerada)
	err = com.Run()
	if err != nil {
		return fmt.Errorf("error al generar la imagen: %v", err)
	}

	return nil
}

func hayExtendida(mbr *structures.MBR) (bool, int) {
	for i, part := range mbr.Mbr_partitions {

		//covnertir pStatus, pType, pFit a char
		pType := rune(part.Part_type[0])
		fmt.Println(i)
		if pType == 'E' {

			return true, int(part.Part_size)
		}

	}
	return false, 0

}

func numeroLogicas(mbr *structures.MBR, pathDes string, ExR bool) (int, string, float64) {
	// se puede ingresar datos desde aqui a la grafica
	porcentajesR := fmt.Sprintf(`
        <TR>
	 `)
	if ExR == true {
		porcentajeLibre := 0.0
		contador := 0
		espacioUsado := 0.0
		for i, part := range mbr.Mbr_partitions {

			pType := rune(part.Part_type[0])
			fmt.Println(i)
			if pType == 'E' {
				ebr := &structures.EBR{}
				err := ebr.DeserializeEBR(pathDes, int(part.Part_start))
				if err != nil {
					fmt.Println("Error al deserializar el EBR:", err)
				}

				porcentajeExtendida := (float64(part.Part_size) * 100) / float64(mbr.Mbr_size)
				porcentajeTL := 0.0

				for {
					//codigo tabla
					contador += 1
					espacioUsado += float64(ebr.Ebr_part_size)
					porcentaje := (float64(ebr.Ebr_part_size) * 100) / float64(mbr.Mbr_size)
					fmt.Println("Porcentaje de la particion logica: ", porcentaje)
					fmt.Println("Porcentaje extendida: ", porcentajeExtendida)
					porcentajeTL += porcentaje
					porcentajesR += fmt.Sprintf(`
  		 			<TD>EBR</TD>
					<TD border="1">Lógica <BR/> %.2f%% del disco</TD>
					`, porcentaje)
					if ebr.Ebr_part_next == 0 {
						break
					}

					siguienteEbr := &structures.EBR{}
					err := siguienteEbr.DeserializeEBR(pathDes, int(ebr.Ebr_part_next))
					if err != nil {
						fmt.Println("Error al deserializar el EBR:", err)
					}

					ebr = siguienteEbr

				}
				porcentajeLibre = porcentajeExtendida - porcentajeTL
				if porcentajeLibre == 0 {
					porcentajesR += "</TR>"
				} else {
					porcentajesR += fmt.Sprintf(`
					<TD border="1">Libre <BR/> %.2f%% del Disco </TD>
					</TR>`, porcentajeLibre)
				}

			}

		}
		return contador, porcentajesR, porcentajeLibre
	}
	return 0, porcentajesR, 0
}

func numeroPrimarias(mbr *structures.MBR, tamanioExtendida int) (int, string) {
	tamanioTotal := mbr.Mbr_size
	contador := 0
	tamanioT := 0
	porcentajesR := ""
	for i, part := range mbr.Mbr_partitions {
		pType := rune(part.Part_type[0])
		fmt.Println(i)
		if pType == 'P' {
			contador += 1
			tamanioT += int(part.Part_size)
			porcentaje := (float64(part.Part_size) * 100) / float64(tamanioTotal)
			fmt.Println("Porcentaje de la particion primaria: ", porcentaje)
			porcentajesR += fmt.Sprintf(`
		  	<TD rowspan="2" border="1">Primaria <BR/> %.2f%% del disco</TD>
            `, porcentaje)
		}

	}
	tamanioSE := int(tamanioTotal) - tamanioExtendida //tamaño sin extendida, solo de primarias
	//espacio libre es espacio sin extendida - particiones primarias
	ELibre := tamanioSE - tamanioT
	if ELibre == 0 {

	} else {
		porcentajesR += fmt.Sprintf(`
		<TD rowspan="2" border="1">Libre <BR/> %.2f%% del Disco </TD>
	`, (float64(ELibre)*100)/float64(tamanioTotal))
	}

	return contador, porcentajesR
}
