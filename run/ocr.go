package run

import (
	"apifinder/ocr"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

func RunOCR(host string, tmpDir string, chromePath string, chromeDriverPath string) {
	filename := "APIInfo-" + host + ".json"

	f, err := os.ReadFile(path.Join(tmpDir, filename))
	if err != nil {
		logrus.Fatal("failed to read json file: ", err)
	}

	var jsonData []data
	if err = json.Unmarshal(f, &jsonData); err != nil {
		logrus.Fatal("failed to parse json file: ", err)
	}

	for index, item := range jsonData {
		fmt.Printf("\rOCR in progress: %d/%d ", index+1, len(jsonData))
		filename, err := ocr.TakeScreenShot(item.URL, host, tmpDir, chromePath, chromeDriverPath)
		if err != nil {
			logrus.Error("failed to take screenshot: ", err)
			continue
		}

		text, err := ocr.Get_OCR(filename)
		fmt.Print("=======OCR text: ", text)
		if err != nil {
			logrus.Error("failed to get OCR: ", err)
			continue
		}

		jsonData[index].Content += "\n\n" + text
	}
}
