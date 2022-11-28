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

	engine.RegistActionRunner("common", &actions.CommonActionRunner{})

	param := make(map[string]interface{})

	data, _ := ioutil.ReadFile(*file)

	flow, err := engine.ParseFlow(string(data))
	if err != nil {
		fmt.Println(err)
		return
	}

	runtime := engine.CreateRuntime(flow, param)

	runner := engine.CommonFlowRunner{}
	engine.Execute(runtime, &runner, *timeout)

	fmt.Println("time used(ms):", runtime.Timeused)

}
