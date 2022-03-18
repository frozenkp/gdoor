package screenshot

import (
	"fmt"
	"github.com/kbinani/screenshot"
	"image/png"
	"os"
	"time"
)

func TakeScreenShot() ([]string, error) {
	n := screenshot.NumActiveDisplays()
	t := time.Now()

	output := make([]string, 0)
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			return nil, err
		}

		fileName := fmt.Sprintf("%d_%s.png", i, t.Format("20060102150405"))
		file, err := os.Create(fileName)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		png.Encode(file, img)
		output = append(output, fileName)
	}

	return output, nil
}
