package model

//Service model
type Service struct {
	ServiceName   string        `bson:"serviceName" json:"serviceName"`
	Category      string        `bson:"category" json:"category"`
	ServiceOnline bool          `bson:"serviceOnline" json:"serviceOnline"`
	IngressRules  []IngressRule `bson:"ingressRules" json:"ingressRules"`
	IngressOnline bool          `bson:"ingressOnline" json:"ingressOnline"`
}

//IngressRule model
type IngressRule struct {
	IngressHost string `bson:"ingressHost" json:"ingressHost"`
	IngressPath string `bson:"ingressPath" json:"ingressPath"`
}

//IsOnline checks if ingress and service are available
func (service *Service) IsOnline() bool {
	return (service.ServiceOnline && service.IngressOnline)
}
