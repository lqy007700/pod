package model

type PodPort struct {
	ID            int64  `json:"id"`
	PodId         int64  `json:"pod_id"`
	ContainerPort int32  `json:"container_port"`
	Protocol      string `json:"protocol"`
}
