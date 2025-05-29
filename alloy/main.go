package main

import (
    "context"
    "fmt"
    "os"
    "time"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/leaderelection"
    "k8s.io/client-go/tools/leaderelection/resourcelock"
)

func main() {
    podName := os.Getenv("POD_NAME")
    if podName == "" {
        fmt.Println("POD_NAME not set")
        os.Exit(1)
    }

    config, _ := rest.InClusterConfig()
    clientset, _ := kubernetes.NewForConfig(config)

    lock := &resourcelock.LeaseLock{
        LeaseMeta: metav1.ObjectMeta{
            Name:      "testing-leader-election",
            Namespace: "default",
        },
        Client: clientset.CoordinationV1(),
        LockConfig: resourcelock.ResourceLockConfig{
            Identity: podName,
        },
    }

    ctx := context.TODO()
    leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
        Lock:          lock,
        LeaseDuration: 15 * time.Second,
        RenewDeadline: 10 * time.Second,
        RetryPeriod:   2 * time.Second,
        Callbacks: leaderelection.LeaderCallbacks{
            OnStartedLeading: func(ctx context.Context) {
                fmt.Printf("[%s] I am the leader, having the lock, doing work...\n", podName)
                for {
                    time.Sleep(5 * time.Second)
                    fmt.Printf("[%s] Working podName...\n", podName)
                }
            },
            OnStoppedLeading: func() {
                fmt.Printf("[%s] Lost the lease-lock, exiting. podName \n", podName)
                os.Exit(0)
            },
        },
    })
}
