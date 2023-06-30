package service

import (
	"context"
	"github.com/zxnlx/pod/domain/model"
	"github.com/zxnlx/pod/domain/repository"
	"github.com/zxnlx/pod/proto/pod"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"strconv"
)

type IPodDataService interface {
	AddPod(*model.Pod) (int64, error)
	DelPod(int64) error
	UpdatePod(*model.Pod) error
	FindPodById(int64) (*model.Pod, error)
	FindAllPod() ([]model.Pod, error)

	CreateToK8s(*pod.PodInfo) error
	DelForK8s(*model.Pod) error
	UpdateForK8s(*pod.PodInfo) error
}

func NewPodDataService(podRepo repository.IPodRepository,
	clientSet *kubernetes.Clientset) IPodDataService {
	return &PodDataService{
		PodRepository: podRepo,
		K8sClientSet:  clientSet,
		deployment:    &v1.Deployment{},
	}
}

type PodDataService struct {
	PodRepository repository.IPodRepository
	K8sClientSet  *kubernetes.Clientset
	deployment    *v1.Deployment
}

func (p *PodDataService) AddPod(pod *model.Pod) (int64, error) {
	return p.PodRepository.CreatePod(pod)
}

func (p *PodDataService) DelPod(id int64) error {
	return p.PodRepository.DeletePodById(id)
}

func (p *PodDataService) UpdatePod(pod *model.Pod) error {
	return p.PodRepository.Update(pod)
}

func (p *PodDataService) FindPodById(id int64) (*model.Pod, error) {
	return p.PodRepository.FindPodById(id)
}

func (p *PodDataService) FindAllPod() ([]model.Pod, error) {
	return p.PodRepository.FindAll()
}

func (p *PodDataService) CreateToK8s(info *pod.PodInfo) error {
	p.setDeployment(info)
	_, err := p.K8sClientSet.AppsV1().Deployments(info.PodNamespace).Get(context.Background(), info.PodName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = p.K8sClientSet.AppsV1().Deployments(info.PodNamespace).Create(context.Background(), p.deployment, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (p *PodDataService) DelForK8s(m *model.Pod) error {
	//TODO implement me
	panic("implement me")
}

func (p *PodDataService) UpdateForK8s(info *pod.PodInfo) error {
	//TODO implement me
	panic("implement me")
}

func (p *PodDataService) setDeployment(info *pod.PodInfo) {
	deploy := &v1.Deployment{}
	deploy.TypeMeta = metav1.TypeMeta{
		Kind:       "deployment",
		APIVersion: "v1",
	}
	deploy.ObjectMeta = metav1.ObjectMeta{
		Name:      info.PodName,
		Namespace: info.PodNamespace,
		Labels: map[string]string{
			"app-name": info.PodName,
			"author":   "zxnl",
		},
	}

	deploy.Name = info.PodName
	deploy.Spec = v1.DeploymentSpec{
		Replicas: &info.PodReplicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app-name": info.PodName,
			},
		},
		Template: v12.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app-name": info.PodName,
				},
			},
			Spec: v12.PodSpec{
				Containers: []v12.Container{
					{
						Name:            info.PodName,
						Image:           info.PodImages,
						Ports:           p.getContainerPort(info),
						Env:             p.getEnv(info),
						Resources:       p.getResource(info),
						ImagePullPolicy: p.getImagePullPolicy(info),
					},
				},
			},
		},
	}
	p.deployment = deploy
}

func (p *PodDataService) getResource(info *pod.PodInfo) v12.ResourceRequirements {
	return v12.ResourceRequirements{
		Limits: v12.ResourceList{
			"cpu":    resource.MustParse(strconv.FormatFloat(float64(info.PodCpuMax), 'f', 6, 64)),
			"memory": resource.MustParse(strconv.FormatFloat(float64(info.PodMemMax), 'f', 6, 64)),
		},
		// todo 最小资源
		Requests: v12.ResourceList{
			"cpu":    resource.MustParse(strconv.FormatFloat(float64(info.PodCpuMax), 'f', 6, 64)),
			"memory": resource.MustParse(strconv.FormatFloat(float64(info.PodMemMax), 'f', 6, 64)),
		},
	}
}

func (p *PodDataService) getContainerPort(info *pod.PodInfo) []v12.ContainerPort {
	res := make([]v12.ContainerPort, 0, len(info.PodPort))
	for _, v := range info.PodPort {
		res = append(res, v12.ContainerPort{
			Name:          "port-" + strconv.FormatInt(int64(v.ContainerPort), 10),
			ContainerPort: v.ContainerPort,
			Protocol:      p.getProtocol(v.Protocol),
			HostIP:        "",
		})
	}
	return res
}

func (p *PodDataService) getEnv(info *pod.PodInfo) []v12.EnvVar {
	res := make([]v12.EnvVar, 0, len(info.PodEnv))
	for _, v := range info.PodEnv {
		res = append(res, v12.EnvVar{
			Name:      v.EnvKey,
			Value:     v.EnvVal,
			ValueFrom: nil,
		})
	}
	return res
}

func (p *PodDataService) getProtocol(protocol string) v12.Protocol {
	switch protocol {
	case "TCP":
		return "TCP"
	case "UDP":
		return "UDP"
	default:
		return "TCP"
	}
}

func (p *PodDataService) getImagePullPolicy(pod *pod.PodInfo) v12.PullPolicy {
	switch pod.PodPullPolicy {
	case "Always":
		return "Always"
	case "Never":
		return "Never"
	case "IfNotPresent":
		return "IfNotPresent"
	default:
		return "Always"
	}
}
