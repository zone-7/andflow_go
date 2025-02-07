package andflow

import (
	"strings"

	"github.com/dop251/goja"
)

var actionRunnerMap map[string]ActionRunner = make(map[string]ActionRunner)

type Result int

const (
	RESULT_SUCCESS Result = 1  //执行通过
	RESULT_REJECT  Result = 0  //没有执行
	RESULT_FAILURE Result = -1 //执行失败
)

type Prop struct {
	Name     string `json:"name" yaml:"name"`
	Default  string `json:"default" yaml:"default"`
	Label    string `json:"label" yaml:"label"`
	Required bool   `json:"required" yaml:"required"`
}

type ActionRunner interface {
	Properties() []Prop
	Execute(s *Session, param *ActionParam, state *ActionStateModel) (Result, error)
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

func GetActionRunners() map[string]ActionRunner {
	return actionRunnerMap
}

type ScriptActionRunner struct {
	funcs map[string]func(s *Session, param *ActionParam, args ...interface{}) interface{}
}

func (a *ScriptActionRunner) SetActionFunc(name string, act func(s *Session, param *ActionParam, args ...interface{}) interface{}) {
	if a.funcs == nil {
		a.funcs = make(map[string]func(s *Session, param *ActionParam, args ...interface{}) interface{})
	}
	a.funcs[name] = act
}

func (a *ScriptActionRunner) Properties() []Prop {
	return []Prop{}
}

func (a *ScriptActionRunner) Execute(s *Session, param *ActionParam, state *ActionStateModel) (Result, error) {

	action := s.GetFlow().GetAction(param.ActionId)

	sc := action.ScriptAfter
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

	SetCommonActionScriptFunc(rts, s, param, state)

	script := "function $exec(){\n" + sc + "\n}\n $exec();\n"
	val, err := rts.RunString(script)
	if err != nil {
		return 0, err
	}

	obj := GetScriptIntResult(val)

	return obj, nil
}
