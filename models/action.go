package models

// import "context"

// type Action interface {
// 	GetName() string
// 	Init(callback interface{})
// 	PrepareMetadata(userid int64, flowCode string, metadata string) string
// 	Filter(ctx context.Context, runtimeId string, preActionId string, actionId string, callback interface{}) (bool, error)
// 	Exec(ctx context.Context, runtimeId string, preActionId string, actionId string, callback interface{}) (interface{}, error)
// }

// type ActionRequest struct {
// 	HostPort string        `json:"host_port"`
// 	Method   string        `json:"method"`
// 	Params   []interface{} `json:"params"`
// }

// type ActionResponse struct {
// 	Results []interface{} `json:"results"`
// 	Code    int           `json:"code"`
// 	Msg     string        `json:"msg"`
// }

// type CallbackRequest struct {
// 	Callbacker string        `json:"callbacker"`
// 	Method     string        `json:"method"`
// 	RuntimeId  string        `json:"runtime_id"`
// 	Params     []interface{} `json:"params"`
// }

// type CallbackResponse struct {
// 	Results []interface{} `json:"results"`
// 	Code    int           `json:"code"`
// 	Msg     string        `json:"msg"`
// }
// type InitCallbacker interface {
// 	GetFlowPluginPath() string
// 	GetFlowActionPath(name string) string
// }

// type ActionCallbacker interface {
// 	GetRuntime() string
// 	GetFlow() string
// 	GetFlowCode() string
// 	GetFlowName() string
// 	GetFlowType() string
// 	GetFlowLogType() string
// 	GetFlowDict(name string) string
// 	GetFlowTimeout() string
// 	GetAction(aid string) string
// 	GetActionName(aid string) string
// 	GetActionTitle(aid string) string
// 	GetActionIcon(aid string) string
// 	GetActionScript(aid string) string
// 	GetActionFilter(aid string) string
// 	GetActionParams(aid string) map[string]string
// 	GetActionParam(aid, name string) string
// 	ShowMessage(aid string, content_type string, content string)
// 	ShowRuntimeState()
// 	ShowText(aid string, content string)
// 	ShowGrid(aid string, content string)
// 	ShowChart(aid string, content string)
// 	ShowHtml(aid string, content string)
// 	ShowForm(aid string, content string)
// 	ShowWeb(aid string, content string)
// 	GetRuntimeActionData(aid, name string) interface{}
// 	GetRuntimeActionDatas(aid string) map[string]interface{}
// 	SetRuntimeActionData(aid, name string, value interface{})
// 	SetRuntimeActionIcon(aid, icon string)
// 	SetRuntimeData(name string, value interface{})
// 	GetRuntimeData(name string) interface{}
// 	GetRuntimeDatas() map[string]interface{}
// 	GetRuntimeParams() map[string]interface{}
// 	GetRuntimeParam(name string) interface{}
// 	SetRuntimeParam(name string, value interface{})
// 	LogRuntimeFlowError(content string)
// 	LogRuntimeFlowInfo(content string)
// 	LogRuntimeFlowWarn(content string)
// 	LogRuntimeActionError(aid string, content string)
// 	LogRuntimeActionInfo(aid string, content string)
// 	LogRuntimeActionWarn(aid string, content string)
// 	LogRuntimeLinkError(sid string, tid string, content string)
// 	LogRuntimeLinkInfo(sid string, tid string, content string)
// 	LogRuntimeLinkWarn(sid string, tid string, content string)
// 	GetRuntimeId() string
// 	GetRuntimeDes() string
// 	GetRuntimeContextId() string
// 	GetRuntimeClientId() string
// 	SetRuntimeIserror(iserror int)
// 	GetRuntimeIserror() int
// 	SetRuntimeFlowState(state int)
// 	GetRuntimeFlowState() int
// 	SetRuntimeFlowBeginTime(t int64)
// 	GetRuntimeFlowBeginTime() int64
// 	SetRuntimeFlowEndTime(t int64)
// 	GetRuntimeFlowEndTime() int64
// 	SaveRuntime()
// 	ExecuteFlowByCode(clientId string, contextId string, code string, params map[string]interface{}) (string, error)
// 	ExecuteFlowByModel(clientId string, contextId string, flow string, params map[string]interface{}) (string, error)
// 	ExtendRuntimeAction(aid string, rt string)
// }

// func ParseActionCallbacker(obj interface{}) ActionCallbacker {
// 	p, ok := obj.(ActionCallbacker)

// 	if ok == false {
// 		return nil
// 	}
// 	return p
// }

// func ParseInitCallbacker(obj interface{}) InitCallbacker {
// 	p, ok := obj.(InitCallbacker)

// 	if ok == false {
// 		return nil
// 	}
// 	return p
// }
