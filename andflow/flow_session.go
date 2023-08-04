package andflow

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

var chan_len = 100

type Session struct {
	Id string

	Ctx           context.Context
	Controller    RuntimeController
	Router        FlowRouter
	Runner        FlowRunner
	ActionChanMap sync.Map //[string]chan *ActionParam
	LinkChanMap   sync.Map //map[string]chan *LinkParam
}

func (s *Session) getActionChan(actionId string) chan *ActionParam {

	c, ok := s.ActionChanMap.Load(actionId)
	if ok && c != nil {
		return c.(chan *ActionParam)
	}

	return nil
}

func (s *Session) setActionChan(actionId string, c chan *ActionParam) {
	s.ActionChanMap.Store(actionId, c)
}

func (s *Session) delActionChan(actionId string) {

	s.ActionChanMap.Delete(actionId)

}

func (s *Session) getLinkChan(sourceId string, targetId string) chan *LinkParam {

	c, ok := s.LinkChanMap.Load(sourceId + "_" + targetId)
	if ok && c != nil {
		return c.(chan *LinkParam)
	}

	return nil
}

func (s *Session) setLinkChan(sourceId string, targetId string, c chan *LinkParam) {
	s.LinkChanMap.Store(sourceId+"_"+targetId, c)
}
func (s *Session) delLinkChan(sourceId string, targetId string) {

	s.LinkChanMap.Delete(sourceId + "_" + targetId)

}

func (s *Session) GetFlow() *FlowModel {
	flow := s.Controller.GetFlow()
	return flow
}

func (s *Session) GetRuntime() *RuntimeModel {
	runtime := s.Controller.GetRuntime()
	return runtime
}

func (s *Session) GetParamMap() map[string]interface{} {
	return s.Controller.GetParamMap()
}
func (s *Session) GetParam(key string) interface{} {
	return s.Controller.GetParam(key)
}
func (s *Session) SetParam(key string, val interface{}) {
	s.Controller.SetParam(key, val)
}

func (s *Session) GetDataMap() map[string]interface{} {
	return s.Controller.GetDataMap()
}
func (s *Session) GetData(key string) interface{} {
	return s.Controller.GetData(key)
}
func (s *Session) SetData(key string, val interface{}) {
	s.Controller.SetData(key, val)
}

func (s *Session) AddLog_flow_error(name, title, content string) {
	s.Controller.AddLog("error", "flow", name, title, content)
}

func (s *Session) AddLog_flow_info(name, title, content string) {
	s.Controller.AddLog("info", "flow", name, title, content)
}

func (s *Session) AddLog_action_error(name, title, content string) {
	s.Controller.AddLog("error", "action", name, title, content)
}

func (s *Session) AddLog_action_info(name, title, content string) {
	s.Controller.AddLog("info", "action", name, title, content)
}

func (s *Session) AddLog_link_error(name, title, content string) {
	s.Controller.AddLog("error", "link", name, title, content)
}

func (s *Session) AddLog_link_info(name, title, content string) {
	s.Controller.AddLog("info", "link", name, title, content)
}

func (s *Session) createActionState(actionId string, preActionId string) *ActionStateModel {
	flow := s.GetFlow()

	action := flow.GetAction(actionId)

	state := &ActionStateModel{}
	state.ActionId = actionId
	state.ActionDes = action.Des
	state.ActionName = action.Name
	state.ActionTitle = action.Title
	state.ActionIcon = action.Icon
	state.PreActionId = preActionId
	content := &ActionContentModel{}

	if action.Content != nil {
		content_str := action.Content["content"]
		content_type := action.Content["content_type"]
		content.ContentType = content_type
		content.Content = content_str
	}
	state.Content = content

	state.BeginTime = time.Now()

	return state
}

func (s *Session) createLinkState(sourceId string, targetId string) *LinkStateModel {
	flow := s.GetFlow()
	link := flow.GetLinkBySourceIdAndTargetId(sourceId, targetId)

	state := &LinkStateModel{}
	state.SourceActionId = link.SourceId
	state.TargetActionId = link.TargetId

	state.BeginTime = time.Now()

	return state
}

