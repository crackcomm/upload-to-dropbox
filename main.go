package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/stacktic/dropbox"
)

var (
	dirName   = flag.String("dir", "", "Directory name (required)")
	doMkdir   = flag.Bool("mkdir", false, "Create the directory")
	fileName  = flag.String("filename", "", "File name to create (required)")
	inputFile = flag.String("input", "", "Input file to upload create (stdin used by default)")
	appId     = flag.String("appId", "", "Dropbox App ID (required, in env: DROPBOX_APP_ID)")
	appSecret = flag.String("appSecret", "", "Dropbox App Secret (required, in env: DROPBOX_APP_SECRET)")
	appToken  = flag.String("token", "", "Dropbox App Token (in env: DROPBOX_APP_TOKEN)")
	chunkSize = flag.Int("chunk-size", 64, "Upload chunk size in megabytes")
)

func main() {
	flag.Parse()
	setFromEnv()

	if *appId == "" || *appSecret == "" || *dirName == "" || *fileName == "" {
		fmt.Println("upload-to-dropbox")
		fmt.Println()
		flag.PrintDefaults()
		return
	}

	dbox := dropbox.NewDropbox()
	dbox.SetAppInfo(*appId, *appSecret)

	if *appToken == "" {
		if err := dbox.Auth(); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		dbox.SetAccessToken(*appToken)
	}

	if *doMkdir {
		if _, err := dbox.CreateFolder(*dirName); err != nil {
			fmt.Printf("Error creating folder %s: %s\n", *dirName, err)
			return
		}
		fmt.Printf("Folder %s successfully created\n", *dirName)
	}

	var err error
	fname := fmt.Sprintf("%s/%s", *dirName, *fileName)
	if *inputFile != "" {
		_, err = dbox.UploadFile(*inputFile, fname, true, "")
	} else {
		_, err = dbox.UploadByChunk(os.Stdin, *chunkSize*1024, fname, false, "")
	}

	if err != nil {
		fmt.Printf("Error uploading file %s: %s\n", fname, err)
		return
	}

	fmt.Printf("File %s successfully created\n", fname)

}

func setFromEnv() {
	if v := os.Getenv("DROPBOX_APP_ID"); *appId == "" && v != "" {
		*appId = v
	}
	if v := os.Getenv("DROPBOX_APP_SECRET"); *appSecret == "" && v != "" {
		*appSecret = v
	}
	if v := os.Getenv("DROPBOX_APP_TOKEN"); *appToken == "" && v != "" {
		*appToken = v
	}
}
