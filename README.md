# APISeeker
A novel tool for proactive detection of shadow APIs in black-box settings based on Go.

# Ethical Considerations
We are committed to ensuring that APISeeker is used responsibly. To mitigate potential negative impacts:

+ Access Control: The core of APISeeker is available only through an application process, ensuring that it is used by authorized users or organizations.
+ Responsible Usage: We encourage users to deploy APISeeker for security research and vulnerability discovery, not for malicious purposes.
+ Community Collaboration: We call for community involvement in developing standards and best practices for using tools like APISeeker to ensure their positive impact on web security.

# Usage scenarios
This software may only be used in defense-related scenarios, including but not limited to cybersecurity research, vulnerability assessment, and defensive security measures.

# License
This project is licensed under our PRIVATE License - see the LICENSE file for details.

# Usage

## Paddle OCR

APISeeker's OCR module is based on Paddle OCR. If you want to use the OCR function, please refer to Paddle OCR's documentation to install the Paddle OCR model.

Here is [Paddle OCR Docs](https://paddlepaddle.github.io/PaddleOCR/main/index.html).

## LLM

To prevent the abuse of the tool, we do not open source the model parameters after fine-tuning. APISeeker supports large language model servers based on the vllm architecture that conform to the OpenAI interface.

## APISeeker

APISeeker is developed in Go language and is configured and run through command-line parameters. The following is a detailed description of each command-line parameter and its usage example:

### Command-line Parameters
- `-url` or `--website`: Specifies the URL of the website to crawl, the default is `https://www.google.com`.
    ```sh
  go run main.go -url=https://example.com
    ```
- `-host` or `--host`: Specifies the hostname of the website, and it cannot be empty.
    ```sh
  go run main.go -host=example.com
    ```
- `-tmp` or `--tmp`: Specifies the temporary directory, the default is `.tmp`. If the input path is not an absolute path, the program's current working directory will be used to concatenate the path.
    ```sh
  go run main.go -tmp=/path/to/tmp
    ```
- `-result` or `--result`: Specifies the result storage directory, the default is `result`. If the input path is not an absolute path, the program's current working directory will be used to concatenate the path.
    ```sh
  go run main.go -result=/path/to/result
    ```
- `-keywords` or `--keywords`: Specifies the service URL of the keyword extraction large model, the default is `http://127.0.0.1:8000/v1/chat/completions`.
    ```sh
  go run main.go -keywords=http://your_keywords_llm:8000/v1/chat/completions
    ```
- `-validations` or `--validations`: Specifies the service URL of the validation large model, the default is `http://127.0.0.1:8000/v1/chat/completions`.
    ```sh
  go run main.go -validations=http://your_validations_llm:8000/v1/chat/completions
    ```
- `-chromedriver` or `--chromedriver`: Specifies the path of ChromeDriver, the default is `/usr/local/bin/chromedriver`.
    ```sh
  go run main.go -chromedriver=/path/to/chromedriver
    ```
- `-chrome` or `--chrome`: Specifies the path of the Chrome browser, the default is `/usr/bin/google-chrome`.
    ```sh
  go run main.go -chrome=/path/to/google-chrome
    ```
- `-cookie` or `--cookie`: You can specify a cookie for possible authentication or other purposes (default is empty).
    ```sh
  go run main.go -cookie="your_cookie_value"
    ```
- `-limit` or `--limit`: Specifies the maximum crawling limit, the default is 500.
    ```sh
  go run main.go -limit=1000
    ```
- `-ocr` or `--ocr`: A boolean flag indicating whether to use OCR to extract text from images, the default is `false`.
    ```sh
  go run main.go -ocr=true
    ```


### Example
Here is a complete example showing how to run the program using all parameters:
```sh
go run main.go -url=https://example.com -host=example.com -tmp=/tmp -result=/results -keywords=http://your_keywords_llm:8000/v1/chat/completions -validations=http://your_validations_llm:8000/v1/chat/completions -chromedriver=/path/to/chromedriver -chrome=/path/to/google-chrome -cookie="your_cookie_value" -limit=1000 -ocr=true
```

### Program Flow
1. **Command-line Parameter Parsing**: The program starts by parsing various parameters input from the command line. If `host` or `website` is empty, an error message will be output and the program will terminate.
2. **Directory Handling**: Get the current working directory of the program, and convert the relative paths of `tmpDir` and `resultDir` to absolute paths. If `tmpDir` does not exist, the directory will be created.
3. **Website Crawling**: Use the `run.CrawlerRun` function to crawl the specified website and store the results in `tmpDir`, using `maxCrawledLimit` and `cookie` for crawling restrictions and authentication.
4. **OCR Processing (Optional)**: If the `-ocr` parameter is enabled, use the `run.RunOCR` function to perform a second crawl using Selenium and extract text using OCR.
5. **Keyword Extraction**: Use the `run.ExtractKeywords` function to extract keywords from `tmpDir`, calling the `keywordsLLMURL` large model service.
6. **Large Model Inference and Validation**: Use the `validation.RunValidation` function for large model inference and validation, and store the results in `resultDir`.


### Notes
- Ensure that the required third-party Go packages are correctly installed. For example:
    ```sh
    go mod tidy
    ```
- Ensure that the paths of ChromeDriver and Chrome browser are correct and compatible with the program's environment.
- Ensure that the URL of the large model service is accessible and the service is running.


Through the above usage method, you can configure the program according to different needs to perform operations such as website crawling, keyword extraction, and large model validation.
