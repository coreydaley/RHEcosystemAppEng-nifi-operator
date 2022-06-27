package nifi

import (
	"context"

	bigdatav1alpha1 "github.com/RHEcosystemAppEng/nifi-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func newService(nifi *bigdatav1alpha1.Nifi) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nifi.Name,
			Namespace: nifi.Namespace,
			Labels:    labelsForNifi(nifi.Name),
		},
	}

}

func (r *Reconciler) reconcileServices(ctx context.Context, req ctrl.Request, nifi *bigdatav1alpha1.Nifi) error {
	svc := newService(nifi)
	svc.Spec = corev1.ServiceSpec{
		Selector: labelsForNifi(nifi.Name),
		Ports: []corev1.ServicePort{
			{
				Name:     nifiConsolePortName,
				Port:     nifiConsolePort,
				Protocol: "TCP",
			},
		},
	}

	// Checking if service already exists
	existingSrv := newService(nifi)
	err := r.Get(ctx, req.NamespacedName, existingSrv)
	if !errors.IsNotFound(err) {
		// if it exists, do nothing
		return nil
	}

	// Set Nifi instance as the owner and controller
	if err := ctrl.SetControllerReference(nifi, svc, r.Scheme); err != nil {
		return err
	}

	return r.Client.Create(ctx, svc)
}
