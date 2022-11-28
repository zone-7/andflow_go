package engine

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
	"github.com/zone-7/andflow_go/models"
)

type FlowRunner interface {
	ExecuteLink(s *Session, param *models.LinkParam) int     //返回三个状态 -1 不通过，1通过，0还没准备好执行
	ExecuteAction(s *Session, param *models.ActionParam) int //返回三个状态 -1 不通过，1通过，0还没准备好执行
}

type CommonFlowRunner struct {
}

func (r *CommonFlowRunner) ExecuteLink(s *Session, param *models.LinkParam) int {
	link := s.Runtime.Flow.GetLinkBySourceIdAndTargetId(param.SourceId, param.TargetId)
	sc := link.Filter
	if len(strings.Trim(sc, " ")) == 0 {
		return 1
	}

	rts := goja.New()
	rts.Set("flow", s.Runtime.Flow)
	rts.Set("link", link)

	SetScriptFun(rts, s, param.SourceId, param.TargetId, true)

	script := "function $exec(){\n" + sc + "\n}\n $exec();\n"
	val, err := rts.RunString(script)
	if err != nil {

		linkState := s.Runtime.GetLastLinkState(param.SourceId, param.TargetId)
		if linkState != nil {
			linkState.IsError = 1
		}

		return 0
	}

	res := GetScriptIntResult(val)

	return res
}

func (r *CommonFlowRunner) ExecuteAction(s *Session, param *models.ActionParam) int {

	action := s.Runtime.Flow.GetAction(param.ActionId)
	name := action.Name

	fmt.Println("开始执行节点：" + action.Name + " " + action.Title)
	defer fmt.Println("结束执行节点：" + action.Name + " " + action.Title)

	runner := GetActionRunner(name)
	if runner == nil {
		fmt.Println("找不到节点" + name + "的执行器")
		return 1
	}
	res, err := runner.Execute(s, param)

	if err != nil {
		actionState := s.Runtime.GetLastActionState(param.ActionId)
		if actionState != nil {
			actionState.IsError = 1
		}

		return 0
	}

	return res
}
