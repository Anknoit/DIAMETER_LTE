package main

import (
	"epiphanyHSS/modules/server"
	"fmt"
	"os"
)

// Variables that will be set at build time
var (
	version   string // The application version
	buildTime string // The time the binary was built
	goVersion string // The Go version used to build the binary
	platform  string // The target platform (OS/Arch)
	customer  string
)

func main() {
	// Check if the --version flag is passed
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		// Display the version information
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Build Time: %s\n", buildTime)

		if goVersion != "" {
			fmt.Printf("Go Version: %s\n", goVersion)
		}
		if platform != "" {
			fmt.Printf("Platform: %s\n", platform)
		}
		if customer != "" {
			fmt.Printf("Customer: %s\n", customer)
		}
		os.Exit(0) // Exit after displaying the version
	}
	// Start the TACACS server
	fmt.Println("Starting the HSS server...")
	server.StartServer()

	// Main application logic can continue here
	fmt.Println("Server is running...")
}
