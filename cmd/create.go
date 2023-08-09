/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"gitd/internal/storage"
	"gitd/internal/transport"
	"github.com/bnb-chain/greenfield-go-sdk/types"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	transport2 "github.com/go-git/go-git/v5/plumbing/transport"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new repo on specify remote",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var url string
		if len(args) == 1 {
			url = args[0]
		} else {
			fmt.Println("Must specify a url, example: gnfd://gnfd-testnet-fullnode-tendermint-us.bnbchain.org:443/<reponame>")
			return
		}
		endpoint, err := transport2.NewEndpoint(url)
		if err != nil {
			fmt.Printf("New endpoint error: %s", err)
			return
		}
		fmt.Println("Endpoint: ", endpoint.String())

		bucketName, _ := strings.CutPrefix(endpoint.Path, "/")
		newStorage, err := storage.NewStorage(
			os.Getenv(transport.EnvChainID),
			"https://"+endpoint.Host+":"+strconv.Itoa(endpoint.Port),
			os.Getenv(transport.EnvPrivateKey),
			bucketName,
		)
		if err != nil {
			fmt.Printf("New storage error: %s", err)
			return
		}
		_, err = newStorage.GnfdClient.HeadBucket(context.Background(), bucketName)
		if err != nil {
			if strings.Contains(err.Error(), "No such bucket") {
				providers, err := newStorage.GnfdClient.ListStorageProviders(context.Background(), true)
				if err != nil {
					fmt.Println("list storage provider error: ", err)
					return
				}
				if len(providers) > 0 {
					_, err := newStorage.GnfdClient.CreateBucket(context.Background(), bucketName, providers[0].OperatorAddress, types.CreateBucketOptions{})
					if err != nil {
						fmt.Println("create bucket error: ", err)
						return
					}
				}
			} else {
				fmt.Println("head bucket error: ", err)
				return
			}
		}

		_, err = git.Init(newStorage, memfs.New())
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
