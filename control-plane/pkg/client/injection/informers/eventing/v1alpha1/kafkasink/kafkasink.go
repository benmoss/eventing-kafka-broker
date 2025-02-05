/*
 * Copyright 2021 The Knative Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Code generated by injection-gen. DO NOT EDIT.

package kafkasink

import (
	context "context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	cache "k8s.io/client-go/tools/cache"
	apiseventingv1alpha1 "knative.dev/eventing-kafka-broker/control-plane/pkg/apis/eventing/v1alpha1"
	versioned "knative.dev/eventing-kafka-broker/control-plane/pkg/client/clientset/versioned"
	v1alpha1 "knative.dev/eventing-kafka-broker/control-plane/pkg/client/informers/externalversions/eventing/v1alpha1"
	client "knative.dev/eventing-kafka-broker/control-plane/pkg/client/injection/client"
	factory "knative.dev/eventing-kafka-broker/control-plane/pkg/client/injection/informers/factory"
	eventingv1alpha1 "knative.dev/eventing-kafka-broker/control-plane/pkg/client/listers/eventing/v1alpha1"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
	logging "knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterInformer(withInformer)
	injection.Dynamic.RegisterDynamicInformer(withDynamicInformer)
}

// Key is used for associating the Informer inside the context.Context.
type Key struct{}

func withInformer(ctx context.Context) (context.Context, controller.Informer) {
	f := factory.Get(ctx)
	inf := f.Eventing().V1alpha1().KafkaSinks()
	return context.WithValue(ctx, Key{}, inf), inf.Informer()
}

func withDynamicInformer(ctx context.Context) context.Context {
	inf := &wrapper{client: client.Get(ctx)}
	return context.WithValue(ctx, Key{}, inf)
}

// Get extracts the typed informer from the context.
func Get(ctx context.Context) v1alpha1.KafkaSinkInformer {
	untyped := ctx.Value(Key{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch knative.dev/eventing-kafka-broker/control-plane/pkg/client/informers/externalversions/eventing/v1alpha1.KafkaSinkInformer from context.")
	}
	return untyped.(v1alpha1.KafkaSinkInformer)
}

type wrapper struct {
	client versioned.Interface

	namespace string
}

var _ v1alpha1.KafkaSinkInformer = (*wrapper)(nil)
var _ eventingv1alpha1.KafkaSinkLister = (*wrapper)(nil)

func (w *wrapper) Informer() cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(nil, &apiseventingv1alpha1.KafkaSink{}, 0, nil)
}

func (w *wrapper) Lister() eventingv1alpha1.KafkaSinkLister {
	return w
}

func (w *wrapper) KafkaSinks(namespace string) eventingv1alpha1.KafkaSinkNamespaceLister {
	return &wrapper{client: w.client, namespace: namespace}
}

func (w *wrapper) List(selector labels.Selector) (ret []*apiseventingv1alpha1.KafkaSink, err error) {
	lo, err := w.client.EventingV1alpha1().KafkaSinks(w.namespace).List(context.TODO(), v1.ListOptions{
		LabelSelector: selector.String(),
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
	if err != nil {
		return nil, err
	}
	for idx := range lo.Items {
		ret = append(ret, &lo.Items[idx])
	}
	return ret, nil
}

func (w *wrapper) Get(name string) (*apiseventingv1alpha1.KafkaSink, error) {
	return w.client.EventingV1alpha1().KafkaSinks(w.namespace).Get(context.TODO(), name, v1.GetOptions{
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
}
