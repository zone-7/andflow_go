package engine

import (
	"sync"
	"time"
)

const (
	EVENT_BEGIN     = "begin"
	EVENT_END       = "end"
	EVENT_MESSAGE   = "message"
	EVENT_FLOWSTATE = "flowstate"
	EVENT_ISERROR   = "iserror"

	EVENT_ACTION_RUNNING_ADD = "action_running_add"
	EVENT_ACTION_RUNNING_DEL = "action_running_del"

	EVENT_ACTION_STATE_ADD = "action_state_add"
	EVENT_ACTION_STATE_SET = "action_state_set"
	EVENT_ACTION_ERROR_SET = "action_error_set"

	EVENT_ACTION_DATA_SET = "action_data_set"
	EVENT_ACTION_ICON_SET = "action_icon_set"

	EVENT_LINK_RUNNING_ADD = "link_running_add"
	EVENT_LINK_RUNNING_DEL = "link_running_del"

	EVENT_LINK_STATE_ADD = "link_state_add"
	EVENT_LINK_STATE_SET = "link_state_set"
	EVENT_LINK_ERROR_SET = "link_error_set"

	EVENT_PARAM_SET = "param_set"
	EVENT_DATA_SET  = "data_set"
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
	AddLog(tp, tag, name, title, content string)
	SetBegin()
	SetEnd()
	SetState(state int)
	GetState() int
	SetError(iserror int)
	GetError() int
	SetMessage(message string)
	GetMessage() string
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
	SetActionError(actionId string, isError int)
	SetActionState(actionId string, state int)

	SetLinkError(sourceId string, targetId string, isError int)
	SetLinkState(sourceId string, targetId string, state int)
}

type CommonRuntimeStore struct {
	Cmd          int
	RuntimeId    string
	Wg           sync.WaitGroup //????????????
	Runtime      *RuntimeModel
	OnChangeFunc func(event string, runtime *RuntimeModel)
	OnSaveFunc   func(runtime *RuntimeModel)
}

func (s *CommonRuntimeStore) Init(runtimeId string) {

	s.RuntimeId = runtimeId
	s.Wg = sync.WaitGroup{}
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

func (s *CommonRuntimeStore) AddLog(tp, tag, name, title, content string) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddLog(tp, tag, name, title, content)
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
		s.OnChangeFunc(EVENT_BEGIN, s.Runtime)
	}
}
func (s *CommonRuntimeStore) SetEnd() {
	if s.Runtime == nil {
		return
	}
	s.Runtime.IsRunning = 0
	s.Runtime.EndTime = time.Now()
	s.Runtime.Timeused = s.Runtime.EndTime.Sub(s.Runtime.BeginTime).Milliseconds()

	for _, ast := range s.Runtime.ActionStates {
		if ast.IsError == 1 {
			s.Runtime.IsError = 1
			break
		}
	}
	for _, lst := range s.Runtime.LinkStates {
		if lst.IsError == 1 {
			s.Runtime.IsError = 1
			break
		}
	}

	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_END, s.Runtime)
	}
}
func (s *CommonRuntimeStore) SetMessage(message string) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.Message = message
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_MESSAGE, s.Runtime)
	}
}

func (s *CommonRuntimeStore) GetMessage() string {
	if s.Runtime == nil {
		return ""
	}
	return s.Runtime.Message
}

func (s *CommonRuntimeStore) SetError(iserror int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.IsError = iserror
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ISERROR, s.Runtime)
	}
}
func (s *CommonRuntimeStore) GetError() int {
	if s.Runtime == nil {
		return 0
	}
	return s.Runtime.IsError
}

func (s *CommonRuntimeStore) SetState(state int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.FlowState = state
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_FLOWSTATE, s.Runtime)
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
		s.OnChangeFunc(EVENT_ACTION_RUNNING_ADD, s.Runtime)
	}
}

func (s *CommonRuntimeStore) DelRunningAction(param *ActionParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.DelRunningAction(param)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_RUNNING_DEL, s.Runtime)
	}
}

func (s *CommonRuntimeStore) AddRunningLink(param *LinkParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddRunningLink(param)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_LINK_RUNNING_ADD, s.Runtime)
	}
}
func (s *CommonRuntimeStore) DelRunningLink(param *LinkParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.DelRunningLink(param)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_LINK_RUNNING_DEL, s.Runtime)
	}
}

func (s *CommonRuntimeStore) AddActionState(state *ActionStateModel) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddActionState(state)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_STATE_ADD, s.Runtime)
	}
}
func (s *CommonRuntimeStore) AddLinkState(state *LinkStateModel) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddLinkState(state)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_LINK_STATE_ADD, s.Runtime)
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
		s.OnChangeFunc(EVENT_PARAM_SET, s.Runtime)
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
		s.OnChangeFunc(EVENT_DATA_SET, s.Runtime)
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
		s.OnChangeFunc(EVENT_ACTION_DATA_SET, s.Runtime)
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
		s.OnChangeFunc(EVENT_ACTION_ICON_SET, s.Runtime)
	}
}
func (s *CommonRuntimeStore) SetActionState(actionId string, state int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetActionState(actionId, state)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_STATE_SET, s.Runtime)
	}
}
func (s *CommonRuntimeStore) SetActionError(actionId string, isError int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetActionError(actionId, isError)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_ERROR_SET, s.Runtime)
	}
}

func (s *CommonRuntimeStore) SetLinkState(sourceId, targetId string, state int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetLinkState(sourceId, targetId, state)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_LINK_STATE_SET, s.Runtime)
	}
}
func (s *CommonRuntimeStore) SetLinkError(sourceId, targetId string, isError int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetLinkError(sourceId, targetId, isError)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_LINK_ERROR_SET, s.Runtime)
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
