PROJECT = k8sportal

build: 
		cd src/main && docker build . -t eneayg/k8sportal-image:0.8 && docker push eneayg/k8sportal-image:0.8
		
		cd src/frontend && docker build . -t eneayg/frontend-k8sportal-image:0.8 && docker push eneayg/frontend-k8sportal-image:0.8
		helm upgrade --install frontend --force src/helm/frontend
		
		kubectl delete pod --selector=app.kubernetes.io/instance=k8sportal
		helm upgrade --install k8sportal --force src/helm/k8sportal
		