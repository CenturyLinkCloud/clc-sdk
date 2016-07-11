package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	clc "github.com/CenturyLinkCloud/clc-sdk"
	"github.com/CenturyLinkCloud/clc-sdk/api"
)

// Version of binary
const Version = "0.1"

func main() {
	un := flag.String("username", "", "clc username")
	pw := flag.String("password", "", "clc password")
	dc := flag.String("dc", "", "datacenter")
	flag.Parse()
	if *un == "" {
		log.Panic("missing flag -username")
	}
	if *pw == "" {
		log.Panic("missing flag -password")
	}
	if *dc == "" {
		log.Panic("missing flag -dc")
	}

	config, _ := api.NewConfig(*un, *pw)
	config.UserAgent = fmt.Sprintf("baremetal-skus: %s", Version)
	client := clc.New(config)
	if err := client.Authenticate(); err != nil {
		log.Panicf("Failed to auth: %v", err)
	}

	resp, err := client.DC.GetBareMetalCapabilities(*dc)
	if err != nil {
		log.Fatalf("Error pulling capabilities: %v", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 3, ' ', 0)
	fmt.Fprintln(w, "OS\tRate\t")
	for _, os := range resp.OperatingSystems {
		formatted := []string{
			os.Type,
			fmt.Sprintf("%.2f", os.HourlyRatePerSocket),
		}
		fmt.Fprintln(w, strings.Join(formatted, "\t"))
	}
	_ = w.Flush()

	w = tabwriter.NewWriter(os.Stdout, 0, 8, 3, ' ', 0)
	fmt.Fprintln(w, "SKU\tRate\tCPU\tMEM\tHD\t")
	for _, sku := range resp.SKUs {
		disk := 0
		for _, s := range sku.Storage {
			disk += s.CapacityInGB
		}
		cores := (sku.Processor.CoresPerSocket * sku.Processor.Sockets)

		formatted := []string{
			sku.ID,
			fmt.Sprintf("%.2f", sku.HourlyRate),
			fmt.Sprintf("%d", cores),
			fmt.Sprintf("%d", sku.Memory[0].CapacityInGB),
			fmt.Sprintf("%d", disk),
		}
		fmt.Fprintln(w, strings.Join(formatted, "\t"))
	}
	_ = w.Flush()

}
