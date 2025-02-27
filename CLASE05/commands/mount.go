package commands

import (
	stores "CLASE05/stores"
	structures "CLASE05/structures"
	utils "CLASE05/utils"
	"errors" // Paquete para manejar errores y crear nuevos errores con mensajes personalizados
	"fmt"    // Paquete para formatear cadenas y realizar operaciones de entrada/salida
	"regexp" // Paquete para trabajar con expresiones regulares, útil para encontrar y manipular patrones en cadenas

	// Paquete para convertir cadenas a otros tipos de datos, como enteros
	"strings" // Paquete para manipular cadenas, como unir, dividir, y modificar contenido de cadenas
)

// MOUNT estructura que representa el comando mount con sus parámetros
type MOUNT struct {
	path string // Ruta del archivo del disco
	name string // Nombre de la partición
}

/*
	mount -path=/home/Disco1.mia -name=Part1 #id=341a
	mount -path=/home/Disco2.mia -name=Part1 #id=342a
	mount -path=/home/Disco3.mia -name=Part2 #id=343a
*/

// CommandMount parsea el comando mount y devuelve una instancia de MOUNT
func ParseMount(tokens []string) (*MOUNT, error) {
	cmd := &MOUNT{} // Crea una nueva instancia de MOUNT

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mount
	re := regexp.MustCompile(`-path="[^"]+"|-path=[^\s]+|-name="[^"]+"|-name=[^\s]+`)
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
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		case "-name":
			// Verifica que el nombre no esté vacío
			if value == "" {
				return nil, errors.New("el nombre no puede estar vacío")
			}
			cmd.name = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -path y -name hayan sido proporcionados
	if cmd.path == "" {
		return nil, errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.name == "" {
		return nil, errors.New("faltan parámetros requeridos: -name")
	}

	// Montamos la partición
	err := commandMount(cmd)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return cmd, nil // Devuelve el comando MOUNT creado
}

func commandMount(mount *MOUNT) error {
	// Crear una instancia de MBR
	var mbr structures.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.DeserializeMBR(mount.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	// Buscar la partición con el nombre especificado
	partition, indexPartition := mbr.GetPartitionByName(mount.name)
	if partition == nil {
		fmt.Println("Error: la partición no existe")
		return errors.New("la partición no existe")
	}

	/* SOLO PARA VERIFICACIÓN */
	// Print para verificar que la partición se encontró correctamente
	fmt.Println("\nPartición disponible:")
	partition.PrintPartition()

	// Generar un id único para la partición
	idPartition, partitionCorrelative, err := generatePartitionID(mount)
	if err != nil {
		fmt.Println("Error generando el id de partición:", err)
		return err
	}

	//  Guardar la partición montada en la lista de montajes globales
	stores.MountedPartitions[idPartition] = mount.path

	// Modificamos la partición para indicar que está montada
	partition.MountPartition(partitionCorrelative, idPartition)

	/* SOLO PARA VERIFICACIÓN */
	// Print para verificar que la partición se haya montado correctamente
	fmt.Println("\nPartición montada (modificada):")
	partition.PrintPartition()

	// Guardar la partición modificada en el MBR
	mbr.Mbr_partitions[indexPartition] = *partition

	// Serializar la estructura MBR en el archivo binario
	err = mbr.SerializeMBR(mount.path)
	if err != nil {
		fmt.Println("Error serializando el MBR:", err)
		return err
	}

	return nil
}

func generatePartitionID(mount *MOUNT) (string, int, error) {
	// Asignar una letra a la partición y obtener el índice
	letter, partitionCorrelative, err := utils.GetLetterAndPartitionCorrelative(mount.path)
	if err != nil {
		fmt.Println("Error obteniendo la letra:", err)
		return "", 0, err
	}

	// Crear id de partición
	idPartition := fmt.Sprintf("%s%d%s", stores.Carnet, partitionCorrelative, letter)

	return idPartition, partitionCorrelative, nil
}
