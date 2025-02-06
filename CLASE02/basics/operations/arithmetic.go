package operations

// Add realiza la suma de dos números de tipo float64
// Es una función exportada que puede ser usada por otros paquetes
func Add(a, b float64) float64 {
	return a + b
}

// Subtract realiza la resta entre dos números float64
// Es una función exportada que puede ser usada por otros paquetes
func Subtract(a, b float64) float64 {
	return a - b
}

// Multiply realiza la multiplicación entre dos números float64
// Es una función exportada que puede ser usada por otros paquetes
func Multiply(a, b float64) float64 {
	return a * b
}

// Divide realiza la división entre dos números float64
// Utiliza la función auxiliar isZero para validar división por cero
// Es una función exportada que puede ser usada por otros paquetes
func Divide(a, b float64) float64 {
	if isZero(b) {
		return 0 // Evita división por cero
	}
	return a / b
}

// Power calcula la potencia usando la función auxiliar calculate
// Es una función exportada que puede ser usada por otros paquetes
func Power(base, exponent float64) float64 {
	return calculate(base, exponent)
}

// isZero es una función auxiliar (no exportada)
// Se usa para validar división por cero en la función Divide
func isZero(num float64) bool {
	return num == 0
}

// calculate es una función auxiliar (no exportada)
// Implementa el cálculo de potencia mediante multiplicación iterativa
// Convierte el exponente a int para poder iterar
func calculate(base, exponent float64) float64 {
	result := 1.0
	for i := 0; i < int(exponent); i++ {
		result *= base
	}
	return result
}
