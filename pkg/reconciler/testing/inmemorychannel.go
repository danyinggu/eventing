/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testing

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/types"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	duckv1alpha1 "knative.dev/eventing/pkg/apis/duck/v1alpha1"
	"knative.dev/eventing/pkg/apis/eventing"
	"knative.dev/eventing/pkg/apis/messaging/v1alpha1"
	"knative.dev/pkg/apis"
)

// InMemoryChannelOption enables further configuration of a InMemoryChannel.
type InMemoryChannelOption func(*v1alpha1.InMemoryChannel)

// NewInMemoryChannel creates an InMemoryChannel with InMemoryChannelOptions.
func NewInMemoryChannel(name, namespace string, imcopt ...InMemoryChannelOption) *v1alpha1.InMemoryChannel {
	imc := &v1alpha1.InMemoryChannel{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.InMemoryChannelSpec{},
	}
	for _, opt := range imcopt {
		opt(imc)
	}
	imc.SetDefaults(context.Background())
	return imc
}

func WithInitInMemoryChannelConditions(imc *v1alpha1.InMemoryChannel) {
	imc.Status.InitializeConditions()
}

func WithInMemoryChannelGeneration(gen int64) InMemoryChannelOption {
	return func(s *v1alpha1.InMemoryChannel) {
		s.Generation = gen
	}
}

func WithInMemoryChannelStatusObservedGeneration(gen int64) InMemoryChannelOption {
	return func(s *v1alpha1.InMemoryChannel) {
		s.Status.ObservedGeneration = gen
	}
}

func WithInMemoryChannelDeleted(imc *v1alpha1.InMemoryChannel) {
	deleteTime := metav1.NewTime(time.Unix(1e9, 0))
	imc.ObjectMeta.SetDeletionTimestamp(&deleteTime)
}

func WithInMemoryChannelSubscribers(subscribers []duckv1alpha1.SubscriberSpec) InMemoryChannelOption {
	return func(imc *v1alpha1.InMemoryChannel) {
		imc.Spec.Subscribable = &duckv1alpha1.Subscribable{Subscribers: subscribers}
	}
}

func WithInMemoryChannelDeploymentFailed(reason, message string) InMemoryChannelOption {
	return func(imc *v1alpha1.InMemoryChannel) {
		imc.Status.MarkDispatcherFailed(reason, message)
	}
}

func WithInMemoryChannelDeploymentUnknown(reason, message string) InMemoryChannelOption {
	return func(imc *v1alpha1.InMemoryChannel) {
		imc.Status.MarkDispatcherUnknown(reason, message)
	}
}

func WithInMemoryChannelDeploymentReady() InMemoryChannelOption {
	return func(imc *v1alpha1.InMemoryChannel) {
		imc.Status.PropagateDispatcherStatus(&appsv1.DeploymentStatus{Conditions: []appsv1.DeploymentCondition{{Type: appsv1.DeploymentAvailable, Status: corev1.ConditionTrue}}})
	}
}

func WithInMemoryChannelServicetNotReady(reason, message string) InMemoryChannelOption {
	return func(imc *v1alpha1.InMemoryChannel) {
		imc.Status.MarkServiceFailed(reason, message)
	}
}

func WithInMemoryChannelServiceReady() InMemoryChannelOption {
	return func(imc *v1alpha1.InMemoryChannel) {
		imc.Status.MarkServiceTrue()
	}
}

func WithInMemoryChannelChannelServiceNotReady(reason, message string) InMemoryChannelOption {
	return func(imc *v1alpha1.InMemoryChannel) {
		imc.Status.MarkChannelServiceFailed(reason, message)
	}
}

func WithInMemoryChannelChannelServiceReady() InMemoryChannelOption {
	return func(imc *v1alpha1.InMemoryChannel) {
		imc.Status.MarkChannelServiceTrue()
	}
}

func WithInMemoryChannelEndpointsNotReady(reason, message string) InMemoryChannelOption {
	return func(imc *v1alpha1.InMemoryChannel) {
		imc.Status.MarkEndpointsFailed(reason, message)
	}
}

func WithInMemoryChannelEndpointsReady() InMemoryChannelOption {
	return func(imc *v1alpha1.InMemoryChannel) {
		imc.Status.MarkEndpointsTrue()
	}
}

func WithInMemoryChannelAddress(a string) InMemoryChannelOption {
	return func(imc *v1alpha1.InMemoryChannel) {
		imc.Status.SetAddress(&apis.URL{
			Scheme: "http",
			Host:   a,
		})
	}
}

func WithInMemoryChannelReady(host string) InMemoryChannelOption {
	return func(imc *v1alpha1.InMemoryChannel) {
		imc.Status.SetAddress(&apis.URL{
			Scheme: "http",
			Host:   host,
		})
		imc.Status.MarkChannelServiceTrue()
		imc.Status.MarkEndpointsTrue()
		imc.Status.MarkServiceTrue()
	}
}

func WithInMemoryChannelReadySubscriber(uid string) InMemoryChannelOption {
	return WithInMemoryChannelReadySubscriberAndGeneration(uid, 0)
}

func WithInMemoryChannelReadySubscriberAndGeneration(uid string, observedGeneration int64) InMemoryChannelOption {
	return func(c *v1alpha1.InMemoryChannel) {
		if c.Status.GetSubscribableTypeStatus() == nil { // Both the SubscribableStatus fields are nil
			c.Status.SetSubscribableTypeStatus(duckv1alpha1.SubscribableStatus{})
		}
		c.Status.SubscribableTypeStatus.AddSubscriberToSubscribableStatus(duckv1alpha1.SubscriberStatus{
			UID:                types.UID(uid),
			ObservedGeneration: observedGeneration,
			Ready:              corev1.ConditionTrue,
		})
	}
}

func WithInMemoryChannelStatusSubscribers(subscriberStatuses []duckv1alpha1.SubscriberStatus) InMemoryChannelOption {
	return func(imc *v1alpha1.InMemoryChannel) {
		imc.Status.SetSubscribableTypeStatus(duckv1alpha1.SubscribableStatus{
			Subscribers: subscriberStatuses})
	}
}

func WithInMemoryScopeAnnotation(value string) InMemoryChannelOption {
	return func(imc *v1alpha1.InMemoryChannel) {
		if imc.Annotations == nil {
			imc.Annotations = make(map[string]string)
		}
		imc.Annotations[eventing.ScopeAnnotationKey] = value
	}
}
