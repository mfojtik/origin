package quota_test

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/resource"
	kapi "k8s.io/kubernetes/pkg/apis/core"

	quotaapi "github.com/openshift/openshift-apiserver/pkg/quota/apis/quota"
)

func TestDeepCopy(t *testing.T) {
	make := func() *quotaapi.ClusterResourceQuota {
		q := resource.Quantity{}
		q.Set(100)
		crq := &quotaapi.ClusterResourceQuota{}
		crq.Status.Namespaces.Insert("ns1", kapi.ResourceQuotaStatus{Hard: kapi.ResourceList{"a": q.DeepCopy()}, Used: kapi.ResourceList{"a": q.DeepCopy()}})
		crq.Status.Namespaces.Insert("ns2", kapi.ResourceQuotaStatus{Hard: kapi.ResourceList{"b": q.DeepCopy()}, Used: kapi.ResourceList{"b": q.DeepCopy()}})
		return crq
	}

	check := make()

	original := make()
	if !reflect.DeepEqual(check, original) {
		t.Error("before mutation of copy, check and original should be identical but are not, likely failure in deepequal")
	}
	if !equality.Semantic.DeepEqual(check, original) {
		t.Error("before mutation of copy, check and original should be identical but are not, likely failure in deepequal")
	}

	copied := original.DeepCopy()
	if !reflect.DeepEqual(copied, original) {
		t.Error("before mutation of copy, copied and original should be identical but are not, likely failure in deepequal")
	}
	if !equality.Semantic.DeepEqual(copied, original) {
		t.Error("before mutation of copy, copied and original should be identical but are not, likely failure in deepequal")
	}

	// Mutate the copy
	for e := copied.Status.Namespaces.OrderedKeys().Front(); e != nil; e = e.Next() {
		k := e.Value.(string)
		ns, _ := copied.Status.Namespaces.Get(k)
		for k2, v2 := range ns.Hard {
			v2.Set(v2.Value() + 2)
			ns.Hard[k2] = v2
		}
		for k2, v2 := range ns.Used {
			v2.Set(v2.Value() + 1)
			ns.Used[k2] = v2
		}
		copied.Status.Namespaces.Insert(k, ns)
	}

	if !reflect.DeepEqual(check, original) {
		t.Error("after mutation of copy, check and original should be identical but are not, likely failure in deep copy (ensure custom DeepCopy is being used)")
	}
	if !equality.Semantic.DeepEqual(check, original) {
		t.Error("after mutation of copy, check and original should be identical but are not, likely failure in deep copy (ensure custom DeepCopy is being used)")
	}

	if reflect.DeepEqual(original, copied) {
		t.Error("after mutation of copy, original and copied should be different but are not, likely failure in deep copy (ensure custom DeepCopy is being used)")
	}
	if equality.Semantic.DeepEqual(original, copied) {
		t.Error("after mutation of copy, original and copied should be different but are not, likely failure in deep copy (ensure custom DeepCopy is being used)")
	}
}
