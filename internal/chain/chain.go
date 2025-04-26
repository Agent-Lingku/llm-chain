package chain

// Result 处理结果
type Result struct {
	Data map[string]interface{}
}

// Chain 责任链
type Chain struct {
	head Handler
	tail Handler
}

// NewChain 创建责任链
func NewChain() *Chain {
	return &Chain{}
}

// AddHandler 添加处理类
func (c *Chain) AddHandler(handler Handler) *Chain {
	if c.head == nil {
		c.head = handler
		c.tail = handler
	} else {
		c.tail.SetNext(handler)
		c.tail = handler
	}
	return c
}

// HandleRequest 处理请求
func (c *Chain) HandleRequest(request *Request) *Result {
	if c.head == nil {
		return &Result{Data: request.Data}
	}
	c.head.Handle(request)
	return &Result{Data: request.Data}
}
