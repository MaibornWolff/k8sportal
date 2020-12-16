package model

//Service model
type Service struct {
	ServiceName string `bson:"serviceName" json:"serviceName"`
	FQDN        string `bson:"fqdn" json:"fqdn"`
	Category    string `bson:"category" json:"category"`
	Online      bool   `bson:"online" json:"online"`
}
