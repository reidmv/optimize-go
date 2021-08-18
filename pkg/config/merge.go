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

// Merge types

// mergeString overwrites s1 with a non-empty value of s2
func mergeString(s1 *string, s2 string) {
	if s2 != "" {
		*s1 = s2
	}
}

// Merge elements

func mergeConfig(c1, c2 *Config) {
	mergeServers(c1, c2.Servers)
	mergeAuthorizations(c1, c2.Authorizations)
	mergeClusters(c1, c2.Clusters)
	mergeControllers(c1, c2.Controllers)
	mergeContexts(c1, c2.Contexts)
	mergeString(&c1.CurrentContext, c2.CurrentContext)
	mergeString(&c1.Environment, c2.Environment)
}

func mergeServer(s1, s2 *Server) {
	mergeString(&s1.Identifier, s2.Identifier)
	mergeString(&s1.API.AccountsEndpoint, s2.API.AccountsEndpoint)
	mergeString(&s1.API.ExperimentsEndpoint, s2.API.ExperimentsEndpoint)
	mergeString(&s1.Authorization.Issuer, s2.Authorization.Issuer)
	mergeString(&s1.Authorization.AuthorizationEndpoint, s2.Authorization.AuthorizationEndpoint)
	mergeString(&s1.Authorization.TokenEndpoint, s2.Authorization.TokenEndpoint)
	mergeString(&s1.Authorization.RevocationEndpoint, s2.Authorization.RevocationEndpoint)
	mergeString(&s1.Authorization.RegistrationEndpoint, s2.Authorization.RegistrationEndpoint)
	mergeString(&s1.Authorization.DeviceAuthorizationEndpoint, s2.Authorization.DeviceAuthorizationEndpoint)
	mergeString(&s1.Authorization.JSONWebKeySetURI, s2.Authorization.JSONWebKeySetURI)
	mergeString(&s1.Application.BaseURL, s2.Application.BaseURL)
	mergeString(&s1.Application.AuthSuccessEndpoint, s2.Application.AuthSuccessEndpoint)
}

func mergeAuthorization(a1, a2 *Authorization) {
	// Do not merge credentials, just shallow copy them wholesale if they are present
	if a2.Credential.TokenCredential != nil && a2.Credential.AccessToken != "" {
		a1.Credential.ClientCredential = nil
		a1.Credential.TokenCredential = new(TokenCredential)
		*a1.Credential.TokenCredential = *a2.Credential.TokenCredential
	}
	if a2.Credential.ClientCredential != nil && a2.Credential.ClientID != "" {
		a1.Credential.TokenCredential = nil
		a1.Credential.ClientCredential = new(ClientCredential)
		*a1.Credential.ClientCredential = *a2.Credential.ClientCredential
	}
}

func mergeCluster(c1, c2 *Cluster) {
	mergeString(&c1.KubeConfig, c2.KubeConfig)
	mergeString(&c1.Context, c2.Context)
	mergeString(&c1.Namespace, c2.Namespace)
	mergeString(&c1.Bin, c2.Bin)
	mergeString(&c1.Controller, c2.Controller)
}

func mergeController(c1, c2 *Controller) {
	mergeString(&c1.Namespace, c2.Namespace)
	mergeString(&c1.RegistrationClientURI, c2.RegistrationClientURI)
	mergeString(&c1.RegistrationAccessToken, c2.RegistrationAccessToken)
	idx := make(map[string]string, len(c2.Env))
	for i := range c2.Env {
		idx[c2.Env[i].Name] = c2.Env[i].Value
	}
	for i := range c1.Env {
		if v := idx[c1.Env[i].Name]; v != "" {
			c1.Env[i].Value = v
			delete(idx, c1.Env[i].Name)
		}
	}
	for k, v := range idx {
		c1.Env = append(c1.Env, ControllerEnvVar{Name: k, Value: v})
	}
}

