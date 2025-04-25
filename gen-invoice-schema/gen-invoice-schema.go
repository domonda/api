package main

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"

	"github.com/domonda/api/golang/domonda"
	"github.com/invopop/jsonschema"
)

func main() {
	// ------------------------------------------------------------------------
	// Configure the jsonschema.Reflector
	// ------------------------------------------------------------------------

	// reflector.AddGoComments needs to be called from the directory of the reflected package
	err := os.Chdir("../golang/domonda")
	if err != nil {
		log.Fatalf("Failed to change directory to golang/domonda: %v", err)
	}
	reflector := &jsonschema.Reflector{
		Namer: func(t reflect.Type) string {
			return path.Base(t.PkgPath()) + "." + t.Name()
		},
		ExpandedStruct: true,
		DoNotReference: true,
	}
	err = reflector.AddGoComments("github.com/domonda/api/golang/domonda", ".")
	if err != nil {
		log.Fatalf("Failed to parse Go comments: %v", err)
	}

	// ------------------------------------------------------------------------
	// Generate the schema
	// ------------------------------------------------------------------------

	schema := reflector.Reflect(domonda.Invoice{})
	schema.ID = "https://raw.githubusercontent.com/domonda/api/refs/heads/master/invoice.schema.json"
	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		log.Fatalf("Failed to generate schema: %v", err)
	}

	// ------------------------------------------------------------------------
	// Write the schema to a file
	// ------------------------------------------------------------------------

	filePath, err := filepath.Abs("../../invoice.schema.json")
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}
	err = os.WriteFile(filePath, schemaJSON, 0644)
	if err != nil {
		log.Fatalf("Failed to write schema: %v", err)
	}
	log.Println("Invoice schema written to", filePath)
}
