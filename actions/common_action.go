package actions

import (
	"strings"

	"github.com/dop251/goja"
	"github.com/zone-7/andflow_go/engine"
	"github.com/zone-7/andflow_go/models"
)

type CommonActionRunner struct {
}

func (a *CommonActionRunner) Execute(s *engine.Session, param *models.ActionParam) (int, error) {

	action := s.GetFlow().GetAction(param.ActionId)

	sc := action.Script
	if len(strings.Trim(sc, " ")) == 0 {
		return 1, nil
	}

	rts := goja.New()
	rts.Set("flow", s.GetFlow())
	rts.Set("action", action)

	engine.SetScriptFun(rts, s, param.PreActionId, param.ActionId, false)

	script := "function $exec(){\n" + sc + "\n}\n $exec();\n"
	val, err := rts.RunString(script)
	if err != nil {
		return 0, err
	}

	obj := engine.GetScriptIntResult(val)

	return obj, nil
}
