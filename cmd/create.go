package cmd

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
)

func init() {
	var Path string
	createCmd.Flags().StringVarP(&Path, "path", "p", "", "cluster path")
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "creates",
	Long:  `creates`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating cluster " + strings.Join(args, " "))
		if err := createCluster(args); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func createCluster(args []string) error {
	clusterName := args[0]
	cfgPath := filepath.Join(os.Getenv("HOME"), ".ckube", "clusters", clusterName)
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		err := os.MkdirAll(cfgPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create %q\n", cfgPath)
			return err
		}
		fmt.Printf("Created cluster in %q\n", cfgPath)
	} else {
		fmt.Printf("Cluster %s already exists in %q\n", clusterName, cfgPath)
	}
	if err := checkAndFetchImages(filepath.Join(os.Getenv("HOME"), ".ckube")); err != nil {
		return err
	}

	return nil
}

func checkAndFetchImages(path string) error {
	imagePath := filepath.Join(path, "images")
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		if err := os.MkdirAll(imagePath, os.ModePerm); err != nil {
			return err
		}
	}
	images := map[string]string{"ckube-cmdline": "1q7bnMuRhakhYmR12TnUSQRY2rKqunKID",
		"ckube-initrd.img": "1xOHfGMeVK3vktMwBjR7xAGvKdtmYuULk",
		"ckube-kernel":     "13ndDo3dAHRTBsq5W7v7uWAQ4LGDY7b8a"}
	//images := map[string]string{"ckube-images.tgz": "1OuRNedVzN1fjShOSpS0FfayrTwc2yanc"}
	for image, id := range images {
		imageFile := filepath.Join(imagePath, image)
		if _, err := os.Stat(imageFile); os.IsNotExist(err) {
			fmt.Printf("Downloading %s to %s\n", image, imageFile)
			output, err := os.Create(imageFile)
			defer output.Close()

			ctx := context.Background()
			driveService, err := drive.NewService(ctx, option.WithAPIKey("AIzaSyBhVLu7Kk8HGmhzSO7Xn6GV2_WBCT9nSHc"))
			if err != nil {
				return err
			}
			getCall := driveService.Files.Get(id)
			response, err := getCall.Download()
			if err != nil {
				fmt.Println("Error while downloading")
				return err
			}
			/*
				response, err := http.Get(url)
				if err != nil {
					fmt.Println("Error while downloading", url, "-", err)
					return err
				}
			*/

			defer response.Body.Close()
			n, err := io.Copy(output, response.Body)
			if err != nil {
				return err
			}
			fmt.Println(n, "bytes downloaded")
		}
	}
	return nil
}

// ExtractTarGz extracts a tgz
func ExtractTarGz(gzipStream io.Reader) {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		log.Fatal("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				log.Fatalf("ExtractTarGz: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
			outFile.Close()

		default:
			log.Fatalf(
				"ExtractTarGz: uknown type: %b in %s",
				header.Typeflag,
				header.Name)
		}

	}
}
