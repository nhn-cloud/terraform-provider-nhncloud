package nhncloud

import (
	"reflect"
	"testing"

	"github.com/nhn/nhncloud.gophercloud/openstack/blockstorage/v3/volumetypes"
)

func TestExpandBlockStorageVolumeTypeV3ExtraSpecs(t *testing.T) {
	raw := map[string]interface{}{
		"foo": "foo",
		"bar": "bar",
	}

	expected := volumetypes.ExtraSpecsOpts{
		"foo": "foo",
		"bar": "bar",
	}

	actual := expandBlockStorageVolumeTypeV3ExtraSpecs(raw)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Results differ. Want: %#v, but got %#v", expected, actual)
	}
}
