package db

import "errors"

func (database *Database) SetData(key string,value []byte) (string,error) {
	set := database.Client.Set(Ctx, key,value,0)
	result, err := set.Result()
	if err != nil {
		return "",err
	}
	return result,nil
}

func (database *Database) GetData(key string) (string, error) {
	get := database.Client.Get(Ctx, key)
	result, err := get.Result()
	if err != nil {
		return "", errors.New("Cannot find workload cluster kubeconfig!")
	}
	if err != nil {
		return "", err
	}
	return result, nil
}
