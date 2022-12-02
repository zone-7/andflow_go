package engine

import (
	"sync"
	"time"

	"github.com/zone-7/andflow_go/models"
)

type RuntimeStore interface {
	Init(runtimeId string)
	SetRuntime(runtime *models.RuntimeModel)
	GetRuntime() *models.RuntimeModel
	GetFlow() *models.FlowModel

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
	RuntimeId string
	Wg        sync.WaitGroup //同步控制
	Runtime   *models.RuntimeModel
}

func (s *CommonRuntimeStore) Init(runtimeId string) {
	s.RuntimeId = runtimeId
	s.Wg = sync.WaitGroup{}
}
func (s *CommonRuntimeStore) SetRuntime(runtime *models.RuntimeModel) {
	s.Runtime = runtime
}

func (s *CommonRuntimeStore) RefreshState() {
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

func (s *CommonRuntimeStore) AddLog(tp, tag, title, content string) {
	s.Runtime.AddLog("flow", "stop", "中断", "执行中断")
}

func (s *CommonRuntimeStore) GetRuntime() *models.RuntimeModel {
	return s.Runtime
}
func (s *CommonRuntimeStore) GetFlow() *models.FlowModel {
	return s.Runtime.Flow
}

func (s *CommonRuntimeStore) SetBegin() {
	s.Runtime.BeginTime = time.Now()
}

func (s *CommonRuntimeStore) GetRunningActions() []*models.ActionParam {
	return s.Runtime.RunningActions
}

func (s *CommonRuntimeStore) GetRunningLinks() []*models.LinkParam {
	return s.Runtime.RunningLinks
}

func (s *CommonRuntimeStore) AddRunningAction(param *models.ActionParam) {
	s.Runtime.AddRunningAction(param)
}

func (s *CommonRuntimeStore) DelRunningAction(param *models.ActionParam) {
	s.Runtime.DelRunningAction(param)
}

func (s *CommonRuntimeStore) AddRunningLink(param *models.LinkParam) {
	s.Runtime.AddRunningLink(param)
}
func (s *CommonRuntimeStore) DelRunningLink(param *models.LinkParam) {
	s.Runtime.DelRunningLink(param)
}

func (s *CommonRuntimeStore) AddActionState(state *models.ActionStateModel) {
	s.Runtime.AddActionState(state)
}
func (s *CommonRuntimeStore) AddLinkState(state *models.LinkStateModel) {
	s.Runtime.AddLinkState(state)
}

func (s *CommonRuntimeStore) GetLastActionState(actionId string) *models.ActionStateModel {
	return s.Runtime.GetLastActionState(actionId)
}

func (s *CommonRuntimeStore) GetLastLinkState(sourceId string, targetId string) *models.LinkStateModel {
	return s.Runtime.GetLastLinkState(sourceId, targetId)
}

func (s *CommonRuntimeStore) SetParam(key string, val interface{}) {
	s.Runtime.SetParam(key, val)
}
func (s *CommonRuntimeStore) GetParam(key string) interface{} {
	return s.Runtime.GetParam(key)
}
func (s *CommonRuntimeStore) GetParamMap() map[string]interface{} {
	return s.Runtime.GetParamMap()
}
func (s *CommonRuntimeStore) SetData(key string, val interface{}) {
	s.Runtime.SetData(key, val)
}
func (s *CommonRuntimeStore) GetData(key string) interface{} {
	return s.Runtime.GetData(key)
}
func (s *CommonRuntimeStore) GetDataMap() map[string]interface{} {
	return s.Runtime.GetDataMap()
}

func (s *CommonRuntimeStore) SetActionData(actionId string, name string, val interface{}) {
	s.Runtime.SetActionData(actionId, name, val)
}
func (s *CommonRuntimeStore) GetActionData(actionId string, name string) interface{} {
	return s.Runtime.GetActionData(actionId, name)
}

func (s *CommonRuntimeStore) GetActionDataMap(actionId string) map[string]interface{} {
	return s.Runtime.GetActionDataMap(actionId)
}
