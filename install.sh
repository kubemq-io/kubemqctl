#! /bin/sh
OS=$(uname -s)
arch=$(arch)
echo "$OS"
echo "$arch"

case $OS in
  MSYS_NT-10.0)
    echo "this script is not supported on Windows"
    echo " - Download kubemqctl from: https://github.com/kubemq-io/kubemqctl/releases/download/latest/kubemqctl.exe"
    echo " - Place the file under e.g. C:\Program Files\Kubemqctl\kubemqctl.exe"
    echo " - Add that directory to your system path to access it from any command prompt"
    exit 1
    ;;
  Darwin)
    filename="kubemqctl_darwin_amd64"
    ;;
  Linux)
    case $arch in
      x86_64)
        filename="kubemqctl_linux_amd64"
        ;;
      i386)
        filename="kubemqctl_linux_386"
        ;;
      *)
        echo "There is no kubemqctl $OS support for $arch" >&2
        exit 1
        ;;
    esac
    ;;
  *)
    echo "There is no kubemqctl support for $OS/$arch (yet...)" >&2
    exit 1
    ;;
esac
#set -eu

#filename="kubemqctl_linux_amd64"
url=https://github.com/kubemq-io/kubemqctl/releases/download/latest/${filename}
(
  echo "Downloading ${filename}..."
  curl -LO "${url}" -o /usr/local/bin/kubemqctl
  echo ""
  echo "Download complete!"
)

(
  chmod +x "/usr/local/bin/kubemqctl"
)

echo "kubemqctl was successfully installed 🎉"

echo "Now run:"
echo "kubemqctl                               # validate that kubemqctl is installed"
echo ""
echo "Looking for more? Visit https://docs.kubemq.io/getting-started/quick-start"
echo ""
