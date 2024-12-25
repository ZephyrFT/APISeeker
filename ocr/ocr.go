package ocr

import (
	"log"
	"os"
	"os/exec"
)

func Get_OCR(imagePath string) (string, error) {

	// 检查图片文件是否存在
	_, err := os.Stat(imagePath)
	if os.IsNotExist(err) {
		log.Fatalf("文件不存在: %v\n", imagePath)
	} else if err != nil {
		log.Fatalf("无法访问文件: %v\n", err)
	}

	// 使用 exec.Command 调用 Python 脚本并传递图片路径
	cmd := exec.Command("python", "ocr.py", imagePath)

	// 获取 Python 脚本的标准输出
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	decodedOutput := string(output)

	// 打印 Python 脚本的输出（OCR 结果）
	// fmt.Println("OCR Output:", decodedOutput)

	return decodedOutput, nil
}

// func main() {
// 	imagePath := "image.png"
// 	fmt.Println(Get_OCR(imagePath))
// }
