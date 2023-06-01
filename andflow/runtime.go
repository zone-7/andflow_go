package andflow

import (
	"encoding/json"
	"time"
)

//节点执行参数
type ActionParam struct {
	RuntimeId   string `json:"runtime_id"`
	ActionId    string `json:"action_id"`     //当前节点ID
	PreActionId string `json:"pre_action_id"` //上一个节点ID
	// Timeout     int64  `json:"timeout"`

	Cmd int `json:"cmd"` //指令  1:停止,2.3.4..
}

//连接线执行参数
type LinkParam struct {
	RuntimeId string `json:"runtime_id"`
	SourceId  string `json:"source_id"` //源ID
	TargetId  string `json:"target_id"` //目的ID
	// Timeout   int64  `json:"timeout"`
	Cmd int `json:"stop"` //指令  1:停止,2.3.4..
}

//节点内容
type ActionContentModel struct {
	ActionId    string `bson:"action_id" json:"action_id"`       //节点ID
	ContentType string `bson:"content_type" json:"content_type"` //内容类型
	Content     string `bson:"content" json:"content"`           //内容信息
}

//节点数据
type ActionDataModel struct {
	Name  string      `bson:"name" json:"name"`   //数据名称
	Value interface{} `bson:"value" json:"value"` //数据值
}

//节点状态
type ActionStateModel struct {
	ActionId    string              `bson:"action_id" json:"action_id"`         //节点ID
	ActionName  string              `bson:"action_name" json:"action_name"`     //节点名称
	ActionTitle string              `bson:"action_title" json:"action_title"`   //节点标题
	ActionDes   string              `bson:"action_des" json:"action_des"`       //节点描述
	ActionIcon  string              `bson:"action_icon" json:"action_icon"`     //节点图标
	PreActionId string              `bson:"pre_action_id" json:"pre_action_id"` //上一个节点ID
	IsError     int                 `bson:"is_error" json:"is_error"`           //是否异常
	State       int                 `bson:"state" json:"state"`                 //状态：1，完成并继续往下，0 未来执行， -1完成并终止
	Data        []*ActionDataModel  `bson:"data" json:"data"`                   //执行结果
	Content     *ActionContentModel `bson:"content" json:"content"`             //界面显示的内容
	BeginTime   time.Time           `bson:"begin_time" json:"begin_time"`       //开始时间
	EndTime     time.Time           `bson:"end_time" json:"end_time"`           //完成时间
	Timeused    int64               `bson:"timeused" json:"timeused"`           //耗时
}

//连接线状态
type LinkStateModel struct {
	SourceActionId string    `bson:"source_action_id" json:"source_action_id"` //源ID
	TargetActionId string    `bson:"target_action_id" json:"target_action_id"` //目的ID
	IsError        int       `bson:"is_error" json:"is_error"`                 //是否异常
	State          int       `bson:"state" json:"state"`                       //状态
	BeginTime      time.Time `bson:"begin_time" json:"begin_time"`             //开始时间
	EndTime        time.Time `bson:"end_time" json:"end_time"`                 //完成时间
	Timeused       int64     `bson:"timeused" json:"timeused"`                 //耗时
}

//日志
type LogModel struct {
	Tp      string    `bson:"tp" json:"tp"`           //类型
	Id      string    `bson:"id" json:"id"`           //id
	Name    string    `bson:"name" json:"name"`       //名称
	Title   string    `bson:"title" json:"title"`     //标题
	Time    time.Time `bson:"time" json:"time"`       //时间
	Tag     string    `bson:"tag" json:"tag"`         //标记
	Content string    `bson:"content" json:"content"` //内容
}

//参数类型信息
type ParamInfo struct {
	TypeName          string `bson:"type_name" json:"type_name"`                   //类型
	Size              int    `bson:"size" json:"size"`                             //大小
	ExpireMillisecond int64  `bson:"expire_millisecond" json:"expire_millisecond"` //超时
}

//运行时数据
type RuntimeDataModel struct {
	Name  string      `bson:"name" json:"name"`   //名称
	Value interface{} `bson:"value" json:"value"` //值
}

