package main

import (
	analyzer "CLASE05/analyzer"
	"bufio" // Importa el paquete "bufio" para operaciones de buffer de entrada/salida
	"fmt"   // Importa el paquete "fmt" para formatear e imprimir texto
	"os"    // Importa el paquete "os" para interactuar con el sistema operativo
)

func main() {
	// Crea un nuevo escáner que lee desde la entrada estándar (teclado)
	scanner := bufio.NewScanner(os.Stdin)

	// Bucle infinito para leer comandos del usuario
	for {
		fmt.Print(">>> ") // Imprime el prompt para el usuario

		// Lee la siguiente línea de entrada del usuario
		if !scanner.Scan() {
			break // Si no hay más líneas para leer, rompe el bucle
		}

		// Obtiene el texto ingresado por el usuario
		input := scanner.Text()

		// Llama a la función Analyzer del paquete analyzer para analizar el comando ingresado
		_, err := analyzer.Analyzer(input)
		if err != nil {
			// Si hay un error al analizar el comando, imprime el error y continúa con el siguiente comando
			fmt.Println("Error:", err)
			continue
		}

		// Comentado: Aquí podrías imprimir el comando analizado
		// fmt.Printf("Parsed Command: %+v\n", cmd)
	}

	// Verifica si hubo algún error al leer la entrada
	if err := scanner.Err(); err != nil {
		// Si hubo un error al leer la entrada, lo imprime
		fmt.Println("Error al leer:", err)
	}
}
