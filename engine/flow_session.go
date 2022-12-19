package engine

import (
	"context"
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
	Store         RuntimeStore
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
	flow := s.Store.GetFlow()
	return flow
}

func (s *Session) GetRuntime() *RuntimeModel {
	runtime := s.Store.GetRuntime()
	return runtime
}

func (s *Session) GetParam(key string) interface{} {
	return s.Store.GetParam(key)
}
func (s *Session) SetParam(key string, val interface{}) {
	s.Store.SetParam(key, val)
}

func (s *Session) GetData(key string) interface{} {
	return s.Store.GetData(key)
}
func (s *Session) SetData(key string, val interface{}) {
	s.Store.SetData(key, val)
}

func (s *Session) AddLog_flow_error(title, content string) {
	s.Store.AddLog("error", "flow", title, content)
}

func (s *Session) AddLog_flow_info(title, content string) {
	s.Store.AddLog("info", "flow", title, content)
}

func (s *Session) AddLog_action_error(title, content string) {
	s.Store.AddLog("error", "action", title, content)
}

func (s *Session) AddLog_action_info(title, content string) {
	s.Store.AddLog("info", "action", title, content)
}

func (s *Session) AddLog_link_error(title, content string) {
	s.Store.AddLog("error", "link", title, content)
}

func (s *Session) AddLog_link_info(title, content string) {
	s.Store.AddLog("info", "link", title, content)
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
			s.Store.AddLog("flow", "stop", "中断", "执行中断")
			log.Println("执行超时，中断退出")
			return
		default:
			if s.Store.GetCmd() == 1 {
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func CreateSession(context context.Context, store RuntimeStore, router FlowRouter, runner FlowRunner) *Session {
	session := &Session{}
	// uid, _ := uuid.NewV4()
	// id := strings.ReplaceAll(uid.String(), "-", "")
	session.Id = store.GetRuntime().Id
	session.Ctx = context
	session.Router = router
	session.Runner = runner

	session.Store = store
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

	s.Store.SetBegin()
	defer s.Store.SetEnd()

	s.Store.SetCmd(0)

	//是否全新执行
	firstRun := true
	runningActions := s.Store.GetRunningActions()

	if runningActions != nil && len(runningActions) > 0 {
		firstRun = false
		for _, param := range runningActions {
			s.ToAction(param)
		}
	}

	runningLinks := s.Store.GetRunningLinks()
	if runningLinks != nil && len(runningLinks) > 0 {
		firstRun = false
		for _, param := range runningLinks {
			s.ToLink(param)
		}
	}

	if firstRun {

		runtimeId := s.Store.GetRuntime().Id

		startIds := s.GetFlow().GetStartActionIds()
		for _, actionId := range startIds {
			param := &ActionParam{RuntimeId: runtimeId, ActionId: actionId, PreActionId: ""}
			s.ToAction(param)
		}

	}

	s.waitComplete()

}

func (s *Session) Stop() {

	s.Store.SetCmd(1)
	log.Println("停止流程")
}

func (s *Session) waitComplete() {
	s.Store.Wait()
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

//节点处理
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
		s.Store.AddRunningAction(param)

		if s.Store.GetCmd() == 1 {
			return
		}
		s.Store.WaitAdd(1)
		s.Router.RouteAction(s, param)
	}
}

func (s *Session) PushAction(param *ActionParam) bool {

	c := s.getActionChan(param.ActionId)
	if c == nil {
		s.Store.WaitDone()
		return false
	}

	c <- param
	return true
}

func (s *Session) ExecuteAction(param *ActionParam) {
	defer s.Store.WaitDone() //确认节点执行完成

	actionState := s.createActionState(param.ActionId, param.PreActionId)
	var err error
	var res Result = SUCCESS
	complete := true
	canNext := true
	if s.Runner != nil {
		res, err = s.Runner.ExecuteAction(s, param, actionState)
		if err != nil {
			s.AddLog_action_error(actionState.ActionName, fmt.Sprintf("执行节点错误：%v", err))
		}
		if res != SUCCESS {
			canNext = false
		}
		if res == REJECT {
			complete = false
		}
	}

	actionState.State = int(res)
	if res == FAILURE {
		actionState.IsError = 1
	}

	if canNext {

		flow := s.GetFlow()

		nextLinks := flow.GetLinkBySourceId(param.ActionId)
		if nextLinks != nil && len(nextLinks) > 0 {
			for _, link := range nextLinks {
				if link.Active == "false" {
					continue
				}

				linkParam := &LinkParam{RuntimeId: param.RuntimeId, SourceId: link.SourceId, TargetId: link.TargetId}
				s.ToLink(linkParam)
			}
		}
	}

	if complete {
		s.Store.DelRunningAction(param) //删除待执行节点中已经执行完成的节点

		actionState.EndTime = time.Now()
		actionState.Timeused = actionState.EndTime.Sub(actionState.BeginTime).Milliseconds()
		s.Store.AddActionState(actionState)
	}

}

//连接线处理
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
		s.Store.AddRunningLink(param)
		//如果停止
		if s.Store.GetCmd() == 1 {
			return
		}

		s.Store.WaitAdd(1)
		s.Router.RouteLink(s, param)
	}
}

func (s *Session) PushLink(param *LinkParam) bool {

	c := s.getLinkChan(param.SourceId, param.TargetId)
	if c == nil {
		s.Store.WaitDone()
		return false
	}

	c <- param
	return true
}

func (s *Session) ExecuteLink(param *LinkParam) {
	defer s.Store.WaitDone()

	linkState := s.createLinkState(param.SourceId, param.TargetId)
	var err error
	var res Result = SUCCESS
	complete := true
	toNext := true
	if s.Runner != nil {
		res, err = s.Runner.ExecuteLink(s, param, linkState)
		if err != nil {
			s.AddLog_link_error(param.SourceId+"->"+param.TargetId, fmt.Sprintf("执行连线错误：%v", err))
		}
		if res != SUCCESS {
			toNext = false
		}
		if res == REJECT {
			complete = false
		}
	}

	//状态
	linkState.State = int(res)
	if res == FAILURE {
		linkState.IsError = 1
	}

	if toNext {

		canNext := true
		flow := s.GetFlow()
		nextAction := flow.GetAction(param.TargetId)

		//如果是汇聚模式，需要等待所有连线都执行通过
		if strings.ToLower(nextAction.Collect) == "true" {
			fromLinks := flow.GetLinkByTargetId(param.TargetId)

			if fromLinks != nil && len(fromLinks) > 1 {
				for _, link := range fromLinks {

					st := s.Store.GetLastLinkState(link.SourceId, link.TargetId)

					if st == nil || st.State != int(SUCCESS) {
						canNext = false
					}
				}
			}
		}

		if canNext {
			actionParam := &ActionParam{RuntimeId: param.RuntimeId, ActionId: param.TargetId, PreActionId: param.SourceId}
			s.ToAction(actionParam)
		}

	}

	if complete {
		s.Store.DelRunningLink(param) //移除待办路径中已经执行的路径
		linkState.EndTime = time.Now()
		linkState.Timeused = linkState.EndTime.Sub(linkState.BeginTime).Milliseconds()

		s.Store.AddLinkState(linkState)
	}

}

func (s *Session) ExecuteFlow(flow *FlowModel, param map[string]interface{}, timeout int64) *RuntimeModel {
	runtime := ExecuteFlow(flow, param, timeout)

	return runtime
}
