package pod

import (
	"flag"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v3"
	opentracing2 "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/opentracing/opentracing-go"
	"github.com/zxnlx/common"
	"github.com/zxnlx/pod/domain/repository"
	service2 "github.com/zxnlx/pod/domain/service"
	hystrix2 "github.com/zxnlx/pod/plugin/hystrix"
	"github.com/zxnlx/pod/proto/pod"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"net"
	"net/http"
	"strconv"
)

var (
	// 注册中心配置
	consulHost       = "127.0.0.1"
	consulPort int64 = 8500

	// 链路
	tracerHost = "127.0.0.1"
	tracerPort = 6381

	// 熔断
	hystrixPort int64 = 9092

	//监控
	prometheusPort = 9192
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
	config, err := common.GetConsulConfig(consulHost, consulPort, "/micro/config")
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
	common.Info(db)
	return db
}

func initTracer() {
	// 链路追踪
	// jaeger
	tracer, i, err := common.NewTracer("base", tracerHost+":"+strconv.Itoa(tracerPort))
	if err != nil {
		common.Fatal(err)
		return
	}
	defer i.Close()
	opentracing.SetGlobalTracer(tracer)
}

func initHystrix() {
	// 熔断
	hystrixHandler := hystrix.NewStreamHandler()
	hystrixHandler.Start()
	go func() {
		//http://192.168.0.112:9092/turbine/turbine.stream
		//看板访问地址 http://127.0.0.1:9002/hystrix，url后面一定要带 /hystrix
		err := http.ListenAndServe(net.JoinHostPort("0.0.0.0", "9092"), hystrixHandler)
		if err != nil {
			common.Fatal(err)
		}
	}()
}

func initK8s() *kubernetes.Clientset {
	//k8s
	var k8sConfig *string
	k8sConfig = flag.String("kubeconfig", "", "/Users/lqy007700/Data/config")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *k8sConfig)
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
	initTracer()
	initHystrix()

	clientSet := initK8s()

	// 日志
	// ./filebeat -e -c filebeat.yml

	// 监控
	common.PrometheusBoot(prometheusPort)

	service := micro.NewService(
		micro.Name("go.micro.service.pod"),
		micro.Version("latest"),
		micro.Registry(c),
		// 链路
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// 熔断
		micro.WrapClient(hystrix2.NewClientHystrixWrapper()),
		// 限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(1000)),
	)

	service.Init()

	// 创建表 执行一次
	err := repository.NewPodRepository(db).InitTable()
	if err != nil {
		common.Fatal(err)
		return
	}

	dataService := service2.NewPodDataService(repository.NewPodRepository(db), clientSet)
	err = pod.RegisterPodHandler(service.Server(), &handler.PodHandler())
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
