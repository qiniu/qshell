package operations

import (
	"context"
	"fmt"
	"os"

	sbClient "github.com/qiniu/qshell/v2/iqshell/sandbox"
)

// ListInfo holds parameters for listing injection rules.
type ListInfo struct {
	Format string
}

// List lists all injection rules.
func List(info ListInfo) {
	client, err := sbClient.NewSandboxClient()
	if err != nil {
		sbClient.PrintError("%v", err)
		return
	}

	rules, err := client.ListInjectionRules(context.Background())
	if err != nil {
		sbClient.PrintError("list injection rules failed: %v", err)
		return
	}

	if info.Format == sbClient.FormatJSON {
		sbClient.PrintJSON(rules)
		return
	}

	if len(rules) == 0 {
		fmt.Println("No injection rules found")
		return
	}

	tw := sbClient.NewTable(os.Stdout)
	fmt.Fprintf(tw, "RULE ID\tNAME\tTYPE\tTARGET\tHEADERS\tCREATED AT\tUPDATED AT\n")
	for _, r := range rules {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			r.RuleID,
			r.Name,
			formatInjectionType(r.Injection),
			formatInjectionTarget(r.Injection),
			formatInjectionHeaders(r.Injection),
			sbClient.FormatTimestamp(r.CreatedAt),
			sbClient.FormatTimestamp(r.UpdatedAt),
		)
	}
	tw.Flush()
}
