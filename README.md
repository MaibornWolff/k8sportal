# k8sportal
A web portal to running services in your k8s cluster

## Cluster Configuration

1. Add `clusterPortalShow : true` to the labels of the services you want to show up on the cluster portal
2. Add `clusterPortalShow : true` to the labels of the ingresses to the services you want to show up on the cluster portal
3. (optional) Add `clusterPortalCategory : <insert the Category here>` to to the labels of a service. 

## Installation

1. Change the database address and the correspinding username and password in `src/main/config.go`. The database and the collection of the portal will be created at first statup and cleaned up after every redeployment. 
2. (Optional) Configure the values for Helm in `src/main/k8sportal/values.yaml`
3. (Optional) Adapt the `src/main/web/web.go` file accordingly to step 2 to fit your helm chart
4. Create the image with the included make file
5. Deploy with helm

