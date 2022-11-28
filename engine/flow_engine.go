package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/zone-7/andflow_go/models"
)

var chan_len = 100

type Session struct {
	Id            string
	Wg            sync.WaitGroup //同步控制
	Cmd           int            //指令 1:停止,2.3.4..
	Ctx           context.Context
	Runtime       *models.RuntimeModel
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

	// c := s.getActionChan(actionId)
	s.ActionChanMap.Delete(actionId)
	// if c != nil {
	// 	close(c)
	// }

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

	// c := s.getLinkChan(sourceId, targetId)
	s.LinkChanMap.Delete(sourceId + "_" + targetId)
	// if c != nil {
	// 	close(c)
	// }

}

func (s *Session) changeState() {

	//如果没有什么待办事项就表示执行完了
	if (s.Runtime.RunningActions == nil || len(s.Runtime.RunningActions) == 0) &&
		(s.Runtime.RunningLinks == nil || len(s.Runtime.RunningLinks) == 0) {
		s.Runtime.FlowState = 1
		s.Runtime.EndTime = time.Now()
		s.Runtime.Timeused = s.Runtime.EndTime.Sub(s.Runtime.BeginTime).Milliseconds()
	}

}

// 监控协程
func (s *Session) watch() {

	for {
		select {
		case <-s.Ctx.Done():
			s.Stop()
			s.Runtime.AddLog("flow", "stop", "中断", "执行中断")
			fmt.Println("执行中断")
			return
		default:
			if s.Cmd == 1 {
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func CreateSession(context context.Context, runtime *models.RuntimeModel, runner FlowRunner) *Session {
	session := &Session{}
	uid, _ := uuid.NewV4()
	id := strings.ReplaceAll(uid.String(), "-", "")
	session.Id = id
	session.Wg = sync.WaitGroup{}
	session.Ctx = context
	session.Runner = runner

	session.Runtime = runtime
	//actionid: RuntimeActionParam
	session.ActionChanMap = sync.Map{}
	//make(map[string]chan *models.ActionParam, 0)

	//sourceId_targetId: RuntimeLinkParam
	session.LinkChanMap = sync.Map{}
	//make(map[string]chan *models.LinkParam, 0)

	flow := runtime.Flow

	w_action := &sync.WaitGroup{}
	for _, a := range flow.Actions {
		w_action.Add(1)
		go session.startProcessAction(a.Id, w_action)
	}

	w_link := &sync.WaitGroup{}
	for _, l := range flow.Links {
		w_link.Add(1)
		go session.startProcessLink(l.SourceId, l.TargetId, w_link)
	}

	w_action.Wait()
	w_link.Wait()

	return session
}

func (s *Session) Execute() {

	go s.watch()

	//是否全新执行
	firstRun := true

	if s.Runtime.RunningActions != nil && len(s.Runtime.RunningActions) > 0 {
		firstRun = false
		for _, param := range s.Runtime.RunningActions {
			s.toAction(param)
		}
	}

	if s.Runtime.RunningLinks != nil && len(s.Runtime.RunningLinks) > 0 {
		firstRun = false
		for _, param := range s.Runtime.RunningLinks {
			s.toLink(param)
		}
	}

	if firstRun {
		s.Runtime.BeginTime = time.Now()

		startIds := s.Runtime.Flow.GetStartActionIds()
		for _, actionId := range startIds {
			param := &models.ActionParam{ActionId: actionId, PreActionId: ""}
			s.toAction(param)
		}

	}

}

func (s *Session) Stop() {
	s.Cmd = 1

	for _, a := range s.Runtime.Flow.Actions {
		s.stopProcessAction(a.Id)
	}
	for _, l := range s.Runtime.Flow.Links {
		s.stopProcessLink(l.SourceId, l.TargetId)
	}

	fmt.Println("stop...............................")
}

func (s *Session) WaitComplete() {
	s.Wg.Wait()
}

func (s *Session) Release() {

	s.Stop()

	for _, a := range s.Runtime.Flow.Actions {
		c := s.getActionChan(a.Id)
		if c != nil {
			close(c)
		}
		s.delActionChan(a.Id)

	}
	for _, l := range s.Runtime.Flow.Links {
		c := s.getLinkChan(l.SourceId, l.TargetId)
		if c != nil {
			close(c)
		}
		s.delLinkChan(l.SourceId, l.TargetId)
	}

}

//节点处理
func (s *Session) startProcessAction(actionId string, w *sync.WaitGroup) {

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
			if param.Stop {
				break EVENT_LOOP
			}

			s.exeAction(param)

		}
	}
}

func (s *Session) stopProcessAction(actionId string) {
	c := s.getActionChan(actionId)
	if c == nil {
		return
	}

	param := &models.ActionParam{Stop: true}

	c <- param
}
func (s *Session) toAction(param *models.ActionParam) {
	s.Runtime.AddRunningAction(param)

	c := s.getActionChan(param.ActionId)
	if c == nil {
		return
	}
	if s.Cmd == 1 {
		return
	}
	s.Wg.Add(1)

	c <- param
}

func (s *Session) exeAction(param *models.ActionParam) {
	defer s.Wg.Done()     //确认节点执行完成
	defer s.changeState() //改变状态

	actionState := s.Runtime.CreateActionState(param.ActionId, param.PreActionId)

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
		nextLinks := s.Runtime.Flow.GetLinkBySourceId(param.ActionId)
		if nextLinks != nil && len(nextLinks) > 0 {
			for _, link := range nextLinks {
				if link.Active == "false" {
					continue
				}

				linkParam := &models.LinkParam{SourceId: link.SourceId, TargetId: link.TargetId}
				s.toLink(linkParam)
			}
		}
	} else {
		actionState.State = -1
	}
	if complete {
		s.Runtime.DelRunningAction(param) //删除待执行节点中已经执行完成的节点

		actionState.EndTime = time.Now()
		actionState.Timeused = actionState.EndTime.Sub(actionState.BeginTime).Milliseconds()
		s.Runtime.AddActionState(actionState)
	}

}

//连接线处理
func (s *Session) startProcessLink(sourceId, targetId string, w *sync.WaitGroup) {
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
			if param.Stop {
				break EVENT_LOOP
			}

			s.exeLink(param)
		}
	}

}

