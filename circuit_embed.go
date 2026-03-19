package dsr

import (
	_ "embed"
	"fmt"

	framework "github.com/dpopsuev/origami"
)

//go:embed circuit.yaml
var defaultCircuitYAML []byte

// DefaultCircuitYAML returns the embedded base DSR circuit definition.
func DefaultCircuitYAML() []byte { return defaultCircuitYAML }

// SchematicResolver returns an AssetResolver that resolves "gnd"
// to the embedded base circuit.
func SchematicResolver() framework.AssetResolver {
	return func(name string) ([]byte, error) {
		if name == "gnd" {
			return defaultCircuitYAML, nil
		}
		return nil, fmt.Errorf("unknown schematic %q", name)
	}
}
