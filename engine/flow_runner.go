package engine

import (
	"fmt"
	"log"
	"strings"

	"github.com/dop251/goja"
)

type FlowRunner interface {
	ExecuteLink(s *Session, param *LinkParam, state *LinkStateModel) (Result, error)       //返回三个状态 -1 不通过，1通过，0还没准备好执行
	ExecuteAction(s *Session, param *ActionParam, state *ActionStateModel) (Result, error) //返回三个状态 -1 不通过，1通过，0还没准备好执行
}

type CommonFlowRunner struct {
	actionScriptFunc func(rts *goja.Runtime, s *Session, param *ActionParam, state *ActionStateModel)
	linkScriptFunc   func(rts *goja.Runtime, s *Session, param *LinkParam, state *LinkStateModel)
}

func (r *CommonFlowRunner) SetActionScript(act func(rts *goja.Runtime, s *Session, param *ActionParam, state *ActionStateModel)) {
	r.actionScriptFunc = act
}
func (r *CommonFlowRunner) SetLinkScript(act func(rts *goja.Runtime, s *Session, param *LinkParam, state *LinkStateModel)) {
	r.linkScriptFunc = act
}
func (r *CommonFlowRunner) ExecuteLink(s *Session, param *LinkParam, state *LinkStateModel) (Result, error) {
	link := s.GetFlow().GetLinkBySourceIdAndTargetId(param.SourceId, param.TargetId)
	sc := link.Filter
	if len(strings.Trim(sc, " ")) == 0 {
		return SUCCESS, nil
	}

	rts := goja.New()
	rts.Set("flow", s.GetFlow())
	rts.Set("link", link)

	r.linkScriptFunc(rts, s, param, state)
	SetCommonScriptFunc(rts, s, param.SourceId, param.TargetId, true)

	script := "function $exec(){\n" + sc + "\n}\n $exec();\n"
	val, err := rts.RunString(script)

	if err != nil {
		return FAILURE, err
	}

	res := GetScriptIntResult(val)

	return res, nil
}

func (r *CommonFlowRunner) ExecuteAction(s *Session, param *ActionParam, state *ActionStateModel) (Result, error) {
	var res Result = SUCCESS
	var err error

	action := s.GetFlow().GetAction(param.ActionId)
	name := action.Name

	log.Println("开始执行节点：" + action.Name + " " + action.Title)
	defer log.Println("结束执行节点：" + action.Name + " " + action.Title)

	//0.准备脚本执行环境
	rts := goja.New()
	rts.Set("flow", s.GetFlow())
	rts.Set("action", action)

	r.actionScriptFunc(rts, s, param, state)
	SetCommonScriptFunc(rts, s, param.PreActionId, param.ActionId, false)

	//1.执行过滤脚本
	if len(strings.Trim(action.Filter, " ")) > 0 {

		script_filter := "function $filter(){\n" + action.Filter + "\n}\n $filter();\n"

		val, err := rts.RunString(script_filter)
		if err != nil {
			log.Println(fmt.Sprintf("执行脚本异常：%v", err))
			return FAILURE, err
		}

		res := GetScriptIntResult(val)

		if res != SUCCESS {
			return res, nil
		}

	}

	//2.执行节点执行器
	runner := GetActionRunner(name)
	if runner != nil {
		res, err = runner.Execute(s, param, state)
		if err != nil {
			return FAILURE, err
		}
		if res != SUCCESS {
			return res, nil
		}
	}

	//3.执行事后脚本
	if len(strings.Trim(action.Script, " ")) > 0 {

		script_filter := "function $exec(){\n" + action.Script + "\n}\n $exec();\n"

		val, err := rts.RunString(script_filter)
		if err != nil {

			log.Println(fmt.Sprintf("执行脚本异常：%v", err))
			return FAILURE, err
		}

		res := GetScriptIntResult(val)

		if res != SUCCESS {
			return res, nil
		}
	}

	return res, nil
}