//运行时
type RuntimeModel struct {
	Id             string              `bson:"_id" json:"id"`                          //ID
	Flow           *FlowModel          `bson:"flow" json:"flow"`                       //流程
	Des            string              `bson:"des" json:"des"`                         //描述
	BeginTime      time.Time           `bson:"begin_time" json:"begin_time"`           //开始时间
	EndTime        time.Time           `bson:"end_time" json:"end_time"`               //完成时间
	Timeused       int64               `bson:"timeused" json:"timeused"`               //耗时
	IsRunning      int                 `bson:"is_running" json:"is_running"`           //是否正在运行
	IsError        int                 `bson:"is_error" json:"is_error"`               //是否异常
	RunningLinks   []*LinkParam        `bson:"running_links" json:"running_links"`     //待执行连接线
	RunningActions []*ActionParam      `bson:"running_actions" json:"running_actions"` //待执行节点
	FlowState      int                 `bson:"flow_state" json:"flow_state"`           //流程执行状态
	ActionStates   []*ActionStateModel `bson:"action_states" json:"action_states"`     //节点状态
	LinkStates     []*LinkStateModel   `bson:"link_states" json:"link_states"`         //连接线状态
	Logs           []*LogModel         `bson:"logs" json:"logs"`                       //日志
	Param          []*RuntimeDataModel `bson:"param" json:"param"`                     //参数
	Data           []*RuntimeDataModel `bson:"data" json:"data"`                       //数据
	Message        string              `bson:"message" json:"message"`                 //信息
	CreateTime     time.Time           `bson:"create_time" json:"create_time"`         //创建时间
	UpdateTime     time.Time           `bson:"update_time" json:"update_time"`         //修改时间
	UserId         int64               `bson:"user_id" json:"user_id"`
}

func (a *ActionStateModel) SetData(name string, value interface{}) {
	if a.Data == nil {
		a.Data = make([]*ActionDataModel, 0)
	}
	exists := false
	for _, d := range a.Data {
		if d.Name == name {
			d.Value = value
			exists = true
		}
	}
	if !exists {
		a.Data = append(a.Data, &ActionDataModel{Name: name, Value: value})
	}

}

func (a *ActionStateModel) GetData(name string) interface{} {
	if a.Data == nil {
		return nil
	}

	for _, d := range a.Data {
		if d.Name == name {
			return d.Value
		}
	}
	return nil
}

func (a *ActionStateModel) DelData(name string) {
	if a.Data == nil {
		return
	}
	index := -1
	for i, d := range a.Data {
		if d.Name == name {
			index = i
		}
	}
	if index >= 0 {
		a.Data = append(a.Data[0:index], a.Data[index+1:]...)
	}
}

func (a *ActionStateModel) GetDataMap() map[string]interface{} {
	res := make(map[string]interface{})
	if a.Data == nil {
		return res
	}

	for _, d := range a.Data {
		res[d.Name] = d.Value
	}
	return res
}

func (a *RuntimeModel) AddActionState(state *ActionStateModel) {
	a.ActionStates = append(a.ActionStates, state)
}
func (a *RuntimeModel) GetActionState(actionId string, preActionId string) *ActionStateModel {
	if a.ActionStates == nil {
		return nil
	}
	i := len(a.ActionStates)
	for i >= 0 {
		state := a.ActionStates[i]
		if state.ActionId == actionId && state.PreActionId == preActionId {
			return state
		}
		i--
	}
	return nil
}
func (a *RuntimeModel) GetLastActionState(actionId string) *ActionStateModel {
	if a.ActionStates == nil {
		return nil
	}
	i := len(a.ActionStates) - 1
	for i >= 0 {
		state := a.ActionStates[i]
		if state.ActionId == actionId {
			return state
		}
		i--
	}
	return nil
}

func (a *RuntimeModel) AddLinkState(state *LinkStateModel) {
	a.LinkStates = append(a.LinkStates, state)

}

func (a *RuntimeModel) GetLastLinkState(sourceId string, targetId string) *LinkStateModel {
	if a.LinkStates == nil {
		return nil
	}
	i := len(a.LinkStates) - 1
	for i >= 0 {
		state := a.LinkStates[i]
		if state.SourceActionId == sourceId && state.TargetActionId == targetId {
			return state
		}
		i--
	}
	return nil
}

func (a *RuntimeModel) SetData(name string, value interface{}) {
	if a.Data == nil {
		a.Data = make([]*RuntimeDataModel, 0)
	}
	exists := false
	for _, d := range a.Data {
		if d.Name == name {
			d.Value = value
			exists = true
		}
	}
	if !exists {
		a.Data = append(a.Data, &RuntimeDataModel{Name: name, Value: value})
	}
}

func (a *RuntimeModel) GetData(name string) interface{} {
	if a.Data == nil {
		return nil
	}

	for _, d := range a.Data {
		if d.Name == name {
			return d.Value
		}
	}
	return nil
}

func (a *RuntimeModel) DelData(name string) {
	if a.Data == nil {
		return
	}
	index := -1
	for i, d := range a.Data {
		if d.Name == name {
			index = i
		}
	}
	if index >= 0 {
		a.Data = append(a.Data[0:index], a.Data[index+1:]...)
	}
}

func (a *RuntimeModel) GetDataMap() map[string]interface{} {
	res := make(map[string]interface{})
	if a.Data == nil {
		return res
	}

	for _, d := range a.Data {
		res[d.Name] = d.Value
	}
	return res
}
func (a *RuntimeModel) SetDataMap(data map[string]interface{}) {

	for k, v := range data {
		a.SetData(k, v)
	}
}

func (a *RuntimeModel) GetParamMap() map[string]interface{} {
	res := make(map[string]interface{})
	if a.Param == nil {
		return res
	}

	for _, d := range a.Param {
		res[d.Name] = d.Value
	}
	return res
}

