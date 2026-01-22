# 🖼️ MCP PNG to ICO Converter

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![MCP](https://img.shields.io/badge/MCP-Compatible-5865F2?style=flat)](https://modelcontextprotocol.io)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

MCP (Model Context Protocol) server for converting PNG files to multi-resolution ICO format, optimized for .NET 10 applications.

## Features

- **High-quality conversion**: Uses LANCZOS filter for maximum quality resampling
- **Multi-resolution support**: Generates ICO with 6 resolutions (256x256, 128x128, 64x64, 48x48, 32x32, 16x16)
- **Alpha transparency**: Preserves RGBA transparency from original PNG
- **Dual operation mode**: Direct CLI or MCP server for Claude Desktop
- **MCP Protocol**: Native integration with Claude Desktop and other MCP tools
- **No CGO dependencies**: Pure Go compilation, no C dependencies

## System Requirements

- Go 1.21 or higher (for compilation)
- Windows, macOS, or Linux
- Claude Desktop (for MCP server usage)

## Installation

### Option 1: Compile from source

```bash
# Clone the repository
git clone https://github.com/scopweb/mcp-go-png-to-ico-convert.git
cd mcp-go-png-to-ico-convert

# Install dependencies
go mod download

# Compile for your platform
go build -o mcp-png-convert .

# On Windows
go build -o mcp-png-convert.exe .
```

### Option 2: Cross-platform compilation

To generate binaries for multiple platforms:

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

## Claude Desktop Configuration

1. Locate your Claude Desktop configuration file:
   - **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
   - **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - **Linux**: `~/.config/Claude/claude_desktop_config.json`

2. Add the following configuration to the JSON file:

```json
{
  "mcpServers": {
    "png-to-ico": {
      "command": "/full/path/to/binary/mcp-png-convert",
      "args": []
    }
  }
}
```

**Windows Example:**
```json
{
  "mcpServers": {
    "png-to-ico": {
      "command": "C:\\Users\\YourUser\\mcp-servers\\mcp-png-convert.exe",
      "args": []
    }
  }
}
```

**macOS/Linux Example:**
```json
{
  "mcpServers": {
    "png-to-ico": {
      "command": "/home/user/mcp-servers/mcp-png-convert",
      "args": []
    }
  }
}
```

3. Restart Claude Desktop

4. Verify the server is connected (a tools icon will appear in the interface)

## Usage

The tool supports **two operation modes**:

### Mode 1: Direct CLI (without AI)

Use it directly from the command line for quick conversions:

```bash
# With explicit input and output
./mcp-png-convert app.png app.ico

# Input only (automatically generates app.ico)
./mcp-png-convert app.png

# On Windows
.\mcp-png-convert.exe logo.png logo.ico

# Or with go run (during development)
go run main.go app.png app.ico
```

**Example output:**
```
🔄 Converting: app.png -> app.ico
✅ ICO generated successfully
📦 File: app.ico
📐 Resolutions: [256x256 128x128 64x64 48x48 32x32 16x16]
💾 Size: 24.5 KB
```

### Mode 2: MCP Server with Claude Desktop

Once configured, you can use the server directly from Claude Desktop:

```
Convert app.png to app.ico
```

```
Convert my-logo.png to icon.ico
```

Claude will automatically use the `convert_png_to_ico` tool from the MCP server.

### Tool Parameters

The `convert_png_to_ico` tool accepts the following parameters:

- **input_path** (required): Path to input PNG file (absolute or relative)
- **output_path** (optional): Path to output ICO file. If not specified, uses the same name as PNG with `.ico` extension

### Response Format

The tool returns a JSON with the result:

**Success:**
```json
{
  "success": true,
  "output_file": "/path/to/file.ico",
  "resolutions": ["256x256", "128x128", "64x64", "48x48", "32x32", "16x16"],
  "file_size": "24.5 KB",
  "message": "ICO generated successfully"
}
```

**Error:**
```json
{
  "success": false,
  "error": "Error description",
  "error_code": "FILE_NOT_FOUND|INVALID_PNG|ENCODING_ERROR"
}
```

## Technical Architecture

### Library Stack

- **[github.com/modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk)**: Official MCP SDK
- **[github.com/sergeymakinen/go-ico](https://github.com/sergeymakinen/go-ico)**: ICO format encoding
- **[github.com/disintegration/imaging](https://github.com/disintegration/imaging)**: High-quality resizing
- **image/png** (built-in): PNG decoding

### Conversion Flow

1. Reads PNG file using `image/png.Decode()`
2. Resizes to each resolution (256, 128, 64, 48, 32, 16) using Lanczos filter
3. Encodes all images into a single ICO file with `ico.EncodeAll()`
4. Returns information about the generated file

### Error Codes

| Code | Description |
|------|-------------|
| `FILE_NOT_FOUND` | PNG file doesn't exist or lacks read permissions |
| `INVALID_PNG` | File is not a valid or is corrupted PNG |
| `ENCODING_ERROR` | Error encoding ICO file or writing output file |

## Development

### Project Structure

```
mcp-go-png-to-ico-convert/
├── main.go                  # Entry point (hybrid CLI/MCP mode)
├── internal/
│   ├── converter/
│   │   ├── converter.go     # PNG->ICO conversion logic
│   │   └── converter_test.go # Unit tests
│   └── mcp/
│       └── server.go        # MCP server
├── tests/
│   ├── security/
│   │   └── security_test.go # Security tests
│   └── integration/         # Integration tests (future)
├── dist/                    # Compiled binaries (gitignored)
├── go.mod
├── go.sum
├── .gitignore
├── CLAUDE.md               # Project specifications
├── CHANGELOG.md            # Version history
├── README.md               # Spanish documentation
└── README.en.md            # English documentation
```

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Run only security tests
go test -v ./tests/security/...

# Run benchmarks
go test -bench=. ./internal/converter/
```

### Test Coverage

- **Converter package**: 88.6% coverage
- Security tests: Path traversal, invalid inputs, permissions, concurrent access
- Unit tests: Valid conversion, error handling, file size formatting

## Troubleshooting

### Server doesn't appear in Claude Desktop

1. Verify the binary path is absolute and correct
2. On Unix/Linux, ensure the binary has execution permissions: `chmod +x mcp-png-convert`
3. Check Claude Desktop logs for errors
4. Restart Claude Desktop after changing configuration

### "FILE_NOT_FOUND" Error

- Verify the PNG file path is correct
- Ensure the file exists and you have read permissions
- Use absolute paths if you have issues with relative paths

### "INVALID_PNG" Error

- Verify the file is a valid PNG
- Try opening the file with an image viewer to confirm it's not corrupted
- Ensure the file has the `.png` extension

### Generated ICO doesn't work on Windows

- Verify the original PNG has a minimum size of 256x256 pixels
- Ensure the PNG has the correct format (RGBA)
- The ICO must contain all standard resolutions (done automatically)

## License

MIT License

Copyright (c) 2026 scopweb

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

## Contributing

Contributions are welcome. Please:

1. Fork the project
2. Create a branch for your feature (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## References

- [Model Context Protocol](https://modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [Claude Desktop](https://claude.ai/desktop)
- [ICO Format Specification](https://en.wikipedia.org/wiki/ICO_(file_format))

## Support

To report bugs or request features, open an issue on the GitHub repository.

---

Developed with ❤️ using Go and the MCP protocol
