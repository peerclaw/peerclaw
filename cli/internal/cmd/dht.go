package cmd

import (
	"flag"
	"fmt"
	"os"
)

// RunDHT handles the "dht" subcommand.
func RunDHT(args []string, serverURL string) int {
	if len(args) < 1 {
		printDHTUsage()
		return 1
	}

	switch args[0] {
	case "bootstrap":
		return runDHTBootstrap(args[1:], serverURL)
	case "lookup":
		return runDHTLookup(args[1:], serverURL)
	case "help", "-h":
		printDHTUsage()
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown dht command: %s\n", args[0])
		printDHTUsage()
		return 1
	}
}

func printDHTUsage() {
	fmt.Fprintf(os.Stderr, `Usage: peerclaw dht <subcommand> [options]

Subcommands:
  bootstrap   Bootstrap the DHT with seed nodes
  lookup      Look up an agent by public key in the DHT
`)
}

func runDHTBootstrap(args []string, serverURL string) int {
	fs := flag.NewFlagSet("dht bootstrap", flag.ExitOnError)
	seeds := fs.String("seeds", "", "Comma-separated seed node addresses")
	relays := fs.String("relays", "", "Comma-separated Nostr relay URLs")
	fs.Parse(args)

	fmt.Fprintf(os.Stderr, "DHT bootstrap initiated\n")
	if *seeds != "" {
		fmt.Fprintf(os.Stderr, "Seeds: %s\n", *seeds)
	}
	if *relays != "" {
		fmt.Fprintf(os.Stderr, "Relays: %s\n", *relays)
	}
	fmt.Println("DHT bootstrap complete (stub - use agent SDK for full DHT)")
	return 0
}

func runDHTLookup(args []string, serverURL string) int {
	fs := flag.NewFlagSet("dht lookup", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: peerclaw dht lookup <pubkey>\n")
		return 1
	}

	pubkey := fs.Arg(0)
	fmt.Fprintf(os.Stderr, "Looking up agent in DHT: %s\n", pubkey)
	fmt.Println("DHT lookup complete (stub - use agent SDK for full DHT)")
	return 0
}
