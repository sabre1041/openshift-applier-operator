package manager

import (
	"context"
	"log"

	copapi "github.com/redhat-cop/openshift-applier-operator/pkg/apis/cop/v1alpha1"
	"github.com/redhat-cop/openshift-applier-operator/pkg/util"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	k8sManager "sigs.k8s.io/controller-runtime/pkg/manager"
)

type ApplierManager struct {
	client client.Client
	scheme *runtime.Scheme
}

func New(manager k8sManager.Manager) (*ApplierManager, error) {
	return &ApplierManager{
		client: manager.GetClient(),
		scheme: manager.GetScheme(),
	}, nil
}

func (c *ApplierManager) LaunchApplierJob(applier *copapi.Applier) error {

	log.Printf("Launching Job from resource in namespace '%s' with name '%s'", applier.Namespace, applier.Name)

	job, err := util.GenerateJob(applier)

	if err != nil {
		return err
	}

	// Set Owner Reference
	controllerutil.SetControllerReference(applier, job, c.scheme)

	return c.client.Create(context.TODO(), job)

}

func (c *ApplierManager) FindApplierResourceByToken(namespace string, name string, token string) (*copapi.Applier, error) {

	applier := &copapi.Applier{}

	err := c.client.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, applier)

	if err != nil {
		return nil, err
	}

	return applier, nil

}

func (c *ApplierManager) FindApplierResource(namespace string, name string) (*copapi.Applier, error) {

	instance := &copapi.Applier{}

	err := c.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, instance)

	if err != nil {
		return nil, err
	}

	return instance, nil

}
