import sys
import logging
logging.disable(logging.DEBUG)  # 关闭DEBUG日志的打印


from paddleocr import PaddleOCR

# 编码设置
sys.stdout.reconfigure(encoding='utf-8')

def get_OCR(img_path):
    # Paddleocr目前支持的多语言语种可以通过修改lang参数进行切换
    # 例如`ch`, `en`, `fr`, `german`, `korean`, `japan`
    ocr = PaddleOCR(use_angle_cls=True, lang="ch")  # need to run only once to download and load model into memory
    OCRstr = ""
    
    result = ocr.ocr(img_path, cls=True)
    
    for idx in range(len(result)):
        res = result[idx]
        for line in res:
            OCRstr += str(line[1][0]) + " "
    OCRstr += "\n"
    return OCRstr

if __name__ == "__main__":
    # 从命令行参数获取图片路径
    if len(sys.argv) != 2:
        print("Usage: python ocr.py <image_path>")
        sys.exit(1)
    image_path = sys.argv[1]
    result = get_OCR(image_path)
    
    # 打印OCR结果到标准输出
    print(result)