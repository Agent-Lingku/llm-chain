package agent

import "learn/internal/agent/prompts"

// Role 定义 Agent 角色
type Role string

const (
	DemandAnalysisRole Role = "需求分析"
	FrontEndRole       Role = "前端工程师"
	TaskExecutionRole  Role = "任务执行"
	AssistanceRole     Role = "协助"
	MonitoringRole     Role = "监控"
	ResultFeedbackRole Role = "结果反馈"
)

// RolePromptMap 存储 Agent 角色和对应的 Prompt 字符串
var RolePromptMap = map[Role]string{
	DemandAnalysisRole: prompts.DemandAnalysisPrompt,
	FrontEndRole:       prompts.FrontEndPrompt,
	TaskExecutionRole:  prompts.TaskExecutionPrompt,
	AssistanceRole:     prompts.AssistancePrompt,
	MonitoringRole:     prompts.MonitoringPrompt,
	ResultFeedbackRole: prompts.ResultFeedbackPrompt,
}

// GetAgentPrompt 获取 Agent 的 Prompt
func GetAgentPrompt(role Role) string {
	return RolePromptMap[role]
}
