package helpers

import (
	"github.com/clarketm/json"
	ign3types "github.com/coreos/ignition/v2/config/v3_1/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilrand "k8s.io/apimachinery/pkg/util/rand"

	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
)

var (
	// MasterSelector returns a label selector for masters nodes
	MasterSelector = metav1.AddLabelToSelector(&metav1.LabelSelector{}, "node-role/master", "")
	// WorkerSelector returns a label selector for workers nodes
	WorkerSelector = metav1.AddLabelToSelector(&metav1.LabelSelector{}, "node-role/worker", "")
	// InfraSelector returns a label selector for infra nodes
	InfraSelector = metav1.AddLabelToSelector(&metav1.LabelSelector{}, "node-role/infra", "")
)

// StrToPtr returns a pointer to a string
func StrToPtr(s string) *string {
	return &s
}

// BoolToPtr returns a pointer to a bool
func BoolToPtr(b bool) *bool {
	return &b
}

// NewMachineConfig returns a basic machine config with supplied labels, osurl & files added
func NewMachineConfig(name string, labels map[string]string, osurl string, files []ign3types.File) *mcfgv1.MachineConfig {
	if labels == nil {
		labels = map[string]string{}
	}
	rawIgnition := MarshalOrDie(
		&ign3types.Config{
			Ignition: ign3types.Ignition{
				Version: ign3types.MaxVersion.String(),
			},
			Storage: ign3types.Storage{
				Files: files,
			},
		},
	)

	return &mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: mcfgv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
			UID:    types.UID(utilrand.String(5)),
		},
		Spec: mcfgv1.MachineConfigSpec{
			OSImageURL: osurl,
			Config: runtime.RawExtension{
				Raw: rawIgnition,
			},
		},
	}
}

// NewMachineConfigPool returns a MCP with supplied mcSelector, nodeSelector and machineconfig
func NewMachineConfigPool(name string, mcSelector, nodeSelector *metav1.LabelSelector, currentMachineConfig string) *mcfgv1.MachineConfigPool {
	return &mcfgv1.MachineConfigPool{
		TypeMeta: metav1.TypeMeta{
			APIVersion: mcfgv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"machineconfiguration.openshift.io/mco-built-in": "",
			},
			UID: types.UID(utilrand.String(5)),
		},
		Spec: mcfgv1.MachineConfigPoolSpec{
			NodeSelector:          nodeSelector,
			MachineConfigSelector: mcSelector,
			Configuration: mcfgv1.MachineConfigPoolStatusConfiguration{
				ObjectReference: corev1.ObjectReference{
					Name: currentMachineConfig,
				},
			},
		},
		Status: mcfgv1.MachineConfigPoolStatus{
			Configuration: mcfgv1.MachineConfigPoolStatusConfiguration{
				ObjectReference: corev1.ObjectReference{
					Name: currentMachineConfig,
				},
			},
			Conditions: []mcfgv1.MachineConfigPoolCondition{
				{
					Type:               mcfgv1.MachineConfigPoolRenderDegraded,
					Status:             corev1.ConditionFalse,
					LastTransitionTime: metav1.Unix(0, 0),
					Reason:             "",
					Message:            "",
				},
			},
		},
	}
}

// CreateMachineConfigFromIgnition returns a MachineConfig object from an Ignition config passed to it
func CreateMachineConfigFromIgnition(ignCfg interface{}) *mcfgv1.MachineConfig {
	return &mcfgv1.MachineConfig{
		Spec: mcfgv1.MachineConfigSpec{
			Config: runtime.RawExtension{
				Raw: MarshalOrDie(ignCfg),
			},
		},
	}
}

// MarshalOrDie returns a marshalled interface or panics
func MarshalOrDie(input interface{}) []byte {
	bytes, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	return bytes
}
