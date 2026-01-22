package converter

import (
	"os"
	"path/filepath"
	"testing"
)

// TestConvertPNGtoICO_ValidInput verifica la conversión con un PNG válido
func TestConvertPNGtoICO_ValidInput(t *testing.T) {
	tmpDir := t.TempDir()

	// Crear un PNG válido de prueba (1x1 pixel RGBA)
	testPNG := filepath.Join(tmpDir, "test.png")
	pngData := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89, 0x00, 0x00, 0x00,
		0x0a, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00, 0x00, 0x00, 0x49,
		0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
	}
	err := os.WriteFile(testPNG, pngData, 0644)
	if err != nil {
		t.Fatalf("Failed to create test PNG: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "output.ico")

	// Ejecutar la conversión
	result, err := ConvertPNGtoICO(testPNG, outputPath)

	// Verificaciones
	if err != nil {
		t.Fatalf("ConvertPNGtoICO() failed: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}

	if !result.Success {
		t.Errorf("Expected Success=true, got false. Error: %s", result.Error)
	}

	if result.OutputFile != outputPath {
		t.Errorf("Expected OutputFile=%s, got %s", outputPath, result.OutputFile)
	}

	// Verificar que el archivo ICO fue creado
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output ICO file was not created")
	}

	// Verificar que se crearon todas las resoluciones
	expectedResolutions := []string{"256x256", "128x128", "64x64", "48x48", "32x32", "16x16"}
	if len(result.Resolutions) != len(expectedResolutions) {
		t.Errorf("Expected %d resolutions, got %d", len(expectedResolutions), len(result.Resolutions))
	}

	for i, expected := range expectedResolutions {
		if i >= len(result.Resolutions) || result.Resolutions[i] != expected {
			t.Errorf("Expected resolution %s at index %d, got %s", expected, i, result.Resolutions[i])
		}
	}
}

// TestConvertPNGtoICO_FileNotFound verifica el manejo de archivos no encontrados
func TestConvertPNGtoICO_FileNotFound(t *testing.T) {
	result, err := ConvertPNGtoICO("nonexistent.png", "output.ico")

	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	if result == nil {
		t.Fatal("Result should not be nil")
	}

	if result.Success {
		t.Error("Expected Success=false for non-existent file")
	}

	if result.ErrorCode != "FILE_NOT_FOUND" {
		t.Errorf("Expected ErrorCode='FILE_NOT_FOUND', got '%s'", result.ErrorCode)
	}
}

// TestConvertPNGtoICO_InvalidPNG verifica el manejo de archivos PNG inválidos
func TestConvertPNGtoICO_InvalidPNG(t *testing.T) {
	tmpDir := t.TempDir()

	// Crear un archivo que no es PNG
	invalidFile := filepath.Join(tmpDir, "invalid.png")
	err := os.WriteFile(invalidFile, []byte("This is not a PNG"), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "output.ico")
	result, err := ConvertPNGtoICO(invalidFile, outputPath)

	if err == nil {
		t.Error("Expected error for invalid PNG, got nil")
	}

	if result == nil {
		t.Fatal("Result should not be nil")
	}

	if result.Success {
		t.Error("Expected Success=false for invalid PNG")
	}

	if result.ErrorCode != "INVALID_PNG" {
		t.Errorf("Expected ErrorCode='INVALID_PNG', got '%s'", result.ErrorCode)
	}
}

// TestFormatFileSize verifica el formato correcto del tamaño de archivo
func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"Less than 1KB", 500, "500 B"},
		{"Exactly 1KB", 1024, "1.0 KB"},
		{"1.5KB", 1536, "1.5 KB"},
		{"1MB", 1048576, "1.0 MB"},
		{"1GB", 1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFileSize(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatFileSize(%d) = %s, want %s", tt.bytes, result, tt.expected)
			}
		})
	}
}

// TestConversionResult_Structure verifica la estructura de ConversionResult
func TestConversionResult_Structure(t *testing.T) {
	result := &ConversionResult{
		Success:     true,
		OutputFile:  "/path/to/file.ico",
		Resolutions: []string{"256x256", "128x128"},
		FileSize:    "1.5 KB",
		Message:     "Success",
		Error:       "",
		ErrorCode:   "",
	}

	if !result.Success {
		t.Error("Expected Success to be true")
	}

	if result.OutputFile != "/path/to/file.ico" {
		t.Error("OutputFile not set correctly")
	}

	if len(result.Resolutions) != 2 {
		t.Error("Resolutions not set correctly")
	}
}

// Benchmark para medir el rendimiento de la conversión
func BenchmarkConvertPNGtoICO(b *testing.B) {
	tmpDir := b.TempDir()

	// Crear un PNG válido de prueba
	testPNG := filepath.Join(tmpDir, "test.png")
	pngData := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89, 0x00, 0x00, 0x00,
		0x0a, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00, 0x00, 0x00, 0x00, 0x49,
		0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
	}
	err := os.WriteFile(testPNG, pngData, 0644)
	if err != nil {
		b.Fatalf("Failed to create test PNG: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outputPath := filepath.Join(tmpDir, "bench_output.ico")
		_, _ = ConvertPNGtoICO(testPNG, outputPath)
	}
}
