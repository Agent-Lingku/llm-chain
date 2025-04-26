package interfaces

type IAgent interface {
	IsFinished() bool               // 是否完成
	HasContext() bool               // 是否有上下文
	EchoRoleInfo() (string, string) // 输出角色信息
	ExecuteTask() (string, error)   // 执行结果
}
