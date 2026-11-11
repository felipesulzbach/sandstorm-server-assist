package netdiag

import (
	"bufio"
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"sandstorm-server-assist/desktop/internal/domain/netcheck"
)

// O roteador é consultado por UPnP (IGD) só para leitura: IP externo e redirecionamentos das portas. O app não
// cria nem apaga redirecionamentos. Só conversa com endereços da rede local.

const (
	ssdpAddr      = "239.255.255.250:1900"
	maxUPnPBody   = 256 << 10
	upnpNoSuchMap = 714 // NoSuchEntryInArray
)

var ssdpTargets = []string{
	"urn:schemas-upnp-org:device:InternetGatewayDevice:1",
	"urn:schemas-upnp-org:device:InternetGatewayDevice:2",
}

// Serviços que sabem dos redirecionamentos, na ordem de preferência.
var wanServices = []string{
	"urn:schemas-upnp-org:service:WANIPConnection:2",
	"urn:schemas-upnp-org:service:WANIPConnection:1",
	"urn:schemas-upnp-org:service:WANPPPConnection:1",
}

// lanClient fala só com o roteador: sem proxy, sem redirecionamento, com prazo curto.
var lanClient = &http.Client{
	Timeout: 4 * time.Second,
	Transport: &http.Transport{
		Proxy:       nil,
		DialContext: (&net.Dialer{Timeout: 3 * time.Second}).DialContext,
	},
	CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
}

// Gateway é o roteador que respondeu por UPnP.
type Gateway struct {
	control *url.URL
	service string
}

// Discover procura o roteador na rede local. Devolve nil se nenhum responder (UPnP desligado ou sem roteador).
func Discover(ctx context.Context) *Gateway {
	for _, loc := range ssdpSearch(ctx) {
		if g, err := describe(ctx, loc); err == nil {
			return g
		}
	}
	return nil
}

// ssdpSearch manda o M-SEARCH e junta os endereços de descrição que vierem do próprio aparelho que respondeu.
func ssdpSearch(ctx context.Context) []*url.URL {
	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return nil
	}
	defer func() { _ = conn.Close() }()
	dst, err := net.ResolveUDPAddr("udp4", ssdpAddr)
	if err != nil {
		return nil
	}
	for _, st := range ssdpTargets {
		msg := "M-SEARCH * HTTP/1.1\r\nHOST: " + ssdpAddr + "\r\nMAN: \"ssdp:discover\"\r\nMX: 2\r\nST: " + st + "\r\n\r\n"
		_, _ = conn.WriteToUDP([]byte(msg), dst)
	}
	deadline := time.Now().Add(2500 * time.Millisecond)
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}
	_ = conn.SetReadDeadline(deadline)

	var out []*url.URL
	seen := map[string]bool{}
	buf := make([]byte, 2048)
	for len(out) < 8 {
		n, src, err := conn.ReadFromUDP(buf)
		if err != nil {
			break
		}
		loc := ssdpLocation(buf[:n])
		u, err := lanURL(loc, src.IP)
		if err != nil || seen[u.String()] {
			continue
		}
		seen[u.String()] = true
		out = append(out, u)
	}
	return out
}

// ssdpLocation lê o cabeçalho LOCATION da resposta.
func ssdpLocation(b []byte) string {
	r := bufio.NewReader(bytes.NewReader(b))
	for {
		line, err := r.ReadString('\n')
		if k, v, ok := strings.Cut(line, ":"); ok && strings.EqualFold(strings.TrimSpace(k), "location") {
			return strings.TrimSpace(v)
		}
		if err != nil {
			return ""
		}
	}
}

// lanURL aceita só http://<IPv4 da rede local>:<porta>/..., e só se o IP for o do aparelho que respondeu.
func lanURL(raw string, from net.IP) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "http" || u.User != nil {
		return nil, errors.New("upnp: invalid address")
	}
	ip := net.ParseIP(u.Hostname())
	if ip == nil || ip.To4() == nil || !ip.IsPrivate() || (from != nil && !ip.Equal(from)) {
		return nil, errors.New("upnp: address outside the local network")
	}
	if p := u.Port(); p != "" {
		if n, err := strconv.Atoi(p); err != nil || n < 1 || n > 65535 {
			return nil, errors.New("upnp: invalid port")
		}
	}
	return u, nil
}

type upnpService struct {
	ServiceType string `xml:"serviceType"`
	ControlURL  string `xml:"controlURL"`
}

type upnpDevice struct {
	Services []upnpService `xml:"serviceList>service"`
	Devices  []upnpDevice  `xml:"deviceList>device"`
}

type upnpRoot struct {
	URLBase string     `xml:"URLBase"`
	Device  upnpDevice `xml:"device"`
}

