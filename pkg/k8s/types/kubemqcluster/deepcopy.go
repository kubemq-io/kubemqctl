package kubemqcluster

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubemqCluster) DeepCopyInto(out *KubemqCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)

}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubemqCluster.
func (in *KubemqCluster) DeepCopy() *KubemqCluster {
	if in == nil {
		return nil
	}
	out := new(KubemqCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KubemqCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubemqClusterList) DeepCopyInto(out *KubemqClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]KubemqCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}

}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubemqClusterList.
func (in *KubemqClusterList) DeepCopy() *KubemqClusterList {
	if in == nil {
		return nil
	}
	out := new(KubemqClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KubemqClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubemqClusterSpec) DeepCopyInto(out *KubemqClusterSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.Volume != nil {
		in, out := &in.Volume, &out.Volume
		*out = new(VolumeConfig)
		**out = **in
	}
	if in.License != nil {
		in, out := &in.License, &out.License
		*out = new(LicenseConfig)
		**out = **in
	}
	if in.Image != nil {
		in, out := &in.Image, &out.Image
		*out = new(ImageConfig)
		**out = **in
	}
	if in.Api != nil {
		in, out := &in.Api, &out.Api
		*out = new(ApiConfig)
		**out = **in
	}
	if in.Rest != nil {
		in, out := &in.Rest, &out.Rest
		*out = new(RestConfig)
		**out = **in
	}
	if in.Grpc != nil {
		in, out := &in.Grpc, &out.Grpc
		*out = new(GrpcConfig)
		**out = **in
	}
	if in.Tls != nil {
		in, out := &in.Tls, &out.Tls
		*out = new(TlsConfig)
		**out = **in
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(ResourceConfig)
		**out = **in
	}
	if in.NodeSelectors != nil {
		in, out := &in.NodeSelectors, &out.NodeSelectors
		*out = (*in).DeepCopy()
	}
	if in.Authentication != nil {
		in, out := &in.Authentication, &out.Authentication
		*out = new(AuthenticationConfig)
		**out = **in
	}
	if in.Authorization != nil {
		in, out := &in.Authorization, &out.Authorization
		*out = new(AuthorizationConfig)
		**out = **in
	}
	if in.Health != nil {
		in, out := &in.Health, &out.Health
		*out = new(HealthConfig)
		**out = **in
	}
	if in.Routing != nil {
		in, out := &in.Routing, &out.Routing
		*out = new(RoutingConfig)
		**out = **in
	}
	if in.Log != nil {
		in, out := &in.Log, &out.Log
		*out = (*in).DeepCopy()
	}
	if in.Notification != nil {
		in, out := &in.Notification, &out.Notification
		*out = new(NotificationConfig)
		**out = **in
	}
	if in.Store != nil {
		in, out := &in.Store, &out.Store
		*out = (*in).DeepCopy()
	}
	if in.Queue != nil {
		in, out := &in.Queue, &out.Queue
		*out = (*in).DeepCopy()
	}
	if in.Gateways != nil {
		in, out := &in.Gateways, &out.Gateways
		*out = (*in).DeepCopy()
	}

}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubemqClusterSpec.
func (in *KubemqClusterSpec) DeepCopy() *KubemqClusterSpec {
	if in == nil {
		return nil
	}
	out := new(KubemqClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubemqClusterStatus) DeepCopyInto(out *KubemqClusterStatus) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubemqClusterStatus.
func (in *KubemqClusterStatus) DeepCopy() *KubemqClusterStatus {
	if in == nil {
		return nil
	}
	out := new(KubemqClusterStatus)
	in.DeepCopyInto(out)
	return out
}
