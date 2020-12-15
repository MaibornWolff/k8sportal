PROJECT = k8sportal

build: 
		cd src/main && docker build . -t k8sportal-image
