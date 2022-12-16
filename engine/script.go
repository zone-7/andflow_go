package engine

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"github.com/dop251/goja"
)

func GetScriptIntResult(val goja.Value) Result {
	if val == goja.Null() || val == goja.NaN() {
		return SUCCESS
	}

	obj := val.Export()
	switch obj.(type) {
	case bool:
		if obj.(bool) {
			return SUCCESS //1
		} else {
			return FAILURE // -1
		}
	case string:
		if obj.(string) == "true" {
			return SUCCESS
		} else if obj.(string) == "false" {
			return FAILURE
		} else if obj.(string) == "1" {
			return SUCCESS
		} else if obj.(string) == "-1" {
			return FAILURE
		} else if obj.(string) == "0" {
			return FAILURE
		}
	case int:
		if obj.(int) > 0 {
			return SUCCESS
		} else if obj.(int) < 0 {
			return -1
		} else {
			return FAILURE
		}
	case int64:
		if obj.(int64) > 0 {
			return SUCCESS
		} else if obj.(int64) < 0 {
			return FAILURE
		} else {
			return FAILURE
		}
	case int32:
		if obj.(int32) > 0 {
			return SUCCESS
		} else if obj.(int32) < 0 {
			return FAILURE
		} else {
			return FAILURE
		}
	case int16:
		if obj.(int16) > 0 {
			return SUCCESS
		} else if obj.(int16) < 0 {
			return FAILURE
		} else {
			return FAILURE
		}
	case float64:
		if obj.(float64) > 0 {
			return SUCCESS
		} else if obj.(float64) < 0 {
			return FAILURE
		} else {
			return FAILURE
		}
	case float32:
		if obj.(float32) > 0 {
			return SUCCESS
		} else if obj.(float32) < 0 {
			return FAILURE
		} else {
			return FAILURE
		}
	default:
		return FAILURE

	}
	return FAILURE
}