func (d upnpDevice) find(service string) string {
	for _, s := range d.Services {
		if strings.TrimSpace(s.ServiceType) == service {
			return strings.TrimSpace(s.ControlURL)
		}
	}
	for _, c := range d.Devices {
		if u := c.find(service); u != "" {
			return u
		}
	}
	return ""
}

// describe lê a descrição do aparelho e acha o serviço de conexão WAN.
func describe(ctx context.Context, loc *url.URL) (*Gateway, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, loc.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := lanClient.Do(req) // #nosec G107 G704 -- endereço conferido por lanURL (IPv4 privado do próprio aparelho)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upnp: HTTP %d", resp.StatusCode)
	}
	var root upnpRoot
	if err := xml.NewDecoder(io.LimitReader(resp.Body, maxUPnPBody)).Decode(&root); err != nil {
		return nil, err
	}
	base := loc
	if b := strings.TrimSpace(root.URLBase); b != "" {
		if bu, err := url.Parse(b); err == nil && bu.Hostname() == loc.Hostname() {
			base = bu
		}
	}
	for _, svc := range wanServices {
		ctrl := root.Device.find(svc)
		if ctrl == "" {
			continue
		}
		ref, err := url.Parse(ctrl)
		if err != nil {
			continue
		}
		u := base.ResolveReference(ref)
		if u.Scheme != "http" || u.Hostname() != loc.Hostname() {
			continue // controle em outro aparelho: recusa
		}
		return &Gateway{control: u, service: svc}, nil
	}
	return nil, errors.New("upnp: no WAN connection service")
}

// call faz uma ação SOAP. Devolve os valores da resposta (pelo nome do elemento) ou o código de erro do UPnP.
func (g *Gateway) call(ctx context.Context, action string, args [][2]string) (map[string]string, int, error) {
	var body bytes.Buffer
	body.WriteString(`<?xml version="1.0"?><s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" ` +
		`s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body><u:` + action + ` xmlns:u="` + g.service + `">`)
	for _, a := range args {
		body.WriteString("<" + a[0] + ">")
		_ = xml.EscapeText(&body, []byte(a[1]))
		body.WriteString("</" + a[0] + ">")
	}
	body.WriteString("</u:" + action + "></s:Body></s:Envelope>")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.control.String(), &body)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", `text/xml; charset="utf-8"`)
	req.Header.Set("SOAPAction", `"`+g.service+"#"+action+`"`)
	resp, err := lanClient.Do(req) // #nosec G107 G704 -- controle no mesmo IPv4 privado da descrição (describe)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = resp.Body.Close() }()
	values, err := soapValues(io.LimitReader(resp.Body, maxUPnPBody))
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode != http.StatusOK {
		code, _ := strconv.Atoi(values["errorCode"])
		return nil, code, fmt.Errorf("upnp: %s failed (HTTP %d, code %d)", action, resp.StatusCode, code)
	}
	return values, 0, nil
}

// soapValues junta o texto de cada elemento folha pelo nome local.
func soapValues(r io.Reader) (map[string]string, error) {
	dec := xml.NewDecoder(r)
	out := map[string]string{}
	var name string
	var text strings.Builder
	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			return out, nil
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			name = t.Name.Local
			text.Reset()
		case xml.CharData:
			if text.Len() < 1024 {
				text.Write(t)
			}
		case xml.EndElement:
			if t.Name.Local == name {
				out[name] = strings.TrimSpace(text.String())
			}
			name = ""
		}
	}
}

// ExternalIP é o IP do lado da internet, segundo o roteador.
func (g *Gateway) ExternalIP(ctx context.Context) string {
	v, _, err := g.call(ctx, "GetExternalIPAddress", nil)
	if err != nil {
		return ""
	}
	if ip := net.ParseIP(v["NewExternalIPAddress"]); ip != nil && ip.To4() != nil {
		return ip.String()
	}
	return ""
}

// Mapping lê o redirecionamento de uma porta (protocolo "UDP" ou "TCP").
func (g *Gateway) Mapping(ctx context.Context, port int, proto string) netcheck.Mapping {
	if port < 1 || port > 65535 {
		return netcheck.Mapping{}
	}
	v, code, err := g.call(ctx, "GetSpecificPortMappingEntry", [][2]string{
		{"NewRemoteHost", ""}, {"NewExternalPort", strconv.Itoa(port)}, {"NewProtocol", proto},
	})
	if err != nil {
		if code == upnpNoSuchMap {
			return netcheck.Mapping{Known: true}
		}
		return netcheck.Mapping{}
	}
	internal, _ := strconv.Atoi(v["NewInternalPort"])
	client := ""
	if ip := net.ParseIP(v["NewInternalClient"]); ip != nil {
		client = ip.String()
	}
	enabled := v["NewEnabled"] == "1" || strings.EqualFold(v["NewEnabled"], "true")
	return netcheck.Mapping{Known: true, Exists: true, Enabled: enabled, Client: client, InternalPort: internal}
}
