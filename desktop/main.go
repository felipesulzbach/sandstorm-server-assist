// Sandstorm Server Assist: app desktop para instalar, configurar e iniciar um servidor dedicado de
// Insurgency: Sandstorm. Veja o README e docs/.
package main

import (
	"embed"
	"log"

	"sandstorm-server-assist/desktop/internal/adapters/desktop"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if err := desktop.Run(assets); err != nil {
		log.Fatal(err)
	}
}
