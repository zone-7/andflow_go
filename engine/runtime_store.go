package engine

import (
	"sync"
	"time"
)

type RuntimeStore interface {
	Init(runtimeId string)
	SetRuntime(runtime *RuntimeModel)
	GetRuntime() *RuntimeModel
	GetFlow() *FlowModel
	SetCmd(c int)
	GetCmd() int
	WaitAdd(int)
	WaitDone()
	Wait()

	Save()

	AddLog(tp, tag, title, content string)

	RefreshState()

	SetBegin()
	SetEnd()
	SetState(state int)
	GetState() int
	GetRunningActions() []*ActionParam
	GetRunningLinks() []*LinkParam
	AddRunningAction(param *ActionParam)
	DelRunningAction(param *ActionParam)
	AddRunningLink(param *LinkParam)
	DelRunningLink(param *LinkParam)

	AddActionState(state *ActionStateModel)
	AddLinkState(state *LinkStateModel)

	GetLastActionState(actionId string) *ActionStateModel
	GetLastLinkState(sourceId string, targetId string) *LinkStateModel

	SetParam(key string, val interface{})
	GetParam(key string) interface{}
	GetParamMap() map[string]interface{}
	SetData(key string, val interface{})
	GetData(key string) interface{}
	GetDataMap() map[string]interface{}

	SetActionData(actionId string, name string, val interface{})
	GetActionData(actionId string, name string) interface{}
	GetActionDataMap(actionId string) map[string]interface{}
	SetActionIcon(actionId string, icon string)
}

type CommonRuntimeStore struct {
	Cmd          int
	RuntimeId    string
	Wg           sync.WaitGroup //同步控制
	Runtime      *RuntimeModel
	OnChangeFunc func(event string, runtime *RuntimeModel)
	OnSaveFunc   func(runtime *RuntimeModel)
}

func (s *CommonRuntimeStore) Init(runtimeId string) {

	s.RuntimeId = runtimeId
	s.Wg = sync.WaitGroup{}
}
func (s *CommonRuntimeStore) RefreshState() {
	if s.Runtime == nil {
		return
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

func (s *CommonRuntimeStore) AddLog(tp, tag, title, content string) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddLog(tp, tag, title, content)
}

func (s *CommonRuntimeStore) SetRuntime(runtime *RuntimeModel) {
	s.Runtime = runtime
}

func (s *CommonRuntimeStore) GetRuntime() *RuntimeModel {
	return s.Runtime
}
func (s *CommonRuntimeStore) GetFlow() *FlowModel {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.Flow
}

func (s *CommonRuntimeStore) SetBegin() {
	if s.Runtime == nil {
		return
	}
	s.Runtime.IsRunning = 1
	s.Runtime.BeginTime = time.Now()
	if s.OnChangeFunc != nil {
		s.OnChangeFunc("begin", s.Runtime)
	}
}
func (s *CommonRuntimeStore) SetEnd() {
	if s.Runtime == nil {
		return
	}
	s.Runtime.IsRunning = 0
	s.Runtime.EndTime = time.Now()
	s.Runtime.Timeused = s.Runtime.EndTime.Sub(s.Runtime.BeginTime).Milliseconds()

	if s.OnChangeFunc != nil {
		s.OnChangeFunc("end", s.Runtime)
	}
}

func (s *CommonRuntimeStore) SetState(state int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.FlowState = state
	if s.OnChangeFunc != nil {
		s.OnChangeFunc("state", s.Runtime)
	}
}

func (s *CommonRuntimeStore) GetState() int {
	if s.Runtime == nil {
		return 0
	}
	return s.Runtime.FlowState
}
func (s *CommonRuntimeStore) GetRunningActions() []*ActionParam {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.RunningActions
}

func (s *CommonRuntimeStore) GetRunningLinks() []*LinkParam {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.RunningLinks
}

func (s *CommonRuntimeStore) AddRunningAction(param *ActionParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddRunningAction(param)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc("action_running_add", s.Runtime)
	}
}

func (s *CommonRuntimeStore) DelRunningAction(param *ActionParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.DelRunningAction(param)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc("action_running_del", s.Runtime)
	}
}

func (s *CommonRuntimeStore) AddRunningLink(param *LinkParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddRunningLink(param)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc("link_running_add", s.Runtime)
	}
}
func (s *CommonRuntimeStore) DelRunningLink(param *LinkParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.DelRunningLink(param)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc("link_running_del", s.Runtime)
	}
}

func (s *CommonRuntimeStore) AddActionState(state *ActionStateModel) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddActionState(state)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc("action_state_add", s.Runtime)
	}
}
func (s *CommonRuntimeStore) AddLinkState(state *LinkStateModel) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddLinkState(state)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc("link_state_add", s.Runtime)
	}
}

func (s *CommonRuntimeStore) GetLastActionState(actionId string) *ActionStateModel {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetLastActionState(actionId)
}

func (s *CommonRuntimeStore) GetLastLinkState(sourceId string, targetId string) *LinkStateModel {
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
	if s.OnChangeFunc != nil {
		s.OnChangeFunc("param_set", s.Runtime)
	}
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
	if s.OnChangeFunc != nil {
		s.OnChangeFunc("data_set", s.Runtime)
	}
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
	if s.OnChangeFunc != nil {
		s.OnChangeFunc("action_data_set", s.Runtime)
	}
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

func (s *CommonRuntimeStore) SetActionIcon(actionId string, icon string) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetActionIcon(actionId, icon)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc("action_icon_set", s.Runtime)
	}
}

func (s *CommonRuntimeStore) Save() {
	if s.OnSaveFunc != nil {
		s.OnSaveFunc(s.Runtime)
	}
}

func (s *CommonRuntimeStore) SetOnChangeFunc(f func(event string, runtime *RuntimeModel)) {
	s.OnChangeFunc = f
}

func (s *CommonRuntimeStore) SetOnSaveFunc(f func(runtime *RuntimeModel)) {
	s.OnSaveFunc = f
}
