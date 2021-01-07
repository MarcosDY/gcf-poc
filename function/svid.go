package function

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/spiffe/go-spiffe/v2/proto/spiffe/workload"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"google.golang.org/protobuf/proto"
)

var (
	// Secret name in format 'projects/*/secrets/*/versions/*'
	secretName = os.Getenv("SECRET_NAME")

	// Hold X509-SVID parsed from secret
	x509SVID = new(workload.X509SVIDResponse)
)

func init() {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to create secretmanager client: %v", err)
	}

	resp, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	})
	if err != nil {
		log.Fatalf("failed to get secret %q: %v", secretName, err)
	}

	if err := proto.Unmarshal(resp.Payload.Data, x509SVID); err != nil {
		log.Fatalf("failed to unmarshal: %v", err)
	}
}

func SvidGet(w http.ResponseWriter, r *http.Request) {
	svid := x509SVID.Svids[0]

	cert, err := parseCerts(svid.X509Svid)
	if err != nil {
		http.Error(w, "Unable to parse certificates", http.StatusBadRequest)
		log.Printf("failed to parse X509-SVID certificates: %v", err)
		return
	}

	bundle, err := parseCerts(svid.Bundle)
	if err != nil {
		http.Error(w, "Unable to parse bundle", http.StatusBadRequest)
		log.Printf("failed to parse X509-SVID bundle: %v", err)
		return
	}

	key, err := parseKey(svid.X509SvidKey)
	if err != nil {
		http.Error(w, "Unable to parse key", http.StatusBadRequest)
		log.Printf("failed to parse X509-SVID key: %v", err)
		return
	}

	type svidJson struct {
		SpiffeID     string `json:"spiffe_id"`
		Certificates string `json:"certificates"`
		Bundle       string `json:"bundle"`
		Key          string `json:"key"`
	}

	s := &svidJson{
		SpiffeID:     svid.SpiffeId,
		Certificates: cert,
		Bundle:       bundle,
		Key:          key,
	}
	data, _ := json.MarshalIndent(s, "", "  ")
	fmt.Fprintf(w, string(data))
}

func parseCerts(raw []byte) (string, error) {
	certs, err := x509.ParseCertificates(raw)
	if err != nil {
		return "", err
	}

	pemData := []byte{}
	for _, cert := range certs {
		// TODO: demonstration purposes only, remove it
		log.Printf("SPIFFE ID: %q, NotAfter: %s", cert.URIs, cert.NotAfter.String())
		b := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		}
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	return string(pemData), nil
}

func parseKey(rawKey []byte) (string, error) {
	privateKey, err := x509.ParsePKCS8PrivateKey(rawKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	data, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	b := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: data,
	}

	return string(pem.EncodeToMemory(b)), nil
}
