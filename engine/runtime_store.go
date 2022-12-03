package engine

import (
	"sync"
	"time"

	"github.com/zone-7/andflow_go/models"
)

type RuntimeStore interface {
	Init(runtimeId string, timeout int64)
	SetRuntime(runtime *models.RuntimeModel)
	GetRuntime() *models.RuntimeModel
	GetFlow() *models.FlowModel
	SetTimeout(timeout int64)
	GetTimeout() int64
	SetCmd(c int)
	GetCmd() int
	WaitAdd(int)
	WaitDone()
	Wait()

	AddLog(tp, tag, title, content string)

	RefreshState()

	SetBegin()

	GetRunningActions() []*models.ActionParam
	GetRunningLinks() []*models.LinkParam
	AddRunningAction(param *models.ActionParam)
	DelRunningAction(param *models.ActionParam)
	AddRunningLink(param *models.LinkParam)
	DelRunningLink(param *models.LinkParam)

	AddActionState(state *models.ActionStateModel)
	AddLinkState(state *models.LinkStateModel)

	GetLastActionState(actionId string) *models.ActionStateModel
	GetLastLinkState(sourceId string, targetId string) *models.LinkStateModel

	SetParam(key string, val interface{})
	GetParam(key string) interface{}
	GetParamMap() map[string]interface{}
	SetData(key string, val interface{})
	GetData(key string) interface{}
	GetDataMap() map[string]interface{}

	SetActionData(actionId string, name string, val interface{})
	GetActionData(actionId string, name string) interface{}
	GetActionDataMap(actionId string) map[string]interface{}
}

type CommonRuntimeStore struct {
	Timeout   int64
	Cmd       int
	RuntimeId string
	Wg        sync.WaitGroup //同步控制
	Runtime   *models.RuntimeModel
}

func (s *CommonRuntimeStore) Init(runtimeId string, timeout int64) {
	s.Timeout = timeout
	s.RuntimeId = runtimeId
	s.Wg = sync.WaitGroup{}
}
func (s *CommonRuntimeStore) RefreshState() {
	if s.Runtime == nil {
		return
	}
	//如果没有什么待办事项就表示执行完了
	if (s.Runtime.RunningActions == nil || len(s.Runtime.RunningActions) == 0) &&
		(s.Runtime.RunningLinks == nil || len(s.Runtime.RunningLinks) == 0) {
		s.Runtime.FlowState = 1
		s.Runtime.EndTime = time.Now()
		s.Runtime.Timeused = s.Runtime.EndTime.Sub(s.Runtime.BeginTime).Milliseconds()
	}
}

func (s *CommonRuntimeStore) WaitAdd(d int) {
	s.Wg.Add(d)
}

func (s *CommonRuntimeStore) WaitDone() {
	s.Wg.Done()
}
func (s *CommonRuntimeStore) Wait() {
	s.Wg.Wait()
}
func (s *CommonRuntimeStore) SetCmd(c int) {
	s.Cmd = c
}
func (s *CommonRuntimeStore) GetCmd() int {
	return s.Cmd
}
func (s *CommonRuntimeStore) SetTimeout(timeout int64) {
	s.Timeout = timeout
}
func (s *CommonRuntimeStore) GetTimeout() int64 {
	return s.Timeout
}

func (s *CommonRuntimeStore) AddLog(tp, tag, title, content string) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddLog("flow", "stop", "中断", "执行中断")
}

func (s *CommonRuntimeStore) SetRuntime(runtime *models.RuntimeModel) {
	s.Runtime = runtime
}

func (s *CommonRuntimeStore) GetRuntime() *models.RuntimeModel {
	return s.Runtime
}
func (s *CommonRuntimeStore) GetFlow() *models.FlowModel {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.Flow
}

func (s *CommonRuntimeStore) SetBegin() {
	if s.Runtime == nil {
		return
	}
	s.Runtime.BeginTime = time.Now()
}

func (s *CommonRuntimeStore) GetRunningActions() []*models.ActionParam {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.RunningActions
}

func (s *CommonRuntimeStore) GetRunningLinks() []*models.LinkParam {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.RunningLinks
}

func (s *CommonRuntimeStore) AddRunningAction(param *models.ActionParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddRunningAction(param)
}

func (s *CommonRuntimeStore) DelRunningAction(param *models.ActionParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.DelRunningAction(param)
}

func (s *CommonRuntimeStore) AddRunningLink(param *models.LinkParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddRunningLink(param)
}
func (s *CommonRuntimeStore) DelRunningLink(param *models.LinkParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.DelRunningLink(param)
}

func (s *CommonRuntimeStore) AddActionState(state *models.ActionStateModel) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddActionState(state)
}
func (s *CommonRuntimeStore) AddLinkState(state *models.LinkStateModel) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddLinkState(state)
}

func (s *CommonRuntimeStore) GetLastActionState(actionId string) *models.ActionStateModel {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetLastActionState(actionId)
}

func (s *CommonRuntimeStore) GetLastLinkState(sourceId string, targetId string) *models.LinkStateModel {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetLastLinkState(sourceId, targetId)
}

func (s *CommonRuntimeStore) SetParam(key string, val interface{}) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetParam(key, val)
}
func (s *CommonRuntimeStore) GetParam(key string) interface{} {
	return s.Runtime.GetParam(key)
}
func (s *CommonRuntimeStore) GetParamMap() map[string]interface{} {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetParamMap()
}
func (s *CommonRuntimeStore) SetData(key string, val interface{}) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetData(key, val)
}
func (s *CommonRuntimeStore) GetData(key string) interface{} {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetData(key)
}
func (s *CommonRuntimeStore) GetDataMap() map[string]interface{} {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetDataMap()
}

func (s *CommonRuntimeStore) SetActionData(actionId string, name string, val interface{}) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetActionData(actionId, name, val)
}
func (s *CommonRuntimeStore) GetActionData(actionId string, name string) interface{} {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetActionData(actionId, name)
}

func (s *CommonRuntimeStore) GetActionDataMap(actionId string) map[string]interface{} {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetActionDataMap(actionId)
}
