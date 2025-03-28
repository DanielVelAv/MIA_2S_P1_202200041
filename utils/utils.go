package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ConvertirTamanioBy(size int, unit string) (int, error) {
	switch unit {
	case "K":
		return size * 1024, nil // Convierte kilobytes a bytes
	case "M":
		return size * 1024 * 1024, nil // Convierte megabytes a bytes
	case "B":
		return size, nil //ya esta en bytes
	default:
		return 0, errors.New("invalid unit") // Devuelve un error si la unidad es inválida
	}
}

var alfabeto = []string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

// mapa para almacenar la asignacion de letras a los paths
var pathToLetter = make(map[string]string)

// indice para siguiente letra disponible en el abecedario
var nextLetterIndex = 0

func GetLetter(path string) (string, error) {
	if _, existe := pathToLetter[path]; !existe {
		if nextLetterIndex < len(alfabeto) {
			pathToLetter[path] = alfabeto[nextLetterIndex]
			nextLetterIndex++
		} else {
			fmt.Println("No hay mas letras para asignar")
			return "", errors.New("no hay mas letras para usar")
		}

	}
	return pathToLetter[path], nil
}

func CrearDirectorios(path string) error {
	directorio := filepath.Dir(path)
	err := os.MkdirAll(directorio, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating director: %v", err)
	}
	return nil
}

func ObtenerNombreArchivo(path string) (string, string) {
	directorio := filepath.Dir(path)
	Nbase := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	nombreDot := filepath.Join(directorio, Nbase+".dot")
	ImagenResultante := path
	return nombreDot, ImagenResultante
}
