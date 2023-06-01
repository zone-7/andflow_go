package andflow

type FlowRouter interface {
	RouteLink(s *Session, param *LinkParam) int
	RouteAction(s *Session, param *ActionParam) int
}

type CommonFlowRouter struct {
}

func (r *CommonFlowRouter) RouteLink(s *Session, param *LinkParam) int {
	s.PushLink(param)
	return 1
}

func (r *CommonFlowRouter) RouteAction(s *Session, param *ActionParam) int {
	s.PushAction(param)
	return 1
}
