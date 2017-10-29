// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var mess = &Messenger{}

func main() {

	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatalf("Please specify a PORT to bind on")
	}
	portStr := fmt.Sprintf(":%s", port)

	log.Println("Listening on port :", port)

	certfile, ok := os.LookupEnv("CERTFILE")
	if !ok {
		log.Fatalf("Please specify a CERTFILE")
	}
	log.Println("Using certfile :", certfile)

	keyfile, ok := os.LookupEnv("KEYFILE")
	if !ok {
		log.Fatalf("Please specify a KEYFILE to bind on")
	}

	log.Println("Using keyfile ", keyfile)

	token, ok := os.LookupEnv("TOKEN")

	if !ok {
		log.Fatalf("Please specify a TOKEN for use with the Facebook Messenger API")
	}

	mess.VerifyToken = token
	mess.AccessToken = token
	log.Println("Using token", token)

	mess.MessageReceived = MessageReceived

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	// Register FB Messenger webhook with the messanger class handler routine
	http.HandleFunc("/webhook/", mess.Handler)

	/* TODO automatic renewal for Let's encrypt

	domain, ok := os.LookupEnv("DOMAIN")
	if !ok {
		log.Fatalf("Please specify a DOMAIN to use with the certManager (to organise TLS for you)")
	}

	log.Printf("Using domain %s for certmanager", domain)

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
		Cache:      autocert.DirCache("certs"),
	}


	server := &http.Server{
		Addr: portStr,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}
	server.ListenAndServeTLS("", "") //key and cert are comming from Let's Encrypt

	*/

	// func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error
	log.Fatal(http.ListenAndServeTLS(portStr, certfile, keyfile, nil))

}

// MessageReceived :Callback to handle when message received.
func MessageReceived(event Event, opts MessageOpts, msg ReceivedMessage) {
	// log.Println("event:", event, " opt:", opts, " msg:", msg)
	profile, err := mess.GetProfile(opts.Sender.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := mess.SendSimpleMessage(opts.Sender.ID, fmt.Sprintf("Hello   , %s %s, %s", profile.FirstName, profile.LastName, msg.Text))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", resp)
}
