package engine

import (
	"log"
	"strings"

	"github.com/dop251/goja"
)

type FlowRunner interface {
	ExecuteLink(s *Session, param *LinkParam) Result     //返回三个状态 -1 不通过，1通过，0还没准备好执行
	ExecuteAction(s *Session, param *ActionParam) Result //返回三个状态 -1 不通过，1通过，0还没准备好执行
}

type CommonFlowRunner struct {
	funcs map[string]func(s *Session, args ...interface{}) interface{}
}

func (r *CommonFlowRunner) SetActionFunc(name string, act func(s *Session, args ...interface{}) interface{}) {
	if r.funcs == nil {
		r.funcs = make(map[string]func(s *Session, args ...interface{}) interface{})
	}
	r.funcs[name] = act
}

func (r *CommonFlowRunner) ExecuteLink(s *Session, param *LinkParam) Result {
	link := s.GetFlow().GetLinkBySourceIdAndTargetId(param.SourceId, param.TargetId)
	sc := link.Filter
	if len(strings.Trim(sc, " ")) == 0 {
		return 1
	}

	rts := goja.New()
	rts.Set("flow", s.GetFlow())
	rts.Set("link", link)

	if r.funcs != nil {
		for name, f := range r.funcs {
			rts.Set(name, func(call goja.FunctionCall) goja.Value {
				values := call.Arguments
				args := make([]interface{}, 0)
				for _, value := range values {
					if value.Equals(goja.Undefined()) || value.Equals(goja.Null()) {
						args = append(args, nil)
					} else {
						args = append(args, value.Export())
					}
				}
				res := f(s, args)
				if res == nil {
					return goja.Null()
				}
				return rts.ToValue(res)
			})
		}
	}

	SetCommonScriptFunc(rts, s, param.SourceId, param.TargetId, true)

	script := "function $exec(){\n" + sc + "\n}\n $exec();\n"
	val, err := rts.RunString(script)
	if err != nil {

		linkState := s.Store.GetLastLinkState(param.SourceId, param.TargetId)
		if linkState != nil {
			linkState.IsError = 1
		}

		return 0
	}

	res := GetScriptIntResult(val)

	return res
}

func (r *CommonFlowRunner) ExecuteAction(s *Session, param *ActionParam) Result {

	action := s.GetFlow().GetAction(param.ActionId)
	name := action.Name

	log.Println("开始执行节点：" + action.Name + " " + action.Title)
	defer log.Println("结束执行节点：" + action.Name + " " + action.Title)

	runner := GetActionRunner(name)
	if runner == nil {
		log.Println("找不到节点" + name + "的执行器,继续执行")
		return SUCCESS
	}
	res, err := runner.Execute(s, param)

	if err != nil {
		actionState := s.Store.GetLastActionState(param.ActionId)
		if actionState != nil {
			actionState.IsError = 1
		}

		return FAILURE
	}

	return res
}
