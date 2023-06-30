package model

/**
- Pending : API Server 已经创建了该Pod，但是Pod中一个或者多个容器镜像还没创建
- Running： Pod内所有的容器已经创建，至少有一个容器处于运行、正在启动、正在重启
- Succeeded：Pod内所有的容器均成功退出，并且不会再重启。
- Failed：Pod内所有的容器均已退出，但至少有一个容器退出失败
- Unknown：由于某种原因无法获取Pod状态，网络不通等
*/

type Pod struct {
	ID           int64  `json:"id"`
	PodName      string `gorm:"unique_index,not_null" json:"pod_name"`
	PodNamespace string `json:"pod_namespace"`
	PodTeamID    int64  `json:"pod_team_id"`
	PodReplicas  int32  `json:"pod_replicas"`

	PodCpuMin float32 `json:"pod_cpu_min"`
	PodCpuMax float32 `json:"pod_cpu_max"`
	PodMemMin float32 `json:"pod_mem_min"`
	PodMemMax float32 `json:"pod_mem_max"`

	PodPort []PodPort `json:"pod_port"`
	PodEnv  []PodEnv  `json:"pod_env"`

	// 镜像拉去策略
	// Always 总是拉取
	// IfNotPresent 本地没有就拉
	// Never 只使用本地镜像，从不拉
	PodPullPolicy string `json:"pod_pull_policy"`

	// 重启策略
	// Always 容器失效时自动重启
	// OnFailure 容器终止且退出码不为0时重启
	// Never 不重启
	PodRestart string `json:"pod_restart"`

	// 发布策略
	// recreate  重建，停止旧版本发布新版本
	// rolling-update 滚动更新，一个一个的发布新版本
	// blue/green 蓝绿发布 新旧版本共存，然后切换流量
	PodType string `json:"pod_type"`

	// 镜像名称+tag
	PodImage string `json:"pod_image"`
}