func (a *RuntimeModel) SetParam(key string, value interface{}) {
	if a.Param == nil {
		a.Param = make([]*RuntimeDataModel, 0)
	}
	exists := false
	for _, d := range a.Param {
		if d.Name == key {
			d.Value = value
			exists = true
		}
	}
	if !exists {
		a.Param = append(a.Param, &RuntimeDataModel{Name: key, Value: value})
	}
}

func (a *RuntimeModel) SetParamMap(param map[string]interface{}) {

	for k, v := range param {
		a.SetParam(k, v)
	}
}

func (a *RuntimeModel) DelParam(key string) {
	if a.Param == nil {
		return
	}
	index := -1
	for i, d := range a.Param {
		if d.Name == key {
			index = i
		}
	}
	if index >= 0 {
		a.Param = append(a.Param[0:index], a.Param[index+1:]...)
	}
}

func (a *RuntimeModel) GetParam(key string) interface{} {
	if a.Param == nil {
		return nil
	}

	for _, d := range a.Param {
		if d.Name == key {
			return d.Value
		}
	}
	return nil
}

func (a *RuntimeModel) SetLinkState(sourceId, targetId string, state int) {
	linkState := a.GetLastLinkState(sourceId, targetId)
	if linkState == nil {
		return
	}
	linkState.State = state
}

func (a *RuntimeModel) SetLinkError(sourceId, targetId string, isError int) {
	linkState := a.GetLastLinkState(sourceId, targetId)
	if linkState == nil {
		return
	}
	linkState.IsError = isError
}
func (a *RuntimeModel) SetActionError(actionId string, isError int) {
	actionState := a.GetLastActionState(actionId)
	if actionState == nil {
		return
	}
	actionState.IsError = isError
}
func (a *RuntimeModel) SetActionState(actionId string, state int) {
	actionState := a.GetLastActionState(actionId)
	if actionState == nil {
		return
	}
	actionState.State = state
}

func (a *RuntimeModel) SetActionIcon(actionId string, icon string) {
	actionState := a.GetLastActionState(actionId)
	if actionState == nil {
		return
	}
	actionState.ActionIcon = icon
}
func (a *RuntimeModel) SetActionData(actionId string, name string, value interface{}) {
	actionState := a.GetLastActionState(actionId)
	if actionState == nil {
		return
	}
	actionState.SetData(name, value)
}

func (a *RuntimeModel) GetActionData(actionId string, name string) interface{} {
	actionState := a.GetLastActionState(actionId)
	if actionState == nil {
		return nil
	}
	return actionState.GetData(name)
}
func (a *RuntimeModel) GetActionDataMap(actionId string) map[string]interface{} {
	actionState := a.GetLastActionState(actionId)
	if actionState == nil {
		return nil
	}

	res := make(map[string]interface{})
	if actionState.Data == nil {
		return res
	}

	for _, d := range actionState.Data {
		res[d.Name] = d.Value
	}
	return res
}
func (a *RuntimeModel) AddLog(tp string, tag string, name string, title string, content string) {
	log := &LogModel{Tp: tp, Tag: tag, Name: name, Title: title, Content: content, Time: time.Now()}
	a.Logs = append(a.Logs, log)
}

func (r *RuntimeModel) AddRunningLink(param *LinkParam) {

	if r.RunningLinks == nil {
		r.RunningLinks = make([]*LinkParam, 0)
	}
	for _, link := range r.RunningLinks {
		if link.SourceId == param.SourceId && link.TargetId == param.TargetId {
			return
		}
	}
	r.RunningLinks = append(r.RunningLinks, param)

}
func (r *RuntimeModel) DelRunningLink(param *LinkParam) {

	if r.RunningLinks == nil {
		r.RunningLinks = make([]*LinkParam, 0)
	}
	for index, link := range r.RunningLinks {
		if link.SourceId == param.SourceId && link.TargetId == param.TargetId {
			r.RunningLinks = append(r.RunningLinks[:index], r.RunningLinks[index+1:]...)
			return
		}
	}

}

func (r *RuntimeModel) AddRunningAction(param *ActionParam) {

	if r.RunningActions == nil {
		r.RunningActions = make([]*ActionParam, 0)
	}
	for _, p := range r.RunningActions {
		if p.ActionId == param.ActionId && p.PreActionId == param.PreActionId {
			return
		}
	}

	r.RunningActions = append(r.RunningActions, param)

}
func (r *RuntimeModel) DelRunningAction(param *ActionParam) {

	if r.RunningActions == nil {
		r.RunningActions = make([]*ActionParam, 0)
	}
	for index, p := range r.RunningActions {
		if p.ActionId == param.ActionId && p.PreActionId == p.PreActionId {
			r.RunningActions = append(r.RunningActions[:index], r.RunningActions[index+1:]...)
			return
		}
	}
}

func (r *RuntimeModel) ToJson() string {
	data, _ := json.Marshal(r)
	return string(data)
}
