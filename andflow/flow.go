package andflow

const (
	// 流程节点样式
	THEME_DEFAULT = "flow_theme_default" //默认样式
	THEME_ICON    = "flow_theme_icon"    //图标
	THEME_ZONE    = "flow_theme_zone"    //右标题
	THEME_BOX     = "flow_theme_box"     //方框

	//连接线样式 Flowchart Straight Bezier StateMachine
	Link_TYPE_FLOWCHART    = "Flowchart"
	LINK_TYPE_STRAIGHT     = "Straight"
	LINK_TYPE_BEZIER       = "Bezier"
	LINK_TYPE_STATEMACHINE = "StateMachine"
)

type ActionContent struct {
	ContentType string `bson:"content_type" json:"content_type"` //类型
	Content     string `bson:"content" json:"content"`           //内容文本
}
type ActionModel struct {
	Id              string            `bson:"id" json:"id"`                       //ID
	Name            string            `bson:"name" json:"name"`                   //名称
	Title           string            `bson:"title" json:"title"`                 //标题
	Icon            string            `bson:"icon" json:"icon"`                   //图标
	Des             string            `bson:"des" json:"des"`                     //描述
	Keywords        string            `bson:"keywords" json:"keywords"`           //关键词
	Left            string            `bson:"left" json:"left"`                   //位置X
	Top             string            `bson:"top" json:"top"`                     //位置Y
	Width           string            `bson:"width" json:"width"`                 //宽度
	Height          string            `bson:"height" json:"height"`               //高度
	Theme           string            `bson:"theme" json:"theme"`                 //样式
	BodyWidth       string            `bson:"body_width" json:"body_width"`       //内容宽度
	BodyHeight      string            `bson:"body_height" json:"body_height"`     //内容高度
	Params          map[string]string `bson:"params" json:"params"`               //运行配置参数
	Collect         string            `bson:"collect" json:"collect"`             //true：所有路线到达节点后才执行；false：任意一个路径到达都执行
	Once            string            `bson:"once" json:"once"`                   //是否在整个过程只执行一次： true，false
	IteratorList    string            `bson:"iterator_list" json:"iterator_list"` //迭代列表参数名，根据列表轮询执行
	IteratorItem    string            `bson:"iterator_item" json:"iterator_item"` //迭代元素参数名
	ScriptBefore    string            `bson:"script_before" json:"script_before"` //执行过滤脚本
	ScriptAfter     string            `bson:"script_after" json:"script_after"`   //执行内容脚本
	ScriptError     string            `bson:"script_error" json:"script_error"`   //执行异常处理脚本
	Content         *ActionContent    `bson:"content" json:"content"`             //显示内容
	BorderColor     string            `bson:"border_color" json:"border_color"`
	BodyColor       string            `bson:"body_color" json:"body_color"`
	BodyTextColor   string            `bson:"body_text_color" json:"body_text_color"`
	HeaderColor     string            `bson:"header_color" json:"header_color"`
	HeaderTextColor string            `bson:"header_text_color" json:"header_text_color"`
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

	m.Content = &ActionContent{ContentType: content_type, Content: content}

}

func (m *ActionModel) GetContent() (string, string) {
	if m.Content == nil {
		return "", ""
	}
	return m.Content.ContentType, m.Content.Content
}

type LinkModel struct {
	Name           string `bson:"name" json:"name"`                       //自定义名称
	Title          string `bson:"title" json:"title"`                     //标题
	Keywords       string `bson:"keywords" json:"keywords"`               //关键词
	SourceId       string `bson:"source_id" json:"source_id"`             //源ID
	SourcePosition string `bson:"source_position" json:"source_position"` //源位置
	TargetId       string `bson:"target_id" json:"target_id"`             //目的ID
	TargetPosition string `bson:"target_position" json:"target_position"` //目的位置
	Filter         string `bson:"filter" json:"filter"`                   //过滤脚本
	Active         string `bson:"active" json:"active"`                   //是否是一个可通的路径，true,false
	LabelSource    string `bson:"label_source" json:"label_source"`       //头标签
	LabelTarget    string `bson:"label_target" json:"label_target"`       //尾标签
	Animation      bool   `bson:"animation" json:"animation"`             //是否动画
	Arrows         []bool `bson:"arrows" json:"arrows"`                   //是否显示头中间和尾箭头
}

