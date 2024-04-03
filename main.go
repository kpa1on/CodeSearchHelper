package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

func main() {
	// 解析命令行参数
	folderPtr := flag.String("f", "", "文件夹路径")
	containsPtr := flag.Bool("c", false, "该参数表示过滤出包含指定字符串的文件的相对路径，不加则过滤出不包含指定字符串的文件的相对路径")
	regexPtr := flag.Bool("r", false, "开启正则表达式搜索模式")
	searchPtr := flag.String("k", "", "要搜索的字符串")
	outputPtr := flag.String("o", "", "运行结果的保存路径，使用的是相对路径")
	extPtr := flag.String("e", "", "指定搜索的文件后缀，多个后缀用逗号分开")
	helpPtr := flag.Bool("h", false, "显示帮助")

	flag.Parse()

	// 如果传入-h参数，则显示帮助信息
	if *helpPtr {
		flag.PrintDefaults()
		fmt.Println("使用方式:")
		fmt.Println("go run main.go -f <文件夹路径> [-c] -k <要搜索的字符串> [-r] [-o <输出文件路径>] [-e <文件后缀:jsp,html...>] [-h]")
		fmt.Println("示例:")
		fmt.Println("go run main.go -f /path/to/folder -c  -k 'exec' -e php -o output.txt")
		return
	}

	// 检查是否指定了文件夹路径
	if *folderPtr == "" {
		flag.PrintDefaults()
		fmt.Println("使用方式:")
		fmt.Println("go run main.go -f <文件夹路径> [-c] -k <要搜索的字符串> [-r] [-o <输出文件路径>] [-e <文件后缀:jsp,html...>] [-h]")
		fmt.Println("示例:")
		fmt.Println("go run main.go -f /path/to/folder -c  -k 'exec' -e php -o output.txt")
		return
	}

	// 递归遍历文件夹
	var mu sync.Mutex
	var wg sync.WaitGroup
	fileCh := make(chan string)
	resultCh := make(chan string)
	var fileContent string
	// 开启多个 goroutine 处理文件
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range fileCh {
				if matchExtensions(filepath.Ext(filePath), *extPtr) {

					if containsInFile(filePath, *searchPtr, *regexPtr) {
						filePath, _ := filepath.Rel(*folderPtr, filePath)
						filePath = strings.ReplaceAll(filePath, "\\", "/")
						if *containsPtr {
							mu.Lock()
							fileContent += filePath + "\n"
							mu.Unlock()
							resultCh <- filePath
						}
					} else {
						filePath, _ := filepath.Rel(*folderPtr, filePath)
						filePath = strings.ReplaceAll(filePath, "\\", "/")
						if !*containsPtr {
							mu.Lock()
							fileContent += filePath + "\n"
							mu.Unlock()
							resultCh <- filePath
						}
					}
				}
			}
		}()
	}

	go func() {
		defer close(fileCh)
		filepath.Walk(*folderPtr, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("遍历文件夹出错:", err)
				return err
			}
			if !info.IsDir() {
				fileCh <- path
			}
			return nil
		})
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// 输出匹配的文件路径到控制台或保存到指定文件中
	for path := range resultCh {
		fmt.Println(path)
	}

	// 保存匹配的文件路径到输出文件
	if *outputPtr != "" {
		absPath, fileErr := filepath.Abs(*outputPtr)
		if fileErr != nil {
			fmt.Println("无法获取绝对路径:", fileErr)
			return
		}

		*outputPtr = absPath
		err := ioutil.WriteFile(*outputPtr, []byte(fmt.Sprintf("%v\n", fileContent)), 0644)
		if err != nil {
			fmt.Println("无法保存输出文件:", err)
			return
		}
		fmt.Println("匹配文件路径已保存到", *outputPtr)
	}
}

// matchExtensions 检查文件的后缀是否匹配给定的后缀列表
func matchExtensions(fileExt, extList string) bool {
	if extList == "" {
		return true
	}
	exts := strings.Split(extList, ",")
	for _, ext := range exts {
		if strings.TrimSpace(fileExt) == "."+strings.TrimSpace(ext) {
			return true
		}
	}
	return false
}

// containsInFile 判断文件的内容是否包含搜索字符串
func containsInFile(filePath, search string, regex bool) bool {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("读取文件出错:", err)
		return false
	}

	if regex {
		matched, err := regexp.MatchString(search, string(fileContent))
		if err != nil {
			fmt.Println("正则表达式错误:", err)
			return false
		}
		return matched
	}

	return strings.Contains(string(fileContent), search)
}
