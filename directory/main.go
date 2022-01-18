package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

//cat 命令的实现
func cat(r *bufio.Reader)  {
	for {
		buf, err := r.ReadBytes('\n') //注意是字节
		if err == io.EOF {
			//退出之前将已读的内容输出
			_,_ = fmt.Fprintf(os.Stdout,"%s",buf)
			break
		}
		_, _ = fmt.Fprintf(os.Stdout, "%s", buf)
	}

}

func main() {
	resp, err := http.Get("https://www.liwenzhou.com/")
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read from resp.Body failed, err:%v\n", err)
		return
	}
	fmt.Print(string(body))


/*
	//使用文件操作相关的知识，模拟实现linux平台cat命令的功能
	flag.Parse()   //解析命令行参数
	if flag.NArg() == 0 {
		//如果没有参数默认从标准输入读取内容
		cat(bufio.NewReader(os.Stdout))
	}
	//依次读取每个文件的内容并打印 到终端
	for i:=0;i<flag.NArg();i++ {
		f,err := os.Open(flag.Arg(i))
		if err != nil {
			_,_ = fmt.Fprintf(os.Stdout,"Reading from %s failed,err:%v\n",flag.Arg(i),err)
			continue
		}
		cat(bufio.NewReader(f))
	}
*/


/*
	file, err := os.OpenFile("xx.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer file.Close()
	str := "hello 沙河"
	file.Write([]byte(str))       //写入字节切片数据
	file.WriteString("hello 小王子") //直接写入字符串数据

*/

/*
	//只读方式打开当前目录下的main.go文件
	file,err := os.Open("./main.go")
	if err != nil {
		fmt.Println("open file failed!,err:",err)
		return
	}
	defer file.Close()
	//使用Read方法读取数据
	var temp= make([]byte,30000)
	n,err := file.Read(temp)
	if err == io.EOF {
		fmt.Println("文件读取完了")
		return
	}
	if err != nil {
		fmt.Println("read file failed,err:",err)
		return
	}
	fmt.Printf("读取了%d字节数据\n",n)
	fmt.Println(string(temp[:n]))
*/
}