type GroupModel struct {
	Id      string   `bson:"id" json:"id"`           //组ID
	Title   string   `bson:"title" json:"title"`     //标题
	Theme   string   `bson:"theme" json:"theme"`     //样式
	Des     string   `bson:"des" json:"des"`         //描述
	Members []string `bson:"members" json:"members"` //成员ID
	Left    string   `bson:"left" json:"left"`       //坐标X
	Top     string   `bson:"top" json:"top"`         //坐标Y
	Width   string   `bson:"width" json:"width"`     //宽度
	Height  string   `bson:"height" json:"height"`   //高度

	BorderColor string `bson:"border_color" json:"border_color"`

	BodyColor     string `bson:"body_color" json:"body_color"`
	BodyTextColor string `bson:"body_text_color" json:"body_text_color"`

	HeaderColor     string `bson:"header_color" json:"header_color"`
	HeaderTextColor string `bson:"header_text_color" json:"header_text_color"`
}

type ListModel struct {
	Id     string   `bson:"id" json:"id"`         //组ID
	Title  string   `bson:"title" json:"title"`   //标题
	Theme  string   `bson:"theme" json:"theme"`   //样式
	Des    string   `bson:"des" json:"des"`       //描述
	Items  []string `bson:"items" json:"items"`   //成员ID
	Left   string   `bson:"left" json:"left"`     //坐标X
	Top    string   `bson:"top" json:"top"`       //坐标Y
	Width  string   `bson:"width" json:"width"`   //宽度
	Height string   `bson:"height" json:"height"` //高度

	BorderColor string `bson:"border_color" json:"border_color"`

	BodyColor     string `bson:"body_color" json:"body_color"`
	BodyTextColor string `bson:"body_text_color" json:"body_text_color"`

	HeaderColor     string `bson:"header_color" json:"header_color"`
	HeaderTextColor string `bson:"header_text_color" json:"header_text_color"`

	ItemColor     string `bson:"item_color" json:"item_color"`
	ItemTextColor string `bson:"item_text_color" json:"item_text_color"`
}

type TipModel struct {
	Id              string `bson:"id" json:"id"`           //组ID
	Title           string `bson:"title" json:"title"`     //标题
	Theme           string `bson:"theme" json:"theme"`     //样式
	Des             string `bson:"des" json:"des"`         //描述
	Content         string `bson:"content" json:"content"` //内容
	Left            string `bson:"left" json:"left"`       //坐标X
	Top             string `bson:"top" json:"top"`         //坐标Y
	Width           string `bson:"width" json:"width"`     //宽度
	Height          string `bson:"height" json:"height"`   //高度
	BorderColor     string `bson:"border_color" json:"border_color"`
	BodyColor       string `bson:"body_color" json:"body_color"`
	BodyTextColor   string `bson:"body_text_color" json:"body_text_color"`
	HeaderColor     string `bson:"header_color" json:"header_color"`
	HeaderTextColor string `bson:"header_text_color" json:"header_text_color"`
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
	Lists               []*ListModel      `bson:"lists" json:"lists"`                                   //LIST
	Tips                []*TipModel       `bson:"tips" json:"tips"`                                     //TIP

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

// 获取作为起点的节点
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

// 获取组内的起点节点
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

func CreateFlowModel(code string, name string) *FlowModel {
	model := &FlowModel{}
	model.Actions = make([]*ActionModel, 0)
	model.Groups = make([]*GroupModel, 0)
	model.Lists = make([]*ListModel, 0)
	model.Tips = make([]*TipModel, 0)
	model.Links = make([]*LinkModel, 0)
	model.Dict = make([]*FlowDictModel, 0)
	model.Params = make([]*FlowParamModel, 0)
	model.ShowActionBody = "true"
	model.ShowActionContent = "true"
	model.ShowActionState = "true"
	model.Theme = THEME_DEFAULT
	model.LinkType = Link_TYPE_FLOWCHART

	model.Code = code
	model.Name = name

	return model
}
