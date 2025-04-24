package commands

import (
	stores "backend/stores"
	structures "backend/structures"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

// MKFS estructura que representa el comando mkfs con sus parámetros
type MKFS struct {
	id  string // ID del disco
	typ string // Tipo de formato (full)
	fs  string // Sistema de archivos (2fs o 3fs)
}

/*
   mkfs -id=vd1 -type=full -fs=2fs
   mkfs -id=vd2 -fs=3fs
   mkfs -id=vd3
*/

func ParseMkfs(tokens []string) (string, error) {
	cmd := &MKFS{} // Crea una nueva instancia de MKFS

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mkfs
	re := regexp.MustCompile(`-id=[^\s]+|-type=[^\s]+|-fs=[23]fs`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Verificar que todos los tokens fueron reconocidos por la expresión regular
	if len(matches) != len(tokens) {
		// Identificar el parámetro inválido
		for _, token := range tokens {
			if !re.MatchString(token) {
				return "", fmt.Errorf("parámetro inválido: %s", token)
			}
		}
	}

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return "", fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar diferentes parámetros
		switch key {
		case "-id":
			// Verifica que el id no esté vacío
			if value == "" {
				return "", errors.New("el id no puede estar vacío")
			}
			cmd.id = value
		case "-type":
			// Verifica que el tipo sea "full"
			if value != "full" {
				return "", errors.New("el tipo debe ser full")
			}
			cmd.typ = value
		case "-fs":
			// Verifica que el sistema de archivos sea "2fs" o "3fs"
			if value != "2fs" && value != "3fs" {
				return "", errors.New("el sistema de archivos debe ser 2fs o 3fs")
			}
			cmd.fs = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que el parámetro -id haya sido proporcionado
	if cmd.id == "" {
		return "", errors.New("faltan parámetros requeridos: -id")
	}

	// Si no se proporcionó el tipo, se establece por defecto a "full"
	if cmd.typ == "" {
		cmd.typ = "full"
	}

	// Si no se proporcionó el sistema de archivos, se establece por defecto a "2fs"
	if cmd.fs == "" {
		cmd.fs = "2fs"
	}

	// Aquí se puede agregar la lógica para ejecutar el comando mkfs con los parámetros proporcionados
	err := commandMkfs(cmd)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("MKFS: Sistema de archivos creado exitosamente\n"+
		"-> ID: %s\n"+
		"-> Tipo: %s\n"+
		"-> Sistema de archivos: %s",
		cmd.id,
		cmd.typ,
		map[string]string{"2fs": "EXT2", "3fs": "EXT3"}[cmd.fs]), nil
}

func commandMkfs(mkfs *MKFS) error {
	fmt.Println("Creando sistema de archivos...", mkfs.fs)

	// Obtener la partición montada
	mountedPartition, partitionPath, err := stores.GetMountedPartition(mkfs.id)
	if err != nil {
		return err
	}

	// Calcular el valor de n
	n := calculateN(mountedPartition, mkfs.fs)

	fmt.Printf("Valor de N: %d\n", n)

	// Inicializar un nuevo superbloque
	superBlock := createSuperBlock(mountedPartition, n, mkfs.fs)

	// Crear los bitmaps
	err = superBlock.CreateBitMaps(partitionPath)
	if err != nil {
		return err
	}

	// Validar que sistema de archivos es
	if superBlock.S_filesystem_type == 3 {
		// Crear archivo users.txt ext3
		err = superBlock.CreateUsersFileExt3(partitionPath, int64(mountedPartition.Part_start+int32(binary.Size(structures.SuperBlock{}))))
		if err != nil {
			return err
		}
	} else {
		// Crear archivo users.txt ext2
		err = superBlock.CreateUsersFileExt2(partitionPath)
		if err != nil {
			return err
		}
	}

	// Serializar el superbloque
	err = superBlock.Serialize(partitionPath, int64(mountedPartition.Part_start))
	if err != nil {
		return err
	}

	return nil
}

func calculateN(partition *structures.Partition, fs string) int32 {
	// Numerador: tamaño de la partición menos el tamaño del superblock
	numerator := int(partition.Part_size) - binary.Size(structures.SuperBlock{})

	// Denominador base: 4 + tamaño de inodos + 3 * tamaño de bloques de archivo
	baseDenominator := 4 + binary.Size(structures.Inode{}) + 3*binary.Size(structures.FileBlock{})

	// Si el sistema de archivos es "3fs", se añade el tamaño del journaling al denominador
	temp := 0
	if fs == "3fs" {
		temp = binary.Size(structures.Journal{})
	}

	// Denominador final
	denominator := baseDenominator + temp

	// Calcular n
	n := math.Floor(float64(numerator) / float64(denominator))

	return int32(n)
}

func createSuperBlock(partition *structures.Partition, n int32, fs string) *structures.SuperBlock {
	// Calcular punteros de las estructuras
	journal_start, bm_inode_start, bm_block_start, inode_start, block_start := calculateStartPositions(partition, fs, n)

	fmt.Printf("Journal Start: %d\n", journal_start)
	fmt.Printf("Bitmap Inode Start: %d\n", bm_inode_start)
	fmt.Printf("Bitmap Block Start: %d\n", bm_block_start)
	fmt.Printf("Inode Start: %d\n", inode_start)
	fmt.Printf("Block Start: %d\n", block_start)

	// Tipo de sistema de archivos
	var fsType int32

	if fs == "2fs" {
		fsType = 2
	} else {
		fsType = 3
	}

	// Crear un nuevo superbloque
	superBlock := &structures.SuperBlock{
		S_filesystem_type:   fsType,
		S_inodes_count:      0,
		S_blocks_count:      0,
		S_free_inodes_count: int32(n),
		S_free_blocks_count: int32(n * 3),
		S_mtime:             float32(time.Now().Unix()),
		S_umtime:            float32(time.Now().Unix()),
		S_mnt_count:         1,
		S_magic:             0xEF53,
		S_inode_size:        int32(binary.Size(structures.Inode{})),
		S_block_size:        int32(binary.Size(structures.FileBlock{})),
		S_first_ino:         inode_start,
		S_first_blo:         block_start,
		S_bm_inode_start:    bm_inode_start,
		S_bm_block_start:    bm_block_start,
		S_inode_start:       inode_start,
		S_block_start:       block_start,
	}
	return superBlock
}

func calculateStartPositions(partition *structures.Partition, fs string, n int32) (int32, int32, int32, int32, int32) {
	superblockSize := int32(binary.Size(structures.SuperBlock{}))
	journalSize := int32(binary.Size(structures.Journal{}))
	inodeSize := int32(binary.Size(structures.Inode{}))

	// Inicializar posiciones
	// EXT2
	journalStart := int32(0)
	bmInodeStart := partition.Part_start + superblockSize
	bmBlockStart := bmInodeStart + n
	inodeStart := bmBlockStart + (3 * n)
	blockStart := inodeStart + (inodeSize * n)

	// Ajustar para EXT3
	if fs == "3fs" {
		journalStart = partition.Part_start + superblockSize
		bmInodeStart = journalStart + (journalSize * n)
		bmBlockStart = bmInodeStart + n
		inodeStart = bmBlockStart + (3 * n)
		blockStart = inodeStart + (inodeSize * n)
	}

	return journalStart, bmInodeStart, bmBlockStart, inodeStart, blockStart
}
