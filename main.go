package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cheggaaa/pb"
	"github.com/mitchellh/ioprogress"
	"github.com/stacktic/dropbox"
)

var (
	dirName   = flag.String("dir", "", "Directory name (required, in env: DROPBOX_DIR)")
	doMkdir   = flag.Bool("mkdir", false, "Create the directory")
	fileName  = flag.String("filename", "", "File name to create (required)")
	inputFile = flag.String("input", "", "Input file to upload create (stdin used by default)")
	appID     = flag.String("appID", "", "Dropbox App ID (required, in env: DROPBOX_APP_ID)")
	appSecret = flag.String("appSecret", "", "Dropbox App Secret (required, in env: DROPBOX_APP_SECRET)")
	appToken  = flag.String("token", "", "Dropbox App Token (in env: DROPBOX_APP_TOKEN)")
	chunkSize = flag.Int("chunk-size", 64, "Upload chunk size in megabytes")
)

func printUsage() {
	fmt.Println("upload-to-dropbox")
	fmt.Println()
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	setFromEnv()

	if *appID == "" || *appSecret == "" {
		fatalf("Application ID and Secret is required.")
	}

	dbox := dropbox.NewDropbox()
	dbox.SetAppInfo(*appID, *appSecret)

	if *appToken == "" {
		if err := dbox.Auth(); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		dbox.SetAccessToken(*appToken)
	}

	if *fileName == "" && *inputFile == "" {
		fmt.Println("Filename is required.")
		os.Exit(1)
		return
	} else if *fileName == "" {
		*fileName = *inputFile
	}

	if *doMkdir && *dirName != "" {
		if _, err := dbox.CreateFolder(*dirName); err != nil && strings.Contains(err.Error(), "already exists") {
			fmt.Printf("Folder %s already exists\n", *dirName)
		} else if err != nil {
			fatalf("Error creating folder %s: %s", *dirName, err)
		}
		fmt.Printf("Folder %s successfully created\n", *dirName)
	}

	fname := strings.Trim(strings.Join([]string{*dirName, *fileName}, "/"), "/")
	fmt.Printf("Uploading file %s => %s\n", *inputFile, fname)

	var input io.ReadCloser
	if *inputFile != "" {
		// _, err = dbox.UploadFile(*inputFile, fname, true, "")
		info, err := os.Stat(*inputFile)
		if err != nil {
			fatalf("%v", err)
		}

		file, err := os.Open(*inputFile)
		if err != nil {
			fatalf("Error opening file %q: %v", *inputFile, err)
		}
		defer file.Close()
		bar := pb.StartNew(int(info.Size()))
		bar.Units = pb.U_BYTES

		input = ioutil.NopCloser(&ioprogress.Reader{
			Reader: file,
			Size:   info.Size(),
			DrawFunc: func(progress int64, total int64) error {
				bar.Set(int(progress))
				if progress >= total {
					bar.FinishPrint("Done")
				}
				return nil
			},
		})
	} else {
		input = os.Stdin
	}

	_, err := dbox.UploadByChunk(input, *chunkSize*1024, fname, false, "")
	if err != nil {
		fatalf("Error uploading file %s: %v", fname, err)
	}

	fmt.Printf("File %s successfully created\n", fname)
}

func fatalf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
	os.Exit(1)
}

func setFromEnv() {
	if v := os.Getenv("DROPBOX_DIR"); *dirName == "" && v != "" {
		*dirName = v
	}
	if v := os.Getenv("DROPBOX_APP_ID"); *appID == "" && v != "" {
		*appID = v
	}
	if v := os.Getenv("DROPBOX_APP_SECRET"); *appSecret == "" && v != "" {
		*appSecret = v
	}
	if v := os.Getenv("DROPBOX_APP_TOKEN"); *appToken == "" && v != "" {
		*appToken = v
	}
}
