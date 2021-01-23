package model

//Service model
type Service struct {
	ServiceName   string        `bson:"serviceName" json:"serviceName"`
	Category      string        `bson:"category" json:"category"`
	ServiceExists bool          `bson:"serviceExists" json:"serviceExists"`
	IngressRules  []IngressRule `bson:"ingressRules" json:"ingressRules"`
	IngressExists bool          `bson:"ingressExists" json:"ingressExists"`
}

//IngressRule model
type IngressRule struct {
	IngressHost string `bson:"ingressHost" json:"ingressHost"`
	IngressPath string `bson:"ingressPath" json:"ingressPath"`
}

//Exists checks if ingress and service are available
func (service *Service) Exists() bool {
	return (service.ServiceExists && service.IngressExists)
}
