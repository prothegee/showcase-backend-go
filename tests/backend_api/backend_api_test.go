/*
IMPORTANT:
- should be on last test
*/
package test_unittest

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"

	backend_api "showcase-backend-go/cmd/backend_api/api"
	backend_api_account "showcase-backend-go/cmd/backend_api/api/account"
	backend_api_auth "showcase-backend-go/cmd/backend_api/api/auth"
	backend_api_game1 "showcase-backend-go/cmd/backend_api/api/game1"
	"showcase-backend-go/pkg"
	config "showcase-backend-go/pkg/configs"
	mw "showcase-backend-go/pkg/middleware"

	"github.com/google/uuid"
)

// --------------------------------------------------------- //

// net global backend_api test data
//
// cfg - server config
//
// server - as in target from backend listener
var (
	cfg pkg.ConfigServer
	server string
)

// end-user temporal global backend_api test data
//
// userId
//
// email
//
// password
//
// authorizationData
//
// stashName
var (
	userId uuid.UUID
	email, password, authorizationData, stashName string
)

// --------------------------------------------------------- //

// func doRequest(method, url string,
// 			   headers map[string]string,
// 		   	   body any) (*http.Response, error) {
// 	var err error
// 	var bodyBytes []byte
// }

func initializeConfig() {
	var err error
	cfg, err = pkg.ConfigServerLoad(config.BACKEND_API_CONFIG_JSON); if err != nil {
		log.Fatal(err.Error())
	}

	server = fmt.Sprintf("http://%s:%s",
		cfg.Listener.BackendApi.Address, strconv.Itoa(int(cfg.Listener.BackendApi.Port)))

	if len(server) <= 0 {
		log.Fatal("server variable still empty")
	}

	genUser, err := pkg.GenRandomNumber(3, 9);
	genDomain, err := pkg.GenRandomNumber(3, 9)
	genTld, err := pkg.GenRandomNumber(3, 6)
	genPassword, err := pkg.GenRandomNumber(6, 9)

	user, _ := pkg.GenRandomAlphanumeric(genUser)
	domain, _ := pkg.GenRandomAlphanumeric(genDomain)
	tld, _ := pkg.GenRandomAlphanumeric(genTld)
	password, _ = pkg.GenRandomAlphanumeric(genPassword)

	email = fmt.Sprintf("%s@%s.%s",
		strings.ToLower(user),
		strings.ToLower(domain),
		strings.ToLower(tld))
}

// --------------------------------------------------------- //

func TestBackendApi_1_check_status(t *testing.T) {
	initializeConfig()

	url := fmt.Sprint(server + backend_api.BackendApiStatusHint)
	body := []byte(nil)

	if len(cfg.Security.WhitelistHost) <= 0 || len(cfg.Security.WhitelistOrigin) <= 0 {
		t.Fatal("whitelist host or whitelist origin can't be empty\n")
	}

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(body)); if err != nil {
		t.Fatalf("fail to make new request; %v\n", err.Error())
	}

	req.Host = cfg.Security.WhitelistHost[0]
	req.Header.Set(pkg.HTTP_HEADER_ORIGIN, cfg.Security.WhitelistOrigin[0])

	client := &http.Client{}

	resp, err := client.Do(req); if err != nil {
		t.Fatalf("can't do client request: %v\n", err.Error())
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expecting 200 but got %d; resp body: %v\n", resp.StatusCode, string(respBody))
	}
}

func TestBackendApi_2_create_account_user(t *testing.T) {
	url := fmt.Sprint(server + backend_api_account.BackendApiAccountUserHint)
	body := map[string]any{
		"email": email,
		"password": password,
	}
	bodyBytes, err := json.Marshal(body); if err != nil {
		t.Fatal("fail to make json marshal\n")
	}

	req, err := http.NewRequest(http.MethodPost,
		url, bytes.NewBuffer(bodyBytes)); if err != nil {
		t.Fatalf("fail to make new request; %v\n", err.Error())
	}

	req.Host = cfg.Security.WhitelistHost[0]
	req.Header.Set(pkg.HTTP_HEADER_ORIGIN, cfg.Security.WhitelistOrigin[0])

	client := &http.Client{}

	resp, err := client.Do(req); if err != nil {
		t.Fatalf("can't do client request; %v\n", err.Error())
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expecting 200 but got %d; resp body: %v\n", resp.StatusCode, string(respBody))
	}
}

