package model

type Status string

const (
	StatusPending   Status = "待执行" // 任务待执行
	StatusRunning   Status = "执行中" // 任务正在执行
	StatusCompleted Status = "已完成" // 任务已成功完成
	StatusFailed    Status = "失败"  // 任务执行失败
	StatusCancelled Status = "已取消" // 任务被取消
)

// Task 表示一个任务
type Task struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status Status `json:"status"` // 任务状态
	// 预留字段
	Flag1 string `json:"flag1"`
	Flag2 string `json:"flag2"`
	Flag3 string `json:"flag3"`
}
