// Package winsys tem o que o app faz no Windows fora do servidor (ADR-022): abrir junto com o Windows, criar a
// regra do Firewall para o servidor (com a confirmação de administrador do Windows) e listar as placas de rede.
package winsys

import (
	"net"
	"path/filepath"
	"sort"
)

// AutostartFlag é o argumento com que o Windows abre o app ao entrar (o app abre minimizado).
const AutostartFlag = "--autostart"

// FirewallFlag é o argumento do app aberto como administrador só para criar a regra do Firewall.
const FirewallFlag = "--firewall-rule"

// RulePrefix é o começo do nome das regras do Firewall criadas pelo app.
const RulePrefix = "Sandstorm Server Assist - "

// ServerPrograms são os executáveis do servidor que precisam de regra (o lançador abre o jogo, que recebe as
// conexões).
func ServerPrograms(serverDir string) []string {
	return []string{
		filepath.Join(serverDir, "InsurgencyServer.exe"),
		filepath.Join(serverDir, "Insurgency", "Binaries", "Win64", "InsurgencyServer-Win64-Shipping.exe"),
	}
}

// RuleName é o nome da regra de um executável.
func RuleName(program string) string {
	name := filepath.Base(program)
	return RulePrefix + name[:len(name)-len(filepath.Ext(name))]
}

// Address é um endereço IPv4 desta máquina.
type Address struct {
	IP   string `json:"ip"`
	Name string `json:"name"`
}

// LocalAddresses lista os IPv4 das placas de rede ligadas (inclui as de programas de VPN), sem o loopback.
func LocalAddresses() []Address {
	var out []Address
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, ifc := range ifaces {
		if ifc.Flags&net.FlagUp == 0 || ifc.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := ifc.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			ipn, ok := a.(*net.IPNet)
			if !ok || ipn.IP.To4() == nil || ipn.IP.IsLoopback() || ipn.IP.IsLinkLocalUnicast() {
				continue
			}
			out = append(out, Address{IP: ipn.IP.String(), Name: ifc.Name})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// IsLocal diz se o IP é de uma placa desta máquina.
func IsLocal(ip string) bool {
	for _, a := range LocalAddresses() {
		if a.IP == ip {
			return true
		}
	}
	return false
}
