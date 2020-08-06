package framework

import (
	"context"
	"fmt"
	"reflect"

	kappsapi "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetDaemonSet gets DaemonSet object by name and namespace.
func GetDaemonSet(c client.Client, name, namespace string) (*kappsapi.DaemonSet, error) {
	key := types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
	d := &kappsapi.DaemonSet{}

	if err := wait.PollImmediate(RetryShort, WaitShort, func() (bool, error) {
		if err := c.Get(context.TODO(), key, d); err != nil {
			klog.Errorf("Error querying api for DaemonSet object %q: %v, retrying...", name, err)
			return false, nil
		}
		return true, nil
	}); err != nil {
		return nil, fmt.Errorf("error getting DaemonSet %q: %v", name, err)
	}
	return d, nil
}

// DeleteDaemonSet deletes the specified DaemonSet
func DeleteDaemonSet(c client.Client, daemonSet *kappsapi.DaemonSet) error {
	return wait.PollImmediate(RetryShort, WaitShort, func() (bool, error) {
		if err := c.Delete(context.TODO(), daemonSet); err != nil {
			klog.Errorf("error querying api for DaemonSet object %q: %v, retrying...", daemonSet.Name, err)
			return false, nil
		}
		return true, nil
	})
}

// UpdateDaemonSet updates the specified DaemonSet
func UpdateDaemonSet(c client.Client, name, namespace string, updated *kappsapi.DaemonSet) error {
	return wait.PollImmediate(RetryShort, WaitMedium, func() (bool, error) {
		d, err := GetDaemonSet(c, name, namespace)
		if err != nil {
			klog.Errorf("Error getting daemonSet: %v", err)
			return false, nil
		}
		if err := c.Patch(context.TODO(), d, client.MergeFrom(updated)); err != nil {
			klog.Errorf("error patching daemonSet object %q: %v, retrying...", name, err)
			return false, nil
		}
		return true, nil
	})
}

// IsDaemonSetAvailable returns true if the daemonSet has one or more availabe replicas
func IsDaemonSetAvailable(c client.Client, name, namespace string) bool {
	if err := wait.PollImmediate(RetryShort, WaitLong, func() (bool, error) {
		d, err := GetDaemonSet(c, name, namespace)
		if err != nil {
			klog.Errorf("Error getting daemonSet: %v", err)
			return false, nil
		}
		if d.Status.NumberReady < 1 {
			klog.Errorf("DaemonSet %q is not available. Status: %s",
				d.Name, daemonSetInfo(d))
			return false, nil
		}
		klog.Infof("DaemonSet %q is available. Status: %s",
			d.Name, daemonSetInfo(d))
		return true, nil
	}); err != nil {
		klog.Errorf("Error checking isDaemonSetAvailable: %v", err)
		return false
	}
	return true
}

// IsDaemonSetSynced returns true if provided daemonSet spec matched one found on cluster
func IsDaemonSetSynced(c client.Client, ds *kappsapi.DaemonSet, name, namespace string) bool {
	d, err := GetDaemonSet(c, name, namespace)
	if err != nil {
		klog.Errorf("Error getting daemonSet: %v", err)
		return false
	}
	if !reflect.DeepEqual(d.Spec, ds.Spec) {
		klog.Errorf("DaemonSet %q is not updated. Spec is not equal to: %v",
			d.Name, ds.Spec)
		return false
	}
	klog.Infof("DaemonSet %q is updated. Spec is matched", d.Name)
	return true
}

func daemonSetInfo(d *kappsapi.DaemonSet) string {
	return fmt.Sprintf("(desired: %d, updated: %d, available: %d, unavailable: %d)",
		d.Status.DesiredNumberScheduled, d.Status.UpdatedNumberScheduled,
		d.Status.NumberAvailable, d.Status.NumberUnavailable)
}
