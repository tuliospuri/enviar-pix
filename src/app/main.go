package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func HealthCheckAction(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("ok"))
}

func createHttpClient() *http.Client {
	// Recuperar o certificado
	rawCert := strings.ReplaceAll(os.Getenv("CERT"), `\n`, "\n")
	rawKey := strings.ReplaceAll(os.Getenv("KEY"), `\n`, "\n")

	// Carregar o certificado
	certificate, err := tls.X509KeyPair([]byte(rawCert), []byte(rawKey))
	if err != nil {
		log.Fatal("Erro ao carregar certificado")
	}

	// Carregar o certificado para o HTTP Client
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{certificate},
		},
	}

	// Criar o HTTP Client
	httpClient := &http.Client{
		Transport: transport,
	}

	return httpClient
}

func createOAuthTokenRequest() *http.Request {
	// Definir o endpoint
	endpoint := fmt.Sprintf("%s/oauth/token", os.Getenv("API"))

	// Criar o corpo da requisição
	body := `{"grant_type": "client_credentials"}`

	// Criar a requisição
	request, err := http.NewRequest("POST", endpoint, strings.NewReader(body))
	if err != nil {
		log.Fatal("Erro ao criar a requisição /oauth/token")
	}

	// Adicionar cabeçalhos
	bearer := fmt.Sprintf("%s:%s", os.Getenv("ID"), os.Getenv("SECRET"))
	bearerB64 := base64.StdEncoding.EncodeToString([]byte(bearer))

	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", bearerB64))

	return request
}

func getAccessTokenFromOAuthTokenResponse(response *http.Response) string {
	// Ler o corpo da resposta
	oauthTokenResponseBody, _ := ioutil.ReadAll(response.Body)

	var oauthTokenJson map[string]interface{}

	if err := json.Unmarshal(oauthTokenResponseBody, &oauthTokenJson); err != nil {
		log.Fatal("Erro ao gerar JSON da resposta /oauth/token")
	}

	// Retornar o access_token
	return oauthTokenJson["access_token"].(string)
}

func createPixSendRequest(accessToken string) *http.Request {
	// Definir o endpoint
	endpoint := fmt.Sprintf("%s/v2/pix", os.Getenv("API"))

	// Criar o corpo da requisição
	body := `{
		"valor": "0.05",
		"pagador": {
			"chave": "62a7436b-b2d6-4da0-89f5-38f7b6935446",
			"infoPagador": "Envio de Pix por voz com a Siri"
		},
		"favorecido": {
			"chave": "41af8590-5fc2-4b75-a99d-5b8e55a61428"
		}
	}`

	// Criar a requisição
	request, err := http.NewRequest("POST", endpoint, strings.NewReader(body))
	if err != nil {
		log.Fatal("Erro ao criar a requisição /v2/pix")
	}

	// Adicionar cabeçalhos
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	return request
}

func IndexAction(response http.ResponseWriter, request *http.Request) {
	// Criar o HTTP Client
	httpClient := createHttpClient()

	// Criar requisição para obter o Access Token da API
	oauthTokenRequest := createOAuthTokenRequest()

	// Executar a requisição
	oauthTokenResponse, err := httpClient.Do(oauthTokenRequest)
	if err != nil {
		log.Fatal("Erro ao executar a requisição /oauth/token")
	}

	// Recuperar access_token da resposta
	accessToken := getAccessTokenFromOAuthTokenResponse(oauthTokenResponse)

	// Criar requisição de Envio de Pix
	pixSendRequest := createPixSendRequest(accessToken)

	// Executar a requisição
	pixSendResponse, err := httpClient.Do(pixSendRequest)
	if err != nil {
		log.Fatal("Erro ao executar a requisição /v2/pix")
	}

	// Ler o corpo da resposta
	pixSendResponseBody, _ := ioutil.ReadAll(pixSendResponse.Body)

	var pixSendJson map[string]interface{}

	if err := json.Unmarshal(pixSendResponseBody, &pixSendJson); err != nil {
		log.Fatal("Erro ao gerar JSON da resposta /v2/pix")
	}

	// Verificar o resultado do Envio de Pix
	if e2eId, ok := pixSendJson["e2eId"].(string); ok && e2eId != "" {
		response.WriteHeader(http.StatusOK)
		response.Write([]byte("ok"))
		return
	}

	response.WriteHeader(http.StatusBadRequest)
	response.Write([]byte("nok"))
	return
}

func main() {
	port := os.Getenv("PORT")

	router := mux.NewRouter()

	router.Path("/healthcheck").HandlerFunc(HealthCheckAction)
	router.Path("/").HandlerFunc(IndexAction)

	http.ListenAndServe(":"+port, router)
}
