package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	fmt.Println("Setting up DNS for Harbor Registry")
	fmt.Println("=================================")

	// Step 1: Add DNS record for Harbor
	fmt.Println("\n1. Adding DNS record for harbor.local")
	addDNSRecord("harbor.local", "192.168.1.100")

	// Step 2: Start DNS server
	fmt.Println("\n2. Starting DNS server")
	startDNSServer()

	// Step 3: Configure system to use DNS server
	fmt.Println("\n3. Configuring system to use local DNS server")
	configureDNS()

	fmt.Println("\nSetup complete! You can now access Harbor at https://harbor.local")
}

func addDNSRecord(domain, ip string) {
	cmd := exec.Command("gotransport", "dns", "add", "--domain", domain, "--type", "A", "--value", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to add DNS record: %v\n%s", err, output)
	}
	fmt.Println(string(output))
}

func startDNSServer() {
	// Start DNS server in the background
	var cmd *exec.Cmd

	// Determine if we need sudo based on the OS and port
	if runtime.GOOS == "windows" {
		// On Windows, use non-privileged port without elevation
		cmd = exec.Command("gotransport", "dns", "serve", "--insecure")
	} else {
		// On Unix-like systems, we need sudo for privileged ports
		fmt.Println("Note: Running DNS server requires administrative privileges")
		cmd = exec.Command("sudo", "gotransport", "dns", "serve")
	}

	// Set up command to run in background
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start DNS server: %v", err)
	}

	fmt.Println("DNS server started successfully")
}

func configureDNS() {
	// OS-specific DNS configuration
	switch runtime.GOOS {
	case "linux":
		configureLinuxDNS()
	case "darwin":
		configureMacDNS()
	case "windows":
		configureWindowsDNS()
	default:
		fmt.Println("Manual DNS configuration required for your OS.")
		fmt.Println("Please add 127.0.0.1 as your primary DNS server.")
	}
}

func configureLinuxDNS() {
	fmt.Println("For Linux systems:")
	fmt.Println("1. Edit /etc/resolv.conf:")
	fmt.Println("   sudo nano /etc/resolv.conf")
	fmt.Println("2. Add the following line at the top:")
	fmt.Println("   nameserver 127.0.0.1")
	fmt.Println("3. Save and exit")
}

func configureMacDNS() {
	fmt.Println("For macOS systems:")
	fmt.Println("1. Open System Preferences > Network")
	fmt.Println("2. Select your active network connection and click 'Advanced'")
	fmt.Println("3. Go to the 'DNS' tab")
	fmt.Println("4. Click '+' and add 127.0.0.1 as the first DNS server")
	fmt.Println("5. Click 'OK' and then 'Apply'")
}

func configureWindowsDNS() {
	fmt.Println("For Windows systems:")
	fmt.Println("1. Open Control Panel > Network and Sharing Center")
	fmt.Println("2. Click on your active connection and then 'Properties'")
	fmt.Println("3. Select 'Internet Protocol Version 4' and click 'Properties'")
	fmt.Println("4. Select 'Use the following DNS server addresses'")
	fmt.Println("5. Enter 127.0.0.1 as the Preferred DNS server")
	fmt.Println("6. Click 'OK' to save changes")
}
