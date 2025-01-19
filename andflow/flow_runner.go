package andflow

import (
	"encoding/json"
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
	OnLinkExecuted(s *Session, param *LinkParam, state *LinkStateModel, res Result, err error)
	OnActionExecuted(s *Session, param *ActionParam, state *ActionStateModel, res Result, err error)
	OnTimeout(s *Session)
}

type CommonFlowRunner struct {
	ActionScriptFunc   func(rts *goja.Runtime, s *Session, param *ActionParam, state *ActionStateModel)
	LinkScriptFunc     func(rts *goja.Runtime, s *Session, param *LinkParam, state *LinkStateModel)
	ActionFailureFunc  func(s *Session, param *ActionParam, state *ActionStateModel, err error)
	LinkFailureFunc    func(s *Session, param *LinkParam, state *LinkStateModel, err error)
	ActionExecutedFunc func(s *Session, param *ActionParam, state *ActionStateModel, res Result, err error)
	LinkExecutedFunc   func(s *Session, param *LinkParam, state *LinkStateModel, res Result, err error)

	TimeoutFunc func(s *Session)
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

func (r *CommonFlowRunner) SetActionExecutedFunc(e func(s *Session, param *ActionParam, state *ActionStateModel, res Result, err error)) {
	r.ActionExecutedFunc = e
}
func (r *CommonFlowRunner) SetLinkExecutedFunc(e func(s *Session, param *LinkParam, state *LinkStateModel, res Result, err error)) {
	r.LinkExecutedFunc = e
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
	if len(strings.Trim(action.ScriptBefore, " ")) > 0 {

		script_before := "function $filter(){\n" + action.ScriptBefore + "\n}\n $filter();\n"

		val, err := rts.RunString(script_before)

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

	//获取迭代器，用于循环执行
	var iteratorList []interface{}

	if len(action.IteratorList) > 0 {
		//先找参数
		iteratorListParam := s.GetParam(action.IteratorList)

		if iteratorListParam != nil {

			switch iteratorListParam.(type) {
			case []interface{}:
				iteratorList = iteratorListParam.([]interface{})
			case string:
				err = json.Unmarshal([]byte(iteratorListParam.(string)), &iteratorList)
				if err != nil {
					iteratorList = append(iteratorList, iteratorListParam.(string))
				}

			case map[string]interface{}:

				for _, v := range iteratorListParam.(map[string]interface{}) {
					iteratorList = append(iteratorList, fmt.Sprintf("%v", v))
				}
			}

		} else {
			//找不到参数就直接从字符串转化

			iteratorListStr, err := ReplaceTemplate(action.IteratorList, "iterator_list", s.GetParamMap())
			if err != nil {
				iteratorListStr = action.IteratorList
			}
			err = json.Unmarshal([]byte(iteratorListStr), &iteratorList)
			if err != nil {
				iteratorList = append(iteratorList, iteratorListParam.(string))
			}
		}

	}

	if len(iteratorList) == 0 {
		iteratorList = append(iteratorList, "0")
	}

	//2.执行节点执行器
	runner := GetActionRunner(name)
	if runner != nil {

		//迭代执行
		for _, item := range iteratorList {

			if len(action.IteratorItem) > 0 {
				s.SetParam(action.IteratorItem, item)
			}

			res, err = runner.Execute(s, param, state)

			//如果异常 或者执行失败
			if err != nil || res == RESULT_FAILURE {

				if err == nil {
					err = errors.New("节点" + action.Name + "," + action.Title + "执行错误")
				}

				s.AddLog_action_error(action.Name, action.Title, err.Error())

				//执行异常处理脚本
				if len(strings.Trim(action.ScriptError, " ")) > 0 {
					script_error := "function $exec(){\n" + action.ScriptError + "\n}\n $exec();\n"
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

	}

	//3.执行事后脚本
	if len(strings.Trim(action.ScriptAfter, " ")) > 0 {

		script_after := "function $exec(){\n" + action.ScriptAfter + "\n}\n $exec();\n"

		val, err := rts.RunString(script_after)
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

func (r *CommonFlowRunner) OnLinkExecuted(s *Session, param *LinkParam, state *LinkStateModel, res Result, err error) {
	if r.LinkExecutedFunc != nil {
		r.LinkExecutedFunc(s, param, state, res, err)
	}
}

func (r *CommonFlowRunner) OnActionExecuted(s *Session, param *ActionParam, state *ActionStateModel, res Result, err error) {
	if r.ActionExecutedFunc != nil {
		r.ActionExecutedFunc(s, param, state, res, err)
	}
}

func (r *CommonFlowRunner) OnTimeout(s *Session) {
	if r.TimeoutFunc != nil {
		r.TimeoutFunc(s)
	}

}
