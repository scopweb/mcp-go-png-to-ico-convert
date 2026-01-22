# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-01-22

### Added
- Initial release of MCP PNG to ICO Converter
- CLI mode for direct conversions from command line
  - Support for explicit input/output paths
  - Auto-generate output filename from input when not specified
- MCP server mode for Claude Desktop integration
  - `convert_png_to_ico` tool with JSON responses
  - Error handling with specific error codes (FILE_NOT_FOUND, INVALID_PNG, ENCODING_ERROR)
- Multi-resolution ICO support with 6 resolutions:
  - 256x256, 128x128, 64x64, 48x48, 32x32, 16x16
- High-quality image resampling using LANCZOS filter
- RGBA transparency preservation from source PNG
- Cross-platform support (Windows, macOS, Linux)
- Pure Go implementation without CGO dependencies

### Technical Details
- Built with Go 1.21+
- Dependencies:
  - `github.com/modelcontextprotocol/go-sdk` - MCP protocol implementation
  - `github.com/sergeymakinen/go-ico` - ICO encoding
  - `github.com/disintegration/imaging` - High-quality image processing
  - `image/png` (stdlib) - PNG decoding

### Documentation
- Comprehensive README.md with usage examples
- CLAUDE.md with project specifications
- MIT License
- .gitignore for Go projects

[unreleased]: https://github.com/scopweb/mcp-go-png-to-ico-convert/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/scopweb/mcp-go-png-to-ico-convert/releases/tag/v1.0.0
