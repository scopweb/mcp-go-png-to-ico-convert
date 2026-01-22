# Conversión de PNG a ICO

Se trata de una aplicación en Go para convertir un archivo PNG y pasarlo a ICO.
Este ICO es para insertarlo en aplicaciones de consola de .NET 10.
Quiero que sea de calidad.

## Requisitos funcionales

### Entrada
- Archivo PNG (app.png)
- Soportar ruta relativa o absoluta

### Salida
- Archivo ICO con múltiples resoluciones incrustadas en un solo archivo

### Resoluciones requeridas
El ICO debe contener las siguientes resoluciones:
- 256x256
- 128x128
- 64x64
- 48x48
- 32x32
- 16x16

### Características de calidad
- Redimensionamiento de alta calidad (usar resampling LANCZOS o equivalente)
- Mantener transparencia alfa (RGBA)
- Generar un ICO estándar válido compatible con Windows y .NET
- Archivo resultante optimizado en tamaño (compresión adecuada)

### Uso esperado
```bash
go run main.go app.png app.ico
# o
./mcp-png-convert app.png app.ico
```

### Validación
- El archivo ICO resultante debe ser válido y reconocido por Windows
- Debe funcionar correctamente como icono de aplicación en .NET 10

## Librerías Go recomendadas

### Librerías principales (Tier 1)

#### 1. **github.com/sergeymakinen/go-ico** ⭐ RECOMENDADA
- Propósito: Encoder/decoder de ICO nativo en Go
- Funcionalidad: Codificar múltiples imágenes en un archivo ICO
- Ventajas:
  - Sin dependencias externas
  - Soporte completo para ICO estándar
  - Función `ico.EncodeAll()` perfecta para múltiples resoluciones
  - Mantenida y confiable
- Desventajas: Menos optimizaciones de calidad

#### 2. **github.com/donnie4w/ico**
- Propósito: Librería de encoding para formato ICO en Go
- Funcionalidad: Codificar imágenes a formato ICO
- Ventajas:
  - Implementación pura en Go
  - Sin dependencias externas
  - Bien documentada

#### 3. **image/png + image** (Go built-in)
- Propósito: Decodificar PNG y procesamiento básico de imágenes
- Funcionalidad: Leer PNG y manipular píxeles
- Ventajas:
  - Librería estándar de Go (incluida)
  - Sin dependencias
  - Soporte completo para RGBA

### Librerías de procesamiento de imagen (Tier 2)

#### 4. **github.com/disintegration/imaging** ⭐ RECOMENDADA PARA RESAMPLING
- Propósito: Procesamiento de imágenes de alta calidad
- Funcionalidad: Redimensionamiento con filtros LANCZOS, Gaussian, etc.
- Ventajas:
  - Excelente calidad de redimensionamiento
  - Soporte para múltiples filtros de resampling
  - Sin dependencias CGO
  - Pura Go
- Perfecto para: Redimensionar PNG a múltiples resoluciones sin pérdida de calidad

#### 5. **github.com/h2non/bimg**
- Propósito: Procesamiento de imágenes rápido y de alto nivel
- Funcionalidad: Crop, resize, rotate, watermark usando libvips
- Ventajas:
  - 4x más rápido que Go image package
  - Muy optimizado
- Desventajas:
  - Requiere libvips instalada (dependencia C)
  - Más complejo de configurar

### Solución recomendada

**Stack principal:**
```
1. github.com/sergeymakinen/go-ico (encoding ICO)
2. github.com/disintegration/imaging (resampling LANCZOS)
3. image/png (decodificación PNG built-in)
```

**Instalación:**
```bash
go get github.com/sergeymakinen/go-ico
go get github.com/disintegration/imaging
```

### Flujo de implementación

1. Leer PNG con `image/png.Decode()`
2. Redimensionar a cada resolución con `imaging.Resize()` usando filter LANCZOS
3. Codificar múltiples imágenes redimensionadas con `ico.EncodeAll()`
4. Guardar archivo ICO resultante

## Requisitos MCP (Claude Code + Claude Desktop)

### Protocolo MCP
Este proyecto debe implementar el **Model Context Protocol (MCP)** para funcionar como servidor:

#### 1. **Estructura de herramientas MCP**
El servidor MCP debe exponer una herramienta llamada `convert_png_to_ico` con los siguientes parámetros:

```json
{
  "name": "convert_png_to_ico",
  "description": "Convierte un archivo PNG a ICO con múltiples resoluciones para aplicaciones .NET",
  "inputSchema": {
    "type": "object",
    "properties": {
      "input_path": {
        "type": "string",
        "description": "Ruta del archivo PNG de entrada (absoluta o relativa)"
      },
      "output_path": {
        "type": "string",
        "description": "Ruta del archivo ICO de salida (por defecto: mismo nombre con extensión .ico)"
      }
    },
    "required": ["input_path"]
  }
}
```

#### 2. **Respuesta estándar**
El servidor debe retornar respuestas JSON estructuradas:

```json
{
  "success": true,
  "output_file": "/ruta/al/archivo.ico",
  "resolutions": ["256x256", "128x128", "64x64", "48x48", "32x32", "16x16"],
  "file_size": "24KB",
  "message": "ICO generado exitosamente"
}
```

O en caso de error:

```json
{
  "success": false,
  "error": "Descripción del error",
  "error_code": "FILE_NOT_FOUND|INVALID_PNG|ENCODING_ERROR"
}
```

#### 3. **Implementación técnica del servidor MCP**
- Usar `github.com/anthropics/mcp-go` (SDK oficial de Anthropic para Go)
- Escuchar en stdio (estándar para MCP)
- Implementar handlers para:
  - `tools/list` - Listar herramientas disponibles
  - `tools/call` - Ejecutar la conversión

#### 4. **Configuración en Claude Desktop**
El archivo de configuración `claude_desktop_config.json` debe incluir:

```json
{
  "mcpServers": {
    "png-to-ico": {
      "command": "/ruta/al/binario/mcp-png-convert",
      "args": []
    }
  }
}
```

#### 5. **Compatibilidad multiplataforma**
El binario compilado debe ser:
- `mcp-png-convert` en macOS/Linux
- `mcp-png-convert.exe` en Windows

Compilar para múltiples plataformas:
```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o dist/mcp-png-convert.exe

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o dist/mcp-png-convert-darwin-amd64

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o dist/mcp-png-convert-darwin-arm64

# Linux
GOOS=linux GOARCH=amd64 go build -o dist/mcp-png-convert
```

#### 6. **Validaciones en el servidor MCP**
- ✅ Validar que el archivo PNG existe
- ✅ Validar que es un PNG válido
- ✅ Validar permisos de lectura/escritura
- ✅ Retornar códigos de error específicos
- ✅ Logging estructurado (stderr)

#### 7. **Documentación README.md requerida**
El proyecto debe incluir:
- Instrucciones de instalación
- Cómo configurar en Claude Desktop
- Ejemplos de uso
- Requisitos del sistema
- Licencia

#### 8. **Testing**
- Pruebas unitarias para conversión
- Pruebas de integración MCP
- Test con archivos PNG de diferentes tamaños/formatos
- Validar que los ICO generados sean válidos con herramientas externas