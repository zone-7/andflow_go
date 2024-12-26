package andflow

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"regexp"
	"strings"
)

func unescapeHTML(s any) template.HTML {

	var str string
	if v, ok := s.(string); ok {
		str = v
	} else {
		switch s.(type) {
		case string:
		case float32:
		case float64:
		case int:
		case int16:
		case int32:
		case int64:
		case bool:
			str = fmt.Sprintf("%v", s)
			break
		default:
			data, err := json.Marshal(s)
			if err == nil {
				str = string(data)
			} else {
				str = fmt.Sprintf("%v", s)
			}
		}
	}

	return template.HTML(str)
}

// 模版替换
func ReplaceTemplate(temp string, name string, params map[string]any) (string, error) {
	if strings.Index(temp, "{{") < 0 || strings.Index(temp, "}}") < 0 {
		return temp, nil
	}
	// 正则表达式,替换前面没有$.或. 的不合规
	need_replaces := make(map[string]string)
	reg := regexp.MustCompile(`\{\{[^{}]*\}\}`)

	items := reg.FindAll([]byte(temp), -1)
	for _, item := range items {

		str := string(item)

		key := strings.Replace(str, "{{", "", -1)
		key = strings.Replace(key, "}}", "", -1)
		key = strings.Trim(key, " ")

		newkey := key
		if strings.Index(key, "$.") != 0 && strings.Index(key, ".") != 0 {
			newkey = "." + key
		}

		need_replaces[str] = "{{ unescapeHTML " + newkey + "}}"

	}
	for o, n := range need_replaces {
		temp = strings.ReplaceAll(temp, o, n)
	}

	//解析模板
	t, err := template.New(name).Funcs(template.FuncMap{"unescapeHTML": unescapeHTML}).Parse(temp)

	if err != nil {
		return temp, err
	}

	b := bytes.NewBuffer(make([]byte, 0))
	bw := bufio.NewWriter(b)

	err = t.Execute(bw, params)
	bw.Flush()

	return string(b.Bytes()), err
}
