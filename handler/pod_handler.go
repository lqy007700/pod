package handler

import (
	"context"
	"errors"
	"github.com/zxnlx/common"
	"github.com/zxnlx/pod/domain/model"
	"github.com/zxnlx/pod/domain/service"
	"github.com/zxnlx/pod/proto/pod"
)

type PodHandler struct {
	PodDataService service.IPodDataService
}

func (p *PodHandler) AddPod(ctx context.Context, info *pod.PodInfo, resp *pod.Response) error {
	common.Info("Add Pod")
	podModel := &model.Pod{}
	err := common.SwapTo(info, podModel)
	if err != nil {
		common.Error(err)
		return err
	}

	err = p.PodDataService.CreateToK8s(info)
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
	resp.Msg = "Add success"
	return nil
}

func (p *PodHandler) DelPod(ctx context.Context, req *pod.PodId, resp *pod.Response) error {
	common.Info("Del Pod")
	info, err := p.PodDataService.FindPodById(req.Id)
	if err != nil {
		common.Error(err)
		resp.Msg = err.Error()
		return err
	}

	if info == nil {
		resp.Msg = "Pod Not Exist"
		common.Error(resp.Msg)
		return errors.New(resp.Msg)
	}

	err = p.PodDataService.DelForK8s(info)
	if err != nil {
		resp.Msg = "Pod Not Exist"
		common.Error(resp.Msg)
		return err
	}
	resp.Msg = "Del success"
	return nil
}

func (p *PodHandler) FindPodById(ctx context.Context, req *pod.PodId, info *pod.PodInfo) error {
	podInfo, err := p.PodDataService.FindPodById(req.Id)
	if err != nil {
		return err
	}
	err = common.SwapTo(podInfo, info)
	if err != nil {
		return err
	}
	return nil
}

func (p *PodHandler) UpdatePod(ctx context.Context, info *pod.PodInfo, res *pod.Response) error {
	err := p.PodDataService.UpdateForK8s(info)
	if err != nil {
		return err
	}

	podInfo, err := p.PodDataService.FindPodById(info.Id)
	if err != nil {
		return err
	}

	err = common.SwapTo(info, podInfo)
	if err != nil {
		return err
	}

	err = p.PodDataService.UpdatePod(podInfo)
	if err != nil {
		return err
	}
	return nil
}

func (p *PodHandler) FindAllPod(ctx context.Context, all *pod.FindAll, list *pod.PodList) error {
	allPod, err := p.PodDataService.FindAllPod()
	if err != nil {
		return err
	}

	for _, v := range allPod {
		podInfo := &pod.PodInfo{}
		err = common.SwapTo(v, podInfo)
		if err != nil {
			return err
		}

		list.PodList = append(list.PodList, podInfo)
	}
	return nil
}
