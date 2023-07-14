package andflow

type ActionModel struct {
	Id         string            `bson:"id" json:"id"`                   //ID
	Name       string            `bson:"name" json:"name"`               //名称
	Title      string            `bson:"title" json:"title"`             //标题
	Icon       string            `bson:"icon" json:"icon"`               //图标
	Des        string            `bson:"des" json:"des"`                 //描述
	Keywords   string            `bson:"keywords" json:"keywords"`       //关键词
	Left       string            `bson:"left" json:"left"`               //位置X
	Top        string            `bson:"top" json:"top"`                 //位置Y
	Width      string            `bson:"width" json:"width"`             //宽度
	Height     string            `bson:"height" json:"height"`           //高度
	Theme      string            `bson:"theme" json:"theme"`             //样式
	BodyWidth  string            `bson:"body_width" json:"body_width"`   //内容宽度
	BodyHeight string            `bson:"body_height" json:"body_height"` //内容高度
	Params     map[string]string `bson:"params" json:"params"`           //参数
	Collect    string            `bson:"collect" json:"collect"`         //true：所有路线到达节点后才执行；false：任意一个路径到达都执行
	Once       string            `bson:"once" json:"once"`               //是否只执行一次： true，false
	Filter     string            `bson:"filter" json:"filter"`           //执行过滤脚本
	Script     string            `bson:"script" json:"script"`           //执行内容脚本
	Error      string            `bson:"error" json:"error"`             //执行异常处理脚本
	Content    map[string]string `bson:"content" json:"content"`         //内容
}

func (m *ActionModel) GetParam(name string) string {
	if m.Params != nil {
		return m.Params[name]
	}
	return ""
}
func (m *ActionModel) SetParam(name string, value string) {
	if m.Params == nil {
		m.Params = make(map[string]string)
	}

	m.Params[name] = value
}
func (m *ActionModel) SetContent(content_type string, content string) {
	m.Content = map[string]string{"content_type": content_type, "content": content}
}
func (m *ActionModel) GetContent() (string, string) {
	if m.Content == nil {
		return "", ""
	}
	return m.Content["content_type"], m.Content["content"]
}

type LinkModel struct {
	Title          string `bson:"title" json:"title"`                     //标题
	Keywords       string `bson:"keywords" json:"keywords"`               //关键词
	SourceId       string `bson:"source_id" json:"source_id"`             //源ID
	SourcePosition string `bson:"source_position" json:"source_position"` //源位置
	TargetId       string `bson:"target_id" json:"target_id"`             //目的ID
	TargetPosition string `bson:"target_position" json:"target_position"` //目的位置
	Filter         string `bson:"filter" json:"filter"`                   //过滤脚本
	Active         string `bson:"active" json:"active"`                   //是否是一个可通的路径，true,false
}

type GroupModel struct {
	Id      string   `bson:"id" json:"id"`           //组ID
	Title   string   `bson:"title" json:"title"`     //标题
	Des     string   `bson:"des" json:"des"`         //描述
	Members []string `bson:"members" json:"members"` //成员ID
	Left    string   `bson:"left" json:"left"`       //坐标X
	Top     string   `bson:"top" json:"top"`         //坐标Y
	Width   string   `bson:"width" json:"width"`     //宽度
	Height  string   `bson:"height" json:"height"`   //高度
}

type FlowParamModel struct {
	Name  string `bson:"name" json:"name"`
	Value string `bson:"value" json:"value"`
}
type FlowDictModel struct {
	Name  string `bson:"name" json:"name"`
	Label string `bson:"label" json:"label"`
}

