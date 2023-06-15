// Copyright The Sigstore Authors.
//
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
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/open-policy-agent/frameworks/constraint/pkg/externaldata"

	"path/filepath"

	"github.com/sigstore/cosign/cmd/cosign/cli/options"
	"github.com/sigstore/cosign/pkg/cosign"
	rekor "github.com/sigstore/rekor/pkg/client"
	"github.com/sigstore/rekor/pkg/generated/client"
	"github.com/sigstore/sigstore/pkg/signature"
)

const (
	apiVersion      = "externaldata.gatekeeper.sh/v1beta1"
	defaultRekorURL = "https://rekor.sigstore.dev"
)

var (
	rekorClient         *client.Rekor
	fulcioRoots         *x509.CertPool
	fulcioIntermediates *x509.CertPool
	//rekorURL            = flag.String("rekor-url", defaultRekorURL, "specify Rekor URL")
	//
	rekorURL        = os.Getenv("SIGSTORE_URL") //"harbor.cudc.aws.comunidad.madrid:443"
	harbosecretpath = "/etc/registry-secret/"
	pubkeypath      = "/cosign-gatekeeper-provider/cosign.pub"
	servercrt       = "/etc/cabundles/server.crt"
	serverkey       = "/etc/cabundles/server.key"
	additionalCerts = "/etc/additionalCerts/"
)

func main() {
	fmt.Println("starting server...")
	http.HandleFunc("/validate", validate)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Crea un cliente HTTP con el transport personalizado
	client := &http.Client{
		Transport: transport,
	}

	// Configura el cliente HTTP global para que use el transport personalizado
	http.DefaultClient = client
	rc, err := rekor.GetRekorClient(rekorURL, rekor.WithUserAgent(options.UserAgent()))
	if err != nil {
		panic(fmt.Sprintf("creating Rekor client: %v", err))
	}
	rekorClient = rc

	fulcioRoots, err = x509.SystemCertPool()
	fulcioIntermediates, err = x509.SystemCertPool()

	files := listFilesWithExtension(additionalCerts)
	for _, file := range files {
		fmt.Println(file)
		certBytes, err := ioutil.ReadFile(file)
		if err != nil {
			panic(fmt.Sprintf("error reading file: %v", err))
		}
		if !fulcioRoots.AppendCertsFromPEM(certBytes) {
			panic("no se pudo agregar el certificado a fulcioRoots")
		}
	}
	// Leer el archivo cosign.pub
	/*
		// Agrega el certificado a fulcioRoots
		if !fulcioRoots.AppendCertsFromPEM(certBytes) {
			panic("no se pudo agregar el certificado a fulcioRoots")
		} */

	srv := &http.Server{
		Addr:              ":8090",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			// Aquí debes especificar la ruta a tu certificado y clave privada
			//Certificates: []tls.Certificate{cert},
		},
	}

	// En lugar de usar ListenAndServe, utiliza ListenAndServeTLS para habilitar HTTPS
	if err := srv.ListenAndServeTLS(servercrt, serverkey); err != nil {
		panic(err)
	}
}

