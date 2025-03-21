package structures

import "fmt"

type Partition struct {
	Part_status      [1]byte  // Estado de la partición
	Part_type        [1]byte  // Tipo de partición
	Part_fit         [1]byte  // Ajuste de la partición
	Part_start       int32    // Byte de inicio de la partición
	Part_size        int32    // Tamaño de la partición
	Part_name        [16]byte // Nombre de la partición
	Part_correlative int32    // Correlativo de la partición
	Part_id          [4]byte  // ID de la partición
}

/*
Part Status:
	N: Disponible
	0: Creado
	1: Montado

Esto queda a su criterio.
*/

// Crear una partición con los parámetros proporcionados
func (p *Partition) CreatePartition(partStart, partSize int, partType, partFit, partName string) {
	// Asignar status de la partición
	p.Part_status[0] = '0' // El valor '0' indica que la partición ha sido creada

	// Asignar el byte de inicio de la partición
	p.Part_start = int32(partStart)

	// Asignar el tamaño de la partición
	p.Part_size = int32(partSize)

	// Asignar el tipo de partición
	if len(partType) > 0 {
		p.Part_type[0] = partType[0]
	}

	// Asignar el ajuste de la partición
	if len(partFit) > 0 {
		p.Part_fit[0] = partFit[0]
	}

	// Asignar el nombre de la partición
	copy(p.Part_name[:], partName)
}

// Montar una partición por el id
func (p *Partition) MountPartition(correlative int, id string) error {
	// Asignar status de la partición
	p.Part_status[0] = '1' // El valor '1' indica que la partición ha sido montada

	// Asignar correlativo a la partición
	p.Part_correlative = int32(correlative)

	// Asignar ID a la partición
	copy(p.Part_id[:], id)

	return nil
}

// Imprimir los valores de la partición
func (p *Partition) PrintPartition() {
	fmt.Printf("Part_status: %c\n", p.Part_status[0])
	fmt.Printf("Part_type: %c\n", p.Part_type[0])
	fmt.Printf("Part_fit: %c\n", p.Part_fit[0])
	fmt.Printf("Part_start: %d\n", p.Part_start)
	fmt.Printf("Part_size: %d\n", p.Part_size)
	fmt.Printf("Part_name: %s\n", string(p.Part_name[:]))
	fmt.Printf("Part_correlative: %d\n", p.Part_correlative)
	fmt.Printf("Part_id: %s\n", string(p.Part_id[:]))
}
