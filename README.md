# CodeSearchHelper
在指定文件夹中搜索包含/不包含某个字符串的文件，并输出它的相对路径，支持正则表达式
# 编译
go build
# 使用
➜ GOLAND  go run .\main.go                                                                        
  -c    该参数表示过滤出包含指定字符串的文件的相对路径，不加则过滤出不包含指定字符串的文件的相对路径
  -e string
        指定搜索的文件后缀，多个后缀用逗号分开
  -f string
        文件夹路径
  -h    显示帮助
  -k string
        要搜索的字符串
  -o string
        运行结果的保存路径，使用的是相对路径
  -r    开启正则表达式搜索模式
使用方式:
go run main.go -f <文件夹路径> [-c] -k <要搜索的字符串> [-r] [-o <输出文件路径>] [-e <文件后缀:jsp,html...>] [-h]
示例:
1、搜索vuln_code文件夹下，所有不包含auth.php字符串的后缀为php的文件，并输出相对路径到output.txt
go run main.go -f /vuln_code/  -k 'auth.php' -e php -o output.txt

2、搜索vuln_code文件夹下，所有包含exec字符串的后缀为php的文件，并输出相对路径到output.txt
go run main.go -f /vuln_code/  -k 'exec' -c -e php -o output.txt

3、开启正则搜索
go run main.go -f /vuln_code/  -k '(exec|Runtime)' -c -r -e php -o output.txt
