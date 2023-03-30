package service

import (
	"github.com/vibeitco/apikey-service/model"
)

const (
	namespaceImage = "apikey_v1"
)

type apikeys []*model.ApiKey

type apikey struct {
	model.ApiKey `bson:",inline"`
}

func (p *apikeys) GetNamespace() string {
	return namespaceImage
}

func (p *apikey) GetNamespace() string {
	return namespaceImage
}

func (o *apikey) SetId(id string)    { o.Id = id }
func (o *apikey) SetCreated(t int64) { o.Created = t }
func (o *apikey) SetUpdated(t int64) { o.Updated = t }
