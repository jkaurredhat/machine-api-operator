package main

import (
	"flag"
	"github.com/golang/glog"
	osconfigv1 "github.com/openshift/api/config/v1"
	"k8s.io/client-go/kubernetes/scheme"
	capiv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	namespace = "openshift-cluster-api"
)

func init() {
	if err := capiv1alpha1.AddToScheme(scheme.Scheme); err != nil {
		glog.Fatal(err)
	}

	if err := osconfigv1.AddToScheme(scheme.Scheme); err != nil {
		glog.Fatal(err)
	}
}

type testConfig struct {
	client client.Client
}

func newClient() (client.Client, error) {
	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	return client.New(cfg, client.Options{})

}

func main() {
	flag.Parse()
	if err := runSuite(); err != nil {
		glog.Fatal(err)
	}
}

func runSuite() error {

	client, err := newClient()
	if err != nil {
		return err
	}
	testConfig := &testConfig{
		client: client,
	}
	glog.Info("RUN: ExpectOperatorAvailable")
	if err := testConfig.ExpectOperatorAvailable(); err != nil {
		glog.Errorf("FAIL: ExpectOperatorAvailable: %v", err)
		return err
	}
	glog.Info("PASS: ExpectOperatorAvailable")

	glog.Info("RUN: ExpectOneClusterObject")
	if err := testConfig.ExpectOneClusterObject(); err != nil {
		glog.Errorf("FAIL: ExpectOneClusterObject: %v", err)
		return err
	}
	glog.Info("PASS: ExpectOneClusterObject")

	glog.Info("RUN: ExpectClusterOperatorStatusAvailable")
	if err := testConfig.ExpectClusterOperatorStatusAvailable(); err != nil {
		glog.Errorf("FAIL: ExpectClusterOperatorStatusAvailable: %v", err)
		return err
	}
	glog.Info("PASS: ExpectClusterOperatorStatusAvailable")

	glog.Info("RUN: ExpectAllMachinesLinkedToANode")
	if err := testConfig.ExpectAllMachinesLinkedToANode(); err != nil {
		glog.Errorf("FAIL: ExpectAllMachinesLinkedToANode: %v", err)
		return err
	}
	glog.Info("PASS: ExpectAllMachinesLinkedToANode")

	glog.Info("RUN: ExpectReconcileControllersDeployment")
	if err := testConfig.ExpectReconcileControllersDeployment(); err != nil {
		glog.Errorf("FAIL: ExpectReconcileControllersDeployment: %v", err)
		return err
	}
	glog.Info("PASS: ExpectReconcileControllersDeployment")

	glog.Info("RUN: ExpectNewNodeWhenDeletingMachine")
	if err := testConfig.ExpectNewNodeWhenDeletingMachine(); err != nil {
		glog.Errorf("FAIL: ExpectNewNodeWhenDeletingMachine: %v", err)
		return err
	}
	glog.Info("PASS: ExpectNewNodeWhenDeletingMachine")
	return nil
}
