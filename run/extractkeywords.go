package run

import (
	"apifinder/keywords"
	"encoding/json"
	"log"
	"os"
	"path"
)

func ExtractKeywords(host string, tmpDir string, LLMURL string) {
	filename := "APIInfo-" + host + ".json"

	f, err := os.ReadFile(path.Join(tmpDir, filename))
	if err != nil {
		log.Fatal(err)
	}

	var jsonData []data
	if err = json.Unmarshal(f, &jsonData); err != nil {
		log.Fatal(err)
	}

	for index, item := range jsonData {
		log.Printf("\rWorking on the %d item...", index+1)
		c, err := keywords.ChatCompletions(LLMURL, "Llama-3.1-8B-Instruct", item.Content, nil)
		if err != nil {
			log.Fatal(err)
			continue
		}
		jsonData[index].Content = c
	}

	log.Println()

	output, err := json.MarshalIndent(jsonData, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	outputPath := path.Join(tmpDir, "extracted-"+filename)
	file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err = file.Write(output); err != nil {
		log.Fatal(err)
	}
}