//设置脚本函数
func SetCommonScriptFunc(rts *goja.Runtime, session *Session, preActionId string, actionId string, islink bool) {
	tp := ""
	title := ""
	if islink {
		tp = "link"
		title = preActionId + "->" + actionId

	} else {
		tp = "action"
		action := session.GetFlow().GetAction(actionId)
		title = action.Name
	}

	//日志
	rts.Set("log", func(call goja.FunctionCall) goja.Value {
		value := call.Argument(0)

		var val string
		if !value.Equals(goja.Undefined()) && !value.Equals(goja.Null()) {
			var bb bool
			var str string

			if value.ExportType() == reflect.TypeOf(bb) {
				if value.Export().(bool) {
					val = "true"
				} else {
					val = "false"
				}

			} else if value.ExportType() == reflect.TypeOf(str) {
				val = value.String()
			} else {
				b, err := json.Marshal(str)
				if err == nil {
					val = string(b)
				}
			}
		}
		if tp == "link" {
			session.AddLog_link_info(title, val)
		} else {
			session.AddLog_action_info(title, val)
		}

		return value
	})

	//打印
	rts.Set("print", func(call goja.FunctionCall) goja.Value {
		value := call.Argument(0)

		var val string
		if !value.Equals(goja.Undefined()) && !value.Equals(goja.Null()) {
			var bb bool
			var str string

			if value.ExportType() == reflect.TypeOf(bb) {
				if value.Export().(bool) {
					val = "true"
				} else {
					val = "false"
				}

			} else if value.ExportType() == reflect.TypeOf(str) {
				val = value.String()
			} else {
				b, err := json.Marshal(str)
				if err == nil {
					val = string(b)
				}
			}

		}
		fmt.Println(val)

		return value
	})
	//执行命令行
	rts.Set("cmd", func(call goja.FunctionCall) goja.Value {
		arg0 := call.Argument(0)
		arg1 := call.Argument(1)

		var command string
		var timeout int64

		if arg0.Equals(goja.Undefined()) || arg0.Equals(goja.Null()) {
			return goja.Null()
		} else {
			command = arg0.ToString().String()
		}

		if arg1.Equals(goja.Undefined()) || arg1.Equals(goja.Null()) {
			timeout = 3000 //3秒
		} else {
			timeout = arg1.ToInteger()
		}

		res, err := cmd(command, timeout)

		if err != nil {
			log.Println(err)
			return goja.Null()
		}

		return rts.ToValue(res)
	})
	//等待
	rts.Set("sleep", func(call goja.FunctionCall) goja.Value {
		arg0 := call.Argument(0)

		var timeout int64

		if arg0.Equals(goja.Undefined()) || arg0.Equals(goja.Null()) {
			return goja.Null()
		} else {
			timeout = arg0.ToInteger()
		}

		time.Sleep(time.Duration(timeout) * time.Millisecond)

		return goja.Null()
	})

	//保存参数，（缓存数据）
	rts.Set("setParam", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0)
		val := call.Argument(1)

		var keyStr string
		if !key.Equals(goja.Undefined()) && !key.Equals(goja.Null()) {
			keyStr = key.String()
		}

		if !val.Equals(goja.Undefined()) && !val.Equals(goja.Null()) {
			obj := val.Export()
			session.Store.SetParam(keyStr, obj)
		} else {
			session.Store.SetParam(keyStr, nil)

		}

		return val
	})

	rts.Set("getParam", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0)

		var value interface{}
		var keyStr string
		if !key.Equals(goja.Undefined()) && !key.Equals(goja.Null()) {
			keyStr = key.String()
			value = session.Store.GetParam(keyStr)

		}
		if value == nil {
			return goja.Null()
		}

		return rts.ToValue(value)

	})

	//保存存数据
	rts.Set("setData", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0)
		val := call.Argument(1)

		var keyStr string
		if !key.Equals(goja.Undefined()) && !key.Equals(goja.Null()) {
			keyStr = key.String()
		}

		if !val.Equals(goja.Undefined()) && !val.Equals(goja.Null()) {
			obj := val.Export()
			session.Store.SetData(keyStr, obj)
		} else {
			session.Store.SetData(keyStr, nil)

		}

		return val
	})

	rts.Set("getData", func(call goja.FunctionCall) goja.Value {
		key := call.Argument(0)

		var value interface{}
		var keyStr string
		if !key.Equals(goja.Undefined()) && !key.Equals(goja.Null()) {
			keyStr = key.String()
			value = session.Store.GetData(keyStr)

		}
		if value == nil {
			return goja.Null()
		}

		return rts.ToValue(value)

	})

	rts.Set("getDatas", func(call goja.FunctionCall) goja.Value {

		value := session.Store.GetDataMap()
		if value == nil {
			return goja.Null()
		}

		return rts.ToValue(value)
	})
	rts.Set("setDatas", func(call goja.FunctionCall) goja.Value {
		datas := call.Argument(0)

		if !datas.Equals(goja.Undefined()) && !datas.Equals(goja.Null()) {
			obj := datas.Export()
			switch obj.(type) {
			case map[string]string:
				for k, v := range obj.(map[string]string) {

					session.Store.SetData(k, v)
				}
			case map[string]interface{}:
				for k, v := range obj.(map[string]interface{}) {

					session.Store.SetData(k, v)
				}
			case map[string]map[string]interface{}:
				for k, v := range obj.(map[string]map[string]interface{}) {

					session.Store.SetData(k, v)
				}
			case map[string][]map[string]interface{}:
				for k, v := range obj.(map[string][]map[string]interface{}) {

					session.Store.SetData(k, v)
				}
			case map[string][]interface{}:
				for k, v := range obj.(map[string][]interface{}) {

					session.Store.SetData(k, v)
				}
			default:
				return goja.Null()

			}

		}

		return datas
	})

	//缓存数据
	rts.Set("setActionData", func(call goja.FunctionCall) goja.Value {
		if len(actionId) == 0 {
			return goja.Null()
		}

		key := call.Argument(0)
		val := call.Argument(1)

		var keyStr string
		if !key.Equals(goja.Undefined()) && !key.Equals(goja.Null()) {
			keyStr = key.String()
		}

		if !val.Equals(goja.Undefined()) && !val.Equals(goja.Null()) {
			obj := val.Export()
			session.Store.SetActionData(actionId, keyStr, obj)

		} else {
			session.Store.SetActionData(actionId, keyStr, nil)

		}

		return val
	})

	rts.Set("getActionData", func(call goja.FunctionCall) goja.Value {
		if len(actionId) == 0 {
			return goja.Null()
		}

		key := call.Argument(0)

		var value interface{}
		var keyStr string
		if !key.Equals(goja.Undefined()) && !key.Equals(goja.Null()) {
			keyStr = key.String()
			value = session.Store.GetActionData(actionId, keyStr)

		}

		if value == nil {
			return goja.Null()
		}

		return rts.ToValue(value)

	})

	rts.Set("getPreActionData", func(call goja.FunctionCall) goja.Value {
		if len(preActionId) == 0 {
			return goja.Null()
		}

		key := call.Argument(0)

		var value interface{}
		var keyStr string
		if !key.Equals(goja.Undefined()) && !key.Equals(goja.Null()) {
			keyStr = key.String()
			value = session.Store.GetActionData(preActionId, keyStr)

		}

		if value == nil {
			return goja.Null()
		}

		return rts.ToValue(value)

	})

	rts.Set("getActionDatas", func(call goja.FunctionCall) goja.Value {
		if len(actionId) == 0 {
			return goja.Null()
		}

		value := session.Store.GetActionDataMap(actionId)

		if value == nil {
			return goja.Null()
		}

		return rts.ToValue(value)

	})
	rts.Set("getPreActionDatas", func(call goja.FunctionCall) goja.Value {
		if len(preActionId) == 0 {
			return goja.Null()
		}

		value := session.Store.GetActionDataMap(preActionId)

		if value == nil {
			return goja.Null()
		}

		return rts.ToValue(value)

	})

	//JSON操作
	jsonObj := rts.NewObject()
	jsonObj.Set("parse", func(call goja.FunctionCall) goja.Value {
		jsonStr := call.Argument(0)
		d := make(map[string]interface{}, 0)

		err := json.Unmarshal([]byte(jsonStr.String()), &d)
		if err != nil {
			return goja.Null()
		}

		return rts.ToValue(d)
	})
	jsonObj.Set("stringfy", func(call goja.FunctionCall) goja.Value {
		jsonParam := call.Argument(0)
		return rts.ToValue(jsonParam.String())

	})
	rts.Set("json", jsonObj)
	rts.Set("JSON", jsonObj)

	//base64操作
	b64 := rts.NewObject()
	b64.Set("encodeToString", func(call goja.FunctionCall) goja.Value {
		option := call.Argument(0)
		obj := option.Export()
		d, ok := obj.([]byte)
		if ok {
			b64 := base64.StdEncoding.EncodeToString(d)
			return rts.ToValue(b64)
		} else {
			s, ok := obj.(string)
			if ok {
				b64 := base64.StdEncoding.EncodeToString([]byte(s))
				return rts.ToValue(b64)
			}
		}

		return goja.Null()
	})

	b64.Set("encodeToByte", func(call goja.FunctionCall) goja.Value {
		option := call.Argument(0)
		obj := option.Export()
		d, ok := obj.([]byte)
		if ok {
			var r []byte
			base64.StdEncoding.Encode(r, d)
			return rts.ToValue(r)
		} else {
			s, ok := obj.(string)
			if ok {
				var r []byte
				base64.StdEncoding.Encode(r, []byte(s))
				return rts.ToValue(r)
			}
		}

		return goja.Null()
	})

	b64.Set("decodeToString", func(call goja.FunctionCall) goja.Value {
		option := call.Argument(0)
		obj := option.Export()
		d, ok1 := obj.([]byte)
		if ok1 {
			var r []byte
			_, err := base64.StdEncoding.Decode(r, d)
			if err == nil {
				return rts.ToValue(string(r))
			}
		} else {
			d, ok2 := obj.(string)
			if ok2 {
				r, err := base64.StdEncoding.DecodeString(d)
				if err == nil {
					return rts.ToValue(string(r))
				}
			}

		}

		return goja.Null()

	})
	b64.Set("decodeToByte", func(call goja.FunctionCall) goja.Value {
		option := call.Argument(0)
		obj := option.Export()
		d, ok1 := obj.([]byte)
		if ok1 {
			var r []byte
			_, err := base64.StdEncoding.Decode(r, d)
			if err == nil {
				return rts.ToValue(r)
			}
		} else {
			d, ok2 := obj.(string)
			if ok2 {
				r, err := base64.StdEncoding.DecodeString(d)
				if err == nil {
					return rts.ToValue(r)
				}
			}

		}

		return goja.Null()
	})
	rts.Set("base64", b64)

	rts.Set("md5", func(call goja.FunctionCall) goja.Value {
		option := call.Argument(0)

		obj := option.Export()
		d, ok := obj.([]byte)
		if ok {
			m := md5.New().Sum(d)
			return rts.ToValue(string(m))
		} else {
			s, ok := obj.(string)
			if ok {
				m := md5.New().Sum([]byte(s))
				return rts.ToValue(string(m))
			}

		}

		return goja.Null()
	})

	stringsObj := rts.NewObject()
	stringsObj.Set("indexOf", func(call goja.FunctionCall) goja.Value {
		arg1 := call.Argument(0)
		arg2 := call.Argument(1)
		if arg1 == nil || arg2 == nil {
			return rts.ToValue(-1)
		}
		obj1 := arg1.String()
		obj2 := arg2.String()

		idx := strings.Index(obj1, obj2)

		return rts.ToValue(idx)
	})
	stringsObj.Set("substr", func(call goja.FunctionCall) goja.Value {
		arg1 := call.Argument(0)
		arg2 := call.Argument(1)
		arg3 := call.Argument(2)
		if arg1 == nil {
			return rts.ToValue("")
		}
		obj1 := arg1.String()
		if arg2 == nil {
			return rts.ToValue("")
		}
		obj2 := arg2.ToInteger()

		obj3 := len(obj1)
		if arg3 != nil {
			obj3 = int(arg3.ToInteger())
		}

		res := obj1[obj2:obj3]

		return rts.ToValue(res)
	})
	stringsObj.Set("trim", func(call goja.FunctionCall) goja.Value {
		arg1 := call.Argument(0)

		if arg1 == nil {
			return rts.ToValue("")
		}
		obj1 := arg1.String()

		res := strings.Trim(obj1, " ")

		return rts.ToValue(res)
	})
	rts.Set("strings", stringsObj)

}

func cmd(command string, timeout int64) (string, error) {
	var err error
	var res string

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	lines := strings.Split(command, "\n")

	for _, line := range lines {

		if len(strings.Trim(line, " ")) == 0 {
			continue
		}

		commandArr := strings.Split(line, " ")

		name := commandArr[0]
		attr := commandArr[1:]

		cmd := exec.CommandContext(ctx, name, attr...)

		var stdout io.ReadCloser

		if stdout, err = cmd.StdoutPipe(); err != nil { //获取输出对象，可以从该对象中读取输出结果
			return "", err
		}

		defer stdout.Close() // 保证关闭输出流

		if err = cmd.Start(); err != nil { // 运行命令
			return "", err
		}

		var opBytes []byte
		if opBytes, err = ioutil.ReadAll(stdout); err != nil { // 读取输出结果
			return "", err
		}

		res = string(opBytes)

	}

	return res, nil
}
