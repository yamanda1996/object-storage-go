package es

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"object-storage-go/api-server/model"
	"strings"
)

type Metadata struct {
	Name 			string
	Version 		int
	Size 			int64
	Hash 			string
}

type SearchResult struct {
	Hits struct{
		Total 		int
		Hits 		[] struct{
			Source 	Metadata `json:"_source"`
		}
	}
}

func GetMetadata(name string, version int) (Metadata, error) {
	// 默认返回最新的元数据
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}


func getMetadata(name string, versionId int) (meta Metadata, err error) {
	url := fmt.Sprintf("http://%s:%d/metadata/objects/%s_%d/_source", model.Config.ElasticSearchConfig.ElasticSearchAddress,
		model.Config.ElasticSearchConfig.ElasticSearchPort, name, versionId)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("get metadata from elastic search failed")
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to get %s_%d", name, versionId)
		return
	}
	result, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(result, &meta)
	return
}

func SearchLatestVersion(name string) (meta Metadata, err error) {
	u := fmt.Sprintf("http://%s:%d/metadata/_search?q=name:%s&size=1&sort=version:desc",
		model.Config.ElasticSearchConfig.ElasticSearchAddress, model.Config.ElasticSearchConfig.ElasticSearchPort, url.PathEscape(name))
	resp, err := http.Get(u)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to search latest version")
		return
	}
	result, err := ioutil.ReadAll(resp.Body)

	var sr SearchResult
	json.Unmarshal(result, &sr)
	if err != nil {
		log.Error("unmarshal search result json failed")
		return
	}
	if len(sr.Hits.Hits) > 0 {
		meta = sr.Hits.Hits[0].Source
	}
	return
}

func PutMetadata(name string, version int, size int64, hash string) error {

	doc := fmt.Sprintf(`{"name":"%s","version":%d,"size":%d,"hash":"%s"}`, name, version, size, hash)
	client := http.Client{}
	url := fmt.Sprintf("http://%s:%d/metadata/objects/%s_%d?op_type=create", model.Config.ElasticSearchConfig.ElasticSearchAddress,
		model.Config.ElasticSearchConfig.ElasticSearchPort, name, version)
	request, _ := http.NewRequest("PUT", url, strings.NewReader(doc))
	request.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(request)

	if err != nil {
		log.Errorf("send request to es failed")
		return err
	}
	defer resp.Body.Close()
	// 如果有多个客户端同时上传同一个元数据,es会返回409 conflict,别的客户端版本+1后上传
	if resp.StatusCode == http.StatusConflict {
		log.Warn("multi client upload file at the same time")
		return PutMetadata(name, version + 1, size, hash)
	}

	if resp.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("fail to put metadata %s_%d, resp %s", name, version, string(result))
	}
	return nil
}

func AddVersion(name, hash string, size int64) (error) {
	metadata , err := SearchLatestVersion(name)
	if err != nil {
		log.Errorf("search [%s] latest version failed", name)
		return fmt.Errorf("search [%s] latest version failed", name)
	}
	err = PutMetadata(name, metadata.Version + 1, size, hash)
	if err != nil {
		log.Errorf("put [%s] metadata failed", name)
		return fmt.Errorf("put [%s] metadata failed", name)
	}
	return nil
}

func SearchAllVersions(name string, from, size int) ([]Metadata, error) {
	url := fmt.Sprintf("http://%s:%d/metadata/_search?sort=name,version&from=%d&size=%d",
		model.Config.ElasticSearchConfig.ElasticSearchAddress, model.Config.ElasticSearchConfig.ElasticSearchPort,
		from, size)
	if len(name) > 0 {
		url += "&q=name:" + name
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	metadatas := make([]Metadata, 0)
	result, _ := ioutil.ReadAll(resp.Body)
	var sr SearchResult
	json.Unmarshal(result, &sr)

	for _, hit := range sr.Hits.Hits {
		metadatas = append(metadatas, hit.Source)
	}
	return metadatas, nil
}