func validate(w http.ResponseWriter, req *http.Request) {
	// only accept POST requests
	if req.Method != http.MethodPost {
		sendResponse(nil, "only POST is allowed", w)
		return
	}

	/* 	body, err := io.ReadAll(req.Body)
	   	if err != nil {
	   		http.Error(w, err.Error(), http.StatusInternalServerError)
	   		return
	   	}
	   	defer req.Body.Close()

	   	// Imprimir el contenido del cuerpo
	   	fmt.Println(string(body))

	   	for key, value := range req.Header {
	   		fmt.Printf("%s: %s\n", key, value)
	   	}
	*/
	// read request body
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		sendResponse(nil, fmt.Sprintf("unable to read request body: %v", err), w)
		return
	}

	// parse request body
	var providerRequest externaldata.ProviderRequest
	err = json.Unmarshal(requestBody, &providerRequest)
	if err != nil {
		sendResponse(nil, fmt.Sprintf("unable to unmarshal request body: %v", err), w)
		return
	}

	results := make([]externaldata.Item, 0)

	ctx := req.Context()

	regUsernameByte, err := ioutil.ReadFile(harbosecretpath + "username")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	regPasswordByte, err := ioutil.ReadFile(harbosecretpath + "password")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	regUsername := string(regUsernameByte)
	regPassword := string(regPasswordByte)

	kc := MyKeychain{
		registry: rekorURL,
		username: regUsername,
		password: regPassword,
	}

	ro := options.RegistryOptions{
		AllowInsecure: true,
		Keychain:      kc,
	}

	co, err := ro.ClientOpts(ctx)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if err != nil {
		sendResponse(nil, fmt.Sprintf("ERROR: %v", err), w)
		return
	}

	// iterate over all keys
	for _, key := range providerRequest.Request.Keys {
		fmt.Println("verify signature for:", key)
		ref, err := name.ParseReference(key)
		if err != nil {
			sendResponse(nil, fmt.Sprintf("ERROR (ParseReference(%q)): %v", key, err), w)
			return
		}

		verifier, err := signature.LoadVerifierFromPEMFile(pubkeypath, crypto.SHA256)

		fmt.Println("cosign  verifyImageSignatures for:", ref)
		checkedSignatures, bundleVerified, err := cosign.VerifyImageSignatures(ctx, ref, &cosign.CheckOpts{
			RekorClient:        rekorClient,
			RegistryClientOpts: co,
			RootCerts:          fulcioRoots,
			IntermediateCerts:  fulcioIntermediates,
			ClaimVerifier:      cosign.SimpleClaimVerifier,
			SigVerifier:        verifier,
		})

		if err != nil {
			fmt.Println(err)
			sendResponse(nil, fmt.Sprintf("VerifyImageSignatures: %v", err), w)
			return
		}

		if bundleVerified {
			fmt.Println("signature verified for:", key)
			fmt.Printf("%d number of valid signatures found for %s, found signatures: %v\n", len(checkedSignatures), key, checkedSignatures)
			results = append(results, externaldata.Item{
				Key:   key,
				Value: key + "_valid",
			})
		} else {
			fmt.Printf("no valid signatures found for: %s\n", key)
			results = append(results, externaldata.Item{
				Key:   key,
				Error: key + "_invalid",
			})
		}
	}

	sendResponse(&results, "", w)
}

// sendResponse sends back the response to Gatekeeper.
func sendResponse(results *[]externaldata.Item, systemErr string, w http.ResponseWriter) {
	response := externaldata.ProviderResponse{
		APIVersion: apiVersion,
		Kind:       "ProviderResponse",
	}

	if results != nil {
		response.Response.Items = *results
	} else {
		response.Response.SystemError = systemErr
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SaveCredentials(registryUrl string, username string, password string) error {
	// Obtenemos la ubicación del directorio de inicio del usuario

	// Creamos una estructura de tipo Credentials
	authConfig := map[string]interface{}{
		"username":      username,
		"password":      password,
		"serveraddress": registryUrl,
	}
	authConfigJson, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}
	encodedAuth := base64.URLEncoding.EncodeToString(authConfigJson)
	config := map[string]interface{}{
		"auths": map[string]interface{}{
			registryUrl: map[string]string{
				"auth": encodedAuth,
			},
		},
	}
	configJson, err := json.Marshal(config)

	configFile, err := os.Create("config.json")
	if err != nil {
		return err
	}
	defer configFile.Close()
	if _, err := configFile.Write(configJson); err != nil {
		return err
	}
	fmt.Println("Credentials saved in", configFile.Name())
	return nil
}

type MyKeychain struct {
	registry string
	username string
	password string
}

func (kc MyKeychain) Resolve(target authn.Resource) (authn.Authenticator, error) {
	creds := authn.FromConfig(authn.AuthConfig{
		Username: kc.username,
		Password: kc.password,
	})
	return creds, nil
}

func listFilesWithExtension(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	var result []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == "crt" {
			result = append(result, filepath.Join(path, file.Name()))
		}
	}

	return result
}