type FlowModel struct {
	Code                string            `bson:"code" json:"code"`
	Name                string            `bson:"name" json:"name"`                                     //流程名称
	FlowType            string            `bson:"flow_type" json:"flow_type"`                           //流程类型
	Timeout             string            `bson:"timeout" json:"timeout"`                               //执行时效
	CacheTimeout        string            `bson:"cache_timeout" json:"cache_timeout"`                   //缓存时效
	Params              []*FlowParamModel `bson:"params" json:"params"`                                 //运行参数列表
	Dict                []*FlowDictModel  `bson:"dict" json:"dict"`                                     //运行字典列表
	LogType             string            `bson:"log_type" json:"log_type"`                             //日志类型
	StoreType           string            `bson:"store_type" json:"store_type"`                         //是否保存运行时内容
	ShowActionState     string            `bson:"show_action_state" json:"show_action_state"`           //是否显示节点运行状态
	ShowActionBody      string            `bson:"show_action_body" json:"show_action_body"`             //是否显示节点Body
	ShowActionContent   string            `bson:"show_action_content" json:"show_action_content"`       //是否显示节点运行消息内容
	ShowActionStateList string            `bson:"show_action_state_list" json:"show_action_state_list"` //是否显示节点运行状态列表
	Theme               string            `bson:"theme" json:"theme"`                                   //皮肤
	LinkType            string            `bson:"link_type" json:"link_type"`                           //连接线类型
	Actions             []*ActionModel    `bson:"actions" json:"actions"`                               //节点列表
	Links               []*LinkModel      `bson:"links" json:"links"`                                   //连线列表
	Groups              []*GroupModel     `bson:"groups" json:"groups"`                                 //组

}

func (t *FlowModel) GetDict(name string) *FlowDictModel {
	if t.Dict == nil || len(t.Dict) == 0 {
		return nil
	}

	for _, d := range t.Dict {
		if d.Name == name {
			return d
		}
	}
	return nil
}

func (t *FlowModel) GetActionModel(id string) *ActionModel {
	for _, node := range t.Actions {
		if node.Id == id {
			return node
		}
	}
	return nil
}
func (t *FlowModel) GetLinkByTargetId(id string) []*LinkModel {
	lks := make([]*LinkModel, 0)
	for _, link := range t.Links {
		if link.TargetId == id {
			lks = append(lks, link)
		}
	}
	return lks
}

func (t *FlowModel) GetLinkBySourceId(id string) []*LinkModel {
	lks := make([]*LinkModel, 0)

	for _, link := range t.Links {
		if link.SourceId == id {
			lks = append(lks, link)
		}
	}
	return lks
}

func (t *FlowModel) GetLinkBySourceIdAndTargetId(sid string, tid string) *LinkModel {

	for _, link := range t.Links {
		if link.SourceId == sid && link.TargetId == tid {
			return link
		}
	}
	return nil
}

func (t *FlowModel) GetGroup(groupId string) *GroupModel {

	for _, g := range t.Groups {

		if g.Id == groupId {
			return g
		}
	}
	return nil
}

func (t *FlowModel) GetAction(actionId string) *ActionModel {

	for _, a := range t.Actions {

		if a.Id == actionId {
			return a
		}
	}
	return nil
}

func (t *FlowModel) GetActionIdsInGroup(groupId string) []string {

	g := t.GetGroup(groupId)
	if g != nil {
		ids := make([]string, 0)
		for _, mid := range g.Members {
			if t.GetAction((mid)) != nil {
				ids = append(ids, mid)
			}
		}
		return ids
	}
	return nil

}

//获取作为起点的节点
func (t *FlowModel) GetStartActionIds() []string {

	actions_start := make([]string, 0)

	for _, action := range t.Actions {
		targets := t.GetLinkByTargetId(action.Id)
		if targets == nil || len(targets) == 0 {
			actions_start = append(actions_start, action.Id)
		}
	}

	return actions_start
}

//获取组内的起点节点
func (t *FlowModel) GetStartActionIdsInGroup(groupId string) []string {

	actions_start := make([]string, 0)

	g := t.GetGroup(groupId)
	if g == nil {
		return nil
	}

	for _, memberId := range g.Members {

		links := t.GetLinkByTargetId(memberId)
		if links == nil || len(links) == 0 {
			actions_start = append(actions_start, memberId)
		}
	}

	return actions_start
}
