package model

//Service model
type Service struct {
	ServiceName   string `bson:"serviceName" json:"serviceName"`
	Category      string `bson:"category" json:"category"`
	ServiceOnline bool   `bson:"serviceOnline" json:"serviceOnline"`
	IngressHost   string `bson:"ingressHost" json:"ingressHost"`
	IngressPath   string `bson:"ingressPath" json:"ingressPath"`
	Fqdn          string `bson:"fqdn" json:"fqdn"`
	IngressOnline bool   `bson:"ingressOnline" json:"ingressOnline"`
}

//IsOnline checks if ingress and service are available
func (service *Service) IsOnline() bool {
	return (service.ServiceOnline && service.IngressOnline)
}

//GetFqdn returns the concatenated Host and Path to access the service via the ingress
func (service *Service) GetFqdn(empty string) string {
	return service.IngressHost + service.IngressPath
}
