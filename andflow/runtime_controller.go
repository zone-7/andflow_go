package andflow

import (
	"sync"
	"time"
)

const (
	//执行过程中的事件
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

	//controller CMD 控制指令
	CMD_STOP  = 1 //停止
	CMD_START = 0 //默认执行
)

type RuntimeController interface {
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

type CommonRuntimeController struct {
	Cmd          int
	RuntimeId    string
	Wg           sync.WaitGroup //同步控制
	Runtime      *RuntimeModel
	OnChangeFunc func(event string, runtime *RuntimeModel)
	OnSaveFunc   func(runtime *RuntimeModel)
}

func (s *CommonRuntimeController) Init(runtimeId string) {

	s.RuntimeId = runtimeId
	s.Wg = sync.WaitGroup{}
}

func (s *CommonRuntimeController) WaitAdd(d int) {
	s.Wg.Add(d)
}

func (s *CommonRuntimeController) WaitDone() {
	s.Wg.Done()
}
func (s *CommonRuntimeController) Wait() {
	s.Wg.Wait()
}
func (s *CommonRuntimeController) SetCmd(c int) {
	s.Cmd = c
}
func (s *CommonRuntimeController) GetCmd() int {
	return s.Cmd
}

func (s *CommonRuntimeController) AddLog(tp, tag, name, title, content string) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddLog(tp, tag, name, title, content)
}

func (s *CommonRuntimeController) SetRuntime(runtime *RuntimeModel) {
	s.Runtime = runtime
}

func (s *CommonRuntimeController) GetRuntime() *RuntimeModel {
	return s.Runtime
}
func (s *CommonRuntimeController) GetFlow() *FlowModel {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.Flow
}

func (s *CommonRuntimeController) SetBegin() {
	if s.Runtime == nil {
		return
	}
	s.Runtime.IsRunning = 1
	s.Runtime.BeginTime = time.Now()
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_BEGIN, s.Runtime)
	}
}
func (s *CommonRuntimeController) SetEnd() {
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
func (s *CommonRuntimeController) SetMessage(message string) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.Message = message
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_MESSAGE, s.Runtime)
	}
}

func (s *CommonRuntimeController) GetMessage() string {
	if s.Runtime == nil {
		return ""
	}
	return s.Runtime.Message
}

func (s *CommonRuntimeController) SetError(iserror int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.IsError = iserror
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ISERROR, s.Runtime)
	}
}
func (s *CommonRuntimeController) GetError() int {
	if s.Runtime == nil {
		return 0
	}
	return s.Runtime.IsError
}

func (s *CommonRuntimeController) SetState(state int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.FlowState = state
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_FLOWSTATE, s.Runtime)
	}
}

func (s *CommonRuntimeController) GetState() int {
	if s.Runtime == nil {
		return 0
	}
	return s.Runtime.FlowState
}
func (s *CommonRuntimeController) GetRunningActions() []*ActionParam {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.RunningActions
}

func (s *CommonRuntimeController) GetRunningLinks() []*LinkParam {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.RunningLinks
}

func (s *CommonRuntimeController) AddRunningAction(param *ActionParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddRunningAction(param)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_RUNNING_ADD, s.Runtime)
	}
}

func (s *CommonRuntimeController) DelRunningAction(param *ActionParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.DelRunningAction(param)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_RUNNING_DEL, s.Runtime)
	}
}

func (s *CommonRuntimeController) AddRunningLink(param *LinkParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddRunningLink(param)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_LINK_RUNNING_ADD, s.Runtime)
	}
}
func (s *CommonRuntimeController) DelRunningLink(param *LinkParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.DelRunningLink(param)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_LINK_RUNNING_DEL, s.Runtime)
	}
}

func (s *CommonRuntimeController) AddActionState(state *ActionStateModel) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddActionState(state)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_STATE_ADD, s.Runtime)
	}
}
func (s *CommonRuntimeController) AddLinkState(state *LinkStateModel) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddLinkState(state)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_LINK_STATE_ADD, s.Runtime)
	}
}

func (s *CommonRuntimeController) GetLastActionState(actionId string) *ActionStateModel {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetLastActionState(actionId)
}

func (s *CommonRuntimeController) GetLastLinkState(sourceId string, targetId string) *LinkStateModel {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetLastLinkState(sourceId, targetId)
}

func (s *CommonRuntimeController) SetParam(key string, val interface{}) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetParam(key, val)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_PARAM_SET, s.Runtime)
	}
}
func (s *CommonRuntimeController) GetParam(key string) interface{} {
	return s.Runtime.GetParam(key)
}
func (s *CommonRuntimeController) GetParamMap() map[string]interface{} {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetParamMap()
}
func (s *CommonRuntimeController) SetData(key string, val interface{}) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetData(key, val)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_DATA_SET, s.Runtime)
	}
}
func (s *CommonRuntimeController) GetData(key string) interface{} {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetData(key)
}
func (s *CommonRuntimeController) GetDataMap() map[string]interface{} {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetDataMap()
}

func (s *CommonRuntimeController) SetActionData(actionId string, name string, val interface{}) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetActionData(actionId, name, val)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_DATA_SET, s.Runtime)
	}
}
func (s *CommonRuntimeController) GetActionData(actionId string, name string) interface{} {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetActionData(actionId, name)
}

func (s *CommonRuntimeController) GetActionDataMap(actionId string) map[string]interface{} {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetActionDataMap(actionId)
}

func (s *CommonRuntimeController) SetActionIcon(actionId string, icon string) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetActionIcon(actionId, icon)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_ICON_SET, s.Runtime)
	}
}
func (s *CommonRuntimeController) SetActionState(actionId string, state int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetActionState(actionId, state)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_STATE_SET, s.Runtime)
	}
}
func (s *CommonRuntimeController) SetActionError(actionId string, isError int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetActionError(actionId, isError)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_ERROR_SET, s.Runtime)
	}
}

func (s *CommonRuntimeController) SetLinkState(sourceId, targetId string, state int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetLinkState(sourceId, targetId, state)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_LINK_STATE_SET, s.Runtime)
	}
}
func (s *CommonRuntimeController) SetLinkError(sourceId, targetId string, isError int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetLinkError(sourceId, targetId, isError)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_LINK_ERROR_SET, s.Runtime)
	}
}

func (s *CommonRuntimeController) Save() {
	if s.OnSaveFunc != nil {
		s.OnSaveFunc(s.Runtime)
	}
}

func (s *CommonRuntimeController) SetOnChangeFunc(f func(event string, runtime *RuntimeModel)) {
	s.OnChangeFunc = f
}

func (s *CommonRuntimeController) SetOnSaveFunc(f func(runtime *RuntimeModel)) {
	s.OnSaveFunc = f
}
