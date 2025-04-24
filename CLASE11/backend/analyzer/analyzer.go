package analyzer

import (
	commands "backend/commands"
	"errors"  // Importa el paquete "errors" para manejar errores
	"fmt"     // Importa el paquete "fmt" para formatear e imprimir texto
	"strings" // Importa el paquete "strings" para manipulación de cadenas
)

// Analyzer analiza el comando de entrada y ejecuta la acción correspondiente
func Analyzer(input string) (string, error) {
	tokens := strings.Fields(input)

	if len(tokens) == 0 {
		return "", errors.New("no se proporcionó ningún comando")
	}

	// Switch para manejar diferentes comandos
	switch tokens[0] {
	case "mkdisk":
		return commands.ParseMkdisk(tokens[1:])
	case "fdisk":
		return commands.ParseFdisk(tokens[1:])
	case "mount":
		return commands.ParseMount(tokens[1:])
	case "mkfs":
		return commands.ParseMkfs(tokens[1:])
	case "login":
		return commands.ParseLogin(tokens[1:])
	case "mkdir":
		return commands.ParseMkdir(tokens[1:])
	case "rep":
		return commands.ParseRep(tokens[1:])
	default:
		return "", fmt.Errorf("comando desconocido: %s", tokens[0])
	}
}