// @note raw control test
func TestBackendApi_3_get_account_user_id(t *testing.T) {
	url := fmt.Sprint(
		server + backend_api_account.BackendApiAccountUserHint + "?email=" + email)

	req, err := http.NewRequest(http.MethodGet, url, nil); if err != nil {
		t.Fatalf("fail to make new request; %v\n", err.Error())
	}

	req.Host = cfg.Security.WhitelistHost[0]
	req.Header.Set(pkg.HTTP_HEADER_ORIGIN, cfg.Security.WhitelistOrigin[0])

	client := &http.Client{}

	resp, err := client.Do(req); if err != nil {
		t.Fatalf("can't do client request; %v\n", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expecting 200 but got %d\n", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body); if err != nil {
		t.Fatalf("fail to read resp body data; %v; body: %s\n", err.Error(), string(respBody))
	}

	var raw map[string]any
	err = json.Unmarshal(respBody, &raw); if err != nil {
		t.Fatalf("unmarshal failed %v\n", err.Error())
	}

	okValue, okExists := raw["ok"]; if !okExists {
		t.Fatal("ok field doesn't exists\n")
	}

	ok := okValue.(bool); if !ok {
		t.Fatalf("ok field value must be true\n")
	}

	dataValue, dataExists := raw["data"]; if !dataExists {
		t.Fatalf("data field doesn't exists\n")
	}

	dataMap, dataMapOk := dataValue.(map[string]any); if !dataMapOk {
		t.Fatalf("data map must be an object; got %T = %v\n", dataMap, dataMap)
	}

	idValue, idExists := dataMap["id"]; if !idExists {
		t.Fatal("expecting id field in data object\n")
	}

	idStr := strings.TrimSpace(idValue.(string))

	idUuid, err := uuid.Parse(idStr); if err != nil {
		t.Fatalf("fail to parse id of uuid; %v\n", err.Error())
	}

	userId = idUuid
	authorizationData = base64.StdEncoding.EncodeToString([]byte(userId.String()))
}

func TestBackendApi_4_create_session(t *testing.T) {
	url := fmt.Sprint(server + backend_api_auth.BackendApiAuthSessionHint)

	req, err := http.NewRequest(http.MethodPost, url, nil); if err != nil {
		t.Fatalf("fail to make new request; %v\n", err.Error())
	}

	req.Host = cfg.Security.WhitelistHost[0]
	req.Header.Set(pkg.HTTP_HEADER_ORIGIN, cfg.Security.WhitelistOrigin[0])
	req.Header.Set(pkg.HTTP_HEADER_AUTHORIZATION,
		mw.AuthorizationHeadKey_bearer+" "+authorizationData)

	client := &http.Client{}

	resp, err := client.Do(req); if err != nil {
		t.Fatalf("can't do client request; %v\n", err.Error())
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expecting 200 but got %d; resp body: %v\n", resp.StatusCode, string(respBody))
	}
}

func TestBackendApi_5_update_email(t *testing.T) {
	genUser, _ := pkg.GenRandomNumber(3, 9)
	genDomain, _ := pkg.GenRandomNumber(3, 9)
	genTld, _ := pkg.GenRandomNumber(3, 6)
	user, _ := pkg.GenRandomAlphanumeric(genUser)
	domain, _ := pkg.GenRandomAlphanumeric(genDomain)
	tld, _ := pkg.GenRandomAlphanumeric(genTld)

	newEmail := fmt.Sprintf("%s@%s.%s",
		strings.ToLower(user),
		strings.ToLower(domain),
		strings.ToLower(tld))

	url := fmt.Sprint(server + backend_api_account.BackendApiAccountUserHint)
	body := map[string]any{
		"id": userId.String(),
		"email": newEmail,
	}
	bodyBytes, err := json.Marshal(body); if err != nil {
		t.Fatal("fail to make json marshal\n")
	}

	req, err := http.NewRequest(http.MethodPatch, url,
		bytes.NewBuffer(bodyBytes)); if err != nil {
		t.Fatalf("fail to make new request; %v\n", err.Error())
	}

	req.Host = cfg.Security.WhitelistHost[0]
	req.Header.Set(pkg.HTTP_HEADER_ORIGIN, cfg.Security.WhitelistOrigin[0])
	req.Header.Set(pkg.HTTP_HEADER_AUTHORIZATION,
		mw.AuthorizationHeadKey_bearer+" "+authorizationData)

	client := &http.Client{}

	resp, err := client.Do(req); if err != nil {
		t.Fatalf("can't do client request; %v\n", err.Error())
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expecting 200 but got %d; resp body: %v\n", resp.StatusCode, string(respBody))
	}
}

func TestBackendApi_6_create_stash(t *testing.T) {
	url := fmt.Sprint(server + backend_api_game1.BackendApiGame1StashHint)
	stashGen, err := pkg.GenRandomAlphanumeric(6)
	stashName = "Stash " + stashGen
	body := map[string]any{
		"name": stashName,
	}
	bodyBytes, err := json.Marshal(body); if err != nil {
		t.Fatal("fail to make json marshal\n")
	}

	req, err := http.NewRequest(http.MethodPost, url,
		bytes.NewBuffer(bodyBytes)); if err != nil {
		t.Fatalf("fail to make new request; %v\n", err.Error())
	}

	req.Host = cfg.Security.WhitelistHost[0]
	req.Header.Set(pkg.HTTP_HEADER_ORIGIN, cfg.Security.WhitelistOrigin[0])
	req.Header.Set(pkg.HTTP_HEADER_AUTHORIZATION,
		mw.AuthorizationHeadKey_bearer+" "+authorizationData)

	client := &http.Client{}

	resp, err := client.Do(req); if err != nil {
		t.Fatalf("can't do client request; %v\n", err.Error())
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expecting 200 but got %d; resp body: %v\n", resp.StatusCode, string(respBody))
	}
}

func TestBackendApi_7_update_stash(t *testing.T) {
	url := fmt.Sprint(server + backend_api_game1.BackendApiGame1StashHint)
	operand, err := pkg.GenRandomNumber(1, 2); if err != nil {
		t.Fatalf("fail to generate random number for operand; %v\n", err.Error())
	}
	itemGenLength, err := pkg.GenRandomNumber(3, 9); if err != nil {
		t.Fatalf("fail to generate random number for itemGenLength; %v\n", err.Error())
	}
	item, err := pkg.GenRandomAlphanumeric(itemGenLength); if err != nil {
		t.Fatalf("fail to generate item name; %v\n", err.Error())
	}
	quantity, err := pkg.GenRandomNumber(1, 3); if err != nil {
		t.Fatalf("fail to genrate random number for quantity; %v\n", err.Error())
	}
	body := map[string]any{
		"name": stashName,
		"operand": operand,
		"item": item,
		"quantity": quantity,
	}
	bodyBytes, err := json.Marshal(body); if err != nil {
		t.Fatal("fail to make json json.Marshal\n")
	}

	req, err := http.NewRequest(http.MethodPatch,
		url, bytes.NewBuffer(bodyBytes)); if err != nil {
		t.Fatalf("fail to make new request; %v\n", err.Error())
	}

	req.Host = cfg.Security.WhitelistHost[0]
	req.Header.Set(pkg.HTTP_HEADER_ORIGIN, cfg.Security.WhitelistOrigin[0])
	req.Header.Set(pkg.HTTP_HEADER_AUTHORIZATION,
		mw.AuthorizationHeadKey_bearer+" "+authorizationData)

	client := &http.Client{}

	resp, err := client.Do(req); if err != nil {
		t.Fatalf("can't do client request; %v\n", err.Error())
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expecting 200 got %d; resp body: %v\n", resp.StatusCode, string(respBody))
	}
}

