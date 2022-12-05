package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/zone-7/andflow_go/actions"
	"github.com/zone-7/andflow_go/engine"
)

func main() {
	file := flag.String("f", "", "流程json文件")
	timeout := flag.Int64("t", 30000, "超时设置默认30s")
	//解析
	flag.Parse()
	if file == nil || len(*file) == 0 {
		flag.Usage()
		return
	}
	//注册执行器
	engine.RegistActionRunner("common", &actions.ScriptActionRunner{})

	param := make(map[string]interface{})

	data, _ := ioutil.ReadFile(*file)

	flow, err := engine.ParseFlow(string(data))
	if err != nil {
		fmt.Println(err)
		return
	}

	runtime := engine.ExecuteFlow(flow, param, *timeout)

	fmt.Println("time used(ms):", runtime.Timeused)

}
