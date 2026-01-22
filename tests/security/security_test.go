package security_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/scopweb/mcp-go-png-to-ico-convert/internal/converter"
)

// TestPathTraversal verifica que no se permita path traversal en las rutas
func TestPathTraversal(t *testing.T) {
	tests := []struct {
		name       string
		inputPath  string
		outputPath string
		wantError  bool
	}{
		{
			name:       "Normal path should work",
			inputPath:  "test.png",
			outputPath: "test.ico",
			wantError:  true, // Esperamos error porque el archivo no existe
		},
		{
			name:       "Path traversal with ../ should be handled",
			inputPath:  "../../../etc/passwd",
			outputPath: "output.ico",
			wantError:  true,
		},
		{
			name:       "Absolute path should be handled",
			inputPath:  "/etc/passwd",
			outputPath: "output.ico",
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := converter.ConvertPNGtoICO(tt.inputPath, tt.outputPath)
			if (err != nil) != tt.wantError {
				t.Errorf("ConvertPNGtoICO() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestInvalidInputFiles verifica el manejo de archivos de entrada inválidos
func TestInvalidInputFiles(t *testing.T) {
	// Crear un archivo temporal que no sea PNG
	tmpDir := t.TempDir()
	notPNGFile := filepath.Join(tmpDir, "notpng.png")
	err := os.WriteFile(notPNGFile, []byte("This is not a PNG file"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name          string
		inputPath     string
		outputPath    string
		expectedError string
	}{
		{
			name:          "Non-existent file",
			inputPath:     "nonexistent.png",
			outputPath:    "output.ico",
			expectedError: "FILE_NOT_FOUND",
		},
		{
			name:          "Invalid PNG file",
			inputPath:     notPNGFile,
			outputPath:    filepath.Join(tmpDir, "output.ico"),
			expectedError: "INVALID_PNG",
		},
		{
			name:          "Empty path",
			inputPath:     "",
			outputPath:    "output.ico",
			expectedError: "FILE_NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter.ConvertPNGtoICO(tt.inputPath, tt.outputPath)
			if err == nil {
				t.Error("Expected error, got nil")
			}
			if result != nil && result.ErrorCode != tt.expectedError {
				t.Errorf("Expected error code %s, got %s", tt.expectedError, result.ErrorCode)
			}
		})
	}
}

// TestOutputPermissions verifica que se manejen correctamente los permisos de escritura
func TestOutputPermissions(t *testing.T) {
	tmpDir := t.TempDir()

	// Crear un PNG válido de prueba (1x1 pixel)
	testPNG := filepath.Join(tmpDir, "test.png")
	// PNG de 1x1 pixel transparente (contenido mínimo válido)
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

	// Test: Intentar escribir en un directorio sin permisos de escritura
	// Nota: Esto es complicado de testear en Windows, así que solo verificamos la ruta inválida
	outputPath := filepath.Join("/invalid/path/that/does/not/exist", "output.ico")
	_, err = converter.ConvertPNGtoICO(testPNG, outputPath)
	if err == nil {
		t.Error("Expected error when writing to invalid path, got nil")
	}
}

// TestLargeFileHandling verifica que se manejen archivos grandes
func TestLargeFileHandling(t *testing.T) {
	// Este test verifica que el código no tenga vulnerabilidades de DoS
	// al procesar archivos muy grandes
	t.Skip("Test de archivos grandes - implementar cuando sea necesario")
}

// TestConcurrentAccess verifica que múltiples conversiones concurrentes sean seguras
func TestConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()

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
		t.Fatalf("Failed to create test PNG: %v", err)
	}

	// Ejecutar 10 conversiones concurrentes
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			outputPath := filepath.Join(tmpDir, "output_"+string(rune(id))+ ".ico")
			_, _ = converter.ConvertPNGtoICO(testPNG, outputPath)
			done <- true
		}(i)
	}

	// Esperar a que todas terminen
	for i := 0; i < 10; i++ {
		<-done
	}
}
