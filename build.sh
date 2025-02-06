OUTPUT_DIR="builds"
mkdir -p $OUTPUT_DIR

# Build for Intel Mac
echo "Building for MacOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -o $OUTPUT_DIR/ftparchive-mac-amd64 ./cmd/ftparchive

# Build for Apple Silicon Mac
echo "Building for MacOS (Arm)..."
GOOS=darwin GOARCH=arm64 go build  -o $OUTPUT_DIR/ftparchive-mac-arm64 ./cmd/ftparchive

# Build for Windows
echo "Building for Windows (x86)..."
GOOS=windows GOARCH=amd64 go build  -o $OUTPUT_DIR/ftparchive-windows-amd64.exe ./cmd/ftparchive

echo "Building for Windows (Arm)..."
GOOS=windows GOARCH=arm64 go build  -o $OUTPUT_DIR/ftparchive-windows-arm.exe ./cmd/ftparchive

# Build for Linux
echo "Building for Linux (x86)..."
GOOS=linux GOARCH=amd64 go build  -o $OUTPUT_DIR/ftparchive-linux-amd64 ./cmd/ftparchive

echo "Building for Linux (ARM)..."
GOOS=linux GOARCH=arm64 go build  -o $OUTPUT_DIR/ftparchive-linux-arm ./cmd/ftparchive

#Next, zip all the builds
cd ./$OUTPUT_DIR

echo "Now zipping builds..."

zip ftparchive-mac-amd64.zip ftparchive-mac-amd64
zip ftparchive-mac-arm64.zip ftparchive-mac-arm64
zip ftparchive-windows-amd64.zip ftparchive-windows-amd64.exe
zip ftparchive-windows-arm.zip ftparchive-windows-arm.exe
zip ftparchive-linux-amd64.zip ftparchive-linux-amd64
zip ftparchive-linux-arm.zip ftparchive-linux-arm