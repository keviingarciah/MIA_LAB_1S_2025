package structures

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

type Journal struct {
	J_count   int32       // 4 bytes
	J_content Information // 110 bytes
	// Total: 114 bytes
}

type Information struct {
	I_operation [10]byte // 10 bytes
	I_path      [32]byte // 32 bytes
	I_content   [64]byte // 64 bytes
	I_date      float32  // 4 bytes
	// Total: 110 bytes
}

// SerializeJournal escribe la estructura Journal en un archivo binario
func (journal *Journal) Serialize(path string, journauling_start int64) error {
	// Calcular la posición en el archivo
	offset := journauling_start + (int64(binary.Size(Journal{})) * int64(journal.J_count))

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mover el puntero del archivo a la posición especificada
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Serializar la estructura Journal directamente en el archivo
	err = binary.Write(file, binary.LittleEndian, journal)
	if err != nil {
		return err
	}

	return nil
}

// DeserializeJournal lee la estructura Journal desde un archivo binario
func (journal *Journal) Deserialize(path string, offset int64) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	// Mover el puntero del archivo a la posición especificada
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Deserializar la estructura Journal directamente desde el archivo
	err = binary.Read(file, binary.LittleEndian, journal)
	if err != nil {
		return err
	}

	return nil
}

// PrintJournal imprime en consola la estructura Journal
func (journal *Journal) Print() {
	// Convertir el tiempo de montaje a una fecha
	date := time.Unix(int64(journal.J_content.I_date), 0)

	fmt.Println("Journal:")
	fmt.Printf("J_count: %d", journal.J_count)
	fmt.Println("Information:")
	fmt.Printf("I_operation: %s", string(journal.J_content.I_operation[:]))
	fmt.Printf("I_path: %s", string(journal.J_content.I_path[:]))
	fmt.Printf("I_content: %s", string(journal.J_content.I_content[:]))
	fmt.Printf("I_date: %s", date.Format(time.RFC3339))
}
