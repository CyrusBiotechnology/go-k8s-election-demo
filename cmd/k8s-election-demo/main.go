package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/CyrusBiotechnology/go-k8s-election-demo/pkg/k8s"
	"github.com/golang/glog"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
)

func main() {
	hostname, _ := os.Hostname()

	electionid := flag.String("election", "go-k8s-election-demo", "Leader election ID (name of configmap)")
	kubeconfig := flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "Absolute path to kubeconfig file")
	ttlseconds := flag.Int("ttl", 15, "TTL for leader election in seconds")

	flag.Parse()

	// TTL of the lock
	ttl := time.Duration(*ttlseconds) * time.Second

	// Get a kubernetes client
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	print(config)
	client := clientset.NewForConfigOrDie(config)

	broadcaster := record.NewBroadcaster()

	recorder := broadcaster.NewRecorder(scheme.Scheme, apiv1.EventSource{
		Component: *electionid,
		Host:      hostname,
	})

	pod, err := k8s.GetPodDetails(client)
	if err != nil {
		panic(err)
	}

	// The lock that we will share
	lock := resourcelock.ConfigMapLock{
		ConfigMapMeta: metav1.ObjectMeta{Namespace: pod.Namespace, Name: *electionid},
		Client:        client.CoreV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity:      pod.Name,
			EventRecorder: recorder,
		},
	}

	// Events from Kubernetes
	callbacks := leaderelection.LeaderCallbacks{
		OnStartedLeading: func(stop <-chan struct{}) {
			// leaderelection will log the event
			// TODO do master work
			// Block until done leading
			<-stop
			glog.Info("stopped leading")
		},
		OnStoppedLeading: func() {
			// leaderelection will log the event
		},
		OnNewLeader: func(identity string) {
			if identity == pod.Name {
				// I just got the lock
				return
			}
			glog.Infof("new leader elected: %v", identity)
		},
	}

	// Configure the leader election
	le, err := leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{
		Lock:          &lock,
		LeaseDuration: ttl,
		RenewDeadline: ttl / 2,
		RetryPeriod:   ttl / 4,
		Callbacks:     callbacks,
	})

	// Start election
	le.Run()
}
