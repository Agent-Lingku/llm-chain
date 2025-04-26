package main

import (
	"fmt"
	"log"

	"learn/internal/chain"
	"learn/internal/config"
)

func main() {
	// 初始化配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}
	log.Printf("应用启动配置: %+v", cfg)

	ch := chain.NewChain()

	// 添加处理类
	ch.AddHandler(chain.NewRequester()).
		AddHandler(chain.NewThinker()).
		AddHandler(chain.NewTaskPublisher()).
		AddHandler(chain.NewTaskExecutor()).
		AddHandler(chain.NewTaskCollector())

	// 创建请求
	request := &chain.Request{
		Message: "写一个学校官网的主页",
		Data:    make(map[string]any),
	}

	// 处理请求
	result := ch.HandleRequest(request)

	// 输出结果
	fmt.Println("处理结果:", result.Data)
}
