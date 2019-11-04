/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "sigs.k8s.io/kube-storage-version-migrator/pkg/apis/migration/v1alpha1"
)

// FakeStorageVersionMigrations implements StorageVersionMigrationInterface
type FakeStorageVersionMigrations struct {
	Fake *FakeMigrationV1alpha1
}

var storageversionmigrationsResource = schema.GroupVersionResource{Group: "migration.k8s.io", Version: "v1alpha1", Resource: "storageversionmigrations"}

var storageversionmigrationsKind = schema.GroupVersionKind{Group: "migration.k8s.io", Version: "v1alpha1", Kind: "StorageVersionMigration"}

// Get takes name of the storageVersionMigration, and returns the corresponding storageVersionMigration object, and an error if there is any.
func (c *FakeStorageVersionMigrations) Get(name string, options v1.GetOptions) (result *v1alpha1.StorageVersionMigration, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(storageversionmigrationsResource, name), &v1alpha1.StorageVersionMigration{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.StorageVersionMigration), err
}

// List takes label and field selectors, and returns the list of StorageVersionMigrations that match those selectors.
func (c *FakeStorageVersionMigrations) List(opts v1.ListOptions) (result *v1alpha1.StorageVersionMigrationList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(storageversionmigrationsResource, storageversionmigrationsKind, opts), &v1alpha1.StorageVersionMigrationList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.StorageVersionMigrationList{ListMeta: obj.(*v1alpha1.StorageVersionMigrationList).ListMeta}
	for _, item := range obj.(*v1alpha1.StorageVersionMigrationList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested storageVersionMigrations.
func (c *FakeStorageVersionMigrations) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(storageversionmigrationsResource, opts))
}

// Create takes the representation of a storageVersionMigration and creates it.  Returns the server's representation of the storageVersionMigration, and an error, if there is any.
func (c *FakeStorageVersionMigrations) Create(storageVersionMigration *v1alpha1.StorageVersionMigration) (result *v1alpha1.StorageVersionMigration, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(storageversionmigrationsResource, storageVersionMigration), &v1alpha1.StorageVersionMigration{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.StorageVersionMigration), err
}

// Update takes the representation of a storageVersionMigration and updates it. Returns the server's representation of the storageVersionMigration, and an error, if there is any.
func (c *FakeStorageVersionMigrations) Update(storageVersionMigration *v1alpha1.StorageVersionMigration) (result *v1alpha1.StorageVersionMigration, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(storageversionmigrationsResource, storageVersionMigration), &v1alpha1.StorageVersionMigration{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.StorageVersionMigration), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeStorageVersionMigrations) UpdateStatus(storageVersionMigration *v1alpha1.StorageVersionMigration) (*v1alpha1.StorageVersionMigration, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(storageversionmigrationsResource, "status", storageVersionMigration), &v1alpha1.StorageVersionMigration{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.StorageVersionMigration), err
}

// Delete takes name of the storageVersionMigration and deletes it. Returns an error if one occurs.
func (c *FakeStorageVersionMigrations) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(storageversionmigrationsResource, name), &v1alpha1.StorageVersionMigration{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeStorageVersionMigrations) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(storageversionmigrationsResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.StorageVersionMigrationList{})
	return err
}

// Patch applies the patch and returns the patched storageVersionMigration.
func (c *FakeStorageVersionMigrations) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.StorageVersionMigration, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(storageversionmigrationsResource, name, pt, data, subresources...), &v1alpha1.StorageVersionMigration{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.StorageVersionMigration), err
}
