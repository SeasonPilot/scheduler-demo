package plugins

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
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

type preFilterState struct {
	framework.Resource
}

func (p preFilterState) Clone() framework.StateData {
	return p
}

func computePodResourceLimit(pod *corev1.Pod) *preFilterState {
	//var res *preFilterState  //fixme: 指针未初始化
	var res = &preFilterState{}
	for _, c := range pod.Spec.Containers {
		res.Add(c.Resources.Limits)
	}
	return res
}

var _ framework.PreFilterPlugin = &Sample{}
var _ framework.FilterPlugin = &Sample{}

const (
	Name              = "sample-plugin"
	preFilterStateKey = "PreFilter" + Name
)

func (s Sample) Name() string {
	return Name
}

func (s Sample) PreFilter(ctx context.Context, state *framework.CycleState, p *corev1.Pod) *framework.Status {
	if klog.V(2).Enabled() {
		klog.InfoS("Start PreFilter Pod", "pod", p.Name)
	}
	state.Write(preFilterStateKey, computePodResourceLimit(p))
	return nil
}

func (s Sample) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

func getPreFilterState(cycleState *framework.CycleState) (*preFilterState, error) {
	read, err := cycleState.Read(preFilterStateKey)
	if err != nil {
		return nil, err
	}
	filterState, ok := read.(*preFilterState)
	if !ok {
		return nil, fmt.Errorf("%+v convert to SamplePlugin preFilterState error", filterState)
	}
	return filterState, nil
}

func (s Sample) Filter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	filterState, err := getPreFilterState(state)
	if err != nil {
		return framework.NewStatus(framework.Error, err.Error())
	}

	if klog.V(2).Enabled() {
		klog.InfoS("Start Filter Pod", "pod", pod.Name, "node", nodeInfo.Node().Name, "preFilterState", filterState)
	}
	// logic
	return framework.NewStatus(framework.Success, "")
}

//type Option func(runtime.Registry) error
//type Registry map[string]PluginFactory
//type PluginFactory = func(configuration runtime.Object, f v1alpha1.FrameworkHandle) (v1alpha1.Plugin, error)

func SampleFactory(configuration runtime.Object, f framework.FrameworkHandle) (framework.Plugin, error) {
	args, err := getSampleArgs(configuration)
	if err != nil {
		return nil, err
	}

	// validate args
	if klog.V(2).Enabled() {
		klog.InfoS("Successfully get plugin config args", "plugin", Name, "args", args)
	}

	return &Sample{
		args:    args,
		handler: f,
	}, nil
}

func getSampleArgs(obj runtime.Object) (*SampleArgs, error) {
	sa := &SampleArgs{}
	err := frameworkruntime.DecodeInto(obj, sa)
	if err != nil {
		return nil, err
	}
	return sa, nil
}
