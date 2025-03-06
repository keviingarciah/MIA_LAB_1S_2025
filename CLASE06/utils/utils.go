package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ConvertToBytes convierte un tamaño y una unidad a bytes
func ConvertToBytes(size int, unit string) (int, error) {
	switch unit {
	case "K":
		return size * 1024, nil // Convierte kilobytes a bytes
	case "M":
		return size * 1024 * 1024, nil // Convierte megabytes a bytes
	default:
		return 0, errors.New("invalid unit") // Devuelve un error si la unidad es inválida
	}
}

// Lista con todo el abecedario
var alphabet = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

// Mapa para almacenar la asignación de letras a los diferentes paths
var pathToLetter = make(map[string]string)

// Mapa para almacenar el contador de particiones por path
var pathToPartitionCount = make(map[string]int)

// Índice para la siguiente letra disponible en el abecedario
var nextLetterIndex = 0

// GetLetter obtiene la letra asignada a un path y el siguiente índice de partición
func GetLetterAndPartitionCorrelative(path string) (string, int, error) {
	// Asignar una letra al path si no tiene una asignada
	if _, exists := pathToLetter[path]; !exists {
		if nextLetterIndex < len(alphabet) {
			pathToLetter[path] = alphabet[nextLetterIndex]
			pathToPartitionCount[path] = 0 // Inicializar el contador de particiones
			nextLetterIndex++
		} else {
			fmt.Println("Error: no hay más letras disponibles para asignar")
			return "", 0, errors.New("no hay más letras disponibles para asignar")
		}
	}

	// Incrementar y obtener el siguiente índice de partición para este path
	pathToPartitionCount[path]++
	nextIndex := pathToPartitionCount[path]

	return pathToLetter[path], nextIndex, nil
}

// createParentDirs crea las carpetas padre si no existen
func CreateParentDirs(path string) error {
	dir := filepath.Dir(path)
	// os.MkdirAll no sobrescribe las carpetas existentes, solo crea las que no existen
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error al crear las carpetas padre: %v", err)
	}
	return nil
}

// getFileNames obtiene el nombre del archivo .dot y el nombre de la imagen de salida
func GetFileNames(path string) (string, string) {
	dir := filepath.Dir(path)
	baseName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	dotFileName := filepath.Join(dir, baseName+".dot")
	outputImage := path
	return dotFileName, outputImage
}
