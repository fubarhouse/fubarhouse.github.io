package main

import "k8s.io/klog/v2"

func log(name string, action string, category string, state bool) {
	if !state {
		klog.Warningf("could not %v %v %v\n", action, category, name)
		return
	}
	klog.Infof("%v %v was %vd\n", name, category, action)
}
