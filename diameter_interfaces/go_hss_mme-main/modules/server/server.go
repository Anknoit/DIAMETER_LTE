package server

import (
	"flag"
	"log"
	"net"
	"net/http"
	"strings"

	"epiphanyHSS/modules/common"
	"epiphanyHSS/modules/go-diameter/diam"
	"epiphanyHSS/modules/go-diameter/diam/datatype"
	"epiphanyHSS/modules/go-diameter/diam/sm"
	"epiphanyHSS/modules/handlers"
)

const (
	VENDOR_3GPP = 10415
)

// ClientAuthenticator handles client authentication for incoming connections
type ClientAuthenticator struct {
	mux *sm.StateMachine
}

// NewClientAuthenticator creates a new authenticator wrapper around the provided handler
func NewClientAuthenticator(mux *sm.StateMachine) *ClientAuthenticator {
	return &ClientAuthenticator{
		mux: mux,
	}
}

// Function to clean Diameter Identity strings
func cleanDiameterIdentity(value string) string {
	// Remove "DiameterIdentity{" prefix and everything after "}"
	start := strings.Index(value, "{")
	end := strings.Index(value, "}")
	if start != -1 && end != -1 && start < end {
		return value[start+1 : end]
	}
	return value
}

// ServeDIAM implements the diam.Handler interface
func (ca *ClientAuthenticator) ServeDIAM(conn diam.Conn, m *diam.Message) {
	// Get client IP address
	remoteAddr := conn.RemoteAddr().String()
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		log.Printf("Error parsing remote address: %v", err)
		conn.Close()
		return
	}

	// Get Origin-Host and Origin-Realm from the message
	var originHost, originRealm string
	if avp, err := m.FindAVP("Origin-Host", 0); err == nil {
		originHost = cleanDiameterIdentity(avp.Data.String())
	}
	if avp, err := m.FindAVP("Origin-Realm", 0); err == nil {
		originRealm = cleanDiameterIdentity(avp.Data.String())
	}

	// Get configured clients
	clients, err := common.GetMMEConfig()
	if err != nil {
		log.Printf("Error getting MME config: %v", err)
		conn.Close()
		return
	}

	// Check if client is authorized
	authorized := false
	for _, client := range clients {
		if client.MMEHost == originHost && client.MMERealm == originRealm && strings.HasPrefix(ip, client.IP) {
			authorized = true
			break
		}
	}

	if !authorized {
		log.Printf("Unauthorized client connection attempt from %s (Host: %s, Realm: %s)",
			ip, originHost, originRealm)
		conn.Close()
		return
	}

	// If authorized, pass to the state machine
	ca.mux.ServeDIAM(conn, m)
}

// printErrors prints error reports from the state machine
func printErrors(ec <-chan *diam.ErrorReport) {
	for err := range ec {
		log.Println(err)
	}
}

// listen starts the server with the specified configuration
func listen(networkType, addr, cert, key string, handler diam.Handler) error {
	if len(cert) > 0 && len(key) > 0 {
		log.Println("Starting secure diameter server on", addr)
		return diam.ListenAndServeNetworkTLS(networkType, addr, cert, key, handler, nil)
	}
	log.Println("Starting diameter server on", addr)
	return diam.ListenAndServeNetwork(networkType, addr, handler, nil)
}

func StartServer() {

	InitializeServer()

	// Define all command line flags
	addr := flag.String("addr", "127.0.0.1:3868", "address in the form of ip:port to listen on")
	ppaddr := flag.String("pprof_addr", ":9000", "address in form of ip:port for the pprof server")
	//Setting valus for HSS Server Host and Realm
	host := flag.String("diam_host", "hss-server", "diameter identity host")
	realm := flag.String("diam_realm", "hss-server-realm", "diameter identity realm")
	certFile := flag.String("cert_file", "", "tls certificate file (optional)")
	keyFile := flag.String("key_file", "", "tls key file (optional)")
	networkType := flag.String("network_type", "tcp", "protocol type tcp/sctp/tcp4/tcp6/sctp4/sctp6")
	flag.Parse()

	settings := &sm.Settings{
		OriginHost:       datatype.DiameterIdentity(*host),
		OriginRealm:      datatype.DiameterIdentity(*realm),
		VendorID:         13,
		ProductName:      "ola-diameter-s6a",
		FirmwareRevision: 1,
	}

	// Create the state machine (mux) and set its message handlers
	mux := sm.New(settings)

	mux.Handle("ULR", handlers.HandleULR(*settings))
	mux.Handle("AIR", handlers.HandleAIR(*settings))

	// Wrap the mux with the client authenticator
	authenticator := NewClientAuthenticator(mux)

	// Print error reports
	go printErrors(mux.ErrorReports())

	if len(*ppaddr) > 0 {
		go func() { log.Fatal(http.ListenAndServe(*ppaddr, nil)) }()
	}

	// Use the authenticator instead of mux directly
	err := listen(*networkType, *addr, *certFile, *keyFile, authenticator)
	if err != nil {
		log.Fatal(err)
	}
}
