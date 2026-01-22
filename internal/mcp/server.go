package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/scopweb/mcp-go-png-to-ico-convert/internal/converter"
)

// ConvertPNGtoICOArgs define los parámetros de entrada para la herramienta
type ConvertPNGtoICOArgs struct {
	InputPath  string `json:"input_path" jsonschema:"Ruta del archivo PNG de entrada (absoluta o relativa)"`
	OutputPath string `json:"output_path,omitempty" jsonschema:"Ruta del archivo ICO de salida (por defecto: mismo nombre con extensión .ico)"`
}

// StartServer inicia el servidor MCP
func StartServer() error {
	// Logging a stderr (para no interferir con stdio del protocolo MCP)
	log.SetOutput(os.Stderr)
	log.Println("Servidor MCP PNG-to-ICO iniciado")

	// Crear el servidor MCP con información del servidor
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "png-to-ico",
		Version: "1.0.0",
	}, nil)

	// Registrar la herramienta convert_png_to_ico
	mcp.AddTool(server, &mcp.Tool{
		Name:        "convert_png_to_ico",
		Description: "Convierte un archivo PNG a ICO con múltiples resoluciones para aplicaciones .NET",
	}, handleConvertPNGtoICO)

	// Ejecutar el servidor en modo stdio (estándar para MCP)
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		return fmt.Errorf("error al ejecutar el servidor MCP: %w", err)
	}

	return nil
}

// handleConvertPNGtoICO maneja las llamadas a la herramienta convert_png_to_ico
func handleConvertPNGtoICO(
	ctx context.Context,
	req *mcp.CallToolRequest,
	args ConvertPNGtoICOArgs,
) (*mcp.CallToolResult, any, error) {
	// Logging de la solicitud
	log.Printf("Recibida solicitud de conversión: input=%s, output=%s", args.InputPath, args.OutputPath)

	// Si no se especifica output_path, usar el mismo nombre con extensión .ico
	outputPath := args.OutputPath
	if outputPath == "" {
		ext := filepath.Ext(args.InputPath)
		outputPath = args.InputPath[:len(args.InputPath)-len(ext)] + ".ico"
	}

	// Convertir el PNG a ICO
	result, err := converter.ConvertPNGtoICO(args.InputPath, outputPath)
	if err != nil {
		log.Printf("Error en la conversión: %v", err)
		// Retornar error con contenido JSON
		errorJSON, _ := json.Marshal(result)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(errorJSON)},
			},
			IsError: true,
		}, nil, nil
	}

	log.Printf("Conversión exitosa: %s -> %s", args.InputPath, result.OutputFile)

	// Serializar el resultado a JSON
	resultJSON, _ := json.Marshal(result)
	log.Printf("Resultado: %s", string(resultJSON))

	// Retornar el resultado exitoso
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(resultJSON)},
		},
		IsError: false,
	}, nil, nil
}
