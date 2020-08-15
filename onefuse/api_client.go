// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

const ApiVersion = "api/v3"
const ApiNamespace = "onefuse"
const NamingResourceType = "customNames"
const WorkspaceResourceType = "workspaces"
const MicrosoftADPolicyResourceType = "microsoftADPolicies"
const ModuleEndpointResourceType = "endpoints"

type OneFuseAPIClient struct {
	config *Config
}

type CustomName struct {
	Id        int
	Version   int
	Name      string
	DnsSuffix string
}

type LinkRef struct {
	Href  string `json:"href,omitempty"`
	Title string `json:"title,omitempty"`
}

type Workspace struct {
	Links *struct {
		Self LinkRef `json:"self,omitempty"`
	} `json:"_links,omitempty"`
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type WorkspacesListResponse struct {
	Embedded struct {
		Workspaces []Workspace `json:"workspaces"`
	} `json:"_embedded"`
}

type EndpointsListResponse struct {
	Embedded struct {
		Endpoints []MicrosoftEndpoint `json:"endpoints"` // TODO: Generalize to Endpoints
	} `json:"_embedded"`
}

type MicrosoftEndpoint struct {
	Links *struct {
		Self       *LinkRef `json:"self,omitempty"`
		Workspace  *LinkRef `json:"workspace,omitempty"`
		Credential *LinkRef `json:"credential,omitempty"`
	} `json:"_links,omitempty"`
	ID               int    `json:"id,omitempty"`
	Type             string `json:"type,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"string,omitempty"`
	Host             string `json:"host,omitempty"`
	Port             int    `json:"port,omitempty"`
	SSL              bool   `json:"ssl,omitempty"`
	MicrosoftVersion string `json:"microsoftVersion,omitempty"`
}

type MicrosoftADPolicy struct {
	Links *struct {
		Self              *LinkRef `json:"self,omitempty"`
		Workspace         *LinkRef `json:"workspace,omitempty"`
		MicrosoftEndpoint *LinkRef `json:"microsoftEndpoint,omitempty"`
	} `json:"_links,omitempty"`
	Name                   string `json:"name,omitempty"`
	ID                     int    `json:"id,omitempty"`
	Description            string `json:"description,omitempty"`
	MicrosoftEndpointID    int    `json:"microsoftEndpointId,omitempty"`
	MicrosoftEndpoint      string `json:"microsoftEndpoint,omitempty"`
	ComputerNameLetterCase string `json:"computerNameLetterCase,omitempty"`
	WorkspaceURL           string `json:"workspace,omitempty"`
	OU                     string `json:"ou,omitempty"`
}

func (c *Config) NewOneFuseApiClient() *OneFuseAPIClient {
	return &OneFuseAPIClient{
		config: c,
	}
}

func (apiClient *OneFuseAPIClient) GenerateCustomName(dnsSuffix string, namingPolicyID string, workspaceID string, templateProperties map[string]interface{}) (*CustomName, error) {
	log.Println("onefuse.apiClient: GenerateCustomName")

	config := apiClient.config
	url := collectionURL(config, NamingResourceType)

	if templateProperties == nil {
		templateProperties = make(map[string]interface{})
	}

	if workspaceID == "" {
		defaultWorkspaceID, err := findDefaultWorkspaceID(config)
		if err != nil {
			return nil, err
		} else {
			workspaceID = defaultWorkspaceID
		}
	}

	postBody := map[string]interface{}{
		"namingPolicy":       fmt.Sprintf("/%s/%s/namingPolicies/%s/", ApiVersion, ApiNamespace, namingPolicyID),
		"templateProperties": templateProperties,
		"workspace":          fmt.Sprintf("/%s/%s/workspaces/%s/", ApiVersion, ApiNamespace, workspaceID),
	}

	jsonBytes, err := json.Marshal(postBody)
	if err != nil {
		return nil, errors.WithMessage(err, "onefuse.apiClient: Unable to marshal request body to JSON")
	}

	requestBody := string(jsonBytes)

	payload := strings.NewReader(requestBody)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Unable to create request POST %s %s", url, requestBody))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request POST %s %s", url, requestBody))
	}

	if err = checkForErrors(res); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Request failed POST %s %s", url, requestBody))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from %s %s", url, requestBody))
	}
	defer res.Body.Close()

	customName := CustomName{}
	if err = json.Unmarshal(body, &customName); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	return &customName, nil
}

