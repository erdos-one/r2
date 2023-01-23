#!/bin/bash
REPO_USER=$(git config --get remote.origin.url | cut -d'/' -f4)
REPO_NAME=$(basename `git rev-parse --show-toplevel`)
TAG_NAME=$(git describe --tags --abbrev=0)
BINARY_NAME="r2"

cat > install.sh << EOF
#!/bin/sh
OS=\`uname -s\`
ARCH=\`uname -m\`

curl -fsSLo \$HOME/${BINARY_NAME}-\${OS}-\${ARCH}.tar.gz \\
	https://github.com/${REPO_USER}/${REPO_NAME}/releases/download/${TAG_NAME}/${BINARY_NAME}-\${OS}-\${ARCH}.tar.gz
	
mkdir -p \$HOME/${BINARY_NAME}-${TAG_NAME}
tar -xzf \$HOME/${BINARY_NAME}-\${OS}-\${ARCH}.tar.gz -C \$HOME/${BINARY_NAME}-${TAG_NAME}
chmod +x \$HOME/${BINARY_NAME}-${TAG_NAME}/${BINARY_NAME}

if [ "\$EUID" -eq 0 ]; then
	mv \$HOME/${BINARY_NAME}-${TAG_NAME}/${BINARY_NAME} /usr/bin/${BINARY_NAME}
	rm -r \$HOME/${BINARY_NAME}-${TAG_NAME}
else
	echo "Couldn't add ${BINARY_NAME} to /usr/bin because script is not running at admin"
	echo "Please run the following commands:"
	echo "  mv \$HOME/${BINARY_NAME}-${TAG_NAME}/${BINARY_NAME}-\${OS}-\${ARCH} /usr/bin/${BINARY_NAME}"
	echo "  rm -r \$HOME/${BINARY_NAME}-${TAG_NAME}"
fi
EOF