package commands

import (
	stores "CLASE06/stores"
	structures "CLASE06/structures"
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
}

/*
   mkfs -id=vd1 -type=full
   mkfs -id=vd2
*/

func ParseMkfs(tokens []string) (*MKFS, error) {
	cmd := &MKFS{} // Crea una nueva instancia de MKFS

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mkfs
	re := regexp.MustCompile(`-id=[^\s]+|-type=[^\s]+`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("formato de parámetro inválido: %s", match)
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
				return nil, errors.New("el id no puede estar vacío")
			}
			cmd.id = value
		case "-type":
			// Verifica que el tipo sea "full"
			if value != "full" {
				return nil, errors.New("el tipo debe ser full")
			}
			cmd.typ = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que el parámetro -id haya sido proporcionado
	if cmd.id == "" {
		return nil, errors.New("faltan parámetros requeridos: -id")
	}

	// Si no se proporcionó el tipo, se establece por defecto a "full"
	if cmd.typ == "" {
		cmd.typ = "full"
	}

	// Aquí se puede agregar la lógica para ejecutar el comando mkfs con los parámetros proporcionados
	err := commandMkfs(cmd)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return cmd, nil // Devuelve el comando MKFS creado
}

func commandMkfs(mkfs *MKFS) error {
	// Obtener la partición montada
	mountedPartition, partitionPath, err := stores.GetMountedPartition(mkfs.id)
	if err != nil {
		return err
	}

	// Verificar la partición montada
	fmt.Println("\nPatición montada:")
	mountedPartition.PrintPartition()

	// Calcular el valor de n
	n := calculateN(mountedPartition)

	// Verificar el valor de n
	fmt.Println("\nValor de n:", n)

	// Inicializar un nuevo superbloque
	superBlock := createSuperBlock(mountedPartition, n)

	// Verificar el superbloque
	fmt.Println("\nSuperBlock:")
	superBlock.Print()

	// Crear los bitmaps
	err = superBlock.CreateBitMaps(partitionPath)
	if err != nil {
		return err
	}

	// Crear archivo users.txt
	err = superBlock.CreateUsersFile(partitionPath)
	if err != nil {
		return err
	}

	// Verificar superbloque actualizado
	fmt.Println("\nSuperBlock actualizado:")
	superBlock.Print()

	// Serializar el superbloque
	err = superBlock.Serialize(partitionPath, int64(mountedPartition.Part_start))
	if err != nil {
		return err
	}

	return nil
}

func calculateN(partition *structures.Partition) int32 {
	/*
		numerador = (partition_montada.size - sizeof(Structs::Superblock)
		denominador base = (4 + sizeof(Structs::Inodes) + 3 * sizeof(Structs::Fileblock))
		n = floor(numerador / denominador)
	*/

	numerator := int(partition.Part_size) - binary.Size(structures.SuperBlock{})
	denominator := 4 + binary.Size(structures.Inode{}) + 3*binary.Size(structures.FileBlock{}) // No importa que bloque poner, ya que todos tienen el mismo tamaño
	n := math.Floor(float64(numerator) / float64(denominator))

	return int32(n)
}

func createSuperBlock(partition *structures.Partition, n int32) *structures.SuperBlock {
	// Calcular punteros de las estructuras
	// Bitmaps
	bm_inode_start := partition.Part_start + int32(binary.Size(structures.SuperBlock{}))
	bm_block_start := bm_inode_start + n // n indica la cantidad de inodos, solo la cantidad para ser representada en un bitmap
	// Inodos
	inode_start := bm_block_start + (3 * n) // 3*n indica la cantidad de bloques, se multiplica por 3 porque se tienen 3 tipos de bloques
	// Bloques
	block_start := inode_start + (int32(binary.Size(structures.Inode{})) * n) // n indica la cantidad de inodos, solo que aquí indica la cantidad de estructuras Inode

	// Crear un nuevo superbloque
	superBlock := &structures.SuperBlock{
		S_filesystem_type:   2,
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
