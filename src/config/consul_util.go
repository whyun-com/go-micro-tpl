package config

import(
	"github.com/hashicorp/consul/api"
)

type ConsulUtil struct {
	Address string
	client *api.Client
	ProductionLine string
	ProjectName string
	keyPrefix string
}

func NewConsulUtil(util *ConsulUtil) (error){

	config := api.DefaultConfig()
	config.Address = util.Address
	client, err := api.NewClient(config)
	if err != nil {
		return  err
	}
	util.keyPrefix = util.ProductionLine + "/" 
	util.client = client
	return  nil
}

func (util *ConsulUtil) GetValue(key string) ([]byte, error) {
	kvPair,_,err := util.client.KV().Get(util.keyPrefix + key, nil)
	if err != nil {
		return nil, err
	}
	if kvPair == nil {
		return nil, nil
	}
	return kvPair.Value, nil

}