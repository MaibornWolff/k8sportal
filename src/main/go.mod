module k8sportal

go 1.15

require github.com/gin-gonic/gin v1.6.3

require k8s.io/api v0.19.4

require k8s.io/client-go v0.19.4

require (
	github.com/rs/zerolog v1.20.0
	github.com/vrischmann/envconfig v1.3.0
	go.mongodb.org/mongo-driver v1.4.4
	k8s.io/apimachinery v0.19.4
)
