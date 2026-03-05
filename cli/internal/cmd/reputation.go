package cmd

import (
	"flag"
	"fmt"
	"os"
)

// RunReputation handles the "reputation" subcommand.
func RunReputation(args []string, serverURL string) int {
	if len(args) < 1 {
		printReputationUsage()
		return 1
	}

	switch args[0] {
	case "show":
		return runReputationShow(args[1:])
	case "list":
		return runReputationList(args[1:])
	case "help", "-h":
		printReputationUsage()
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown reputation command: %s\n", args[0])
		printReputationUsage()
		return 1
	}
}

func printReputationUsage() {
	fmt.Fprintf(os.Stderr, `Usage: peerclaw reputation <subcommand> [options]

Subcommands:
  show <pubkey>   Show reputation score for a peer
  list            List all reputation entries
`)
}

func runReputationShow(args []string) int {
	fs := flag.NewFlagSet("reputation show", flag.ExitOnError)
	storePath := fs.String("store", "reputation.json", "Path to reputation store file")
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: peerclaw reputation show <pubkey>\n")
		return 1
	}

	pubkey := fs.Arg(0)
	fmt.Printf("Reputation for %s (store: %s):\n", pubkey, *storePath)
	fmt.Println("  Score: (use agent SDK to query)")
	fmt.Println("  Status: (use agent SDK to query)")
	return 0
}

func runReputationList(args []string) int {
	fs := flag.NewFlagSet("reputation list", flag.ExitOnError)
	storePath := fs.String("store", "reputation.json", "Path to reputation store file")
	fs.Parse(args)

	fmt.Printf("Reputation entries (store: %s):\n", *storePath)
	fmt.Println("  (use agent SDK to list entries)")
	return 0
}
