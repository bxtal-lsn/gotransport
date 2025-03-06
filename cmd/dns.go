package cmd

import (
	"fmt"
	"os"
	"strconv" // Add this

	"github.com/AlecAivazis/survey/v2" // Add this
	"github.com/bxtal-lsn/gotransport/internal/dnsrecords"
	"github.com/bxtal-lsn/gotransport/pkg/dns"
	"github.com/fatih/color"            // Add this
	"github.com/olekukonko/tablewriter" // Add this
	"github.com/spf13/cobra"
)

var (
	// DNS server flags
	dnsAddress     string
	dnsPort        int
	dnsStoragePath string

	// DNS record flags
	dnsDomain   string
	dnsType     string
	dnsValue    string
	dnsTTL      uint32
	dnsInsecure bool
)

// Keep your existing init function unchanged
func init() {
	// Main DNS command
	dnsCmd := &cobra.Command{
		Use:   "dns",
		Short: "DNS server commands",
		Long:  `Commands for managing the DNS server and DNS records`,
	}

	// DNS serve command
	dnsServeCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start DNS server",
		Long:  `Start a DNS server that responds to requests for configured domains`,
		RunE:  runDNSServe,
	}

	// DNS add command
	dnsAddCmd := &cobra.Command{
		Use:   "add",
		Short: "Add DNS record",
		Long:  `Add a new DNS record to the storage`,
		RunE:  runDNSAdd,
	}

	// DNS list command
	dnsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List DNS records",
		Long:  `List all DNS records in the storage`,
		RunE:  runDNSList,
	}

	// DNS remove command
	dnsRemoveCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove DNS record",
		Long:  `Remove a DNS record from the storage`,
		RunE:  runDNSRemove,
	}

	// Add your existing flags here...
	// Add flags to DNS serve command
	dnsServeCmd.Flags().StringVarP(&dnsAddress, "address", "a", "0.0.0.0", "address to listen on")
	dnsServeCmd.Flags().IntVarP(&dnsPort, "port", "p", 53, "port to listen on")
	dnsServeCmd.Flags().StringVarP(&dnsStoragePath, "storage", "s", getDefaultStoragePath(), "path to DNS records storage file")
	dnsServeCmd.Flags().BoolVar(&dnsInsecure, "insecure", false, "run server on non-privileged port (5353) without root")

	// Add flags to DNS add command
	dnsAddCmd.Flags().StringVarP(&dnsDomain, "domain", "d", "", "domain name")
	dnsAddCmd.Flags().StringVarP(&dnsType, "type", "t", "A", "record type (A or CNAME)")
	dnsAddCmd.Flags().StringVar(&dnsValue, "value", "", "record value (IP address for A, domain for CNAME)")
	dnsAddCmd.Flags().Uint32Var(&dnsTTL, "ttl", 3600, "record time to live in seconds")
	dnsAddCmd.Flags().StringVarP(&dnsStoragePath, "storage", "s", getDefaultStoragePath(), "path to DNS records storage file")
	dnsAddCmd.MarkFlagRequired("domain")
	dnsAddCmd.MarkFlagRequired("value")

	// Add flags to DNS list command
	dnsListCmd.Flags().StringVarP(&dnsStoragePath, "storage", "s", getDefaultStoragePath(), "path to DNS records storage file")

	// Add flags to DNS remove command
	dnsRemoveCmd.Flags().StringVarP(&dnsDomain, "domain", "d", "", "domain name")
	dnsRemoveCmd.Flags().StringVarP(&dnsStoragePath, "storage", "s", getDefaultStoragePath(), "path to DNS records storage file")
	dnsRemoveCmd.MarkFlagRequired("domain")

	// Add commands to DNS command
	dnsCmd.AddCommand(dnsServeCmd)
	dnsCmd.AddCommand(dnsAddCmd)
	dnsCmd.AddCommand(dnsListCmd)
	dnsCmd.AddCommand(dnsRemoveCmd)

	// Add DNS command to root command
	rootCmd.AddCommand(dnsCmd)
}

// Keep your existing getDefaultStoragePath function
func getDefaultStoragePath() string {
	return "dns.json"
}

