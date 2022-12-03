package test

import (
	"fmt"
	"io/ioutil"
	"path"
	"testing"
	"time"

	"github.com/zone-7/andflow_go/actions"
	"github.com/zone-7/andflow_go/engine"
)

var demo_path = "./demos"
var i int = 0

//执行所有流程
func TestAll(t *testing.T) {
	param := make(map[string]interface{})

	files, _ := ioutil.ReadDir(demo_path)
	for _, file := range files {
		fmt.Println("执行", file.Name())

		data, _ := ioutil.ReadFile(path.Join(demo_path, file.Name()))

		flow, err := engine.ParseFlow(string(data))
		if err != nil {
			fmt.Println(err)
		}

		t1 := time.Now()

		engine.ExecuteFlow(flow, param, 10000)

		t2 := time.Now().Sub(t1).Milliseconds()
		fmt.Println(file.Name(), "time used(ms):", t2)
		fmt.Println("--------------------------------------------")
	}

}

//测试分步执行
func TestDemo1ByStep(t *testing.T) {
	param := make(map[string]interface{})

	data, _ := ioutil.ReadFile(demo_path + "/1简单流程.json")

	flow, err := engine.ParseFlow(string(data))
	if err != nil {
		fmt.Println(err)
	}
	runtime := engine.CreateRuntime(flow, param)

	step := 0

	engine.ExecuteRuntime(runtime, 10000)
	fmt.Println("执行步骤", step)
	step = step + 1

	runtime.SetParam("name", "zgq")

	for runtime.FlowState == 0 {
		engine.ExecuteRuntime(runtime, 10000)
		fmt.Println("执行步骤", step)
		step = step + 1
	}

	fmt.Println("time used(ms):", runtime.Timeused)
	fmt.Println("--------------------------------------------")

}

//测试超时退出机制
func TestDemo1WithTimeout(t *testing.T) {
	param := make(map[string]interface{})

	data, _ := ioutil.ReadFile(demo_path + "/1简单流程.json")

	flow, err := engine.ParseFlow(string(data))
	if err != nil {
		fmt.Println(err)
	}
	runtime := engine.CreateRuntime(flow, param)

	runtime.SetParam("name", "zgq")

	engine.ExecuteRuntime(runtime, 2000)

	fmt.Println("time used(ms):", runtime.Timeused)
	fmt.Println("--------------------------------------------")

}

//测试通用流程处理器
func TestFlowRunner(t *testing.T) {
	var timeout int64 = 1000
	engine.RegistActionRunner("common", &actions.CommonActionRunner{})

	param := make(map[string]interface{})

	data, _ := ioutil.ReadFile(demo_path + "/4执行脚本.json")

	flow, err := engine.ParseFlow(string(data))
	if err != nil {
		fmt.Println(err)
	}

	runtime := engine.CreateRuntime(flow, param)

	engine.ExecuteRuntime(runtime, timeout)

	runtime.SetParam("name", "zgq")

	engine.ExecuteRuntime(runtime, timeout)

	fmt.Println("time used(ms):", runtime.Timeused)
	fmt.Println("--------------------------------------------")

}

//测试在javascript中执行 系统命令
func TestScriptCmd(t *testing.T) {
	var timeout int64 = 1000
	engine.RegistActionRunner("common", &actions.CommonActionRunner{})

	param := make(map[string]interface{})

	data, _ := ioutil.ReadFile(demo_path + "/5执行命令.json")

	flow, err := engine.ParseFlow(string(data))
	if err != nil {
		fmt.Println(err)
	}

	runtime := engine.CreateRuntime(flow, param)

	engine.ExecuteRuntime(runtime, timeout)

	fmt.Println("time used(ms):", runtime.Timeused)
	fmt.Println("--------------------------------------------")

}
