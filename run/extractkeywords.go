package run

import (
	"apifinder/keywords"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

func ExtractKeywords(host string, tmpDir string, LLMURL string) {
	filename := "APIInfo-" + host + ".json"

	f, err := os.ReadFile(path.Join(tmpDir, filename))
	if err != nil {
		logrus.Fatal(err)
	}

	var jsonData []data
	if err = json.Unmarshal(f, &jsonData); err != nil {
		logrus.Fatal(err)
	}

	for index, item := range jsonData {
		fmt.Printf("\rWorking on the %d item...", index+1)
		c, err := keywords.ChatCompletions(LLMURL, "Llama-3.1-8B-Instruct", item.Content, nil)
		if err != nil {
			logrus.Error(err)
			continue
		}
		jsonData[index].Content = c
	}

	fmt.Println()

	output, err := json.MarshalIndent(jsonData, "", "\t")
	if err != nil {
		logrus.Fatal(err)
	}

	outputPath := path.Join(tmpDir, "extracted-"+filename)
	file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		logrus.Fatal(err)
	}
	defer file.Close()

	if _, err = file.Write(output); err != nil {
		logrus.Fatal(err)
	}
}
