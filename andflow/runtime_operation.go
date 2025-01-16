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

	//operation CMD 控制指令
	CMD_STOP  = 1 //停止
	CMD_START = 0 //执行,默认
)

type RuntimeOperation interface {
	Init(runtime *RuntimeModel)
	GetRuntime() *RuntimeModel
	GetFlow() *FlowModel
	SetCmd(c int)
	GetCmd() int
	WaitAdd(int)
	WaitDone()
	Wait()
	Save()
	AddLog(tp, tag, name, title, content string)
	SetRequestId(id string)
	GetRequestId() string
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

	SetActionData(actionId string, name string, val interface{})
	GetActionData(actionId string, name string) interface{}
	GetActionDataMap(actionId string) map[string]interface{}
	SetActionIcon(actionId string, icon string)
	SetActionError(actionId string, isError int)
	SetActionState(actionId string, state int)

	SetLinkError(sourceId string, targetId string, isError int)
	SetLinkState(sourceId string, targetId string, state int)
}

type CommonRuntimeOperation struct {
	Cmd          int
	Wg           sync.WaitGroup //同步控制
	Runtime      *RuntimeModel
	OnChangeFunc func(event string, runtime *RuntimeModel)
	OnSaveFunc   func(runtime *RuntimeModel)
}

func (s *CommonRuntimeOperation) Init(runtime *RuntimeModel) {
	s.Runtime = runtime
	s.Wg = sync.WaitGroup{}
}

func (s *CommonRuntimeOperation) WaitAdd(d int) {
	s.Wg.Add(d)
}

func (s *CommonRuntimeOperation) WaitDone() {
	s.Wg.Done()
}
func (s *CommonRuntimeOperation) Wait() {
	s.Wg.Wait()
}
func (s *CommonRuntimeOperation) SetCmd(c int) {
	s.Cmd = c
}
func (s *CommonRuntimeOperation) GetCmd() int {
	return s.Cmd
}

func (s *CommonRuntimeOperation) AddLog(tp, tag, name, title, content string) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddLog(tp, tag, name, title, content)
}

func (s *CommonRuntimeOperation) GetRuntime() *RuntimeModel {
	return s.Runtime
}
func (s *CommonRuntimeOperation) GetFlow() *FlowModel {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.Flow
}

func (s *CommonRuntimeOperation) SetRequestId(id string) {
	s.Runtime.RequestId = id
}
func (s *CommonRuntimeOperation) GetRequestId() string {
	return s.Runtime.RequestId
}

func (s *CommonRuntimeOperation) SetBegin() {
	if s.Runtime == nil {
		return
	}
	s.Runtime.IsRunning = 1
	s.Runtime.BeginTime = time.Now()
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_BEGIN, s.Runtime)
	}
}
func (s *CommonRuntimeOperation) SetEnd() {
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
func (s *CommonRuntimeOperation) SetMessage(message string) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.Message = message
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_MESSAGE, s.Runtime)
	}
}

func (s *CommonRuntimeOperation) GetMessage() string {
	if s.Runtime == nil {
		return ""
	}
	return s.Runtime.Message
}

func (s *CommonRuntimeOperation) SetError(iserror int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.IsError = iserror
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ISERROR, s.Runtime)
	}
}
func (s *CommonRuntimeOperation) GetError() int {
	if s.Runtime == nil {
		return 0
	}
	return s.Runtime.IsError
}

func (s *CommonRuntimeOperation) SetState(state int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.FlowState = state
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_FLOWSTATE, s.Runtime)
	}
}

func (s *CommonRuntimeOperation) GetState() int {
	if s.Runtime == nil {
		return 0
	}
	return s.Runtime.FlowState
}
func (s *CommonRuntimeOperation) GetRunningActions() []*ActionParam {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.RunningActions
}

func (s *CommonRuntimeOperation) GetRunningLinks() []*LinkParam {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.RunningLinks
}

func (s *CommonRuntimeOperation) AddRunningAction(param *ActionParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddRunningAction(param)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_RUNNING_ADD, s.Runtime)
	}
}

func (s *CommonRuntimeOperation) DelRunningAction(param *ActionParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.DelRunningAction(param)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_RUNNING_DEL, s.Runtime)
	}
}

func (s *CommonRuntimeOperation) AddRunningLink(param *LinkParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddRunningLink(param)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_LINK_RUNNING_ADD, s.Runtime)
	}
}
func (s *CommonRuntimeOperation) DelRunningLink(param *LinkParam) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.DelRunningLink(param)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_LINK_RUNNING_DEL, s.Runtime)
	}
}

func (s *CommonRuntimeOperation) AddActionState(state *ActionStateModel) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddActionState(state)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_STATE_ADD, s.Runtime)
	}
}
func (s *CommonRuntimeOperation) AddLinkState(state *LinkStateModel) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.AddLinkState(state)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_LINK_STATE_ADD, s.Runtime)
	}
}

func (s *CommonRuntimeOperation) GetLastActionState(actionId string) *ActionStateModel {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetLastActionState(actionId)
}

func (s *CommonRuntimeOperation) GetLastLinkState(sourceId string, targetId string) *LinkStateModel {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetLastLinkState(sourceId, targetId)
}

func (s *CommonRuntimeOperation) SetParam(key string, val interface{}) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetParam(key, val)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_PARAM_SET, s.Runtime)
	}
}
func (s *CommonRuntimeOperation) GetParam(key string) interface{} {
	return s.Runtime.GetParam(key)
}
func (s *CommonRuntimeOperation) GetParamMap() map[string]interface{} {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetParamMap()
}

func (s *CommonRuntimeOperation) SetActionData(actionId string, name string, val interface{}) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetActionData(actionId, name, val)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_DATA_SET, s.Runtime)
	}
}
func (s *CommonRuntimeOperation) GetActionData(actionId string, name string) interface{} {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetActionData(actionId, name)
}

func (s *CommonRuntimeOperation) GetActionDataMap(actionId string) map[string]interface{} {
	if s.Runtime == nil {
		return nil
	}
	return s.Runtime.GetActionDataMap(actionId)
}

func (s *CommonRuntimeOperation) SetActionIcon(actionId string, icon string) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetActionIcon(actionId, icon)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_ICON_SET, s.Runtime)
	}
}
func (s *CommonRuntimeOperation) SetActionState(actionId string, state int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetActionState(actionId, state)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_STATE_SET, s.Runtime)
	}
}
func (s *CommonRuntimeOperation) SetActionError(actionId string, isError int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetActionError(actionId, isError)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_ACTION_ERROR_SET, s.Runtime)
	}
}

func (s *CommonRuntimeOperation) SetLinkState(sourceId, targetId string, state int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetLinkState(sourceId, targetId, state)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_LINK_STATE_SET, s.Runtime)
	}
}
func (s *CommonRuntimeOperation) SetLinkError(sourceId, targetId string, isError int) {
	if s.Runtime == nil {
		return
	}
	s.Runtime.SetLinkError(sourceId, targetId, isError)
	if s.OnChangeFunc != nil {
		s.OnChangeFunc(EVENT_LINK_ERROR_SET, s.Runtime)
	}
}

func (s *CommonRuntimeOperation) Save() {
	if s.OnSaveFunc != nil {
		s.OnSaveFunc(s.Runtime)
	}
}

func (s *CommonRuntimeOperation) SetOnChangeFunc(f func(event string, runtime *RuntimeModel)) {
	s.OnChangeFunc = f
}

func (s *CommonRuntimeOperation) SetOnSaveFunc(f func(runtime *RuntimeModel)) {
	s.OnSaveFunc = f
}
