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

func (r *CommonFlowRunner) SetScriptFunc(name string, act func(s *Session, args ...interface{}) interface{}) {
	if r.funcs == nil {
		r.funcs = make(map[string]func(s *Session, args ...interface{}) interface{})
	}
	r.funcs[name] = act
}

func (r *CommonFlowRunner) ExecuteLink(s *Session, param *LinkParam) Result {
	link := s.GetFlow().GetLinkBySourceIdAndTargetId(param.SourceId, param.TargetId)
	sc := link.Filter
	if len(strings.Trim(sc, " ")) == 0 {
		return SUCCESS
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
		return FAILURE
	}

	res := GetScriptIntResult(val)

	return res
}

func (r *CommonFlowRunner) ExecuteAction(s *Session, param *ActionParam) Result {
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
				r := f(s, args)
				if r == nil {
					return goja.Null()
				}
				return rts.ToValue(r)
			})
		}
	}

	SetCommonScriptFunc(rts, s, param.PreActionId, param.ActionId, false)

	//1.执行过滤脚本
	if len(strings.Trim(action.Filter, " ")) > 0 {

		script_filter := "function $filter(){\n" + action.Filter + "\n}\n $filter();\n"

		val, err := rts.RunString(script_filter)
		if err != nil {
			return FAILURE
		}

		res := GetScriptIntResult(val)

		if res != SUCCESS {
			return res
		}

	}

	//2.执行节点执行器
	runner := GetActionRunner(name)
	if runner != nil {
		res, err = runner.Execute(s, param)
		if err != nil {
			return FAILURE
		}
		if res != SUCCESS {
			return res
		}
	}

	//3.执行事后脚本
	if len(strings.Trim(action.Script, " ")) > 0 {

		script_filter := "function $exec(){\n" + action.Script + "\n}\n $exec();\n"

		val, err := rts.RunString(script_filter)
		if err != nil {
			return FAILURE
		}

		res := GetScriptIntResult(val)

		if res != SUCCESS {
			return res
		}
	}

	return res
}

func (r *CommonFlowRunner) exeActionScript(s *Session, param *ActionParam, sc string) Result {
	action := s.GetFlow().GetAction(param.ActionId)

	rts := goja.New()
	rts.Set("flow", s.GetFlow())
	rts.Set("action", action)

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

	SetCommonScriptFunc(rts, s, param.PreActionId, param.ActionId, false)

	script := "function $exec(){\n" + sc + "\n}\n $exec();\n"
	val, err := rts.RunString(script)
	if err != nil {

		return FAILURE
	}

	res := GetScriptIntResult(val)

	return res
}
