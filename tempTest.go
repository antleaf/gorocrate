package main

import (
	"fmt"
	"os"
)

func readAndWriteExistingROCRate(inputSampleROCratePath,outputSampleROCratePath string) {
	err := os.MkdirAll(outputSampleROCratePath, os.ModePerm)
	roCrate := NewROCrate(inputSampleROCratePath)
	err = roCrate.Unmarshall()
	if err != nil {
		fmt.Println("error:", err)
	}
	//roCrate.GetRootDataEntity().SetProperty("datePublished","1968-04-27")
	roCrate.RootFolderPath = outputSampleROCratePath
	err = roCrate.Marshall()
	if err != nil {
		fmt.Println("error:", err)
	}
}

func createROCrateForExistingFolder(rOCrateRootPath string) {
	roCrate := NewROCrate(rOCrateRootPath)
	roCrate.PopulateROCrateFromFilesystem()
	err := roCrate.Marshall()
	if err != nil {
		fmt.Println("error:", err)
	}
}


func main() {
	//readAndWriteExistingROCRate("/Users/paulwalk/Dropbox/Reference/_ou/NIMS/NIMS RO-Crate/third-party-example-ro-crates/ro-crate-example2","/Users/paulwalk/_temp/ro-crate-sample")
	createROCrateForExistingFolder("/Users/paulwalk/Dropbox/Reference/_ou/NIMS/NIMS RO-Crate/nims_data/sample-mdr-x-data")
}
