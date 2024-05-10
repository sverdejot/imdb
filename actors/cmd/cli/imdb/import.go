package imdb

import (
	"github.com/spf13/cobra"
	"github.com/sverdejot/imdb/actors/pkg/imdb"
)

var connectionString string

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import imdb actors dataset",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		imdb.Import(args[0], connectionString)
	},
}

func init() {
	importCmd.Flags().StringVarP(
		&connectionString,
		"connection-string",
		"c",
		"",
		"database's connection string where the datase will be imported")

	rootCmd.AddCommand(importCmd)
}