func (apiClient *OneFuseAPIClient) GetCustomName(id int) (*CustomName, error) {
	log.Println("onefuse.apiClient: GetCustomName")

	config := apiClient.config

	url := itemURL(config, NamingResourceType, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request GET %s", url))
	}

	setHeaders(req, config)

	client := getHttpClient(config)
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request GET %s", url))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from GET %s", url))
	}
	defer res.Body.Close()

	customName := CustomName{}
	if err = json.Unmarshal(body, &customName); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	return &customName, nil
}

func (apiClient *OneFuseAPIClient) DeleteCustomName(id int) error {
	log.Println("onefuse.apiClient: DeleteCustomName")

	config := apiClient.config

	url := itemURL(config, NamingResourceType, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request DELETE %s", url))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request DELETE %s", url))
	}

	return checkForErrors(res)
}

func (apiClient *OneFuseAPIClient) CreateMicrosoftEndpoint(newEndpoint MicrosoftEndpoint) (*MicrosoftEndpoint, error) {
	log.Println("onefuse.apiClient: CreateMicrosoftEndpoint")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) GetMicrosoftEndpoint(id int) (*MicrosoftEndpoint, error) {
	log.Println("onefuse.apiClient: GetMicrosoftEndpoint")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) GetMicrosoftEndpointByName(name string) (*MicrosoftEndpoint, error) {
	log.Println("onefuse.apiClient: GetMicrosoftEndpointByName")

	config := apiClient.config
	url := fmt.Sprintf("%s?filter=name:%s;type:microsoft", collectionURL(config, ModuleEndpointResourceType), name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request GET %s", url))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request GET %s", url))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from GET %s", url))
	}
	defer res.Body.Close()

	endpoints := EndpointsListResponse{}
	err = json.Unmarshal(body, &endpoints)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	if len(endpoints.Embedded.Endpoints) < 1 {
		return nil, errors.New(fmt.Sprintf("onefuse.apiClient: Could not find Microsoft Endpoint '%s'!", name))
	}

	endpoint := endpoints.Embedded.Endpoints[0]

	return &endpoint, err
}

func (apiClient *OneFuseAPIClient) UpdateMicrosoftEndpoint(id int, updatedEndpoint MicrosoftEndpoint) (*MicrosoftEndpoint, error) {
	log.Println("onefuse.apiClient: UpdateMicrosoftEndpoint")

	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) DeleteMicrosoftEndpoint(id int) error {
	log.Println("onefuse.apiClient: DeleteMicrosoftEndpoint")

	return errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) CreateMicrosoftADPolicy(newPolicy *MicrosoftADPolicy) (*MicrosoftADPolicy, error) {
	log.Println("onefuse.apiClient: CreateMicrosoftADPolicy")

	config := apiClient.config

	// Default workspace if it was not provided
	if newPolicy.WorkspaceURL == "" {
		workspaceID, err := findDefaultWorkspaceID(config)
		if err != nil {
			return nil, errors.WithMessage(err, "onefuse.apiClient: Failed to find default workspace")
		}
		workspaceIDInt, err := strconv.Atoi(workspaceID)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to convert Workspace ID '%s' to integer", workspaceID))
		}

		newPolicy.WorkspaceURL = itemURL(config, WorkspaceResourceType, workspaceIDInt)
	}

	// Construct a URL we are going to POST to
	// /api/v3/onefuse/microsoftADPolicies/
	url := collectionURL(config, MicrosoftADPolicyResourceType)

	jsonBytes, err := json.Marshal(newPolicy)
	if err != nil {
		return nil, errors.WithMessage(err, "onefuse.apiClient: Failed to marshal request body to JSON")
	}

	requestBody := string(jsonBytes)
	payload := strings.NewReader(requestBody)

	// Create the create request
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Unable to create request POST %s %s", url, requestBody))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	// Make the create request
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request POST %s %s", url, requestBody))
	}

	if err = checkForErrors(res); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Request failed POST %s %s", url, requestBody))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from POST %s %s", url, requestBody))
	}
	defer res.Body.Close()

	policy := MicrosoftADPolicy{}
	if err = json.Unmarshal(body, &policy); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	return &policy, nil
}

func (apiClient *OneFuseAPIClient) GetMicrosoftADPolicy(id int) (*MicrosoftADPolicy, error) {
	log.Println("onefuse.apiClient: GetMicrosoftADPolicy")

	config := apiClient.config

	url := itemURL(config, MicrosoftADPolicyResourceType, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request GET %s", url))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request GET %s", url))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from GET %s", url))
	}
	defer res.Body.Close()

	policy := MicrosoftADPolicy{}
	if err = json.Unmarshal(body, &policy); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	return &policy, err
}

