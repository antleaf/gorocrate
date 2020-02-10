package main

import (
	"fmt"
	"os"
)

const InputSampleROCratePath = "/Users/paulwalk/Dropbox/Reference/_ou/NIMS/NIMS RO-Crate/ro-crate-example2"
const OutputSampleROCratePath = "/Users/paulwalk/_temp/ro-crate-sample"
const ROCrateSpecificationPermalingPrefix = "https://w3id.org/ro/crate/"
const ROCRateMetadataFileName = "ro-crate-metadata.jsonld"
const ROCrateSpecificationVersionURI = "https://w3id.org/ro/crate/1.0"


func main() {
	err := os.MkdirAll(OutputSampleROCratePath, os.ModePerm)
	roCrate := NewROCrate()
	//roCrate.AddEntity(NewRootDataEntity(
	//	"My New ROCrate",
	//	"Description of My New ROCrate",
	//	"",
	//	"my new license",
	//))

	err = roCrate.Unmarshall(InputSampleROCratePath)
	if err != nil {
		fmt.Println("error:", err)
	}

	err = roCrate.Marshall(OutputSampleROCratePath)
	if err != nil {
		fmt.Println("error:", err)
	}
	rootDataEntity := roCrate.GetRootDataEntity()
	fmt.Printf("RootDataEntityID = %s",rootDataEntity.GetID())
	fmt.Println()
	//for _,v := range roCrate.Graph {
	//	fmt.Printf("Graph entity = %v",v["@id"])
	//	fmt.Println()
	//}
}
