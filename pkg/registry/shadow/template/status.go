/*
Copyright 2021 The Clusternet Authors.

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

package template

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/rest"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

// StatusREST implements a StatusREST Storage for Shadow API
type StatusREST struct {
	store *REST
}

func (r *StatusREST) Categories() []string {
	return r.store.Categories()
}

func (r *StatusREST) GroupVersionKind(gv schema.GroupVersion) schema.GroupVersionKind {
	return r.store.GroupVersionKind(gv)
}

// New returns empty Deployment object.
func (r *StatusREST) New() runtime.Object {
	return r.store.New()
}

// Get retrieves the object from the storage. It is required to support Patch.
func (r *StatusREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return r.store.Get(ctx, name, options)
}

// Update alters the status subset of an object.
func (r *StatusREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	// We are explicitly setting forceAllowCreate to false in the call to the underlying storage because
	// subresources should never allow create on update.
	return r.store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}

// GetResetFields implements rest.ResetFieldsStrategy
func (r *StatusREST) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return map[fieldpath.APIVersion]*fieldpath.Set{
		fieldpath.APIVersion(r.store.GroupVersion().String()): fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
		),
	}
}

var _ rest.GroupVersionKindProvider = &StatusREST{}
var _ rest.CategoriesProvider = &StatusREST{}
var _ rest.ResetFieldsStrategy = &StatusREST{}

// NewStatusREST returns a StatusREST Storage object that will work against API services.
func NewStatusREST(store *REST) *StatusREST {
	return &StatusREST{
		store: store,
	}
}
