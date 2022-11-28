package test

import (
	"fmt"
	"io/ioutil"
	"path"
	"testing"
	"time"

	"github.com/zone-7/andflow_go/actions"
	"github.com/zone-7/andflow_go/engine"
	"github.com/zone-7/andflow_go/models"
)

var demo_path = "./demos"
var i int = 0

type MyRunner1 struct {
}
type MyRunner2 struct {
}

func (r *MyRunner1) ExecuteLink(s *engine.Session, param *models.LinkParam) int {
	link := s.Runtime.Flow.GetLinkBySourceIdAndTargetId(param.SourceId, param.TargetId)

	fmt.Println("执行link:", link.SourceId, "->", link.TargetId)
	return 1
}

func (r *MyRunner1) ExecuteAction(s *engine.Session, param *models.ActionParam) int {

	action := s.Runtime.Flow.GetAction(param.ActionId)
	// time.Sleep(time.Millisecond * 100)
	fmt.Println("执行action:", action.Id, " ", action.Title)

	return 1
}

func (r *MyRunner2) ExecuteLink(s *engine.Session, param *models.LinkParam) int {
	link := s.Runtime.Flow.GetLinkBySourceIdAndTargetId(param.SourceId, param.TargetId)
	n := s.Runtime.GetParam("name")
	if n == nil {
		return 0
	}

	fmt.Println("执行link:", link.SourceId, "->", link.TargetId)
	return 1
}

func (r *MyRunner2) ExecuteAction(s *engine.Session, param *models.ActionParam) int {

	action := s.Runtime.Flow.GetAction(param.ActionId)
	time.Sleep(time.Millisecond * 2000)
	fmt.Println("执行action:", action.Id, " ", action.Title)

	return 1
}

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
		runtime := engine.CreateRuntime(flow, param)

		runner := MyRunner1{}

		t1 := time.Now()

		engine.Execute(runtime, &runner, 10000)

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

	runner := MyRunner2{}

	step := 0

	engine.Execute(runtime, &runner, 10000)
	fmt.Println("执行步骤", step)
	step = step + 1

	runtime.SetParam("name", "zgq")

	for runtime.FlowState == 0 {
		engine.Execute(runtime, &runner, 10000)
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

	runner := MyRunner2{}

	runtime.SetParam("name", "zgq")

	engine.Execute(runtime, &runner, 2000)

	fmt.Println("time used(ms):", runtime.Timeused)
	fmt.Println("--------------------------------------------")

}

//测试通用流程处理器
func TestFlowRunner(t *testing.T) {

	engine.RegistActionRunner("common", &actions.CommonActionRunner{})

	param := make(map[string]interface{})

	data, _ := ioutil.ReadFile(demo_path + "/4执行脚本.json")

	flow, err := engine.ParseFlow(string(data))
	if err != nil {
		fmt.Println(err)
	}

	runtime := engine.CreateRuntime(flow, param)

	runner := engine.CommonFlowRunner{}

	engine.Execute(runtime, &runner, 10000)

	runtime.SetParam("name", "zgq")

	engine.Execute(runtime, &runner, 10000)

	fmt.Println("time used(ms):", runtime.Timeused)
	fmt.Println("--------------------------------------------")

}

//测试在javascript中执行 系统命令
func TestScriptCmd(t *testing.T) {

	engine.RegistActionRunner("common", &actions.CommonActionRunner{})

	param := make(map[string]interface{})

	data, _ := ioutil.ReadFile(demo_path + "/5执行命令.json")

	flow, err := engine.ParseFlow(string(data))
	if err != nil {
		fmt.Println(err)
	}

	runtime := engine.CreateRuntime(flow, param)

	runner := engine.CommonFlowRunner{}

	engine.Execute(runtime, &runner, 10000)

	fmt.Println("time used(ms):", runtime.Timeused)
	fmt.Println("--------------------------------------------")

}
