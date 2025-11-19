package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadJSON lee un archivo JSON y lo deserializa en la variable destino (struct o map).
func ReadJSON(filename string, dest interface{}) error {
	file, err := os.Open(filename + ".json")

	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error al abrir el archivo %s: %w", filename, err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(dest); err != nil {
		return fmt.Errorf("error al decodificar JSON %s: %w", filename, err)
	}

	return nil
}

// WriteJSON guarda una estructura o mapa en un archivo JSON, formateado bonito.
func WriteJSON(filename string, data interface{}) error {
	file, err := os.Create(filename + ".json")
	if err != nil {
		return fmt.Errorf("error al crear el archivo %s: %w", filename, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("error al codificar JSON: %w", err)
	}

	return nil
}

// UpdateJSON actualiza un campo específico dentro de un JSON (sin perder el resto de los datos).
func UpdateJSON(filename string, key string, newValue interface{}) error {
	var data map[string]interface{}

	// Si el archivo existe, leerlo
	if _, err := os.Stat(filename + ".json"); err == nil {
		if err := ReadJSON(filename, &data); err != nil {
			return err
		}
	} else {
		data = make(map[string]interface{})
	}

	// Actualizar campo
	data[key] = newValue

	return WriteJSON(filename, data)
}
