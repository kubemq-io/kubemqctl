package operator

import (
	"github.com/ghodss/yaml"
	v1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

var crdKubemqConnector = `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: kubemqconnectors.core.k8s.kubemq.io
spec:
  group: core.k8s.kubemq.io
  names:
    kind: KubemqConnector
    listKind: KubemqConnectorList
    plural: kubemqconnectors
    singular: kubemqconnector
  scope: Namespaced
  versions:
    - name: v1beta1
      served: true
      storage: true
      subresources:
        scale:
          labelSelectorPath: .status.selector
          specReplicasPath: .spec.replicas
          statusReplicasPath: .status.replicas
        status: { }
      additionalPrinterColumns:
        - jsonPath: .status.type
          name: Type
          type: string
        - jsonPath: .status.replicas
          name: Replicas
          type: string
        - jsonPath: .status.image
          name: Image
          type: string
        - jsonPath: .status.api
          name: API
          type: string
        - jsonPath: .status.status
          name: Status
          type: string
      schema:
        openAPIV3Schema:
          description: KubemqConnector is the Schema for the kubemqconnectors API
          properties:
            apiVersion:
              description: 'APIVersion defines the versioned schema of this representation
                of an object. Servers should convert recognized schemas to the latest
                internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
              type: string
            kind:
              description: 'Kind is a string value representing the REST resource this
                object represents. Servers may infer this from the endpoint the client
                submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
              type: string
            metadata:
              type: object
            spec:
              description: KubemqConnectorSpec defines the desired state of KubemqConnector
              properties:
                config:
                  type: string
                image:
                  type: string
                node_port:
                  format: int32
                  type: integer
                replicas:
                  format: int32
                  minimum: 0
                  type: integer
                service_type:
                  type: string
                type:
                  type: string
              required:
                - config
                - type
              type: object
            status:
              description: KubemqConnectorStatus defines the observed state of KubemqConnector
              properties:
                api:
                  type: string
                image:
                  type: string
                replicas:
                  format: int32
                  type: integer
                status:
                  type: string
                type:
                  type: string
              required:
                - api
                - image
                - replicas
                - status
                - type
              type: object
          type: object
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kubemq-operator
rules:
  - apiGroups:
      - ""
    resources:
      - pods
      - services
      - endpoints
      - persistentvolumeclaims
      - events
      - configmaps
      - serviceaccounts
      - secrets
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - apiextensions.k8s.io
    resources:
      - customresourcedefinitions
    verbs:
      - patch
      - update
      - create
      - get
  - apiGroups:
      - apps
    resources:
      - deployments
      - replicasets
      - statefulsets
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - monitoring.coreos.com
    resources:
      - servicemonitors
    verbs:
      - get
      - create
  - apiGroups:
      - core.k8s.kubemq.io
    resources:
      - "*"
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
`

type KubemqConnectorCRD struct {
	crd *v1beta1.CustomResourceDefinition
}

func CreateKubemqConnectorCRD() *KubemqConnectorCRD {
	return &KubemqConnectorCRD{
		crd: nil,
	}
}

func (sa *KubemqConnectorCRD) Spec() ([]byte, error) {
	t := NewTemplate(crdKubemqConnector, sa)
	return t.Get()
}

func (sa *KubemqConnectorCRD) Get() (*v1beta1.CustomResourceDefinition, error) {
	if sa.crd != nil {
		return sa.crd, nil
	}
	crd := &v1beta1.CustomResourceDefinition{}
	data, err := sa.Spec()
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, crd)
	if err != nil {
		return nil, err
	}
	sa.crd = crd
	return crd, nil
}
