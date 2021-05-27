package cmd

import (
	"os"
	"strconv"

	"github.com/0chain/gosdk/zboxcore/fileref"
	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/zboxcli/util"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all buckets (or allocations)",
	Args:  cobra.MaximumNArgs(1),
	Long:  `AWS s3 compatible version of the listallocations command. It lists out all the buckets (or allocations.)`,
	Run: func(cmd *cobra.Command, args []string) {

		os.Setenv("ALLOC", "8f1c82f173932fb40e241c3c63fd42e479fe4cc932f7b5f28ab97640d43c275e")
		os.Setenv("remotepath", "/")

		cmd.Flag("allocation").Value.Set(os.Getenv("ALLOC"))
		cmd.Flag("remotepath").Value.Set(os.Getenv("remotepath"))

		allocationID := cmd.Flag("allocation").Value.String()
		allocationObj, err := sdk.GetAllocation(allocationID)
		if err != nil {
			PrintError("Error fetching the allocation", err)
			os.Exit(1)
		}
		remotepath := cmd.Flag("remotepath").Value.String()
		ref, err := allocationObj.ListDir(remotepath)
		if err != nil {
			PrintError(err.Error())
			os.Exit(1)
		}
		header := []string{"Type", "Name", "Path", "Size", "Num Blocks", "Lookup Hash", "Is Encrypted", "Downloads payer"}
		data := make([][]string, len(ref.Children))
		for idx, child := range ref.Children {
			size := strconv.FormatInt(child.Size, 10)
			if child.Type == fileref.DIRECTORY {
				size = ""
			}
			isEncrypted := ""
			if child.Type == fileref.FILE {
				if len(child.EncryptionKey) > 0 {
					isEncrypted = "YES"
				} else {
					isEncrypted = "NO"
				}
			}
			data[idx] = []string{
				child.Type,
				child.Name,
				child.Path,
				size,
				strconv.FormatInt(child.NumBlocks, 10),
				child.LookupHash,
				isEncrypted,
				child.Attributes.WhoPaysForReads.String(),
			}
		}
		util.WriteTable(os.Stdout, header, []string{}, data)
		return
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.PersistentFlags().String("allocation", "", "Allocation ID")
	lsCmd.PersistentFlags().String("remotepath", "", "Remote path to list from")
	listCmd.MarkFlagRequired("allocation")
}