func (apiClient *OneFuseAPIClient) UpdateMicrosoftADPolicy(id int, updatedPolicy *MicrosoftADPolicy) (*MicrosoftADPolicy, error) {
	log.Println("onefuse.apiClient: UpdateMicrosoftADPolicy")

	config := apiClient.config

	url := itemURL(config, MicrosoftADPolicyResourceType, id)

	if updatedPolicy.Name == "" {
		return nil, errors.New("onefuse.apiClient: Microsoft AD Policy Updates Require a Name")
	}

	if updatedPolicy.WorkspaceURL == "" {
		workspaceID, err := findDefaultWorkspaceID(config)
		if err != nil {
			return nil, errors.WithMessage(err, "onefuse.apiClient: Failed to find default workspace")
		}
		workspaceIDInt, err := strconv.Atoi(workspaceID)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to convert Workspace ID '%s' to integer", workspaceID))
		}

		updatedPolicy.WorkspaceURL = itemURL(config, WorkspaceResourceType, workspaceIDInt)
	}

	jsonBytes, err := json.Marshal(updatedPolicy)
	if err != nil {
		return nil, errors.WithMessage(err, "onefuse.apiClient: Failed to marshal request body to JSON")
	}

	requestBody := string(jsonBytes)
	payload := strings.NewReader(requestBody)

	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request PUT %s %s", url, requestBody))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request PUT %s %s", url, requestBody))
	}

	checkForErrors(res)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from PUT %s %s", url, requestBody))
	}
	defer res.Body.Close()

	policy := MicrosoftADPolicy{}
	if err = json.Unmarshal(body, &policy); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarhsal response %s", string(body)))
	}

	return &policy, nil
}

func (apiClient *OneFuseAPIClient) DeleteMicrosoftADPolicy(id int) error {
	log.Println("onefuse.apiClient: DeleteMicrosoftADPolicy")

	config := apiClient.config

	url := itemURL(config, MicrosoftADPolicyResourceType, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request DELETE %s", url))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request DELETE %s", url))
	}

	return checkForErrors(res)
}

func findDefaultWorkspaceID(config *Config) (workspaceID string, err error) {
	fmt.Println("onefuse.findDefaultWorkspaceID")

	filter := "filter=name.exact:Default"
	url := fmt.Sprintf("%s?%s", collectionURL(config, WorkspaceResourceType), filter)

	req, clientErr := http.NewRequest("GET", url, nil)
	if clientErr != nil {
		err = errors.WithMessage(clientErr, fmt.Sprintf("onefuse.findDefaultWorkspaceID: Failed to make request GET %s", url))
		return
	}

	setHeaders(req, config)

	client := getHttpClient(config)
	res, clientErr := client.Do(req)
	if clientErr != nil {
		err = errors.WithMessage(clientErr, fmt.Sprintf("onefuse.findDefaultWorkspaceID: Failed to do request GET %s", url))
		return
	}

	checkForErrors(res)

	body, clientErr := ioutil.ReadAll(res.Body)
	if clientErr != nil {
		err = errors.WithMessage(clientErr, fmt.Sprintf("onefuse.findDefaultWorkspaceID: Failed to do read body from GET %s", url))
		return
	}
	defer res.Body.Close()

	var data WorkspacesListResponse
	json.Unmarshal(body, &data)

	workspaces := data.Embedded.Workspaces
	if len(workspaces) == 0 {
		err = errors.WithMessage(clientErr, "onefuse.findDefaultWorkspaceID: Failed to find default workspace!")
		return
	}
	workspaceID = strconv.Itoa(workspaces[0].ID)
	return
}

func getHttpClient(config *Config) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.verifySSL},
	}
	return &http.Client{Transport: tr}
}

func checkForErrors(res *http.Response) error {
	if res.StatusCode >= 500 {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		return errors.New(string(b))
	} else if res.StatusCode >= 400 {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		return errors.New(string(b))
	}
	return nil
}

func setStandardHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("accept-encoding", "gzip, deflate")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")
}

func setHeaders(req *http.Request, config *Config) {
	setStandardHeaders(req)
	req.Header.Add("Host", fmt.Sprintf("%s:%s", config.address, config.port))
	req.Header.Add("SOURCE", "Terraform")
	req.SetBasicAuth(config.user, config.password)
}

func collectionURL(config *Config, resourceType string) string {
	baseURL := fmt.Sprintf("%s://%s:%s", config.scheme, config.address, config.port)
	endpoint := path.Join(ApiVersion, ApiNamespace, resourceType)
	return fmt.Sprintf("%s/%s/", baseURL, endpoint)
}

func itemURL(config *Config, resourceType string, id int) string {
	idString := strconv.Itoa(id)
	baseURL := collectionURL(config, resourceType)
	return fmt.Sprintf("%s%s/", baseURL, idString)
}
