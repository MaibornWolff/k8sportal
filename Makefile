PROJECT = k8sportal

build:
		cd src/main && docker build . -t k8sportal-image:0.9

		cd src/frontend && docker build . -t frontend-k8sportal-image:0.9
		helm upgrade --install frontend --force src/helm/frontend

		kubectl delete pod --selector=app.kubernetes.io/instance=k8sportal
		helm upgrade --install k8sportal --force src/helm/k8sportal
