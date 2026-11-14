package main

import (
	"avan/internal/config"
	"avan/wire"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Define command-line flags with descriptions
	runningMode := flag.String("running-mode", "", "Run the application in server/API or consumer mode.")
	configPath := flag.String("config", "/etc/config.yaml", "Path to the configuration file.")
	logLevel := flag.String("log-level", "info", "Set log level (debug, info, warn, error).")
	showHelp := flag.Bool("help", false, "Display help information.")

	// Customize the usage function to provide better CLI help
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExample usage:")
		fmt.Fprintln(os.Stderr, "  Run in web mode:     avan --running-mode=server")
		fmt.Fprintln(os.Stderr, "  Run in sender mode:  avan --running-mode=consumer")
		fmt.Fprintln(os.Stderr, "  Run with config:     avan --running-mode=web --config=/path/to/config.yaml")
		fmt.Fprintln(os.Stderr, "  Run with log level:  avan --running-mode=web --log-level=debug")
	}

	// Parse the command-line flags
	flag.Parse()

	// Display help if requested
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Determine the role
	role := *runningMode
	if role == "" {
		role = "both"
	}

	// Pass parsed flags into the config struct
	flagCfg := config.AppFlags{
		Role:       role,
		LogLevel:   logLevel,
		ConfigPath: configPath,
	}

	// Initialize the application
	app := wire.InitializeApp(flagCfg)

	// Logging startup messages
	if role == "both" || role == "server" {
		app.Server.Logger.Info(fmt.Sprintf("running server on port: %d", app.Config.ServerPort))
	}
	if role == "both" || role == "consumer" {
		app.Server.Logger.Info(fmt.Sprintf("running with %d consumer", app.Config.KafkaConsumerCount))
	}

	// Setup default services
	setupDefaultServices(&app)

	// Start gossip service in the background
	//go app.GossipService.PrintMembers()

	// Start the application based on mode
	if role == "web" {
		runWebMode(&app)
	}

	if role == "both" {
		go runWebMode(&app)
	}

	if role == "consumer" || role == "both" {
		app.Notifier.Start()
	}

	app.Server.Logger.Fatal(fmt.Sprintf("please specify running mode with --running-mode=%s|%s", "web", "consumer"))
}

// setupDefaultServices initializes default services
func setupDefaultServices(app *wire.App) {
	app.UserService.CreateAdminUser()
	app.WorkspaceService.CreateDefaultWorkspace()
	app.ProviderService.CreateDefaultProviders()
}

// runWebMode starts the web server
func runWebMode(app *wire.App) {
	app.Router.GenerateRoutes()
	app.Server.Start()
}
