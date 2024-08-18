//go:build !ignore_autogenerated

/*
Copyright 2024.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvironmentVariable) DeepCopyInto(out *EnvironmentVariable) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvironmentVariable.
func (in *EnvironmentVariable) DeepCopy() *EnvironmentVariable {
	if in == nil {
		return nil
	}
	out := new(EnvironmentVariable)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FrontendDeploy) DeepCopyInto(out *FrontendDeploy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FrontendDeploy.
func (in *FrontendDeploy) DeepCopy() *FrontendDeploy {
	if in == nil {
		return nil
	}
	out := new(FrontendDeploy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FrontendDeploy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FrontendDeployList) DeepCopyInto(out *FrontendDeployList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FrontendDeploy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FrontendDeployList.
func (in *FrontendDeployList) DeepCopy() *FrontendDeployList {
	if in == nil {
		return nil
	}
	out := new(FrontendDeployList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FrontendDeployList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FrontendDeploySpec) DeepCopyInto(out *FrontendDeploySpec) {
	*out = *in
	if in.EnvironmentVarialbles != nil {
		in, out := &in.EnvironmentVarialbles, &out.EnvironmentVarialbles
		*out = make([]EnvironmentVariable, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FrontendDeploySpec.
func (in *FrontendDeploySpec) DeepCopy() *FrontendDeploySpec {
	if in == nil {
		return nil
	}
	out := new(FrontendDeploySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FrontendDeployStatus) DeepCopyInto(out *FrontendDeployStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FrontendDeployStatus.
func (in *FrontendDeployStatus) DeepCopy() *FrontendDeployStatus {
	if in == nil {
		return nil
	}
	out := new(FrontendDeployStatus)
	in.DeepCopyInto(out)
	return out
}
