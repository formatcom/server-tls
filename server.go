package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"io/ioutil"
	"os"
)

func main() {

	// REF: https://gist.github.com/xjdrew/97be3811966c8300b724deabc10e38e2
	caCertPEM, err := ioutil.ReadFile("key/ca.pem")
	checkError(err)

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caCertPEM)
	if !ok {
		panic("failed to parse root certificate")
	}

	cert, err := tls.LoadX509KeyPair("key/server.pem", "key/server.key.pem")
	checkError(err)
	config := tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    roots,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			fmt.Println("VerifyPeerCertificate", len(rawCerts))
			for i, raw := range rawCerts {
				fmt.Println("RAW")
				cert, err := x509.ParseCertificate(raw)
				checkError(err)
				subject := cert.Subject
				issuer := cert.Issuer //Alias de la empresa(issuer)
				fmt.Printf(" %d s:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s\n", i, subject.Country, subject.Province, subject.Locality, subject.Organization, subject.OrganizationalUnit, subject.CommonName)
				fmt.Printf("   i:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s\n", issuer.Country, issuer.Province, issuer.Locality, issuer.Organization, issuer.OrganizationalUnit, issuer.CommonName)
			}
			for _, chain := range verifiedChains {
				fmt.Println("VERIFIED")
				for i, cert := range chain {
					subject := cert.Subject
					issuer := cert.Issuer //Alias de la empresa(issuer)
					fmt.Printf(" %d s:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s\n", i, subject.Country, subject.Province, subject.Locality, subject.Organization, subject.OrganizationalUnit, subject.CommonName)
					fmt.Printf("   i:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s\n", issuer.Country, issuer.Province, issuer.Locality, issuer.Organization, issuer.OrganizationalUnit, issuer.CommonName)
				}
			}
			return nil
		},
	}

	service := "0.0.0.0:1200"

	listener, err := tls.Listen("tcp", service, &config)
	checkError(err)
	fmt.Println("Listening")
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("Accepted")
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var buf [512]byte
	for {
		fmt.Println("Trying to read")
		n, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(buf[0:n]))

		_, err2 := conn.Write([]byte("HTTP/1.1 200 OK\nContent-Length: 13\nContent-Type: text/html\n\nHello world!\n"))
		if err2 != nil {
			return
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
