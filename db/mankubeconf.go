package db

import (
	"bytes"
	b64 "encoding/base64"
	"io"
	"mime/multipart"
)

type Kubeconfig struct {
	Kubeconfig string `json:"kubeconfig" binding:"required"`
}

func (database *Database) SetDataWithBase64(key string,kubeconfig []byte) (string,error) {
	sEnc := b64.StdEncoding.EncodeToString(kubeconfig)
	set := database.Client.Set(Ctx, key,sEnc,0)
	result, err := set.Result()
	if err != nil {
		return "",err
	}
	return result,nil
}

func (database *Database) GetDataWithBase64(key string) (string, error) {
	get := database.Client.Get(Ctx, key)
	result, err := get.Result()
	if err != nil {
		return "", ErrNil
	}
	sDec, err := b64.StdEncoding.DecodeString(result)
	if err != nil {
		return "", err
	}
	return string(sDec), nil
}

func (database *Database) GetManagementKubeconfig(key string) (string,error) {
	kubeConfig, err := database.GetDataWithBase64(key)
	if err != nil {
		return "",err
	}
	return kubeConfig,nil
}

func (database *Database) SetManagementKubeconfig(key string,file *multipart.FileHeader) error {
	fileTest,err := file.Open()
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, fileTest)
	if err != nil {
		return err
	}

	err = fileTest.Close()
	if err != nil {
		return err
	}

	_, err = database.SetDataWithBase64(key,buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}