func (s *Session) stopProcessLink(sourceId, targetId string) {
	c := s.getLinkChan(sourceId, targetId)
	if c == nil {
		return
	}

	param := &models.LinkParam{Stop: true}

	c <- param
}

func (s *Session) toLink(param *models.LinkParam) {
	//准备执行的路径
	s.Runtime.AddRunningLink(param)

	c := s.getLinkChan(param.SourceId, param.TargetId)
	if c == nil {
		return
	}

	if s.Cmd == 1 {
		return
	}

	s.Wg.Add(1)

	c <- param
}

func (s *Session) exeLink(param *models.LinkParam) {
	defer s.Wg.Done()
	defer s.changeState() //改变状态

	linkState := s.Runtime.CreateLinkState(param.SourceId, param.TargetId)

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

		nextAction := s.Runtime.Flow.GetAction(param.TargetId)

		//如果是汇聚模式，需要等待所有连线都执行通过
		if strings.ToLower(nextAction.Collect) == "true" {
			fromLinks := s.Runtime.Flow.GetLinkByTargetId(param.TargetId)
			if fromLinks != nil && len(fromLinks) > 1 {
				for _, link := range fromLinks {
					st := s.Runtime.GetLastLinkState(link.SourceId, link.TargetId)
					if st == nil || st.State != 1 {
						canNext = false
					}
				}
			}
		}

		if canNext {
			actionParam := &models.ActionParam{ActionId: param.TargetId, PreActionId: param.SourceId}
			s.toAction(actionParam)
		}

	} else {
		linkState.State = -1
	}

	if complete {
		s.Runtime.DelRunningLink(param) //移除待办路径中已经执行的路径
		linkState.EndTime = time.Now()
		linkState.Timeused = linkState.EndTime.Sub(linkState.BeginTime).Milliseconds()
		s.Runtime.AddLinkState(linkState)
	}

}

func ParseFlow(content string) (*models.FlowModel, error) {

	flowModel := models.FlowModel{}
	err := json.Unmarshal([]byte(content), &flowModel)
	if err != nil {
		return nil, err
	}
	return &flowModel, nil
}

//创建运行时
func CreateRuntime(flow *models.FlowModel, param map[string]interface{}) *models.RuntimeModel {
	runtime := models.RuntimeModel{}
	uid, _ := uuid.NewV4()
	id := strings.ReplaceAll(uid.String(), "-", "")
	runtime.Id = id
	runtime.Flow = flow

	runtime.FlowState = 0
	runtime.Des = flow.Name
	//初始化状态
	if runtime.ActionStates == nil {
		runtime.ActionStates = make([]*models.ActionStateModel, 0)
	}
	//初始化下一步执行的连接线
	if runtime.RunningLinks == nil {
		runtime.RunningLinks = make([]*models.LinkParam, 0)
	}
	if runtime.RunningActions == nil {
		runtime.RunningActions = make([]*models.ActionParam, 0)
	}

	//初始化日志
	if runtime.Logs == nil {
		runtime.Logs = make([]*models.LogModel, 0)
	}

	//初始化数据
	if runtime.Data == nil {
		runtime.Data = make([]*models.RuntimeDataModel, 0)
	}
	//初始化参数
	if param != nil {
		for k, v := range param {
			runtime.SetParam(k, v)
		}
	}

	runtime.CreateTime = time.Now()
	return &runtime
}

func Execute(runtime *models.RuntimeModel, runner FlowRunner, timeout int64) {
	context, _ := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(timeout))

	session := CreateSession(context, runtime, runner)

	session.Execute()
	session.WaitComplete()
	session.Release()

}
