package apiconv

import (
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

func ToYAML(obj client.Object) *string {
	raw, err := yaml.Marshal(removeUnnecessaryFields(obj))
	if err != nil || raw == nil {
		return nil
	}
	return ptr.To(string(raw))
}

func DecodeYAML[T client.Object](raw string, obj T) error {
	return yaml.Unmarshal([]byte(raw), obj)
}

func removeUnnecessaryFields(obj client.Object) client.Object {
	newObj := obj.DeepCopyObject().(client.Object)
	newObj.SetManagedFields(nil)
	return newObj
}
