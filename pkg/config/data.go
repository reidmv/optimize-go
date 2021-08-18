/*
Copyright 2020 GramLabs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/jwt"
)

// NOTE: Configuration JSON names in and below Server and Authorization use snake_case for compatibility with OAuth 2.0 specifications

var (
	// DecodeJWT can be temporarily set to true to make marshalling decode JWT without verification. This is not thread safe.
	DecodeJWT bool
)

// Config is the top level configuration structure for Optimize
type Config struct {
	// Servers is a named list of server configurations
	Servers []NamedServer `json:"servers,omitempty"`
	// Authorizations is a named list of authorizations configurations
	Authorizations []NamedAuthorization `json:"authorizations,omitempty"`
	// Clusters is a named list of cluster configurations
	Clusters []NamedCluster `json:"clusters,omitempty"`
	// Controllers is a named list of controller configurations
	Controllers []NamedController `json:"controllers,omitempty"`
	// Contexts is a named list of context configurations
	Contexts []NamedContext `json:"contexts,omitempty"`
	// CurrentContext is the name of the default context
	CurrentContext string `json:"current-context,omitempty"`
	// Environment identifies the current execution environment
	Environment string `json:"env,omitempty"`
}

// Server contains information about how to communicate with a StormForge API Server
type Server struct {
	// Identifier is a URI used to identify a common set of endpoints making up a StormForge API Server. The identifier
	// may be used to resolve ".well-known" locations, used as an authorization audience, or used as a common base URL
	// when determining default endpoint addresses. The URL must not have any query or fragment components.
	Identifier string `json:"identifier"`
	// API contains the API server metadata necessary to access the programmatic interface.
	API APIServer `json:"api"`
	// Authorization contains the authorization server metadata necessary to access this server
	Authorization AuthorizationServer `json:"authorization"`
	// Application contains information about the public facing user interface.
	Application ApplicationServer `json:"application"`
}

// APIServer is the API server metadata
type APIServer struct {
	// ApplicationsEndpoint is the URL of the applications endpoint
	ApplicationsEndpoint string `json:"applications_endpoint,omitempty"`
	// ExperimentsEndpoint is the URL of the experiments endpoint
	ExperimentsEndpoint string `json:"experiments_endpoint,omitempty"`
	// AccountsEndpoint is the URL of the accounts endpoint
	AccountsEndpoint string `json:"accounts_endpoint,omitempty"`
}

// ApplicationServer is the user facing application.
type ApplicationServer struct {
	// BaseURL is the main entrypoint to the application.
	BaseURL string `json:"base_url,omitempty"`
	// AuthSuccessEndpoint is URL to direct the user to after a successful login.
	AuthSuccessEndpoint string `json:"auth_success_endpoint,omitempty"`
}

// NOTE: AuthorizationServer is defined by https://tools.ietf.org/html/rfc8414 do not add non-standard fields!

// AuthorizationServer is the authorization server metadata
type AuthorizationServer struct {
	// Issuer is the authorization server's identifier, it must be an "https" URL with no query or fragment
	Issuer string `json:"issuer"`
	// AuthorizationEndpoint is the URL of the authorization endpoint
	AuthorizationEndpoint string `json:"authorization_endpoint,omitempty"`
	// TokenEndpoint is the URL of the token endpoint
	TokenEndpoint string `json:"token_endpoint,omitempty"`
	// RevocationEndpoint is the URL of the revocation endpoint
	RevocationEndpoint string `json:"revocation_endpoint,omitempty"`
	// RegistrationEndpoint is the URL of the dynamic client registration endpoint
	RegistrationEndpoint string `json:"registration_endpoint,omitempty"`
	// DeviceAuthorizationEndpoint is the URL of the device flow authorization endpoint
	DeviceAuthorizationEndpoint string `json:"device_authorization_endpoint,omitempty"`
	// JSONWebKeySetURI is URL of the JSON Web Key Set
	JSONWebKeySetURI string `json:"jwks_uri,omitempty"`
}

// Authorization contains information about remote server authorizations
type Authorization struct {
	// Credential is the information that must be presented to prove authorization
	Credential Credential `json:"credential"`
}

// TokenCredential represents a token based credential
type TokenCredential struct {
	// AccessToken is presented to the service being authenticated to
	AccessToken string `json:"access_token"`
	// TokenType is the type of the access token (i.e. "bearer")
	TokenType string `json:"token_type,omitempty"`
	// RefreshToken is presented to the authorization server when the access token expires
	RefreshToken string `json:"refresh_token,omitempty"`
	// Expiry is the time at which the access token expires (or 0 if the token does not expire)
	Expiry time.Time `json:"expiry,omitempty"`
}

// ClientCredential represents a machine-to-machine credential
type ClientCredential struct {
	// ClientID is the client identifier
	ClientID string `json:"client_id"`
	// ClientSecret is the client secret
	ClientSecret string `json:"client_secret"`
	// Scope is the space delimited list of allowable scopes for the client
	Scope string `json:"scope"`
}

// Cluster contains information about communicating with a Kubernetes cluster
type Cluster struct {
	// KubeConfig is the path to a kubeconfig file to use; leave blank to get the default file
	KubeConfig string `json:"kubeconfig,omitempty"`
	// Context is the kubeconfig context to use for the cluster; leave blank to get the current kubeconfig context
	Context string `json:"context"`
	// Namespace overrides the kubeconfig namespace to use for the cluster; leave blank to get the current kubeconfig namespace
	Namespace string `json:"namespace"`
	// Bin is the path to the kubectl binary to use
	Bin string `json:"bin,omitempty"`
	// Controller is the reference to a controller section to use when configuring this cluster
	Controller string `json:"controller,omitempty"`
}

// Controller contains additional controller configuration when working with Optimize on a specific cluster
type Controller struct {
	// DeploymentName is the name of the controller deployment object
	DeploymentName string `json:"deploymentName,omitempty"`
	// Namespace overrides the default namespace to use during configuration
	Namespace string `json:"namespace,omitempty"`
	// RegistrationClientURI is the fully qualified URL of the client configuration endpoint for the controller's client
	RegistrationClientURI string `json:"registration_client_uri,omitempty"`
	// RegistrationAccessToken is the access token to be used for the client configuration endpoint
	RegistrationAccessToken string `json:"registration_access_token,omitempty"`
	// Env defines additional environment variables to load into the controller during authorization
	Env []ControllerEnvVar `json:"env,omitempty"`
}

// ControllerEnvVar is used to specify additional environment variables for a controller during authorization
type ControllerEnvVar struct {
	// Name of the environment variable
	Name string `json:"name"`
	// Value of the environment variable
	Value string `json:"value"`
}

// Context references a remote server...
type Context struct {
	// Server is the name of the remote server to connect to
	Server string `json:"server,omitempty"`
	// Authorization is the name of authorization configuration to use
	Authorization string `json:"authorization,omitempty"`
	// Cluster is the name of the Kubernetes cluster to connect to; it is a name in THIS configuration and does not correspond to the kubeconfig name
	Cluster string `json:"cluster,omitempty"`
}

// NamedServer associates a name to a server configuration
type NamedServer struct {
	// Name is the referencable name for the server
	Name string `json:"name"`
	// Server is the server configuration
	Server Server `json:"server"`
}

// NamedAuthorization associates a name to an authorization configuration
type NamedAuthorization struct {
	// Name is the referencable name for the authorization
	Name string `json:"name"`
	// Authorization is the authorization configuration
	Authorization Authorization `json:"authorization"`
}

// NamedCluster associates a name to cluster configuration
type NamedCluster struct {
	// Name is the referencable name for the cluster
	Name string `json:"name"`
	// Cluster is the cluster configuration
	Cluster Cluster `json:"cluster"`
}

// NamedController associates a name to a controller configuration
type NamedController struct {
	// Name is the referencable name for the controller
	Name string `json:"name"`
	// Controller is the cluster configuration
	Controller Controller `json:"controller"`
}

// NamedContext associates a name to context configuration
type NamedContext struct {
	// Name is the referencable name for the context
	Name string `json:"name"`
	// Context is the context configuration
	Context Context `json:"context"`
}

// Credential is use to represent a credential
type Credential struct {
	// TokenCredential is used to prove authorization using a token that has already been obtained
	*TokenCredential
	// ClientCredential is used to obtain a new token for authorization using the credential information
	*ClientCredential
}

// UnmarshalJSON determines which type of credential is being used
func (c *Credential) UnmarshalJSON(data []byte) error {
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	switch {
	case len(m) == 0:
		return nil
	case m["access_token"] != "":
		c.TokenCredential = &TokenCredential{}
		if err := json.Unmarshal(data, c.TokenCredential); err != nil {
			return nil
		}
	case m["client_id"] != "":
		c.ClientCredential = &ClientCredential{}
		if err := json.Unmarshal(data, c.ClientCredential); err != nil {
			return nil
		}
	default:
		return fmt.Errorf("unknown credential")
	}
	return nil
}

// MarshalJSON ensures token expiry is persisted in UTC
func (c *Credential) MarshalJSON() ([]byte, error) {
	if c.TokenCredential != nil {
		// Override the access token with the decoded JWT claims
		accessToken := interface{}(c.TokenCredential.AccessToken)
		if DecodeJWT {
			if tok, err := jwt.ParseString(c.TokenCredential.AccessToken); err == nil {
				accessToken = tok
			}
		}

		// Override the expiry with the UTC formatted time or "0"
		expiry := "0"
		if !c.TokenCredential.Expiry.IsZero() {
			expiry = c.TokenCredential.Expiry.UTC().Format(time.RFC3339)
		}

		// http://choly.ca/post/go-json-marshalling/
		type TC TokenCredential
		return json.Marshal(&struct {
			*TC
			AccessToken interface{} `json:"access_token,omitempty"`
			Expiry      string      `json:"expiry,omitempty"`
		}{TC: (*TC)(c.TokenCredential), AccessToken: accessToken, Expiry: expiry})
	} else if c.ClientCredential != nil {
		return json.Marshal(c.ClientCredential)
	}
	return []byte("{}"), nil
}

// MarshalJSON omits empty structs
func (srv *Server) MarshalJSON() ([]byte, error) {
	type S Server
	az := &srv.Authorization
	if (AuthorizationServer{}) == srv.Authorization {
		az = nil
	}
	api := &srv.API
	if (APIServer{}) == srv.API {
		api = nil
	}
	app := &srv.Application
	if (ApplicationServer{}) == srv.Application {
		app = nil
	}
	return json.Marshal(&struct {
		*S
		Authorization *AuthorizationServer `json:"authorization,omitempty"`
		API           *APIServer           `json:"api,omitempty"`
		Application   *ApplicationServer   `json:"application,omitempty"`
	}{S: (*S)(srv), Authorization: az, API: api, Application: app})
}
