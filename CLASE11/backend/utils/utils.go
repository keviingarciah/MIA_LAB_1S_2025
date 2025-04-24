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

// GetParentDirectories obtiene las carpetas padres y el directorio de destino
func GetParentDirectories(path string) ([]string, string) {
	// Normalizar el path
	path = filepath.Clean(path)

	// Dividir el path en sus componentes
	components := strings.Split(path, string(filepath.Separator))

	// Lista para almacenar las rutas de las carpetas padres
	var parentDirs []string

	// Construir las rutas de las carpetas padres, excluyendo la última carpeta
	for i := 1; i < len(components)-1; i++ {
		parentDirs = append(parentDirs, components[i])
	}

	// La última carpeta es la carpeta de destino
	destDir := components[len(components)-1]

	return parentDirs, destDir
}

// First devuelve el primer elemento de un slice
func First[T any](slice []T) (T, error) {
	if len(slice) == 0 {
		var zero T
		return zero, errors.New("el slice está vacío")
	}
	return slice[0], nil
}

// RemoveElement elimina un elemento de un slice en el índice dado
func RemoveElement[T any](slice []T, index int) []T {
	if index < 0 || index >= len(slice) {
		return slice // Índice fuera de rango, devolver el slice original
	}
	return append(slice[:index], slice[index+1:]...)
}

// splitStringIntoChunks divide una cadena en partes de tamaño chunkSize y las almacena en una lista
func SplitStringIntoChunks(s string) []string {
	var chunks []string
	for i := 0; i < len(s); i += 64 {
		end := i + 64
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}
