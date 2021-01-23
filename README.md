# k8sportal
A web portal that shows a dynamic link list to installed and labeled services in your Kubernetes cluster


[![Portal](resources/k8sportal_small.png)](resources/k8sportal_small.png)

## Cluster Configuration

In order to make services show up on the portal website you have to edit your services and the corresponding ingresses. Therefore follow step 1 and 2. 

1. Add `clusterPortalShow : true` to the labels of the service you want to show up on the cluster portal
2. Add `clusterPortalShow : true` to the labels of the ingress that point to the service from step 1
3. (optional) Add `clusterPortalCategory : <insert the Category here>` to to the labels of a service.

## Installation

1. Change the database address and the correspinding username and password in `src/main/config.go`. The database and the collection of the portal will be created at first statup and cleaned up after every redeployment. 
2. (Optional) Configure the values for Helm in `src/helm/k8sportal/values.yaml`
3. (Optional) Adapt the `src/main/web/web.go` file accordingly to step 2 to fit your helm chart
4. Create the image with the included make file
5. Deploy with helm

