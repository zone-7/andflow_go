package engine

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/zone-7/andflow_go/models"
)

var chan_len = 100

type Session struct {
	Id string

	Ctx           context.Context
	Store         RuntimeStore
	Router        FlowRouter
	Runner        FlowRunner
	ActionChanMap sync.Map //[string]chan *models.ActionParam
	LinkChanMap   sync.Map //map[string]chan *models.LinkParam
}

func (s *Session) getActionChan(actionId string) chan *models.ActionParam {

	c, ok := s.ActionChanMap.Load(actionId)
	if ok && c != nil {
		return c.(chan *models.ActionParam)
	}

	return nil
}

func (s *Session) setActionChan(actionId string, c chan *models.ActionParam) {
	s.ActionChanMap.Store(actionId, c)
}

func (s *Session) delActionChan(actionId string) {

	s.ActionChanMap.Delete(actionId)

}

func (s *Session) getLinkChan(sourceId string, targetId string) chan *models.LinkParam {

	c, ok := s.LinkChanMap.Load(sourceId + "_" + targetId)
	if ok && c != nil {
		return c.(chan *models.LinkParam)
	}

	return nil
}

func (s *Session) setLinkChan(sourceId string, targetId string, c chan *models.LinkParam) {
	s.LinkChanMap.Store(sourceId+"_"+targetId, c)
}
func (s *Session) delLinkChan(sourceId string, targetId string) {

	s.LinkChanMap.Delete(sourceId + "_" + targetId)

}

func (s *Session) GetFlow() *models.FlowModel {
	flow := s.Store.GetFlow()
	return flow
}

func (s *Session) createActionState(actionId string, preActionId string) *models.ActionStateModel {
	flow := s.GetFlow()

	action := flow.GetAction(actionId)

	state := &models.ActionStateModel{}
	state.ActionId = actionId
	state.ActionDes = action.Des
	state.ActionName = action.Name
	state.ActionTitle = action.Title
	state.ActionIcon = action.Icon
	state.PreActionId = preActionId
	content := &models.ActionContentModel{}

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

func (s *Session) createLinkState(sourceId string, targetId string) *models.LinkStateModel {
	flow := s.GetFlow()
	link := flow.GetLinkBySourceIdAndTargetId(sourceId, targetId)

	state := &models.LinkStateModel{}
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
	uid, _ := uuid.NewV4()
	id := strings.ReplaceAll(uid.String(), "-", "")
	session.Id = id
	session.Ctx = context
	session.Router = router
	session.Runner = runner

	session.Store = store
	//actionid: RuntimeActionParam
	session.ActionChanMap = sync.Map{}
	//make(map[string]chan *models.ActionParam, 0)

	//sourceId_targetId: RuntimeLinkParam
	session.LinkChanMap = sync.Map{}
	//make(map[string]chan *models.LinkParam, 0)

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
			param := &models.ActionParam{RuntimeId: runtimeId, ActionId: actionId, PreActionId: "", Timeout: s.Store.GetTimeout()}
			s.ToAction(param)
		}

	}

	s.waitComplete()

}

func (s *Session) Stop() {

	s.Store.SetCmd(1)
	log.Println("停止执行")
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
		c = make(chan *models.ActionParam, chan_len)
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

	param := &models.ActionParam{Cmd: 1}

	c <- param
}
func (s *Session) ToAction(param *models.ActionParam) {
	if s.Router != nil {
		s.Store.AddRunningAction(param)

		if s.Store.GetCmd() == 1 {
			return
		}
		s.Store.WaitAdd(1)
		s.Router.RouteAction(s, param)
	}
}

func (s *Session) PushAction(param *models.ActionParam) bool {

	c := s.getActionChan(param.ActionId)
	if c == nil {
		s.Store.WaitDone()
		return false
	}

	c <- param
	return true
}

func (s *Session) ExecuteAction(param *models.ActionParam) {
	defer s.Store.WaitDone()     //确认节点执行完成
	defer s.Store.RefreshState() //改变状态

	actionState := s.createActionState(param.ActionId, param.PreActionId)

	complete := true
	canNext := true
	if s.Runner != nil {
		res := s.Runner.ExecuteAction(s, param)
		if res != 1 {
			canNext = false
		}
		if res == 0 {
			complete = false
		}
	}

	if canNext {
		actionState.State = 1
		flow := s.GetFlow()

		nextLinks := flow.GetLinkBySourceId(param.ActionId)
		if nextLinks != nil && len(nextLinks) > 0 {
			for _, link := range nextLinks {
				if link.Active == "false" {
					continue
				}

				linkParam := &models.LinkParam{RuntimeId: param.RuntimeId, SourceId: link.SourceId, TargetId: link.TargetId, Timeout: param.Timeout}
				s.ToLink(linkParam)
			}
		}
	} else {
		actionState.State = -1
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
		c = make(chan *models.LinkParam, chan_len)
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

	param := &models.LinkParam{Cmd: 1}

	c <- param
}

func (s *Session) ToLink(param *models.LinkParam) {
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

func (s *Session) PushLink(param *models.LinkParam) bool {

	c := s.getLinkChan(param.SourceId, param.TargetId)
	if c == nil {
		s.Store.WaitDone()
		return false
	}

	c <- param
	return true
}

func (s *Session) ExecuteLink(param *models.LinkParam) {
	defer s.Store.WaitDone()
	defer s.Store.RefreshState() //改变状态

	linkState := s.createLinkState(param.SourceId, param.TargetId)

	complete := true
	toNext := true
	if s.Runner != nil {
		res := s.Runner.ExecuteLink(s, param)
		if res != 1 {
			toNext = false
		}

		if res == 0 {
			complete = false
		}
	}

	if toNext {
		linkState.State = 1

		canNext := true
		flow := s.GetFlow()
		nextAction := flow.GetAction(param.TargetId)

		//如果是汇聚模式，需要等待所有连线都执行通过
		if strings.ToLower(nextAction.Collect) == "true" {
			fromLinks := flow.GetLinkByTargetId(param.TargetId)

			if fromLinks != nil && len(fromLinks) > 1 {
				for _, link := range fromLinks {

					st := s.Store.GetLastLinkState(link.SourceId, link.TargetId)

					if st == nil || st.State != 1 {
						canNext = false
					}
				}
			}
		}

		if canNext {
			actionParam := &models.ActionParam{RuntimeId: param.RuntimeId, ActionId: param.TargetId, PreActionId: param.SourceId, Timeout: param.Timeout}
			s.ToAction(actionParam)
		}

	} else {
		linkState.State = -1
	}

	if complete {
		s.Store.DelRunningLink(param) //移除待办路径中已经执行的路径
		linkState.EndTime = time.Now()
		linkState.Timeused = linkState.EndTime.Sub(linkState.BeginTime).Milliseconds()

		s.Store.AddLinkState(linkState)
	}

}
func (s *Session) ExecuteFlow(flow *models.FlowModel, param map[string]interface{}, timeout int64) {
	runtime := ExecuteFlow(flow, param, timeout)
	if runtime != nil {
		ds := runtime.GetDataMap()
		if ds != nil {
			for k, v := range ds {
				s.Store.SetData(k, v)
			}
		}
	}
}