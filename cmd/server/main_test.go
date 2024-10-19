package main

// Working, but in development

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type HealthCheckResponse struct {
	Status string `json:"status"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type AboutMeResponse struct {
	Login string `json:"login"`
	Skins []skin
}

type RemoveSkinResponse struct {
	Status string `json:"status"`
}

type skin struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"`
	Type string `json:"Type"`
	Src  string `json:"Src"`
}

const (
	addr             string = "http://localhost:8081/api/v1"
	TestUserLogin    string = "User2347"
	TestUserPassword string = "234"
	AuthHeader       string = "Authorization"
)

func TestApi(t *testing.T) {
	err := healthCheckTest(t)
	if err != nil {
		fmt.Printf("Error test #1: %v", err)
	} else {
		fmt.Println("Succesful test #1")
	}

	err = registerUserTest(t)
	if err != nil {
		fmt.Printf("Error test #2: %v", err)
	} else {
		fmt.Println("Succesful test #2")
	}

	token, err := loginTest(t)
	if err != nil {
		fmt.Printf("Error test #3: %v", err)
	} else {
		fmt.Println("Succesful test #3")
	}

	expected := &AboutMeResponse{
		Login: TestUserLogin,
		Skins: nil,
	}
	err = aboutMeTest(t, token, expected)
	if err != nil {
		fmt.Printf("Error test #4: %v", err)
	} else {
		fmt.Println("Succesful test #4")
	}

	id, err := addSkinTest(t, token)
	if err != nil {
		fmt.Printf("Error test #5: %v", err)
	} else {
		fmt.Println("Succesful test #5")
	}

	err = getSkinsTest(t, token, id)
	if err != nil {
		fmt.Printf("Error test #6: %v", err)
	} else {
		fmt.Println("Succesful test #6")
	}

	err = GetSkinTest(t, token, id)
	if err != nil {
		fmt.Printf("Error test #7: %v", err)
	} else {
		fmt.Println("Succesful test #7")
	}

	expected = &AboutMeResponse{
		Login: TestUserLogin,
		Skins: []skin{{Id: id, Name: "MySkin", Type: "Classic", Src: "Aid"}},
	}
	err = aboutMeTest(t, token, expected)
	if err != nil {
		fmt.Printf("Error test #8: %v", err)
	} else {
		fmt.Println("Succesful test #8")
	}

	err = RemoveSkinTest(t, token, id)
	if err != nil {
		fmt.Printf("Error test #9: %v", err)
	} else {
		fmt.Println("Succesful test #9")
	}
}

func healthCheckTest(t *testing.T) error {

	expected := "ok"
	expectedCode := 200

	resp, err := http.Get(addr + "/")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	assert.Equal(t, expectedCode, resp.StatusCode)

	var hCheckResp *HealthCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&hCheckResp); err != nil {
		fmt.Printf("Error decoding Health Check response: %v", err)
		return err
	}

	assert.Equal(t, expected, hCheckResp.Status)

	return nil
}

