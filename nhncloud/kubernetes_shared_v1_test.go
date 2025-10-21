package nhncloud

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhn-cloud/nhncloud.gophercloud/nhncloud/kubernetes/v1/clustertemplates"
)

func TestUnitExpandKubernetesV1LabelsMap(t *testing.T) {
	labels := map[string]interface{}{
		"foo": "bar",
		"bar": "baz",
	}

	expectedLabels := map[string]string{
		"foo": "bar",
		"bar": "baz",
	}

	actualLabels, err := expandKubernetesV1LabelsMap(labels)
	assert.Equal(t, err, nil)
	assert.Equal(t, expectedLabels, actualLabels)
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
