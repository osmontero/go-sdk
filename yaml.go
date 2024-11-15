package go_sdk

import (
	"os"

	"github.com/threatwinds/logger"
	"gopkg.in/yaml.v3"
	k8syaml "sigs.k8s.io/yaml"
)

// ReadPbYaml reads a YAML file, converts its content to JSON, and returns the JSON bytes.
// If an error occurs while reading the file or converting its content, it returns a logger.Error.
//
// Parameters:
//   - f: The file path of the YAML file to be read.
//
// Returns:
//   - []byte: The JSON bytes converted from the YAML file.
//   - *logger.Error: An error object if an error occurs, otherwise nil.
func ReadPbYaml(f string) ([]byte, *logger.Error) {
	content, err := os.ReadFile(f)
	if err != nil {
		return nil, Logger().ErrorF("error opening file '%s': %s", f, err.Error())
	}

	bytes, err := k8syaml.YAMLToJSON(content)
	if err != nil {
		return nil, Logger().ErrorF("error converting YAML file '%s' to JSON: %s", f, err.Error())
	}

	return bytes, nil
}

// ReadYaml reads a YAML file and unmarshals its content into a specified type.
// The function can also handle JSON mode if specified.
//
// Type Parameters:
//   t: The type into which the YAML content will be unmarshaled.
//
// Parameters:
//   f: The file path to the YAML file.
//   jsonMode: A boolean flag indicating whether to use JSON mode for unmarshaling.
//
// Returns:
//   *t: A pointer to the unmarshaled content of type t.
//   *logger.Error: A pointer to an error object if an error occurs, otherwise nil.
func ReadYaml[t any](f string, jsonMode bool) (*t, *logger.Error) {
	content, err := os.ReadFile(f)
	if err != nil {
		return nil, Logger().ErrorF("error opening file '%s': %s", f, err.Error())
	}

	var value = new(t)
	if jsonMode {
		err = k8syaml.Unmarshal(content, value)
		if err != nil {
			return nil, Logger().ErrorF("error decoding YAML file '%s': %s", f, err.Error())
		}
	} else {
		err = yaml.Unmarshal(content, value)
		if err != nil {
			return nil, Logger().ErrorF("error decoding YAML file '%s': %s", f, err.Error())
		}
	}

	return value, nil
}
