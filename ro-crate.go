package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type ROCrate struct {
	Context string `json:"@context"`
	Graph []map[string]interface{} `json:"@graph"`
}

func NewROCrate() *ROCrate {
	roCrate := ROCrate{
		Context: "https://w3id.org/ro/crate/1.0/context",
	}
	metadataFileDescriptor := NewEntity("CreativeWork","ro-crate-metadata.jsonld")
	metadataFileDescriptor.SetProperty("about",map[string]interface{}{"@id": "./"})
	metadataFileDescriptor.SetProperty("conformsTo",map[string]interface{}{"@id": ROCrateSpecificationVersionURI})
	roCrate.AddEntity(metadataFileDescriptor)
	return &roCrate
}

func (roCrate *ROCrate) AddEntity(e entity) {
	roCrate.Graph = append(roCrate.Graph,e.Properties)
}

func (roCrate *ROCrate) Unmarshall(folderPath string) error {
	filePath := filepath.Join(folderPath,ROCRateMetadataFileName)
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(jsonData),roCrate)
	return err
}

func (roCrate *ROCrate) Marshall(folderPath string) error {
	filePath := filepath.Join(folderPath,ROCRateMetadataFileName)
	b, err := json.Marshal(roCrate)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, b, 0644)
	return err
}

func (roCrate *ROCrate) GetMetadataFileDescriptorEntity() entity {
	var e entity
	for _, entityMap := range roCrate.Graph {
		if entityMap["conformsTo"] != nil {
			conformsTo := entityMap["conformsTo"].(map[string]interface{})
			if strings.HasPrefix(conformsTo["@id"].(string),ROCrateSpecificationPermalingPrefix) {
				e = NewEntityFromMap(entityMap)
				break
			}
		}
	}
	return e
}

func (roCrate *ROCrate) GetRootDataEntity() entity {
	var e entity
	metadataFileDescriptorEntity := roCrate.GetMetadataFileDescriptorEntity()
	aboutMap := metadataFileDescriptorEntity.GetProperty("about").(map[string]interface{})
	if aboutMap != nil {
		id := aboutMap["@id"].(string)
		for _, entityMap := range roCrate.Graph {
			e2 := NewEntityFromMap(entityMap)
			if e2.GetID() == id {
				return e2
			}
		}
	}
	return e
}


