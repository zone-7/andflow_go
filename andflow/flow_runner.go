package andflow

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/dop251/goja"
)

type FlowRunner interface {
	ExecuteLink(s *Session, param *LinkParam, state *LinkStateModel) (Result, error)       //返回三个状态 -1 不通过，1通过，0还没准备好执行
	ExecuteAction(s *Session, param *ActionParam, state *ActionStateModel) (Result, error) //返回三个状态 -1 不通过，1通过，0还没准备好执行
	OnLinkFailure(s *Session, param *LinkParam, state *LinkStateModel, err error)
	OnActionFailure(s *Session, param *ActionParam, state *ActionStateModel, err error)
	OnTimeout(s *Session)
}

type CommonFlowRunner struct {
	ActionScriptFunc  func(rts *goja.Runtime, s *Session, param *ActionParam, state *ActionStateModel)
	LinkScriptFunc    func(rts *goja.Runtime, s *Session, param *LinkParam, state *LinkStateModel)
	ActionFailureFunc func(s *Session, param *ActionParam, state *ActionStateModel, err error)
	LinkFailureFunc   func(s *Session, param *LinkParam, state *LinkStateModel, err error)
	TimeoutFunc       func(s *Session)
}

func (r *CommonFlowRunner) SetActionScriptFunc(act func(rts *goja.Runtime, s *Session, param *ActionParam, state *ActionStateModel)) {
	r.ActionScriptFunc = act
}
func (r *CommonFlowRunner) SetLinkScriptFunc(act func(rts *goja.Runtime, s *Session, param *LinkParam, state *LinkStateModel)) {
	r.LinkScriptFunc = act
}

func (r *CommonFlowRunner) SetActionFailureFunc(e func(s *Session, param *ActionParam, state *ActionStateModel, err error)) {
	r.ActionFailureFunc = e
}
func (r *CommonFlowRunner) SetLinkFailureFunc(e func(s *Session, param *LinkParam, state *LinkStateModel, err error)) {
	r.LinkFailureFunc = e
}
func (r *CommonFlowRunner) SetTimeoutFunc(e func(s *Session)) {
	r.TimeoutFunc = e
}

func (r *CommonFlowRunner) ExecuteLink(s *Session, param *LinkParam, state *LinkStateModel) (Result, error) {
	link := s.GetFlow().GetLinkBySourceIdAndTargetId(param.SourceId, param.TargetId)
	sc := link.Filter
	if len(strings.Trim(sc, " ")) == 0 {
		return RESULT_SUCCESS, nil
	}

	rts := goja.New()
	rts.Set("flow", s.GetFlow())
	rts.Set("link", link)
	if r.LinkScriptFunc != nil {
		r.LinkScriptFunc(rts, s, param, state)
	}

	SetCommonLinkScriptFunc(rts, s, param, state)

	script := "function $exec(){\n" + sc + "\n}\n $exec();\n"
	val, err := rts.RunString(script)

	if err != nil {
		return RESULT_FAILURE, err
	}

	res := GetScriptIntResult(val)

	return res, nil
}

func (r *CommonFlowRunner) ExecuteAction(s *Session, param *ActionParam, state *ActionStateModel) (Result, error) {
	var res Result = RESULT_SUCCESS
	var err error

	action := s.GetFlow().GetAction(param.ActionId)
	name := action.Name

	log.Println("action start: " + action.Name + " " + action.Title)
	defer log.Println("action end: " + action.Name + " " + action.Title)

	//0.准备脚本执行环境
	rts := goja.New()
	rts.Set("flow", s.GetFlow())
	rts.Set("action", action)

	if r.ActionScriptFunc != nil {
		r.ActionScriptFunc(rts, s, param, state)
	}

	SetCommonActionScriptFunc(rts, s, param, state)

	//1.执行过滤脚本
	if len(strings.Trim(action.Filter, " ")) > 0 {

		script_filter := "function $filter(){\n" + action.Filter + "\n}\n $filter();\n"

		val, err := rts.RunString(script_filter)

		if res == RESULT_FAILURE || err != nil {
			if err == nil {
				err = errors.New("执行过滤脚本返回错误")
			}
			log.Println(fmt.Sprintf("script exception：%v", err))

			return RESULT_FAILURE, err
		}

		res := GetScriptIntResult(val)

		if res != RESULT_SUCCESS {
			return res, nil
		}

	}

	//2.执行节点执行器
	runner := GetActionRunner(name)
	if runner != nil {
		res, err = runner.Execute(s, param, state)

		if err != nil || res == RESULT_FAILURE {

			if err == nil {
				err = errors.New("节点" + action.Name + "," + action.Title + "执行错误")
			}

			s.AddLog_action_error(action.Name, action.Title, err.Error())
			//执行异常处理脚本
			if len(strings.Trim(action.Error, " ")) > 0 {
				script_error := "function $exec(){\n" + action.Error + "\n}\n $exec();\n"
				_, err_err := rts.RunString(script_error)
				if err_err != nil {
					log.Println(fmt.Sprintf("script exception：%v", err_err))
				}
			}

			return RESULT_FAILURE, err
		}

		if res != RESULT_SUCCESS {
			return res, nil
		}
	}

	//3.执行事后脚本
	if len(strings.Trim(action.Script, " ")) > 0 {

		script_filter := "function $exec(){\n" + action.Script + "\n}\n $exec();\n"

		val, err := rts.RunString(script_filter)
		if err != nil {

			log.Println(fmt.Sprintf("script exception：%v", err))
			return RESULT_FAILURE, err
		}

		res := GetScriptIntResult(val)

		if res != RESULT_SUCCESS {
			return res, nil
		}
	}

	return res, nil
}

func (r *CommonFlowRunner) OnActionFailure(s *Session, param *ActionParam, state *ActionStateModel, err error) {
	if r.ActionFailureFunc != nil {
		r.ActionFailureFunc(s, param, state, err)
	}
}
func (r *CommonFlowRunner) OnLinkFailure(s *Session, param *LinkParam, state *LinkStateModel, err error) {
	if r.LinkFailureFunc != nil {
		r.LinkFailureFunc(s, param, state, err)
	}

}

func (r *CommonFlowRunner) OnTimeout(s *Session) {
	if r.TimeoutFunc != nil {
		r.TimeoutFunc(s)
	}

}
