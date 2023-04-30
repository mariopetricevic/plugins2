package plugins

import (
	"fmt"
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
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
	if nodeInfo.Node().Name == "masternode"{
		fmt.Println("Inside filter method");
		return framework.NewStatus(framework.Success)
	}
	return framework.NewStatus(framework.Unschedulable, "Node is not masternode")
}

func (s *customFilterPlugin) PreBind(ctx context.Context, pod *v1.Pod, nodeName string) *framework.Status {

	return framework.NewStatus(framework.Success, "")
}

func New(obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	return &customFilterPlugin{}, nil
}

