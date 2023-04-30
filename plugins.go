package plugins

import (
	"fmt"
	"time"
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"github.com/go-ping/ping"
)
type customFilterPlugin struct{
	handle framework.Handle
}

const (
	// Name : name of plugin used in the plugin registry and configurations.
	Name = "CustomFilterPlugin"
)


var _  = framework.FilterPlugin(&customFilterPlugin{})

func (p *customFilterPlugin) Name() string {
	return Name
}


func (s *customFilterPlugin) PreFilter(ctx context.Context, pod *v1.Pod) *framework.Status {
	return framework.NewStatus(framework.Success, "")
}


func (p *customFilterPlugin) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	// Implementirajte logiku filtriranja čvorova ovdje
	// Učitajte sve čvorove
	nodes, err := p.handle.SnapshotSharedLister().NodeInfos().List()
	if err != nil {
		return framework.NewStatus(framework.Error, "Failed to list nodes")
	}

	// Inicijalizirajte najbliži čvor i minimalnu latenciju
	var closestNode string
	var minLatency time.Duration = time.Duration(1<<63 - 1) // Postavljanje maksimalne vrijednosti za usporedbu

	// Pronađite najbliži čvor koristeći ping
	for _, node := range nodes {
		ip := node.Node().Status.Addresses[0].Address
		latency, err := pingNode(ip)
		if err != nil {
			return framework.NewStatus(framework.Error, "Failed to ping node")
		}

		if latency < minLatency {
			fmt.Println("node:", node.Node().Name + " latency: ", latency );
			minLatency = latency
			closestNode = node.Node().Name
		}
	}

	// Ako je trenutni čvor najbliži, postavite status na Success
	if nodeInfo.Node().Name == closestNode {
		return framework.NewStatus(framework.Success)
	}

	return framework.NewStatus(framework.Unschedulable, "Node cannot be scheduled")
}

func (s *customFilterPlugin) PreBind(ctx context.Context, pod *v1.Pod, nodeName string) *framework.Status {

	return framework.NewStatus(framework.Success, "")
}

func New(obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	return &customFilterPlugin{}, nil
}


func pingNode(ip string) (time.Duration, error) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return 0, err
	}
	pinger.Count = 3
	pinger.Timeout = time.Second * 5
	pinger.SetPrivileged(true)
	pinger.Run()
	stats := pinger.Statistics()
	return stats.AvgRtt, nil
}


