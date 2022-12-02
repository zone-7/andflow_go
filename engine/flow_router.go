package engine

import (
	"github.com/zone-7/andflow_go/models"
)

type FlowRouter interface {
	RouteLink(s *Session, param *models.LinkParam) int
	RouteAction(s *Session, param *models.ActionParam) int
}

type CommonFlowRouter struct {
}

func (r *CommonFlowRouter) RouteLink(s *Session, param *models.LinkParam) int {
	s.PushLink(param)
	return 1
}

func (r *CommonFlowRouter) RouteAction(s *Session, param *models.ActionParam) int {
	s.PushAction(param)
	return 1
}
