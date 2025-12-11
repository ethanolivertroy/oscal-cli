# OSCAL CLI

A fast, lightweight command-line tool for working with [Open Security Controls Assessment Language](https://pages.nist.gov/OSCAL/) (OSCAL) documents.

## Features

- **Format Conversion**: Convert OSCAL documents between XML, JSON, and YAML formats
- **Auto-detection**: Automatically detects input format from file extension or content
- **Fast**: Converts the full NIST 800-53 catalog (10MB, 1000+ controls) in ~0.1 seconds
- **Validation**: JSON Schema validation against official OSCAL 1.1.3 schemas
- **Lightweight**: Single static binary (~6MB with embedded schemas), no runtime dependencies
- **Cross-platform**: Builds for macOS, Linux, and Windows

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/ethantroy/oscal-cli.git
cd oscal-cli

# Build
make build

# Or install to your GOPATH/bin
make install
```

### Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/ethantroy/oscal-cli/releases).

## Usage

### Validate Command

Validate OSCAL documents against official NIST JSON schemas:

```bash
# Validate a catalog
oscal validate catalog.json

# Validate an XML document (auto-converted internally)
oscal validate ssp.xml

# Quiet mode (only output errors)
oscal validate catalog.json --quiet
```

The embedded schemas correspond to OSCAL version 1.1.3.

### Convert Command

Convert OSCAL documents between formats:

```bash
# Convert JSON to XML
oscal convert catalog.json --to xml

# Convert XML to YAML with output file
oscal convert catalog.xml --to yaml --output catalog.yaml

# Convert and redirect to file
oscal convert catalog.xml --to json > catalog.json

# Overwrite existing output file
oscal convert catalog.json --to xml --output catalog.xml --overwrite
```

### Supported Document Types

- Catalog
- Profile
- System Security Plan (SSP)
- Component Definition
- Assessment Plan
- Assessment Results
- Plan of Action and Milestones (POA&M)

### Version

```bash
oscal version
```

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for convenience commands)

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Clean build artifacts
make clean
```

### Project Structure

```
oscal-cli/
├── cmd/oscal/          # CLI entry point
├── internal/cli/       # CLI commands (cobra)
├── pkg/oscal/
│   ├── model/          # OSCAL document types
│   ├── io/             # Format detection, load/save
│   └── validate/       # Validation (future)
└── testdata/           # Test fixtures
```

## Roadmap

- [x] `validate` command with JSON Schema validation
- [ ] Profile resolution (`profile resolve`)
- [ ] HTML rendering (`catalog render`)
- [ ] OSCAL constraint validation (beyond schema)

## License

This project is in the public domain. See [LICENSE.md](LICENSE.md) for details.

## Acknowledgments

This project is a Go rewrite of NIST's [oscal-cli](https://github.com/usnistgov/oscal-cli), originally written in Java. Thanks to NIST and the OSCAL team for their work on the original tool and the OSCAL specification.

## Related Projects

- [NIST OSCAL](https://pages.nist.gov/OSCAL/) - Official OSCAL specification
- [oscal-cli (Java)](https://github.com/usnistgov/oscal-cli) - NIST's original Java CLI (this project's inspiration)
- [liboscal-java](https://github.com/usnistgov/liboscal-java) - NIST's Java OSCAL library
- [go-oscal](https://github.com/defenseunicorns/go-oscal) - Defense Unicorns Go OSCAL library
