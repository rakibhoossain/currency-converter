#!/bin/bash

# Currency Converter API - Multi-platform Build Script

set -e

# Build configuration
APP_NAME="currency-converter"
VERSION=${VERSION:-"1.0.0"}
BUILD_DIR="dist"

# Target platforms
TARGETS=(
    "linux/amd64/"
    "windows/amd64/.exe"
    "darwin/amd64/"
    "darwin/arm64/"
)

echo "🚀 Building Currency Converter API v${VERSION}"
echo "=============================================="

# Clean previous builds
if [ -d "$BUILD_DIR" ]; then
    echo "🧹 Cleaning previous builds..."
    rm -rf "$BUILD_DIR"
fi

# Create build directory
mkdir -p "$BUILD_DIR"

# Copy .env.example to build directory
echo "📄 Copying configuration files..."
cp .env.example "$BUILD_DIR/.env.example"

# Create README for distribution
cat > "$BUILD_DIR/README.txt" << 'EOF'
Currency Converter API - Distribution Package
=============================================

Quick Start:
1. Copy .env.example to .env
2. Edit .env and set your OpenExchangeRates App ID and Auth Token:
   - OXR_APP_ID=your_app_id_here
   - AUTH_TOKEN=your_secure_token_here
3. Run the appropriate binary for your platform
4. API will be available at http://localhost:3000

Authentication:
All API endpoints require Bearer token authentication.
Include in header: Authorization: Bearer your_auth_token_here

Endpoints:
- GET /api/health - Health check
- GET /api/currencies - Get currency symbols
- GET /api/rates - Get exchange rates
- POST /api/convert - Convert currencies

For full documentation, visit:
https://github.com/rakibhoossain/currency-converter

EOF

# Build for each target platform
for target in "${TARGETS[@]}"; do
    # Parse target
    IFS='/' read -r os arch extension <<< "$target"
    
    # Set output filename
    output_name="${APP_NAME}-${os}-${arch}${extension}"
    output_path="${BUILD_DIR}/${output_name}"
    
    echo "🔨 Building for ${os}/${arch}..."
    
    # Set environment variables and build
    GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build \
        -ldflags "-s -w -X main.Version=${VERSION}" \
        -o "$output_path" \
        .
    
    # Create platform-specific package
    platform_dir="${BUILD_DIR}/${APP_NAME}-${VERSION}-${os}-${arch}"
    mkdir -p "$platform_dir"
    
    # Copy binary
    cp "$output_path" "$platform_dir/"
    
    # Copy configuration files
    cp .env.example "$platform_dir/"
    cp "$BUILD_DIR/README.txt" "$platform_dir/"
    
    # Create platform-specific start script
    if [ "$os" = "windows" ]; then
        cat > "$platform_dir/start.bat" << 'EOF'
@echo off
echo Starting Currency Converter API...
echo Make sure you have configured .env file
echo Press Ctrl+C to stop the server
currency-converter-windows-amd64.exe
pause
EOF
    else
        cat > "$platform_dir/start.sh" << EOF
#!/bin/bash
echo "Starting Currency Converter API..."
echo "Make sure you have configured .env file"
echo "Press Ctrl+C to stop the server"
./${output_name}
EOF
        chmod +x "$platform_dir/start.sh"
    fi
    
    # Create archive
    echo "📦 Creating archive for ${os}/${arch}..."
    cd "$BUILD_DIR"
    
    if [ "$os" = "windows" ]; then
        zip -r "${APP_NAME}-${VERSION}-${os}-${arch}.zip" "${APP_NAME}-${VERSION}-${os}-${arch}/" > /dev/null
    else
        tar -czf "${APP_NAME}-${VERSION}-${os}-${arch}.tar.gz" "${APP_NAME}-${VERSION}-${os}-${arch}/"
    fi
    
    cd ..
    
    # Clean up directory (keep archive)
    rm -rf "$platform_dir"
    
    echo "✅ Built: ${output_name}"
done

# Create checksums
echo "🔐 Generating checksums..."
cd "$BUILD_DIR"
if command -v sha256sum >/dev/null 2>&1; then
    sha256sum *.tar.gz *.zip 2>/dev/null > checksums.txt || true
elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 *.tar.gz *.zip 2>/dev/null > checksums.txt || true
fi
cd ..

echo ""
echo "✅ Build completed successfully!"
echo "📁 Build artifacts are in the '${BUILD_DIR}' directory:"
echo ""
ls -la "$BUILD_DIR"
echo ""
echo "🚀 Ready for distribution!"
