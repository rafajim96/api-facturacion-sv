package main

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/bootstrap"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"os"
	"fmt"
)

func main() {
	// Crear e inicializar la aplicación
	app := bootstrap.NewApplication()
	if err := app.Initialize(); err != nil {
		// SAFEGUARD: check if logger is nil
		if logs.Logger != nil {
			logs.Fatal("Failed to initialize application", map[string]interface{}{"error": err.Error()})
		} else {
			// Fallback to fmt if logger isn't ready
			fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
		}
		os.Exit(1)
	}

	// Iniciar la aplicación
	if err := app.Start(); err != nil {
		logs.Fatal("Application error", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}
}
