package main

import (
	"fmt" // Importa el paquete fmt para entrada/salida
	"os"  // Importa el paquete os para funciones del sistema operativo

	operations "basics/operations" // Importa nuestro paquete de operaciones matemáticas
)

func main() {
	// Bucle infinito para mantener el programa ejecutándose
	for {
		// Muestra el menú principal
		fmt.Println("\n=== Calculadora Aritmética ===")
		fmt.Println("1. Suma")
		fmt.Println("2. Resta")
		fmt.Println("3. Multiplicación")
		fmt.Println("4. División")
		fmt.Println("5. Potencia")
		fmt.Println("6. Salir")

		// Declara variable para almacenar la opción del usuario
		var choice int
		fmt.Print("Seleccione una operación (1-6): ")
		fmt.Scan(&choice) // Lee la opción del usuario

		// Verifica si el usuario quiere salir
		if choice == 6 {
			fmt.Println("¡Hasta luego!")
			os.Exit(0) // Termina el programa con código 0 (éxito)
		}

		// Declara variables para los operandos
		var a, b float64
		fmt.Print("Ingrese primer número: ")
		fmt.Scan(&a) // Lee el primer número
		fmt.Print("Ingrese segundo número: ")
		fmt.Scan(&b) // Lee el segundo número

		// Estructura switch para ejecutar la operación seleccionada
		switch choice {
		case 1:
			fmt.Printf("Resultado: %.2f\n", operations.Add(a, b)) // Suma
		case 2:
			fmt.Printf("Resultado: %.2f\n", operations.Subtract(a, b)) // Resta
		case 3:
			fmt.Printf("Resultado: %.2f\n", operations.Multiply(a, b)) // Multiplicación
		case 4:
			// Verifica división por cero
			if b == 0 {
				fmt.Println("Error: No se puede dividir entre cero")
				continue // Vuelve al inicio del bucle
			}
			fmt.Printf("Resultado: %.2f\n", operations.Divide(a, b)) // División
		case 5:
			fmt.Printf("Resultado: %.2f\n", operations.Power(a, b)) // Potencia
		default:
			fmt.Println("Opción inválida") // Manejo de opciones no válidas
		}
	}
}
