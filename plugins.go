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
	 // Provjerite je li nodeInfo nil
    	if nodeInfo == nil || nodeInfo.Node() == nil {
        	return framework.NewStatus(framework.Error, "NodeInfo or Node is nil")
    	}
	// Implementirajte logiku filtriranja čvorova ovdje
	// Učitajte sve čvorove
	nodes, err := p.handle.SnapshotSharedLister().NodeInfos().List()
	if err != nil {
		return framework.NewStatus(framework.Error, "Failed to list nodes")
	}

	fmt.Println("Ovdje1" );
	// Inicijalizirajte najbliži čvor i minimalnu latenciju
	var closestNode string
	var minLatency time.Duration = time.Duration(1<<63 - 1) // Postavljanje maksimalne vrijednosti za usporedbu

	// Pronađite najbliži čvor koristeći ping
	for _, node := range nodes {
		fmt.Println("Ovdje2" );
		if node.Node() == nil || len(node.Node().Status.Addresses) == 0 {
			fmt.Println("Ovdje3" );
            		continue
        	}
		fmt.Println("Ovdje4" );
		ip := node.Node().Status.Addresses[0].Address
		fmt.Println("Ovdje5" );
		latency, err := pingNode(ip)
		fmt.Println("Ovdje6" );
		if err != nil {
			fmt.Println("Ovdje7" );
			return framework.NewStatus(framework.Error, "Failed to ping node")
		}

		fmt.Println("Ovdje8" );
		if latency < minLatency {
			fmt.Println("Ovdje9" );
			fmt.Println("node:", node.Node().Name + " latency: ", latency );
			minLatency = latency
			closestNode = node.Node().Name
		}
	}
	fmt.Println("Ovdje10" );
	// Ako je trenutni čvor najbliži, postavite status na Success
	if nodeInfo.Node().Name == closestNode {
		fmt.Println("Ovdje11" );
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


