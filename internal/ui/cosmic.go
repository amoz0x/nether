// Package ui provides cosmic-themed user interface elements for nether
package ui

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Colors and effects
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	
	// Cosmic colors
	Purple    = "\033[35m"
	Cyan      = "\033[36m"
	Blue      = "\033[34m"
	Magenta   = "\033[95m"
	White     = "\033[97m"
	Gray      = "\033[90m"
	
	// Special effects
	Blink     = "\033[5m"
	Lightning = "⚡"
	Star      = "✦"
	Comet     = "☄"
	Galaxy    = "🌌"
	Butterfly = "🦋"
)

// CosmicBanner displays a fast, static cosmic banner
func CosmicBanner() {
	if isQuietMode() {
		return
	}
	
	// Fast static banner - no animation delays
	fmt.Print(Purple + Bold)
	fmt.Print("        ✦ ･ ｡ﾟ☆: *.☽ .* :☆ﾟ. ✦\n")
	fmt.Print(Reset)
	
	// Static nether logo
	fmt.Print(Cyan + Bold)
	fmt.Print("    ███▄    █ ▓█████▄▄▄█████▓ ██░ ██ ▓█████  ██▀███  \n")
	fmt.Print(Magenta + Bold)
	fmt.Print("    ██ ▀█   █ ▓█   ▀▓  ██▒ ▓▒▓██░ ██▒▓█   ▀ ▓██ ▒ ██▒\n")
	fmt.Print(Blue + Bold)
	fmt.Print("   ▓██  ▀█ ██▒▒███  ▒ ▓██░ ▒░▒██▀▀██░▒███   ▓██ ░▄█ ▒\n")
	fmt.Print(Purple + Bold)
	fmt.Print("   ▓██▒  ▐▌██▒▒▓█  ▄░ ▓██▓ ░ ░▓█ ░██ ▒▓█  ▄ ▒██▀▀█▄  \n")
	fmt.Print(Cyan + Bold)
	fmt.Print("   ▒██░   ▓██░░▒████▒ ▒██▒ ░ ░▓█▒░██▓░▒████▒░██▓ ▒██▒\n")
	fmt.Print(White + Bold)
	fmt.Print("   ░ ▒░   ▒ ▒ ░░ ▒░ ░ ▒ ░░    ▒ ░░▒░▒░░ ▒░ ░░ ▒▓ ░▒▓░\n")
	fmt.Print(Gray)
	fmt.Print("   ░ ░░   ░ ▒░ ░ ░  ░   ░     ▒ ░▒░ ░ ░ ░  ░  ░▒ ░ ▒░\n")
	fmt.Print("      ░   ░ ░    ░    ░       ░  ░░ ░   ░     ░░   ░ \n")
	fmt.Print("            ░    ░  ░         ░  ░  ░   ░  ░   ░     \n")
	fmt.Print(Reset)
	
	// Fast subtitle
	fmt.Print(Cyan + Bold + "        ☄ Decentralized Subdomain Enumeration ☄\n" + Reset)
	fmt.Print(Magenta + Dim + "           ✦ v1.0.0-enterprise ✦\n" + Reset)
	fmt.Print("\n")
}

// ThunderResults displays results with thunder-style formatting
func ThunderResults(domain string, count int, newCount int, elapsed time.Duration) {
	if isQuietMode() {
		return
	}
	
	fmt.Print("\n" + strings.Repeat("─", 60) + "\n")
	
	// Thunder header
	fmt.Print(Purple + Bold + Lightning + " COSMIC SCAN COMPLETE " + Lightning + Reset + "\n\n")
	
	// Domain info with cosmic styling
	fmt.Printf("%s%s Target:%s %s%s%s\n", 
		Cyan, Star, Reset, 
		Bold + White, domain, Reset)
	
	// Results with lightning effects
	fmt.Printf("%s%s Total Subdomains:%s %s%s%d%s\n", 
		Blue, Lightning, Reset,
		Bold + Purple, Lightning, count, Reset)
	
	if newCount > 0 {
		fmt.Printf("%s%s New This Scan:%s %s%s%d%s\n", 
			Magenta, Comet, Reset,
			Bold + Cyan, Star, newCount, Reset)
	}
	
	// Cosmic timing with thunder
	fmt.Printf("%s%s Elapsed Time:%s %s%s%v%s\n", 
		Yellow(), Lightning, Reset,
		Bold + White, Lightning, elapsed.Round(time.Millisecond), Reset)
	
	fmt.Print("\n" + strings.Repeat("─", 60) + "\n")
}

// CosmicStatus displays status with cosmic theming
func CosmicStatus(status string, detail string, isGood bool) {
	if isQuietMode() {
		return
	}
	
	icon := Lightning
	color := Cyan
	
	if isGood {
		icon = Star
		color = Purple
	} else {
		icon = Comet
		color = Yellow()
	}
	
	fmt.Printf("%s%s %s:%s %s%s%s\n", 
		color, icon, status, Reset,
		Dim, detail, Reset)
}

// CosmicProgress shows fast progress without animation
func CosmicProgress(message string) {
	if isQuietMode() {
		return
	}
	
	fmt.Printf("%s%s %s%s\n", Cyan, Comet, message, Reset)
}

// Helper functions
func animateText(text string, delayMs int) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}
}

func addSparkles() {
	sparkles := []string{Star, "✧", "⋆", "✦"}
	sparkle := sparkles[rand.Intn(len(sparkles))]
	fmt.Print(" " + Purple + sparkle + Reset)
}

func getCosmicColor(index int) string {
	colors := []string{Purple, Cyan, Blue, Magenta, White}
	return colors[index%len(colors)] + Bold
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func isQuietMode() bool {
	// Check if quiet mode is enabled via environment or flags
	for _, arg := range os.Args {
		if arg == "-q" || arg == "--quiet" {
			return true
		}
	}
	// Also check for NETHER_QUIET environment variable
	return os.Getenv("NETHER_QUIET") == "1"
}

func Yellow() string {
	return "\033[33m"
}

// PrintCosmicList prints a beautifully formatted list of subdomains - FAST
func PrintCosmicList(subdomains []string, domain string) {
	if isQuietMode() {
		// In quiet mode, just print plain list
		for _, sub := range subdomains {
			fmt.Println(sub)
		}
		return
	}
	
	fmt.Print("\n" + Purple + Bold + Galaxy + " DISCOVERED SUBDOMAINS " + Galaxy + Reset + "\n\n")
	
	maxLen := 0
	for _, sub := range subdomains {
		if len(sub) > maxLen {
			maxLen = len(sub)
		}
	}
	
	for i, sub := range subdomains {
		// Alternate colors for visual appeal
		color := Cyan
		if i%2 == 1 {
			color = Blue
		}
		
		// Add cosmic bullet point
		bullet := Star
		if i%3 == 0 {
			bullet = Lightning
		} else if i%3 == 1 {
			bullet = Comet
		}
		
		fmt.Printf("  %s%s%s %s%-*s%s\n", 
			color, bullet, Reset,
			Dim, maxLen, sub, Reset)
	}
	
	fmt.Print("\n")
}

// CosmicHeader displays the cosmic header for commands
func CosmicHeader(command string) {
	if isQuietMode() {
		return
	}
	
	fmt.Print("\n" + Purple + Bold)
	fmt.Print("  " + strings.Repeat("✦", 20) + "\n")
	fmt.Printf("   %s%s %s %s%s\n", 
		Butterfly, Lightning, command, Lightning, Butterfly)
	fmt.Print("  " + strings.Repeat("✦", 20))
	fmt.Print(Reset + "\n\n")
}
