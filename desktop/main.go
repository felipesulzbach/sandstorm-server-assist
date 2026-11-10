// Sandstorm Server Assist: app desktop para instalar, configurar e iniciar um servidor dedicado de
// Insurgency: Sandstorm. Veja o README e docs/.
package main

import (
	"embed"
	"log"
	"os"

	"sandstorm-server-assist/desktop/internal/adapters/desktop"
	"sandstorm-server-assist/desktop/internal/adapters/winsys"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Cópia aberta como administrador só para criar a regra do Firewall (ADR-022).
	if len(os.Args) == 2 && os.Args[1] == winsys.FirewallFlag {
		os.Exit(desktop.FirewallHelper())
	}
	if err := desktop.Run(assets); err != nil {
		log.Fatal(err)
	}
}