// 监控协程
func (s *Session) watch() {

	for {

		select {
		case <-s.Ctx.Done():
			s.Stop()
			s.AddLog_flow_info(s.GetFlow().Code, s.GetFlow().Name, "timeout")
			log.Println("timeout, suspend")
			s.Runner.OnTimeout(s)
			return
		default:
			if s.Controller.GetCmd() == CMD_STOP {
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func CreateSession(context context.Context, controller RuntimeController, router FlowRouter, runner FlowRunner) *Session {
	session := &Session{}
	// uid, _ := uuid.NewV4()
	// id := strings.ReplaceAll(uid.String(), "-", "")
	session.Id = controller.GetRuntime().Id
	session.Ctx = context
	session.Router = router
	session.Runner = runner

	session.Controller = controller
	//actionid: RuntimeActionParam
	session.ActionChanMap = sync.Map{}
	//make(map[string]chan *ActionParam, 0)

	//sourceId_targetId: RuntimeLinkParam
	session.LinkChanMap = sync.Map{}
	//make(map[string]chan *LinkParam, 0)

	return session
}

func (s *Session) startChannel() {

	go s.watch()

	flow := s.GetFlow()

	w_action := &sync.WaitGroup{}
	for _, a := range flow.Actions {
		w_action.Add(1)
		go s.runProcessAction(a.Id, w_action)
	}

	w_link := &sync.WaitGroup{}
	for _, l := range flow.Links {
		w_link.Add(1)
		go s.runProcessLink(l.SourceId, l.TargetId, w_link)
	}

	w_action.Wait()
	w_link.Wait()

}

func (s *Session) Execute() {

	s.startChannel()
	defer s.stopChannel()

	s.Controller.SetBegin()
	defer s.Controller.SetEnd()

	s.Controller.SetCmd(CMD_START)

	//是否全新执行
	firstRun := true
	runningActions := s.Controller.GetRunningActions()

	if runningActions != nil && len(runningActions) > 0 {
		firstRun = false
		for _, param := range runningActions {
			s.ToAction(param)
		}
	}

	runningLinks := s.Controller.GetRunningLinks()
	if runningLinks != nil && len(runningLinks) > 0 {
		firstRun = false
		for _, param := range runningLinks {
			s.ToLink(param)
		}
	}

	//第一次执行或者从头执行。
	if firstRun {
		//清空状态
		s.Controller.GetRuntime().ActionStates = make([]*ActionStateModel, 0)
		s.Controller.GetRuntime().LinkStates = make([]*LinkStateModel, 0)

		runtimeId := s.Controller.GetRuntime().Id

		startIds := s.GetFlow().GetStartActionIds()
		for _, actionId := range startIds {
			param := &ActionParam{RuntimeId: runtimeId, ActionId: actionId, PreActionId: ""}
			s.ToAction(param)
		}

	}

	s.waitComplete()

}

func (s *Session) Stop() {
	s.Controller.SetCmd(CMD_STOP)
	log.Println("stop")
}

func (s *Session) waitComplete() {
	s.Controller.Wait()
}

func (s *Session) stopChannel() {

	s.Stop()
	flow := s.GetFlow()

	for _, a := range flow.Actions {
		s.stopProcessAction(a.Id)
	}
	for _, l := range flow.Links {
		s.stopProcessLink(l.SourceId, l.TargetId)
	}

	for _, a := range flow.Actions {
		c := s.getActionChan(a.Id)
		if c != nil {
			close(c)
		}
		s.delActionChan(a.Id)

	}
	for _, l := range flow.Links {
		c := s.getLinkChan(l.SourceId, l.TargetId)
		if c != nil {
			close(c)
		}
		s.delLinkChan(l.SourceId, l.TargetId)
	}

}

// 节点处理
func (s *Session) runProcessAction(actionId string, w *sync.WaitGroup) {

	c := s.getActionChan(actionId)
	if c == nil {
		c = make(chan *ActionParam, chan_len)
		s.setActionChan(actionId, c)
	}
	w.Done()

EVENT_LOOP:
	for {
		select {
		// 消息通道中有消息则执行，否则堵塞
		case param := <-c:
			if param.Cmd == 1 {
				break EVENT_LOOP
			}

			s.ExecuteAction(param)

		}
	}
}

func (s *Session) stopProcessAction(actionId string) {
	c := s.getActionChan(actionId)
	if c == nil {
		return
	}

	param := &ActionParam{Cmd: 1}

	c <- param
}
func (s *Session) ToAction(param *ActionParam) {
	if s.Router != nil {
		s.Controller.AddRunningAction(param)

		if s.Controller.GetCmd() == CMD_STOP {
			return
		}
		s.Controller.WaitAdd(1)
		s.Router.RouteAction(s, param)
	}
}

func (s *Session) PushAction(param *ActionParam) bool {

	c := s.getActionChan(param.ActionId)
	if c == nil {
		s.Controller.WaitDone()
		return false
	}

	c <- param
	return true
}

func (s *Session) ExecuteAction(param *ActionParam) {
	defer s.Controller.WaitDone() //确认节点执行完成

	actionState := s.createActionState(param.ActionId, param.PreActionId)
	var err error
	var res Result = RESULT_SUCCESS

	if s.Runner != nil {
		res, err = s.Runner.ExecuteAction(s, param, actionState)
		if err != nil || res == RESULT_FAILURE {
			if err == nil {
				err = errors.New("节点返回错误")
			}

			//错误回调
			s.Runner.OnActionFailure(s, param, actionState, err)

			s.AddLog_action_error(actionState.ActionName, actionState.ActionTitle, err.Error())
		}

	}

	actionState.State = int(res)
	if res == RESULT_FAILURE {
		actionState.IsError = 1
	}

	if res == RESULT_SUCCESS {

		flow := s.GetFlow()

		//在执行过程中指定的后续可执行节点
		nextActionIds := actionState.NextActionIds

		//所有后续连线
		nextLinks := flow.GetLinkBySourceId(param.ActionId)
		if nextLinks != nil && len(nextLinks) > 0 {
			for _, link := range nextLinks {
				if link.Active == "false" {
					continue
				}
				//如果有特别指定后续节点，就进行判断,如果不在指定的后续节点当中就不执行。
				if nextActionIds != nil && len(nextActionIds) > 0 {
					if arrayIndexOf(nextActionIds, link.TargetId) < 0 {
						continue
					}
				}

				linkParam := &LinkParam{RuntimeId: param.RuntimeId, SourceId: link.SourceId, TargetId: link.TargetId}
				s.ToLink(linkParam)
			}
		}
	}

	if res == RESULT_SUCCESS || res == RESULT_FAILURE {
		s.Controller.DelRunningAction(param) //删除待执行节点中已经执行完成的节点

		actionState.EndTime = time.Now()
		actionState.Timeused = actionState.EndTime.Sub(actionState.BeginTime).Milliseconds()
		s.Controller.AddActionState(actionState)
	}

}

// 连接线处理
func (s *Session) runProcessLink(sourceId, targetId string, w *sync.WaitGroup) {
	c := s.getLinkChan(sourceId, targetId)
	if c == nil {
		c = make(chan *LinkParam, chan_len)
		s.setLinkChan(sourceId, targetId, c)
	}
	w.Done()

EVENT_LOOP:
	for {
		select {
		// 消息通道中有消息则执行，否则堵塞
		case param := <-c:
			if param.Cmd == 1 {
				break EVENT_LOOP
			}

			s.ExecuteLink(param)
		}
	}

}

func (s *Session) stopProcessLink(sourceId, targetId string) {
	c := s.getLinkChan(sourceId, targetId)
	if c == nil {
		return
	}

	param := &LinkParam{Cmd: 1}

	c <- param
}

func (s *Session) ToLink(param *LinkParam) {
	if s.Router != nil {
		//准备执行的路径
		s.Controller.AddRunningLink(param)
		//如果停止
		if s.Controller.GetCmd() == CMD_STOP {
			return
		}

		s.Controller.WaitAdd(1)
		s.Router.RouteLink(s, param)
	}
}

func (s *Session) PushLink(param *LinkParam) bool {

	c := s.getLinkChan(param.SourceId, param.TargetId)
	if c == nil {
		s.Controller.WaitDone()
		return false
	}

	c <- param
	return true
}

func (s *Session) ExecuteLink(param *LinkParam) {
	defer s.Controller.WaitDone()

	linkState := s.createLinkState(param.SourceId, param.TargetId)
	sourceAction := s.GetFlow().GetAction(param.SourceId)
	targetAction := s.GetFlow().GetAction(param.TargetId)

	var err error
	var res Result = RESULT_SUCCESS

	if s.Runner != nil {
		res, err = s.Runner.ExecuteLink(s, param, linkState)
		if res == RESULT_FAILURE || err != nil {
			if err == nil {
				err = errors.New("连接线返回错误")
			}
			//错误回调
			s.Runner.OnLinkFailure(s, param, linkState, err)

			s.AddLog_link_error(param.SourceId+"->"+param.TargetId, sourceAction.Title+"->"+targetAction.Title, fmt.Sprintf("执行连线错误：%v", err))
		}
	}

	//状态
	linkState.State = int(res)
	if res == RESULT_FAILURE {
		linkState.IsError = 1
	}

	if res == RESULT_SUCCESS {

		canNext := true
		flow := s.GetFlow()
		nextAction := flow.GetAction(param.TargetId)

		//如果是汇聚模式，需要等待所有连线都执行通过
		if strings.ToLower(nextAction.Collect) == "true" {
			fromLinks := flow.GetLinkByTargetId(param.TargetId)
			if fromLinks != nil && len(fromLinks) > 1 {

				activeLinks := make([]*LinkModel, 0)
				for _, link := range fromLinks {
					if link.Active != "false" {
						activeLinks = append(activeLinks, link)
					}
				}

				passCount := 0
				for _, link := range activeLinks {
					st := s.Controller.GetLastLinkState(link.SourceId, link.TargetId)
					if st != nil && st.State == int(RESULT_SUCCESS) {
						passCount = passCount + 1
					}
				}
				if passCount+1 < len(activeLinks) {
					canNext = false
				}
			}
		}

		if canNext {
			actionParam := &ActionParam{RuntimeId: param.RuntimeId, ActionId: param.TargetId, PreActionId: param.SourceId}
			s.ToAction(actionParam)
		}

	}

	if res == RESULT_SUCCESS || res == RESULT_FAILURE {
		s.Controller.DelRunningLink(param) //移除待办路径中已经执行的路径
		linkState.EndTime = time.Now()
		linkState.Timeused = linkState.EndTime.Sub(linkState.BeginTime).Milliseconds()

		s.Controller.AddLinkState(linkState)
	}

}

func (s *Session) ExecuteFlow(flow *FlowModel, param map[string]interface{}, timeout int64) *RuntimeModel {
	runtime := ExecuteFlow(flow, param, timeout)

	return runtime
}

func arrayIndexOf(array []string, val string) (index int) {

	for i := 0; i < len(array); i++ {
		if array[i] == val {
			return i
		}
	}
	return -1
}
