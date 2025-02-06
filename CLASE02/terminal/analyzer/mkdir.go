package analyzer

import (
	"errors"
	"fmt"
	"strings"
)

// MKDIR representa la estructura para el comando mkdir
// Contiene los parámetros necesarios para crear directorios
type MKDIR struct {
	path string // Ruta donde se creará el directorio
	p    bool   // Flag para crear directorios padres si no existen (-p)
}

// Errores comunes
var (
	ErrEmptyPath     = errors.New("la ruta no puede estar vacía")
	ErrPathRequired  = errors.New("el parámetro -path es obligatorio")
	ErrInvalidFormat = "formato inválido: %s (debe ser -param=valor)"
	ErrUnknownParam  = "parámetro desconocido: %s"
	ErrSpacesInPath  = "la ruta sin comillas no puede contener espacios: %s"
)

/*
   Ejemplos de uso del comando mkdir:
   mkdir -path=/home/user/docs
   mkdir -p -path="/home/user/docs 1"
   mkdir -path=/home/user/docs -p
*/

// validParam verifica si un parámetro es válido para el comando mkdir
// param: string con el nombre del parámetro a validar
// retorna: bool indicando si es válido
func validParam(param string) bool {
	validParams := map[string]bool{
		"-path": true,
		"-p":    true,
	}
	return validParams[strings.ToLower(param)]
}

// handleQuotedPath procesa una ruta que contiene comillas
// Retorna la ruta procesada y si se completó el procesamiento
func handleQuotedPath(token string, quotedPath string, inQuotes bool) (string, bool, bool) {
	// Si estamos procesando una ruta con comillas
	if inQuotes {
		if strings.HasSuffix(token, "\"") {
			return quotedPath + " " + strings.TrimSuffix(token, "\""), true, false
		}
		return quotedPath + " " + token, false, true
	}

	// Inicio de una ruta con comillas
	if strings.HasPrefix(token, "-path=\"") {
		path := strings.TrimPrefix(token, "-path=\"")
		if strings.HasSuffix(path, "\"") {
			return strings.TrimSuffix(path, "\""), false, false
		}
		return path, false, true
	}

	return quotedPath, false, false
}

// ParseMkdir analiza los tokens del comando y retorna una estructura MKDIR
// tokens: slice de strings con los parámetros del comando
// Retorna: puntero a MKDIR y error si hay problemas en el parsing
func ParseMkdir(tokens []string) (*MKDIR, error) {
	cmd := &MKDIR{}
	var quotedPath string
	inQuotes := false

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		// Manejo de flag -p
		if strings.EqualFold(token, "-p") {
			cmd.p = true
			continue
		}

		// Manejo de rutas con comillas
		if inQuotes || strings.HasPrefix(token, "-path=\"") {
			processed, done, continuing := handleQuotedPath(token, quotedPath, inQuotes)
			if continuing {
				quotedPath = processed
				inQuotes = true
				continue
			}
			if done {
				cmd.path = processed
				continue
			}
		}

		// Manejo normal de parámetros sin comillas
		parts := strings.SplitN(token, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf(ErrInvalidFormat, token)
		}

		param := strings.ToLower(parts[0])
		value := parts[1]

		if !validParam(param) {
			return nil, fmt.Errorf(ErrUnknownParam, param)
		}

		switch param {
		case "-path":
			if value == "" {
				return nil, ErrEmptyPath
			}
			if strings.Contains(value, " ") {
				return nil, fmt.Errorf(ErrSpacesInPath, value)
			}
			cmd.path = value
		}
	}

	if cmd.path == "" {
		return nil, ErrPathRequired
	}

	// Ejecutar el comando mkdir

	return cmd, nil
}