func registerUserTest(t *testing.T) error {

	expected := "Success"
	expectedCode := 200
	expectedCodeIfAlrRegistered := 400
	expectedIfAlrRegistered := "This user is already registered"

	body := gin.H{
		"login":    TestUserLogin,
		"password": TestUserPassword,
	}
	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(addr+"/user/register", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	var registerResp *RegisterResponse
	if resp.StatusCode == expectedCodeIfAlrRegistered {
		if err := json.NewDecoder(resp.Body).Decode(&registerResp); err != nil {
			fmt.Printf("Error decoding Register response: %v", err)
			return err
		}
		assert.Equal(t, expectedIfAlrRegistered, registerResp.Error)
		return nil
	}

	assert.Equal(t, expectedCode, resp.StatusCode)

	if err := json.NewDecoder(resp.Body).Decode(&registerResp); err != nil {
		fmt.Printf("Error decoding Register response: %v", err)
		return err
	}

	assert.Equal(t, expected, registerResp.Message)

	return nil

}

func loginTest(t *testing.T) (string, error) {

	expectedCode := 200

	body := gin.H{
		"login":    TestUserLogin,
		"password": TestUserPassword,
	}
	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(addr+"/user/login", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer resp.Body.Close()

	assert.Equal(t, expectedCode, resp.StatusCode)

	var loginResp *LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		fmt.Printf("Error decoding Login response: %v", err)
		return "", err
	}
	return loginResp.Token, nil
}

func aboutMeTest(t *testing.T, token string, expected *AboutMeResponse) error {

	expectedCode := 200

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, addr+"/user/me", nil)
	if err != nil {
		log.Fatal(err)
		return err
	}
	req.Header.Add(AuthHeader, "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	assert.Equal(t, expectedCode, resp.StatusCode)

	var aboutMeResp *AboutMeResponse
	if err := json.NewDecoder(resp.Body).Decode(&aboutMeResp); err != nil {
		fmt.Printf("Error decoding AboutMe response: %v", err)
		return err
	}

	assert.Equal(t, expected.Login, aboutMeResp.Login)
	assert.Equal(t, expected.Skins, aboutMeResp.Skins)

	return nil

}

func addSkinTest(t *testing.T, token string) (int, error) {

	expected := &skin{
		Name: "MySkin",
		Type: "Classic",
		Src:  "Aid",
	}
	expectedCode := 201

	body := gin.H{
		"skinname": "MySkin",
		"skintype": "Classic",
		"skinsrc":  "Aid",
	}
	jsonBody, _ := json.Marshal(body)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, addr+"/skins/add", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	req.Header.Add(AuthHeader, "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer resp.Body.Close()

	assert.Equal(t, expectedCode, resp.StatusCode)

	var skinResp *skin
	if err := json.NewDecoder(resp.Body).Decode(&skinResp); err != nil {
		fmt.Printf("Error decoding AddSkin response: %v", err)
		return 0, err
	}

	assert.Equal(t, expected.Name, skinResp.Name)
	assert.Equal(t, expected.Type, skinResp.Type)
	assert.Equal(t, expected.Src, skinResp.Src)

	return skinResp.Id, nil
}

func getSkinsTest(t *testing.T, token string, skinId int) error {

	expected := []skin{
		{
			Id:   skinId,
			Name: "MySkin",
			Type: "Classic",
			Src:  "Aid",
		},
	}
	expectedCode := 200

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, addr+"/skins", nil)
	if err != nil {
		log.Fatal(err)
		return err
	}
	req.Header.Add(AuthHeader, "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	assert.Equal(t, expectedCode, resp.StatusCode)

	var skinsResp []skin
	if err := json.NewDecoder(resp.Body).Decode(&skinsResp); err != nil {
		fmt.Printf("Error decoding GetSkins response: %v", err)
		return err
	}

	assert.Equal(t, expected, skinsResp)

	return nil
}

func GetSkinTest(t *testing.T, token string, skinId int) error {

	expected := &skin{
		Id:   skinId,
		Name: "MySkin",
		Type: "Classic",
		Src:  "Aid",
	}
	expectedCode := 200

	client := &http.Client{}

	reqUrl := fmt.Sprintf("/skins/%d", skinId)
	req, err := http.NewRequest(http.MethodGet, addr+reqUrl, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}
	req.Header.Add(AuthHeader, "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	assert.Equal(t, expectedCode, resp.StatusCode)

	var skinResp *skin
	if err := json.NewDecoder(resp.Body).Decode(&skinResp); err != nil {
		fmt.Printf("Error decoding AddSkin response: %v", err)
		return err
	}

	assert.Equal(t, expected, skinResp)
	return nil
}

func RemoveSkinTest(t *testing.T, token string, skinId int) error {
	expected := &RemoveSkinResponse{Status: "Success"}
	expectedCode := 200

	client := &http.Client{}

	reqUrl := fmt.Sprintf("/skins/%d", skinId)
	req, err := http.NewRequest(http.MethodGet, addr+reqUrl, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}
	req.Header.Add(AuthHeader, "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	assert.Equal(t, expectedCode, resp.StatusCode)

	var skinResp *RemoveSkinResponse
	if err := json.NewDecoder(resp.Body).Decode(&skinResp); err != nil {
		fmt.Printf("Error decoding AddSkin response: %v", err)
		return err
	}

	assert.Equal(t, expected, skinResp)
	return nil
}
