package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/yourusername/mcp-go-png-convert/internal/converter"
	"github.com/yourusername/mcp-go-png-convert/internal/mcp"
)

func main() {
	// Modo CLI: Si hay exactamente 2 argumentos (input y output)
	if len(os.Args) == 3 {
		inputPath := os.Args[1]
		outputPath := os.Args[2]
		runCLIMode(inputPath, outputPath)
		return
	}

	// Modo CLI: Si hay 1 argumento (solo input, output = input.ico)
	if len(os.Args) == 2 {
		inputPath := os.Args[1]
		// Generar nombre de salida automáticamente
		outputPath := inputPath[:len(inputPath)-4] + ".ico"
		runCLIMode(inputPath, outputPath)
		return
	}

	// Sin argumentos → Modo servidor MCP
	if err := mcp.StartServer(); err != nil {
		log.Fatalf("Error al iniciar el servidor MCP: %v", err)
		os.Exit(1)
	}
}

// runCLIMode ejecuta la conversión en modo línea de comandos
func runCLIMode(inputPath, outputPath string) {
	fmt.Printf("🔄 Convirtiendo: %s -> %s\n", inputPath, outputPath)

	// Ejecutar la conversión
	result, err := converter.ConvertPNGtoICO(inputPath, outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Mostrar resultado exitoso
	if result.Success {
		fmt.Printf("✅ %s\n", result.Message)
		fmt.Printf("📦 Archivo: %s\n", result.OutputFile)
		fmt.Printf("📐 Resoluciones: %v\n", result.Resolutions)
		fmt.Printf("💾 Tamaño: %s\n", result.FileSize)
		os.Exit(0)
	} else {
		// Mostrar error en formato JSON
		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		fmt.Fprintf(os.Stderr, "❌ Error en conversión:\n%s\n", string(resultJSON))
		os.Exit(1)
	}
}
