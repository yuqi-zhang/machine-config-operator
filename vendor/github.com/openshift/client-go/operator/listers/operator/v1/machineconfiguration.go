// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/openshift/api/operator/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// MachineConfigurationLister helps list MachineConfigurations.
// All objects returned here must be treated as read-only.
type MachineConfigurationLister interface {
	// List lists all MachineConfigurations in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.MachineConfiguration, err error)
	// Get retrieves the MachineConfiguration from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.MachineConfiguration, error)
	MachineConfigurationListerExpansion
}

// machineConfigurationLister implements the MachineConfigurationLister interface.
type machineConfigurationLister struct {
	indexer cache.Indexer
}

// NewMachineConfigurationLister returns a new MachineConfigurationLister.
func NewMachineConfigurationLister(indexer cache.Indexer) MachineConfigurationLister {
	return &machineConfigurationLister{indexer: indexer}
}

// List lists all MachineConfigurations in the indexer.
func (s *machineConfigurationLister) List(selector labels.Selector) (ret []*v1.MachineConfiguration, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.MachineConfiguration))
	})
	return ret, err
}

// Get retrieves the MachineConfiguration from the index for a given name.
func (s *machineConfigurationLister) Get(name string) (*v1.MachineConfiguration, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("machineconfiguration"), name)
	}
	return obj.(*v1.MachineConfiguration), nil
}