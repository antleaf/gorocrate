package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ROCrateJSONConverter struct {
	Context string `json:"@context"`
	Graph []map[string]interface{} `json:"@graph"`
}

type ROCrate struct {
	Context string
	RootFolderPath string
	Entities []*entity
}

func NewROCrate(rootFolderPath string) *ROCrate {
	roCrate := ROCrate{
		Context: "https://w3id.org/ro/crate/1.0/context",
		RootFolderPath: rootFolderPath,
	}
	roCrate.Entities = make([]*entity,0)
	metadataFileDescriptor := NewEntity("CreativeWork","ro-crate-metadata.jsonld")
	metadataFileDescriptor.SetProperty("about",map[string]interface{}{"@id": "./"})
	metadataFileDescriptor.SetProperty("conformsTo",map[string]interface{}{"@id": ROCrateSpecificationVersionURI})
	roCrate.AddEntity(metadataFileDescriptor)
	rootDataEntity := NewEntity("Dataset","./")
	rootDataEntity.SetProperty("datePublished",time.Now().Format("2006-01-02"))
	roCrate.AddEntity(rootDataEntity)
	return &roCrate
}

func (roCrate *ROCrate) AddEntity(e entity) {
	roCrate.Entities = append(roCrate.Entities,&e)
}

func (roCrate *ROCrate) Unmarshall() error {
	filePath := filepath.Join(roCrate.RootFolderPath,ROCRateMetadataFileName)
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	var rOCrateJSONConverter ROCrateJSONConverter
	err = json.Unmarshal([]byte(jsonData),&rOCrateJSONConverter)
	if err != nil {
		return err
	}
	roCrate.Context = rOCrateJSONConverter.Context
	roCrate.Entities = nil
	roCrate.Entities = make([]*entity,0)
	for _,v := range rOCrateJSONConverter.Graph {
		roCrate.AddEntity(NewEntityFromMap(v))
	}
	return err
}

func (roCrate *ROCrate) Marshall() error {
	filePath := filepath.Join(roCrate.RootFolderPath,ROCRateMetadataFileName)
	rOCrateJSONConverter := ROCrateJSONConverter {
		Context: roCrate.Context,
	}
	rOCrateJSONConverter.Graph = make([]map[string]interface{},0)
	for _,v := range roCrate.Entities {
		rOCrateJSONConverter.Graph = append(rOCrateJSONConverter.Graph,v.Properties)
	}
	b, err := json.Marshal(rOCrateJSONConverter)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, b, 0644)
	return err
}

func (roCrate *ROCrate) GetMetadataFileDescriptorEntityPtr() *entity {
	var e *entity
	for _, entityPtr := range roCrate.Entities {
		if entityPtr.GetID() == ROCRateMetadataFileName {
			e = entityPtr
			break
		}
	}
	return e
}

//func (roCrate *ROCrate) GetMetadataFileDescriptorEntityPtr() *entity {
//	var e *entity
//	for _, entityPtr := range roCrate.Entities {
//		if entityPtr.GetProperty("conformsTo") != nil {
//			conformsTo := entityPtr.GetProperty("conformsTo").(map[string]interface{})
//			if strings.HasPrefix(conformsTo["@id"].(string),ROCrateSpecificationPermalingPrefix) {
//				e = entityPtr
//				break
//			}
//		}
//	}
//	return e
//}

func (roCrate *ROCrate) GetRootDataEntity() *entity {
	var e *entity
	metadataFileDescriptorEntity := roCrate.GetMetadataFileDescriptorEntityPtr()
	aboutMap := metadataFileDescriptorEntity.GetProperty("about").(map[string]interface{})
	if aboutMap != nil {
		id := aboutMap["@id"].(string)
		for _, entityPtr := range roCrate.Entities {
			if entityPtr.GetID() == id {
				e = entityPtr
			}
		}
	}
	return e
}

func (roCrate *ROCrate) PopulateROCrateFromFilesystem() {
	rootDataEntity := roCrate.GetRootDataEntity()
	if rootDataEntity.GetProperty("hasPart") == nil {
		rootDataEntity.SetProperty("hasPart",make([]map[string]string,0))
	}
	hasPartArray := rootDataEntity.GetProperty("hasPart").([]map[string]string)
	err := filepath.Walk(roCrate.RootFolderPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			pathAsId := strings.TrimPrefix(path,roCrate.RootFolderPath)
			if strings.HasPrefix(pathAsId,string(os.PathSeparator)) {
				pathAsId = "." + pathAsId
			}
			if info.IsDir() {
				pathAsId += string(os.PathSeparator)
				newDirEntity := NewEntity("Dataset",pathAsId)
				roCrate.AddEntity(newDirEntity)
			} else { // is a file
				newFileEntity := NewEntity("File",pathAsId)
				newFileEntity.SetProperty("contentSize",info.Size())
				f, err := os.Open(path)
				if err != nil {
				//	do something about it
				}
				defer f.Close()
				contentType, err := GetFileContentType(f)
				newFileEntity.SetProperty("encodingFormat",contentType)
				roCrate.AddEntity(newFileEntity)
			}
			hasPartArray = append(hasPartArray,map[string]string{"@id": pathAsId})
			return nil
		})
	rootDataEntity.SetProperty("hasPart",hasPartArray)
	if err != nil {
		log.Println(err)
	}
}

func GetFileContentType(out *os.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}
	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}