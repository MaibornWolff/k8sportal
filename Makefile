PROJECT = k8sportal

build: 
		cd src/main && docker build . -t k8sportal-image
		kubectl delete pod --selector=app.kubernetes.io/instance=k8sportal
		helm upgrade --install k8sportal --force src/helm/k8sportal