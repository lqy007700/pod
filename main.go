package main

import (
	"fmt"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/zxnlx/common"
	"github.com/zxnlx/pod/domain/repository"
	service2 "github.com/zxnlx/pod/domain/service"
	"github.com/zxnlx/pod/handler"
	"github.com/zxnlx/pod/proto/pod"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"strconv"
)

var (
	serviceHost = "host.docker.internal"
	servicePort = "8083"

	// 注册中心配置
	consulHost       = serviceHost
	consulPort int64 = 8500
)

// 注册中心
func initRegistry() registry.Registry {
	return consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			consulHost + ":" + strconv.FormatInt(consulPort, 10),
		}
	})
}

func initConfig() *gorm.DB {
	// 配置中心
	config, err := common.GetConsulConfig(consulHost, consulPort, "/base/micro/config")
	if err != nil {
		common.Fatal(err)
		return nil
	}

	mysqlConf, err := common.GetMysqlFormConsul(config, "mysql")
	if err != nil {
		common.Fatal(err)
		return nil
	}

	// 连接mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlConf.User, mysqlConf.Pwd, mysqlConf.Host, mysqlConf.Port, mysqlConf.Database)
	common.Info(dsn)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		common.Fatal(err)
		return nil
	}
	return db
}

func initK8s() *kubernetes.Clientset {
	//k8s
	//var k8sConfig *string
	//k8sConfig = flag.String("kubeconfig", "", "/Users/lqy007700/Data/config")
	//flag.Parse()
	//common.Info(*k8sConfig)

	//config, err := clientcmd.BuildConfigFromFlags("", "/Users/lqy007700/Data/config")
	config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		common.Fatal(err)
		return nil
	}
	//
	//config, err := rest.InClusterConfig()
	//if err != nil {
	//	return
	//}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		common.Fatal(err)
		return nil
	}
	return clientset
}

func main() {
	c := initRegistry()
	db := initConfig()

	clientSet := initK8s()

	// 日志
	// ./filebeat -e -c filebeat.yml

	service := micro.NewService(
		micro.Server(server.NewServer(func(options *server.Options) {
			options.Advertise = "host.docker.internal:" + servicePort
		})),
		micro.Name("go.micro.service.pod"),
		micro.Version("latest"),
		micro.Registry(c),
		micro.Address(":"+servicePort),
	)

	service.Init()

	dataService := service2.NewPodDataService(repository.NewPodRepository(db), clientSet)
	err := pod.RegisterPodHandler(service.Server(), &handler.PodHandler{PodDataService: dataService})
	if err != nil {
		common.Fatal(err)
		return
	}

	err = service.Run()
	if err != nil {
		common.Fatal(err)
		return
	}
}
