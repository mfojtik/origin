package build

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestLogOptionsDrift(t *testing.T) {
	popts := reflect.TypeOf(corev1.PodLogOptions{})
	bopts := reflect.TypeOf(BuildLogOptions{})

	for i := 0; i < popts.NumField(); i++ {
		// Verify name
		name := popts.Field(i).Name
		boptsField, found := bopts.FieldByName(name)
		if !found {
			t.Errorf("buildLogOptions drifting from podLogOptions! Field %q wasn't found!", name)
		}
		// Verify type
		if should, is := popts.Field(i).Type, boptsField.Type; is != should {
			t.Errorf("buildLogOptions drifting from podLogOptions! Field %q should be a %s but is %s!", name, should.String(), is.String())
		}
	}
}
