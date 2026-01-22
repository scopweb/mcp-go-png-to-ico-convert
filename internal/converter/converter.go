package converter

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
	"github.com/sergeymakinen/go-ico"
)

// Resoluciones requeridas para el ICO (de mayor a menor)
var requiredResolutions = []int{256, 128, 64, 48, 32, 16}

// ConversionResult contiene el resultado de la conversión
type ConversionResult struct {
	Success     bool     `json:"success"`
	OutputFile  string   `json:"output_file,omitempty"`
	Resolutions []string `json:"resolutions,omitempty"`
	FileSize    string   `json:"file_size,omitempty"`
	Message     string   `json:"message,omitempty"`
	Error       string   `json:"error,omitempty"`
	ErrorCode   string   `json:"error_code,omitempty"`
}

// ConvertPNGtoICO convierte un archivo PNG a ICO con múltiples resoluciones
func ConvertPNGtoICO(inputPath, outputPath string) (*ConversionResult, error) {
	// Validar que el archivo PNG existe
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return &ConversionResult{
			Success:   false,
			Error:     fmt.Sprintf("El archivo PNG no existe: %s", inputPath),
			ErrorCode: "FILE_NOT_FOUND",
		}, fmt.Errorf("archivo no encontrado: %s", inputPath)
	}

	// Abrir el archivo PNG
	file, err := os.Open(inputPath)
	if err != nil {
		return &ConversionResult{
			Success:   false,
			Error:     fmt.Sprintf("No se puede abrir el archivo: %v", err),
			ErrorCode: "FILE_NOT_FOUND",
		}, err
	}
	defer file.Close()

	// Decodificar el PNG
	img, err := png.Decode(file)
	if err != nil {
		return &ConversionResult{
			Success:   false,
			Error:     fmt.Sprintf("No es un archivo PNG válido: %v", err),
			ErrorCode: "INVALID_PNG",
		}, err
	}

	// Generar las imágenes redimensionadas para cada resolución
	images := make([]image.Image, 0, len(requiredResolutions))
	resolutionStrings := make([]string, 0, len(requiredResolutions))

	for _, size := range requiredResolutions {
		// Redimensionar usando filtro Lanczos (alta calidad)
		// imaging.Lanczos es el mejor filtro para downsampling
		resized := imaging.Resize(img, size, size, imaging.Lanczos)
		images = append(images, resized)
		resolutionStrings = append(resolutionStrings, fmt.Sprintf("%dx%d", size, size))
	}

	// Crear el archivo ICO de salida
	outFile, err := os.Create(outputPath)
	if err != nil {
		return &ConversionResult{
			Success:   false,
			Error:     fmt.Sprintf("No se puede crear el archivo de salida: %v", err),
			ErrorCode: "ENCODING_ERROR",
		}, err
	}
	defer outFile.Close()

	// Codificar todas las imágenes en un solo archivo ICO
	err = ico.EncodeAll(outFile, images)
	if err != nil {
		return &ConversionResult{
			Success:   false,
			Error:     fmt.Sprintf("Error al codificar el archivo ICO: %v", err),
			ErrorCode: "ENCODING_ERROR",
		}, err
	}

	// Obtener el tamaño del archivo resultante
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		return &ConversionResult{
			Success:   false,
			Error:     fmt.Sprintf("Error al obtener información del archivo: %v", err),
			ErrorCode: "ENCODING_ERROR",
		}, err
	}

	// Formatear el tamaño del archivo
	fileSize := formatFileSize(fileInfo.Size())

	// Retornar resultado exitoso
	return &ConversionResult{
		Success:     true,
		OutputFile:  outputPath,
		Resolutions: resolutionStrings,
		FileSize:    fileSize,
		Message:     "ICO generado exitosamente",
	}, nil
}

// formatFileSize formatea el tamaño del archivo a una cadena legible
func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
