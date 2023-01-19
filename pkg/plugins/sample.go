package plugins

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	frameworkruntime "k8s.io/kubernetes/pkg/scheduler/framework/runtime"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

type SampleArgs struct {
	FavoriteColor  string `json:"favorite_color,omitempty"`
	FavoriteNumber int    `json:"favorite_number,omitempty"`
	ThanksTo       string `json:"thanks_to,omitempty"`
}

type Sample struct {
	args    *SampleArgs
	handler framework.FrameworkHandle
}

var _ framework.PreFilterPlugin = &Sample{}
var _ framework.FilterPlugin = &Sample{}

const SamplePlugin = "sample-plugin"

func (s Sample) Name() string {
	//TODO implement me
	panic("implement me")
}

func (s Sample) PreFilter(ctx context.Context, state *framework.CycleState, p *corev1.Pod) *framework.Status {
	//TODO implement me
	panic("implement me")
}

func (s Sample) PreFilterExtensions() framework.PreFilterExtensions {
	//TODO implement me
	panic("implement me")
}

func (s Sample) Filter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	//TODO implement me
	panic("implement me")
}

//type Option func(runtime.Registry) error
//type Registry map[string]PluginFactory
//type PluginFactory = func(configuration runtime.Object, f v1alpha1.FrameworkHandle) (v1alpha1.Plugin, error)

func SampleFactory(configuration runtime.Object, f framework.FrameworkHandle) (framework.Plugin, error) {
	args := getSampleArgs(configuration)
	return &Sample{
		args:    args,
		handler: f,
	}, nil
}

func getSampleArgs(obj runtime.Object) *SampleArgs {
	sa := &SampleArgs{}
	err := frameworkruntime.DecodeInto(obj, sa)
	if err != nil {
		return nil
	}
	return sa
}
