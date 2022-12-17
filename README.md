这个项目使用golang开发，用于执行andflow的流程设计结果json文件。

<img src="./doc/arch.png"/>

#### 目前开源版本支持单机执行，需要支持集群并行计算版本请联系作者。
<img src="./doc/wx.png"/>


## 入门例子
```

import (
	"fmt"
	"io/ioutil"

	"github.com/zone-7/andflow_go/actions"
	"github.com/zone-7/andflow_go/engine"
)

func main() {
	file := "4执行脚本.json"

	engine.RegistActionRunner("common", &actions.CommonActionRunner{})

	param := make(map[string]interface{})

	data, _ := ioutil.ReadFile(file)

	flow, err := engine.ParseFlow(string(data))
	if err != nil {
		fmt.Println(err)
		return
	}

	runtime := engine.CreateRuntime(flow, param)

	runner := engine.CommonFlowRunner{}
	engine.Execute(runtime, &runner, 10000)

	fmt.Println("time used(ms):", runtime.Timeused)

}


```

## 执行步骤如下：
```
    
    //1. 设置参数
	param := make(map[string]interface{})
    //读取流程json文件 
	data, _ := ioutil.ReadFile(demo_path + "/5执行命令.json")
    //2. 解析json文件
	flow, err := engine.ParseFlow(string(data))
	if err != nil {
		fmt.Println(err)
	}
 
    //3. 创建一个运行时
	runtime := engine.CreateRuntime(flow, param)
     
    //4. 执行流程，超时时间是10000毫秒
	engine.Execute(runtime, 10000)
      

```

