package cmd

import (
	"fmt"

	"github.com/slackernews/cli/pkg/api"
	"github.com/slackernews/cli/pkg/formatters"
	"github.com/spf13/cobra"
)

func printLinks(cmd *cobra.Command, links []api.RenderableLink, asJSON bool) error {
	if len(links) == 0 {
		if asJSON {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "[]")
			return nil
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No links found")
		return nil
	}

	if asJSON {
		out, err := formatters.FormatJSON(links)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), string(out))
		return nil
	}

	formatters.FormatTable(cmd.OutOrStdout(), links)
	return nil
}
