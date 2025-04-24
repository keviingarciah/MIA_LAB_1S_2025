package reports

import (
	structures "backend/structures"
	utils "backend/utils"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// ReportInode genera un reporte de un inodo y lo guarda en la ruta especificada
func ReportInode(superblock *structures.SuperBlock, diskPath string, path string) error {
	// Crear las carpetas padre si no existen
	err := utils.CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Obtener el nombre base del archivo sin la extensión
	dotFileName, outputImage := utils.GetFileNames(path)

	// Iniciar el contenido DOT
	dotContent := `digraph G {
        node [shape=plaintext]
    `

	// Iterar sobre cada inodo
	for i := int32(0); i < superblock.S_inodes_count; i++ {
		inode := &structures.Inode{}
		// Deserializar el inodo
		err := inode.Deserialize(diskPath, int64(superblock.S_inode_start+(i*superblock.S_inode_size)))
		if err != nil {
			return err
		}

		// Convertir tiempos a string
		atime := time.Unix(int64(inode.I_atime), 0).Format(time.RFC3339)
		ctime := time.Unix(int64(inode.I_ctime), 0).Format(time.RFC3339)
		mtime := time.Unix(int64(inode.I_mtime), 0).Format(time.RFC3339)

		// Definir el contenido DOT para el inodo actual
		dotContent += fmt.Sprintf(`inode%d [label=<
            <table border="0" cellborder="1" cellspacing="0">
                <tr><td colspan="2"> REPORTE INODO %d </td></tr>
                <tr><td>i_uid</td><td>%d</td></tr>
                <tr><td>i_gid</td><td>%d</td></tr>
                <tr><td>i_size</td><td>%d</td></tr>
                <tr><td>i_atime</td><td>%s</td></tr>
                <tr><td>i_ctime</td><td>%s</td></tr>
                <tr><td>i_mtime</td><td>%s</td></tr>
                <tr><td>i_type</td><td>%c</td></tr>
                <tr><td>i_perm</td><td>%s</td></tr>
                <tr><td colspan="2">BLOQUES DIRECTOS</td></tr>
            `, i, i, inode.I_uid, inode.I_gid, inode.I_size, atime, ctime, mtime, rune(inode.I_type[0]), string(inode.I_perm[:]))

		// Agregar los bloques directos a la tabla hasta el índice 11
		for j, block := range inode.I_block {
			if j > 11 {
				break
			}
			dotContent += fmt.Sprintf("<tr><td>%d</td><td>%d</td></tr>", j+1, block)
		}

		// Agregar los bloques indirectos a la tabla
		dotContent += fmt.Sprintf(`
                <tr><td colspan="2">BLOQUE INDIRECTO</td></tr>
                <tr><td>%d</td><td>%d</td></tr>
                <tr><td colspan="2">BLOQUE INDIRECTO DOBLE</td></tr>
                <tr><td>%d</td><td>%d</td></tr>
                <tr><td colspan="2">BLOQUE INDIRECTO TRIPLE</td></tr>
                <tr><td>%d</td><td>%d</td></tr>
            </table>>];
        `, 13, inode.I_block[12], 14, inode.I_block[13], 15, inode.I_block[14])

		// Agregar enlace al siguiente inodo si no es el último
		if i < superblock.S_inodes_count-1 {
			dotContent += fmt.Sprintf("inode%d -> inode%d;\n", i, i+1)
		}
	}

	// Cerrar el contenido DOT
	dotContent += "}"

	// Crear el archivo DOT
	dotFile, err := os.Create(dotFileName)
	if err != nil {
		return err
	}
	defer dotFile.Close()

	// Escribir el contenido DOT en el archivo
	_, err = dotFile.WriteString(dotContent)
	if err != nil {
		return err
	}

	// Generar la imagen con Graphviz
	cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", outputImage)
	err = cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println("Imagen de los inodos generada:", outputImage)
	return nil
}