func mergeContext(c1, c2 *Context) {
	mergeString(&c1.Server, c2.Server)
	mergeString(&c1.Authorization, c2.Authorization)
	mergeString(&c1.Cluster, c2.Cluster)
}

// Merge lists

func mergeServers(data *Config, servers []NamedServer) {
	idx := make(map[string]*Server, len(servers))
	for i := range servers {
		idx[servers[i].Name] = &servers[i].Server
	}
	for i := range data.Servers {
		if svr := idx[data.Servers[i].Name]; svr != nil {
			mergeServer(&data.Servers[i].Server, svr)
			delete(idx, data.Servers[i].Name)
		}
	}
	for k, v := range idx {
		data.Servers = append(data.Servers, NamedServer{Name: k, Server: *v})
	}
}

func mergeAuthorizations(data *Config, authorizations []NamedAuthorization) {
	idx := make(map[string]*Authorization, len(authorizations))
	for i := range authorizations {
		idx[authorizations[i].Name] = &authorizations[i].Authorization
	}
	for i := range data.Authorizations {
		if az := idx[data.Authorizations[i].Name]; az != nil {
			mergeAuthorization(&data.Authorizations[i].Authorization, az)
			delete(idx, data.Authorizations[i].Name)
		}
	}
	for k, v := range idx {
		data.Authorizations = append(data.Authorizations, NamedAuthorization{Name: k, Authorization: *v})
	}
}

func mergeClusters(data *Config, clusters []NamedCluster) {
	idx := make(map[string]*Cluster, len(clusters))
	for i := range clusters {
		idx[clusters[i].Name] = &clusters[i].Cluster
	}
	for i := range data.Clusters {
		if cstr := idx[data.Clusters[i].Name]; cstr != nil {
			mergeCluster(&data.Clusters[i].Cluster, cstr)
			delete(idx, data.Clusters[i].Name)
		}
	}
	for k, v := range idx {
		data.Clusters = append(data.Clusters, NamedCluster{Name: k, Cluster: *v})
	}
}

func mergeControllers(data *Config, controllers []NamedController) {
	idx := make(map[string]*Controller, len(controllers))
	for i := range controllers {
		idx[controllers[i].Name] = &controllers[i].Controller
	}
	for i := range data.Controllers {
		if ctrl := idx[data.Controllers[i].Name]; ctrl != nil {
			mergeController(&data.Controllers[i].Controller, ctrl)
			delete(idx, data.Controllers[i].Name)
		}
	}
	for k, v := range idx {
		data.Controllers = append(data.Controllers, NamedController{Name: k, Controller: *v})
	}
}

func mergeContexts(data *Config, contexts []NamedContext) {
	idx := make(map[string]*Context, len(contexts))
	for i := range contexts {
		idx[contexts[i].Name] = &contexts[i].Context
	}
	for i := range data.Contexts {
		if ctx := idx[data.Contexts[i].Name]; ctx != nil {
			mergeContext(&data.Contexts[i].Context, ctx)
			delete(idx, data.Contexts[i].Name)
		}
	}
	for k, v := range idx {
		data.Contexts = append(data.Contexts, NamedContext{Name: k, Context: *v})
	}
}

// Find elements

func findServer(l []NamedServer, name string) *Server {
	for i := range l {
		if l[i].Name == name {
			return &l[i].Server
		}
	}
	return nil
}

func findAuthorization(l []NamedAuthorization, name string) *Authorization {
	for i := range l {
		if l[i].Name == name {
			return &l[i].Authorization
		}
	}
	return nil
}

func findCluster(l []NamedCluster, name string) *Cluster {
	for i := range l {
		if l[i].Name == name {
			return &l[i].Cluster
		}
	}
	return nil
}

func findController(l []NamedController, name string) *Controller {
	for i := range l {
		if l[i].Name == name {
			return &l[i].Controller
		}
	}
	return nil
}

func findContext(l []NamedContext, name string) *Context {
	for i := range l {
		if l[i].Name == name {
			return &l[i].Context
		}
	}
	return nil
}
