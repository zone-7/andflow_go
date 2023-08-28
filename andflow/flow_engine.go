package andflow

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

var runnings map[string]*Session = make(map[string]*Session)

func ParseFlow(content string) (*FlowModel, error) {

	flowModel := FlowModel{}
	err := json.Unmarshal([]byte(content), &flowModel)
	if err != nil {
		return nil, err
	}
	return &flowModel, nil
}

// 创建运行时
func CreateRuntime(flow *FlowModel, param map[string]interface{}) *RuntimeModel {
	runtime := RuntimeModel{}
	uid, _ := uuid.NewV4()
	id := strings.ReplaceAll(uid.String(), "-", "")
	runtime.Id = id
	runtime.Flow = flow
	runtime.FlowState = 0
	runtime.Des = flow.Name
	//初始化状态
	if runtime.ActionStates == nil {
		runtime.ActionStates = make([]*ActionStateModel, 0)
	}
	//初始化下一步执行的连接线
	if runtime.RunningLinks == nil {
		runtime.RunningLinks = make([]*LinkParam, 0)
	}
	if runtime.RunningActions == nil {
		runtime.RunningActions = make([]*ActionParam, 0)
	}

	//初始化日志
	if runtime.Logs == nil {
		runtime.Logs = make([]*LogModel, 0)
	}

	//初始化数据
	if runtime.Data == nil {
		runtime.Data = make([]*RuntimeDataModel, 0)
	}
	//初始化参数
	if param != nil {
		for k, v := range param {
			runtime.SetParam(k, v)
		}
	}

	runtime.CreateTime = time.Now()
	return &runtime
}

func GetSessions() map[string]*Session {
	return runnings
}
func GetSession(runtimeId string) *Session {
	return runnings[runtimeId]
}

func Execute(controller RuntimeOperation, router FlowRouter, runner FlowRunner, timeout int64) {

	ctx := context.Background()
	if timeout > 0 {
		ctx, _ = context.WithTimeout(ctx, time.Millisecond*time.Duration(timeout))
	}

	session := CreateSession(ctx, controller, router, runner)

	runnings[session.Id] = session
	defer delete(runnings, session.Id)

	session.Execute()

}

func ExecuteRuntime(runtime *RuntimeModel, timeout int64) *RuntimeModel {
	runner := &CommonFlowRunner{}
	router := &CommonFlowRouter{}
	operation := &CommonRuntimeOperation{}
	operation.Init(runtime)
	// controller.SetRuntime(runtime)

	Execute(operation, router, runner, timeout)
	runtime = operation.GetRuntime()

	return runtime

}

func ExecuteFlow(flow *FlowModel, param map[string]interface{}, timeout int64) *RuntimeModel {
	runtime := CreateRuntime(flow, param)

	runner := &CommonFlowRunner{}
	router := &CommonFlowRouter{}
	operation := &CommonRuntimeOperation{}
	operation.Init(runtime)
	Execute(operation, router, runner, timeout)
	runtime = operation.GetRuntime()
	return runtime

}
