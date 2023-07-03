package handler

import (
	"context"
	"github.com/zxnlx/common"
	"github.com/zxnlx/pod/domain/model"
	"github.com/zxnlx/pod/domain/service"
	"github.com/zxnlx/pod/proto/pod"
)

type PodHandler struct {
	PodDataService service.IPodDataService
}

func (p *PodHandler) AddPod(ctx context.Context, info *pod.PodInfo, resp *pod.Response) error {
	common.Info("添加Pod")
	podModel := &model.Pod{}
	err := common.SwapTo(info, podModel)
	if err != nil {
		common.Fatal(err)
		return err
	}

	err := p.PodDataService.CreateToK8s(info)
	if err != nil {
		common.Error(err)
		resp.Msg = err.Error()
		return err
	}

	addPod, err := p.PodDataService.AddPod(podModel)
	if err != nil {
		common.Error(err)
		resp.Msg = err.Error()
		return err
	}

	common.Info(addPod)
	resp.Msg = "添加成功"
	return nil
}

func (p *PodHandler) DelPod(ctx context.Context, id *pod.PodId, response *pod.Response) error {
	//TODO implement me
	panic("implement me")
}

func (p *PodHandler) FindPodById(ctx context.Context, id *pod.PodId, info *pod.PodInfo) error {
	//TODO implement me
	panic("implement me")
}

func (p *PodHandler) UpdatePod(ctx context.Context, info *pod.PodInfo, response *pod.Response) error {
	//TODO implement me
	panic("implement me")
}

func (p *PodHandler) FindAllPod(ctx context.Context, all *pod.FindAll, list *pod.PodList) error {
	//TODO implement me
	panic("implement me")
}
