package engine

import (
	"strings"

	"github.com/dop251/goja"
)

var actionRunnerMap map[string]ActionRunner = make(map[string]ActionRunner)

type Result int

const (
	SUCCESS Result = 1  //执行通过
	REJECT  Result = 0  //没有执行
	FAILURE Result = -1 //执行失败
)

type ActionRunner interface {
	Execute(s *Session, param *ActionParam) (Result, error)
}

func RegistActionRunner(name string, runner ActionRunner) {
	actionRunnerMap[name] = runner
}
func GetActionRunner(name string) ActionRunner {
	runner := actionRunnerMap[name]
	if runner == nil {
		runner = actionRunnerMap["common"]
	}
	if runner == nil {
		runner = actionRunnerMap["*"]
	}
	if runner == nil {
		runner = actionRunnerMap[""]
	}

	return runner
}

type ScriptActionRunner struct {
	funcs map[string]func(s *Session, param *ActionParam, args ...interface{}) interface{}
}

func (a *ScriptActionRunner) SetScriptFunc(name string, act func(s *Session, param *ActionParam, args ...interface{}) interface{}) {
	if a.funcs == nil {
		a.funcs = make(map[string]func(s *Session, param *ActionParam, args ...interface{}) interface{})
	}
	a.funcs[name] = act
}

func (a *ScriptActionRunner) Execute(s *Session, param *ActionParam) (Result, error) {

	action := s.GetFlow().GetAction(param.ActionId)

	sc := action.Script
	if len(strings.Trim(sc, " ")) == 0 {
		return 1, nil
	}

	rts := goja.New()

	if a.funcs != nil {
		for name, f := range a.funcs {
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

				res := f(s, param, args)
				if res == nil {
					return goja.Null()
				}
				return rts.ToValue(res)
			})
		}
	}

	rts.Set("flow", s.GetFlow())
	rts.Set("action", action)

	SetCommonScriptFunc(rts, s, param.PreActionId, param.ActionId, false)

	script := "function $exec(){\n" + sc + "\n}\n $exec();\n"
	val, err := rts.RunString(script)
	if err != nil {
		return 0, err
	}

	obj := GetScriptIntResult(val)

	return obj, nil
}
