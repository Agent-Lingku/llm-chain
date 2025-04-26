package model

type Agent struct {
	ID     int    `json:"id"`
	Role   string `json:"role"`   // 角色：搜集家、任务发布家、执行者、作家
	Params string `json:"params"` // 参数
}
