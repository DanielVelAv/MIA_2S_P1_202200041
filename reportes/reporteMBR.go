package reportes

import (
	structures "MIA_2S_P1_202200041/estructuras"
	utils "MIA_2S_P1_202200041/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func ReporteMBR(mbr *structures.MBR, path string, pathDes string) error {
	err := utils.CrearDirectorios(path)
	if err != nil {
		return err
	}

	nombreArchivo, ImagenGenerada := utils.ObtenerNombreArchivo(path)

	contenido := fmt.Sprintf(`digraph G {
	node [shape=plaintext]
	tabla [label=<
      <table border="0" cellborder="1" cellspacing="0">
        <tr><td colspan="2" bgcolor="#46235a"><font color="white">Reporte MBR</font></td></tr>
        <tr><td>MBR_tamano</td><td>%d</td></tr>
        <tr><td bgcolor="#e7d7ea">MBR_fecha_creacion</td><td bgcolor="#e7d7ea">%s</td></tr>
	    <tr><td>MBR_disk_signature</td><td>%x</td></tr>
    	`, mbr.Mbr_size, time.Unix(int64(mbr.Mbr_creation_date), 0), mbr.Mbr_disk_signature)

	//particiones
	for i, part := range mbr.Mbr_partitions {

		//convertir a string
		PartNombre := strings.TrimRight(string(part.Part_name[:]), "\x00")
		//covnertir pStatus, pType, pFit a char
		pStatus := rune(part.Part_status[0])
		pType := rune(part.Part_type[0])
		pFit := rune(part.Part_fit[0])

		if part.Part_start != -1 {
			contenido += fmt.Sprintf(`
			<tr><td colspan="2" bgcolor="#46235a"><font color="white">Particion %d</font></td></tr>
			<tr><td>Part_status</td><td>%c</td></tr>
			<tr><td bgcolor="#e7d7ea">Part_type</td><td bgcolor="#e7d7ea">%c</td></tr>
			<tr><td>Part_fit</td><td>%c</td></tr>
			<tr><td bgcolor="#e7d7ea">Part_start</td><td bgcolor="#e7d7ea">%d</td></tr>
			<tr><td>Part_size</td><td>%d</td></tr>
			<tr><td bgcolor="#e7d7ea">Part_name</td><td bgcolor="#e7d7ea">%s</td></tr>
			`, i+1, pStatus, pType, pFit, part.Part_start, part.Part_size, PartNombre)

			if pType == 'E' {
				ebr := &structures.EBR{}
				err = ebr.DeserializeEBR(pathDes, int(part.Part_start))
				if err != nil {
					fmt.Println("Error al deserializar el EBR:", err)
					return err
				}

				contenido += `<tr><td colspan="2" bgcolor="#46235a"><font color="white">EBRs</font></td></tr>`

				for {
					//codigo tabla
					EBRNombre := strings.TrimRight(string(ebr.Ebr_part_name[:]), "\x00")
					contenido += fmt.Sprintf(`
						<tr><td colspan="2" bgcolor="#f07d7d"><font color="white">Particion Logica</font></td></tr>
						<tr><td>part_status</td><td>%s</td></tr>
						<tr><td bgcolor="#f4b5af">part_next</td><td bgcolor="#f4b5af">%d</td></tr>
						<tr><td>part_fit</td><td>%s</td></tr>	
						<tr><td bgcolor="#f4b5af">part_start</td><td bgcolor="#f4b5af">%d</td></tr>
						<tr><td>part_size</td><td>%d</td></tr>
						<tr><td bgcolor="#f4b5af">part_name</td><td bgcolor="#f4b5af">%s</td></tr>
					`, string(ebr.Ebr_part_moun[:]), ebr.Ebr_part_next, string(ebr.Ebr_part_fit[:]), ebr.Ebr_part_start, ebr.Ebr_part_size, EBRNombre)

					if ebr.Ebr_part_next == 0 {
						break
					}
					siguienteEbr := &structures.EBR{}
					err := siguienteEbr.DeserializeEBR(pathDes, int(ebr.Ebr_part_next))
					if err != nil {
						return err
					}
					ebr = siguienteEbr
				}

			}

		}

	}
	contenido += `</table>>]}`

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