// Keep your existing runDNSServe function, or update it if needed
func runDNSServe(cmd *cobra.Command, args []string) error {
	// If insecure flag is set, use non-privileged port
	if dnsInsecure && dnsPort == 53 {
		dnsPort = 5353
		fmt.Println("Running in insecure mode on port 5353")
	}

	// Check if running as root when using privileged port
	if dnsPort < 1024 && os.Getuid() != 0 {
		return fmt.Errorf("must run as root to bind to port %d. Try using --insecure flag or sudo", dnsPort)
	}

	// Create DNS server
	server, err := dns.NewServer(dnsAddress, dnsPort, dnsStoragePath)
	if err != nil {
		return err
	}

	// Start server and handle signals
	return server.StartWithSignalHandling()
}

// Function to handle DNS remove command
func runDNSRemove(cmd *cobra.Command, args []string) error {
	// Create storage
	storage, err := dnsrecords.NewStorage(dnsStoragePath)
	if err != nil {
		return err
	}

	// Remove record
	if err := storage.Remove(dnsDomain); err != nil {
		return err
	}

	fmt.Printf("Removed DNS record: %s\n", dnsDomain)
	return nil
}

// Update this function
func runDNSList(cmd *cobra.Command, args []string) error {
	// Create storage
	storage, err := dnsrecords.NewStorage(dnsStoragePath)
	if err != nil {
		return err
	}

	// List records
	records := storage.List()

	if len(records) == 0 {
		color.Yellow("No DNS records found")
		return nil
	}

	// Display records
	color.Cyan("DNS Records:")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Domain", "Type", "Value", "TTL"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
	)
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
		tablewriter.Colors{tablewriter.FgYellowColor},
		tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgHiWhiteColor},
	)

	for _, record := range records {
		table.Append([]string{
			record.Domain,
			string(record.Type),
			record.Value,
			fmt.Sprintf("%d", record.TTL),
		})
	}

	table.Render()
	return nil
}

// Update this function
func runDNSAdd(cmd *cobra.Command, args []string) error {
	// If domain or value not provided, prompt the user
	if dnsDomain == "" {
		prompt := &survey.Input{
			Message: "Domain name:",
		}
		if err := survey.AskOne(prompt, &dnsDomain); err != nil {
			return err
		}
	}

	// Select record type
	if dnsType == "" {
		prompt := &survey.Select{
			Message: "Record type:",
			Options: []string{"A", "CNAME"},
			Default: "A",
		}
		if err := survey.AskOne(prompt, &dnsType); err != nil {
			return err
		}
	}

	// Prompt for value
	if dnsValue == "" {
		prompt := &survey.Input{
			Message: fmt.Sprintf("Value for %s record:", dnsType),
			Help:    "IP address for A record, domain name for CNAME",
		}
		if err := survey.AskOne(prompt, &dnsValue); err != nil {
			return err
		}
	}

	// Prompt for TTL
	if !cmd.Flags().Changed("ttl") {
		ttlStr := "3600"
		prompt := &survey.Input{
			Message: "TTL (seconds):",
			Default: ttlStr,
		}
		if err := survey.AskOne(prompt, &ttlStr); err != nil {
			return err
		}
		parsedTTL, err := strconv.ParseUint(ttlStr, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid TTL: %w", err)
		}
		dnsTTL = uint32(parsedTTL)
	}

	// Confirm before adding
	var confirm bool
	confirmPrompt := &survey.Confirm{
		Message: fmt.Sprintf("Add DNS record %s -> %s (%s)?", dnsDomain, dnsValue, dnsType),
		Default: true,
	}
	if err := survey.AskOne(confirmPrompt, &confirm); err != nil {
		return err
	}
	if !confirm {
		printWarning("Operation cancelled")
		return nil
	}

	// Create storage and add record
	storage, err := dnsrecords.NewStorage(dnsStoragePath)
	if err != nil {
		return err
	}

	recordType := dnsrecords.RecordType(dnsType)
	if err := storage.Add(dnsDomain, recordType, dnsValue, dnsTTL); err != nil {
		return err
	}

	printSuccess("Added DNS record: %s -> %s (%s)", dnsDomain, dnsValue, dnsType)
	return nil
}

// Add the runDNSRemove function if needed

