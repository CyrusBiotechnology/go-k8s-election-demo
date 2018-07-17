# Go Kubernetes Leader Election Demo

This code demonstrates coordinated locking on the kubernetes platform for go 
applications.

# Running

Notes:
1) The app needs a unique POD_NAME environment variable, as well as a common 
POD_NAMESPACE variable.
2) The app talks to the Kubernetes cluster in the currently configured context, 
so be sure to point it at a local or non-production cluster. 

To see the app in action, run the following three commands in separate 
terminals. `cd` to a clone of this repo and run:


    POD_NAME=test-1 POD_NAMESPACE=default go run cmd/k8s-election-demo/main.go -logtostderr
   
This "pod" should successfully acquire the lock right away (as long as there is not a 
lock that is currently valid).


    POD_NAME=test-2 POD_NAMESPACE=default go run cmd/k8s-election-demo/main.go -logtostderr
    POD_NAME=test-3 POD_NAMESPACE=default go run cmd/k8s-election-demo/main.go -logtostderr

Now kill the first command. After a short time you should see either test-2 or test-3 
acquire the lock, with the other process seeing and printing the new owner.

