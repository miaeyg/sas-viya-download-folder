# sas-viya-download-folder

![Static Badge](https://img.shields.io/badge/license-MIT-blue)

This project downloads a SAS Content folder to an on-disk folder. 

Code accepts as input a base folder and then downloads all .sas files (can be modified) in base folder and all subfolders recursively.

Thanks to this blog https://communities.sas.com/t5/SAS-Communities-Library/Go-Viya-First-steps-with-Go-language-and-SAS-Viya/ta-p/704659 for providing base code.

### Prerequisites

- [Go](https://golang.org/) 1.22.2 (or later)

### Installation

#### Option 1

- Using `go build`: <br> Clone the project from the GitHub repository. Then,
  from the project root, run the following command:
  ```
  go build -o {executable file name} main.go
  ```