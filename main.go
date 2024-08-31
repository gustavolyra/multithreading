package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type BrasilApi struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Unidade     string `json:"unidade"`
	Complemento string `json:"complemento"`

	Ibge  string `json:"ibge"`
	Gia   string `json:"gia"`
	Ddd   string `json:"ddd"`
	Siafi string `json:"siafi"`
}

func main() {
	ch1 := make(chan BrasilApi)
	ch2 := make(chan ViaCep)

	var cep string
	fmt.Printf("CEP (apenas numeros): ")
	_, err := fmt.Scanln(&cep)
	if err != nil {
		panic(err)
	}

	go getBrasil(ch1, cep)
	go getViaCep(ch2, cep)

	select {
	case json1 := <-ch1:
		fmt.Println("--> BrasilApi <--")
		fmt.Println("Cidade: ", json1.City)
		fmt.Println("Bairro: ", json1.Neighborhood)
		fmt.Println("Logradouro: ", json1.Street)
		fmt.Println("Estado: ", json1.State)
		fmt.Println("Cep: ", json1.Cep)
		fmt.Println("Servico: ", json1.Service)

	case json2 := <-ch2:
		fmt.Println("--> ViaCep <--")
		fmt.Println("Cidade: ", json2.Localidade)
		fmt.Println("Bairro: ", json2.Bairro)
		fmt.Println("Logradouro: ", json2.Logradouro)
		fmt.Println("Complemento: ", json2.Complemento)
		fmt.Println("Unidade: ", json2.Unidade)
		fmt.Println("Estado: ", json2.Uf)
		fmt.Println("Cep: ", json2.Cep)
		fmt.Println("Ibge: ", json2.Ibge)
		fmt.Println("Gia: ", json2.Gia)
		fmt.Println("Ddd: ", json2.Ddd)
		fmt.Println("Siafi: ", json2.Siafi)

	case <-time.After(time.Second * 1):
		fmt.Printf("!!! Timeout !!!")
	}
}

func getBrasil(ch chan BrasilApi, cep string) {
	req, err := http.NewRequest("GET", "https://brasilapi.com.br/api/cep/v1/"+cep, nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var jsonBrasil BrasilApi
	err = json.Unmarshal([]byte(body), &jsonBrasil)
	if err != nil {
		panic(err)
	}
	ch <- jsonBrasil
}

func getViaCep(ch chan ViaCep, cep string) {
	req, err := http.NewRequest("GET", "https://viacep.com.br/ws/"+cep+"/json/", nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var jsonViaCep ViaCep
	err = json.Unmarshal([]byte(body), &jsonViaCep)
	if err != nil {
		panic(err)
	}
	ch <- jsonViaCep
}
