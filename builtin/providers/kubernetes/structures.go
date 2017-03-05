package kubernetes

import (
	"fmt"

	api "k8s.io/kubernetes/pkg/api/v1"
)

func expandMetadata(in []interface{}) api.ObjectMeta {
	meta := api.ObjectMeta{}
	m := in[0].(map[string]interface{})

	meta.Annotations = expandStringMap(m["annotations"].(map[string]interface{}))
	meta.Labels = expandStringMap(m["labels"].(map[string]interface{}))

	if v, ok := m["generate_name"]; ok {
		meta.GenerateName = v.(string)
	}
	if v, ok := m["name"]; ok {
		meta.Name = v.(string)
	}
	if v, ok := m["namespace"]; ok {
		meta.Namespace = v.(string)
	}

	return meta
}

func expandStringMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = v.(string)
	}
	return result
}

func flattenMetadata(meta api.ObjectMeta) []map[string]interface{} {
	m := make(map[string]interface{})
	m["annotations"] = meta.Annotations
	m["generate_name"] = meta.GenerateName
	m["labels"] = meta.Labels
	m["name"] = meta.Name
	m["resource_version"] = meta.ResourceVersion
	m["self_link"] = meta.SelfLink
	m["uid"] = fmt.Sprintf("%v", meta.UID)
	m["generation"] = meta.Generation

	if meta.Namespace != "" {
		m["namespace"] = meta.Namespace
	}

	return []map[string]interface{}{m}
}
