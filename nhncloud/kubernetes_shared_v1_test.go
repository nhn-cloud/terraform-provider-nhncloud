package nhncloud

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/kubernetes/v1/clustertemplates"
)

func TestUnitExpandKubernetesV1LabelsMap(t *testing.T) {
	labels := map[string]interface{}{
		"foo":                   "bar",
		"bar":                   "baz",
		"pods_network_subnet":   "24",
		"extra_volumes":         `["vol1","vol2"]`,
		"extra_security_groups": `["sg-123","sg-456"]`,
	}

	actualLabels, err := expandKubernetesV1LabelsMap(labels)
	assert.NoError(t, err)

	assert.Equal(t, "bar", actualLabels["foo"])
	assert.Equal(t, "baz", actualLabels["bar"])
	assert.Equal(t, "24", actualLabels["pods_network_subnet"])

	extraVolumes, ok := actualLabels["extra_volumes"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(extraVolumes))
	assert.Equal(t, "vol1", extraVolumes[0])
	assert.Equal(t, "vol2", extraVolumes[1])

	extraSgs, ok := actualLabels["extra_security_groups"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(extraSgs))
	assert.Equal(t, "sg-123", extraSgs[0])
	assert.Equal(t, "sg-456", extraSgs[1])
}

func TestUnitExpandKubernetesV1LabelsMapWithVariousTypes(t *testing.T) {
	labels := map[string]interface{}{
		"string_value":  "test",
		"integer_value": 123,
		"float_value":   45.67,
		"bool_value":    true,
		"array_value":   `["a","b","c"]`,
		"object_value":  `{"key":"value"}`,
	}

	actualLabels, err := expandKubernetesV1LabelsMap(labels)
	assert.NoError(t, err)

	assert.Equal(t, "test", actualLabels["string_value"])
	assert.Equal(t, 123, actualLabels["integer_value"])
	assert.Equal(t, 45.67, actualLabels["float_value"])
	assert.Equal(t, true, actualLabels["bool_value"])

	arrayVal, ok := actualLabels["array_value"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 3, len(arrayVal))
	assert.Equal(t, "a", arrayVal[0])

	// Objects remain as JSON strings
	assert.Equal(t, `{"key":"value"}`, actualLabels["object_value"])
}

func TestUnitExpandKubernetesV1LabelsString(t *testing.T) {
	labels := map[string]interface{}{
		"foo": "bar",
		"bar": "baz",
	}

	expectedLabels1 := "{'foo':'bar','bar':'baz'}"
	expectedLabels2 := "{'bar':'baz','foo':'bar'}"

	actualLabels, err := expandKubernetesV1LabelsString(labels)
	assert.Equal(t, err, nil)

	if actualLabels != expectedLabels1 && actualLabels != expectedLabels2 {
		t.Fatalf("Unexpected labels. Got %s, expected %s or %s",
			actualLabels, expectedLabels1, expectedLabels2)
	}
}

func TestUnitKubernetesClusterTemplateV1AppendUpdateOpts(t *testing.T) {
	actualUpdateOpts := []clustertemplates.UpdateOptsBuilder{}

	expectedUpdateOpts := []clustertemplates.UpdateOptsBuilder{
		clustertemplates.UpdateOpts{
			Op:    clustertemplates.ReplaceOp,
			Path:  "/master_lb_enabled",
			Value: "True",
		},
		clustertemplates.UpdateOpts{
			Op:    clustertemplates.ReplaceOp,
			Path:  "/registry_enabled",
			Value: "True",
		},
	}

	actualUpdateOpts = kubernetesClusterTemplateV1AppendUpdateOpts(
		actualUpdateOpts, "master_lb_enabled", "True")

	actualUpdateOpts = kubernetesClusterTemplateV1AppendUpdateOpts(
		actualUpdateOpts, "registry_enabled", "True")

	assert.Equal(t, expectedUpdateOpts, actualUpdateOpts)
}
