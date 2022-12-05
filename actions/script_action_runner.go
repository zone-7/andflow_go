package actions

import (
	"strings"

	"github.com/dop251/goja"
	"github.com/zone-7/andflow_go/engine"
	"github.com/zone-7/andflow_go/models"
)

type ScriptActionRunner struct {
	funcs map[string]func(args ...interface{}) interface{}
}

func (a *ScriptActionRunner) SetActionFunc(name string, act func(args ...interface{}) interface{}) {
	if a.funcs == nil {
		a.funcs = make(map[string]func(args ...interface{}) interface{})
	}
	a.funcs[name] = act
}

func (a *ScriptActionRunner) Execute(s *engine.Session, param *models.ActionParam) (int, error) {

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
				res := f(args)
				return rts.ToValue(res)
			})
		}
	}

	rts.Set("flow", s.GetFlow())
	rts.Set("action", action)

	engine.SetScriptFunc(rts, s, param.PreActionId, param.ActionId, false)

	script := "function $exec(){\n" + sc + "\n}\n $exec();\n"
	val, err := rts.RunString(script)
	if err != nil {
		return 0, err
	}

	obj := engine.GetScriptIntResult(val)

	return obj, nil
}
