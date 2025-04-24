package commands

import (
	stores "backend/stores"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// LOGIN estructura que representa el comando login con sus parámetros
type LOGIN struct {
	user string // Usuario
	pass string // Contraseña
	id   string // ID del disco
}

/*
	login -user=root -pass=123 -id=062A3E2D
*/

func ParseLogin(tokens []string) (string, error) {
	cmd := &LOGIN{} // Crea una nueva instancia de LOGIN

	// Unir tokens en una sola cadena y luego dividir por espacios, respetando las comillas
	args := strings.Join(tokens, " ")
	// Expresión regular para encontrar los parámetros del comando mkfs
	re := regexp.MustCompile(`-user=[^\s]+|-pass=[^\s]+|-id=[^\s]+`)
	// Encuentra todas las coincidencias de la expresión regular en la cadena de argumentos
	matches := re.FindAllString(args, -1)

	// Verificar que todos los tokens fueron reconocidos por la expresión regular
	if len(matches) != len(tokens) {
		// Identificar el parámetro inválido
		for _, token := range tokens {
			if !re.MatchString(token) {
				return "", fmt.Errorf("parámetro inválido: %s", token)
			}
		}
	}

	// Itera sobre cada coincidencia encontrada
	for _, match := range matches {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return "", fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar diferentes parámetros
		switch key {
		case "-user":
			if value == "" {
				return "", errors.New("el usuario no puede estar vacío")
			}
			cmd.user = value
		case "-pass":
			if value == "" {
				return "", errors.New("la contraseña no puede estar vacía")
			}
			cmd.pass = value
		case "-id":
			// Verifica que el id no esté vacío
			if value == "" {
				return "", errors.New("el id no puede estar vacío")
			}
			cmd.id = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return "", fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que el parámetro -id haya sido proporcionado
	if cmd.id == "" {
		return "", errors.New("faltan parámetros requeridos: -id")
	}

	// Si no se proporcionó el tipo, se establece por defecto a "full"
	if cmd.user == "" {
		return "", errors.New("faltan parámetros requeridos: -user")
	}

	// Si no se proporcionó el tipo, se establece por defecto a "full"
	if cmd.pass == "" {
		return "", errors.New("faltan parámetros requeridos: -pass")
	}

	// Aquí se puede agregar la lógica para ejecutar el comando mkfs con los parámetros proporcionados
	err := commandLogin(cmd)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("LOGIN: Usuario: %s, Contraseña: %s, ID: %s", cmd.user, cmd.pass, cmd.id), nil
}

func commandLogin(login *LOGIN) error {
	// Obtener la partición montada
	partitionSuperblock, _, partitionPath, err := stores.GetMountedPartitionSuperblock(login.id)
	if err != nil {
		return fmt.Errorf("error al obtener la partición montada: %w", err)
	}

	// Obtener el bloque de usuarios
	usersBlock, err := partitionSuperblock.GetUsersBlock(partitionPath)
	if err != nil {
		return fmt.Errorf("error al obtener el bloque de usuarios: %w", err)
	}

	fmt.Println(usersBlock)

	// Convertir el contenido del bloque a string y separar por líneas
	content := strings.Trim(string(usersBlock.B_content[:]), "\x00")
	lines := strings.Split(content, "\n")

	fmt.Println(content)

	// Variables para almacenar la información del usuario
	var foundUser bool
	var userPassword string

	// Buscar el usuario en las líneas
	for _, line := range lines {
		// Dividir la línea en campos
		fields := strings.Split(line, ",")
		// Limpiar espacios en blanco de cada campo
		for i := range fields {
			fields[i] = strings.TrimSpace(fields[i])
		}

		// Verificar si es una línea de usuario (tipo U)
		if len(fields) == 5 && fields[1] == "U" {
			// Comparar el nombre de usuario (campo 3)
			if strings.EqualFold(fields[3], login.user) {
				foundUser = true
				userPassword = fields[4]
				break
			}
		}
	}

	// Verificar si se encontró el usuario
	if !foundUser {
		return fmt.Errorf("el usuario %s no existe", login.user)
	}

	// Verificar la contraseña
	if !strings.EqualFold(userPassword, login.pass) {
		return fmt.Errorf("la contraseña no coincide")
	}

	// If validation succeeds, set the auth state
	stores.Auth.Login(login.user, login.pass, login.id)

	return nil
}
