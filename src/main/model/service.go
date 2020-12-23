package model

//Service model
type Service struct {
	ServiceName   string `bson:"serviceName" json:"serviceName"`
	Category      string `bson:"category" json:"category"`
	ServiceOnline bool   `bson:"serviceOnline" json:"serviceOnline"`
	IngressHost   string `bson:"ingressHost" json:"ingressHost"`
	IngressPath   string `bson:"ingressPath" json:"ingressPath"`
	IngressOnline bool   `bson:"ingressOnline" json:"ingressOnline"`
}

//IsOnline checks if ingress and service are available
func IsOnline(service Service) bool {
	return (service.ServiceOnline && service.IngressOnline)
}

//GetFqdn returns the concatenated Host and Path to access the service via the ingress
func GetFqdn(service Service) string {
	return service.IngressHost + service.IngressPath
}
