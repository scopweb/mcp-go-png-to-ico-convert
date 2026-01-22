# MCP PNG to ICO Converter

Servidor MCP (Model Context Protocol) para convertir archivos PNG a ICO con múltiples resoluciones, optimizado para aplicaciones .NET 10.

## Características

- **Conversión de alta calidad**: Utiliza filtro LANCZOS para redimensionamiento de máxima calidad
- **Múltiples resoluciones**: Genera ICO con 6 resoluciones (256x256, 128x128, 64x64, 48x48, 32x32, 16x16)
- **Transparencia alfa**: Mantiene transparencia RGBA del PNG original
- **Protocolo MCP**: Integración nativa con Claude Desktop y otras herramientas MCP
- **Sin dependencias CGO**: Compilación pura en Go, sin dependencias de C

## Requisitos del sistema

- Go 1.21 o superior (para compilación)
- Windows, macOS, o Linux
- Claude Desktop (para uso como servidor MCP)

## Instalación

### Opción 1: Compilar desde código fuente

```bash
# Clonar el repositorio
git clone https://github.com/yourusername/mcp-go-png-convert.git
cd mcp-go-png-convert

# Instalar dependencias
go mod download

# Compilar para tu plataforma
go build -o mcp-png-convert .

# En Windows
go build -o mcp-png-convert.exe .
```

### Opción 2: Compilación multiplataforma

Para generar binarios para múltiples plataformas:

```bash
# Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -o dist/mcp-png-convert.exe

# macOS Intel (64-bit)
GOOS=darwin GOARCH=amd64 go build -o dist/mcp-png-convert-darwin-amd64

# macOS Apple Silicon (ARM64)
GOOS=darwin GOARCH=arm64 go build -o dist/mcp-png-convert-darwin-arm64

# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o dist/mcp-png-convert
```

## Configuración en Claude Desktop

1. Localiza tu archivo de configuración de Claude Desktop:
   - **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
   - **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - **Linux**: `~/.config/Claude/claude_desktop_config.json`

2. Agrega la siguiente configuración al archivo JSON:

```json
{
  "mcpServers": {
    "png-to-ico": {
      "command": "/ruta/completa/al/binario/mcp-png-convert",
      "args": []
    }
  }
}
```

**Ejemplo en Windows:**
```json
{
  "mcpServers": {
    "png-to-ico": {
      "command": "C:\\Users\\TuUsuario\\mcp-servers\\mcp-png-convert.exe",
      "args": []
    }
  }
}
```

**Ejemplo en macOS/Linux:**
```json
{
  "mcpServers": {
    "png-to-ico": {
      "command": "/home/usuario/mcp-servers/mcp-png-convert",
      "args": []
    }
  }
}
```

3. Reinicia Claude Desktop

4. Verifica que el servidor esté conectado (aparecerá un icono de herramientas en la interfaz)

## Uso

### Desde Claude Desktop

Una vez configurado, puedes usar el servidor directamente desde Claude Desktop:

```
Convierte app.png a app.ico
```

```
Convierte mi-logo.png a icono.ico
```

Claude automáticamente utilizará la herramienta `convert_png_to_ico` del servidor MCP.

### Parámetros de la herramienta

La herramienta `convert_png_to_ico` acepta los siguientes parámetros:

- **input_path** (requerido): Ruta del archivo PNG de entrada (absoluta o relativa)
- **output_path** (opcional): Ruta del archivo ICO de salida. Si no se especifica, se usa el mismo nombre del PNG con extensión `.ico`

### Formato de respuesta

La herramienta retorna un JSON con el resultado:

**Éxito:**
```json
{
  "success": true,
  "output_file": "/ruta/al/archivo.ico",
  "resolutions": ["256x256", "128x128", "64x64", "48x48", "32x32", "16x16"],
  "file_size": "24.5 KB",
  "message": "ICO generado exitosamente"
}
```

**Error:**
```json
{
  "success": false,
  "error": "Descripción del error",
  "error_code": "FILE_NOT_FOUND|INVALID_PNG|ENCODING_ERROR"
}
```

## Arquitectura técnica

### Stack de librerías

- **[github.com/modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk)**: SDK oficial de MCP
- **[github.com/sergeymakinen/go-ico](https://github.com/sergeymakinen/go-ico)**: Encoding de formato ICO
- **[github.com/disintegration/imaging](https://github.com/disintegration/imaging)**: Redimensionamiento de alta calidad
- **image/png** (built-in): Decodificación de PNG

### Flujo de conversión

1. Lee el archivo PNG usando `image/png.Decode()`
2. Redimensiona a cada resolución (256, 128, 64, 48, 32, 16) usando filtro Lanczos
3. Codifica todas las imágenes en un solo archivo ICO con `ico.EncodeAll()`
4. Retorna información del archivo generado

### Códigos de error

| Código | Descripción |
|--------|-------------|
| `FILE_NOT_FOUND` | El archivo PNG no existe o no tiene permisos de lectura |
| `INVALID_PNG` | El archivo no es un PNG válido o está corrupto |
| `ENCODING_ERROR` | Error al codificar el archivo ICO o escribir el archivo de salida |

## Desarrollo

### Estructura del proyecto

```
mcp-go-png-convert/
├── main.go                  # Punto de entrada
├── internal/
│   ├── converter/
│   │   └── converter.go     # Lógica de conversión PNG->ICO
│   └── mcp/
│       └── server.go        # Servidor MCP
├── dist/                    # Binarios compilados
├── go.mod
├── go.sum
├── CLAUDE.md               # Especificaciones del proyecto
└── README.md
```

### Testing manual

Para probar el servidor manualmente:

1. Crea un archivo PNG de prueba (`test.png`)
2. Ejecuta el servidor en modo debug (opcional: redirige logs a un archivo)
3. Usa un cliente MCP para llamar a la herramienta

## Solución de problemas

### El servidor no aparece en Claude Desktop

1. Verifica que la ruta al binario sea absoluta y correcta
2. En Unix/Linux, asegúrate de que el binario tenga permisos de ejecución: `chmod +x mcp-png-convert`
3. Revisa los logs de Claude Desktop para ver errores
4. Reinicia Claude Desktop después de cambiar la configuración

### Error "FILE_NOT_FOUND"

- Verifica que la ruta del archivo PNG sea correcta
- Asegúrate de que el archivo existe y tienes permisos de lectura
- Usa rutas absolutas si tienes problemas con rutas relativas

### Error "INVALID_PNG"

- Verifica que el archivo sea un PNG válido
- Intenta abrir el archivo con un visor de imágenes para confirmar que no está corrupto
- Asegúrate de que el archivo tenga la extensión `.png`

### El ICO generado no funciona en Windows

- Verifica que el PNG original tenga un tamaño mínimo de 256x256 píxeles
- Asegúrate de que el PNG tenga el formato correcto (RGBA)
- El ICO debe contener todas las resoluciones estándar (esto se hace automáticamente)

## Licencia

MIT License

Copyright (c) 2026 [Your Name]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## Contribuciones

Las contribuciones son bienvenidas. Por favor:

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/amazing-feature`)
3. Commit tus cambios (`git commit -m 'Add amazing feature'`)
4. Push a la rama (`git push origin feature/amazing-feature`)
5. Abre un Pull Request

## Referencias

- [Model Context Protocol](https://modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [Claude Desktop](https://claude.ai/desktop)
- [ICO Format Specification](https://en.wikipedia.org/wiki/ICO_(file_format))

## Soporte

Para reportar bugs o solicitar features, abre un issue en el repositorio de GitHub.

---

Desarrollado con ❤️ usando Go y el protocolo MCP
