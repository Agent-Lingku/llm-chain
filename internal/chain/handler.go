package chain

import (
	"fmt"
	"learn/internal/agent"
	"learn/internal/gen"
	"learn/internal/util"
	"log"
	"os"
)

// Request 请求上下文
type Request struct {
	Message string
	Data    map[string]any
}

// Handler 处理接口
type Handler interface {
	SetNext(Handler) Handler
	Handle(request *Request) *Request
	GetName() string
}

// BaseHandler 基础处理类
type BaseHandler struct {
	next Handler
	name string
}

func (h *BaseHandler) SetNext(next Handler) Handler {
	h.next = next
	return next
}

func (h *BaseHandler) GetName() string {
	return h.name
}

func (h *BaseHandler) Handle(request *Request) *Request {
	if h.next != nil {
		return h.next.Handle(request)
	}
	return request
}

// NewBaseHandler 创建基础处理类
func NewBaseHandler(name string) *BaseHandler {
	return &BaseHandler{name: name}
}

// Requester 创建Agent链条
type Requester struct {
	BaseHandler
}

// NewRequester 新建一个需求分析者
func NewRequester() *Requester {
	return &Requester{BaseHandler: *NewBaseHandler("Requester")}
}

var (
	// 新建模型
	ollama = gen.NewLocalLargeModelClient()
)

func (h *Requester) Handle(request *Request) *Request {
	app := agent.NewAgent(
		agent.WithTaskID("1"),
		agent.WithAgentName("需求分析者"),
		agent.WithModel("qwen2.5-coder:1.5b"),
		agent.WithRole(agent.DemandAnalysisRole),
		agent.WithUserPrompt(request.Message),
	)

	fmt.Println(h.GetName(), "处理请求:", request.Message)

	toolCalls, result, err := app.ExecuteTask(ollama.Client)

	request.Data["Requester"] = map[string]interface{}{
		"tool_calls": toolCalls,
		"data":       result,
	}

	if err != nil {
		if reqMap, ok := request.Data["Requester"].(map[string]interface{}); ok {
			reqMap["err"] = err.Error()
		} else {
			log.Printf("Requester数据格式错误: %T", request.Data["Requester"])
		}
	}

	log.Printf("处理完成: %s\n", request.Message)
	return h.BaseHandler.Handle(request)
}

type Thinker struct {
	BaseHandler
}

func NewThinker() *Thinker {
	return &Thinker{BaseHandler: *NewBaseHandler("Thinker")}
}

func (h *Thinker) Handle(request *Request) *Request {
	app := agent.NewAgent(
		agent.WithTaskID("2"),
		agent.WithAgentName("前端工程师"),
		agent.WithModel("qwen2.5-coder:1.5b"),
		agent.WithRole(agent.FrontEndRole),
		agent.WithUserPrompt("请给我完整代码，不允许省略。"),
	)

	toolCalls, result, err := app.ExecuteTask(ollama.Client, util.AppendUserPrompt(
		request.Data["Requester"].(map[string]interface{})["data"].(string),
	))

	request.Data["Thinker"] = map[string]interface{}{
		"tool_calls": toolCalls,
		"data":       result,
	}

	if err != nil {
		if thinkMap, ok := request.Data["Thinker"].(map[string]interface{}); ok {
			thinkMap["err"] = err.Error()
		} else {
			log.Printf("Thinker数据格式错误: %T", request.Data["Thinker"])
		}
	}

	fmt.Println("前端工程师", "处理结果:", result)
	codeBlocks := util.ExtractCodeBlocks(result)
	if len(codeBlocks) > 0 {
		err := os.WriteFile("demo.html", []byte(codeBlocks[0]), 0666)
		if err != nil {
			log.Printf("写入文件失败: %s\n", err)
		}
	} else {
		log.Println("未找到代码块")
	}
	return h.BaseHandler.Handle(request)
}

type TaskPublisher struct {
	BaseHandler
}

func NewTaskPublisher() *TaskPublisher {
	return &TaskPublisher{BaseHandler: *NewBaseHandler("TaskPublisher")}
}

func (h *TaskPublisher) Handle(request *Request) *Request {
	fmt.Println(h.GetName(), "处理请求:", request.Message)
	request.Data["TaskPublisher"] = "处理成功"
	return h.BaseHandler.Handle(request)
}

type TaskExecutor struct {
	BaseHandler
}

func NewTaskExecutor() *TaskExecutor {
	return &TaskExecutor{BaseHandler: *NewBaseHandler("TaskExecutor")}
}

func (h *TaskExecutor) Handle(request *Request) *Request {
	fmt.Println(h.GetName(), "处理请求:", request.Message)
	request.Data["TaskExecutor"] = "处理成功"
	return h.BaseHandler.Handle(request)
}

type TaskCollector struct {
	BaseHandler
}

func NewTaskCollector() *TaskCollector {
	return &TaskCollector{BaseHandler: *NewBaseHandler("TaskCollector")}
}

func (h *TaskCollector) Handle(request *Request) *Request {
	fmt.Println(h.GetName(), "处理请求:", request.Message)
	request.Data["TaskCollector"] = "处理成功"
	return h.BaseHandler.Handle(request)
}
