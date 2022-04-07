/*
Copyright 2021 GramLabs, Inc.

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

package v2

// ApplicationName represents a name token used to identify an application.
type ApplicationName string

func (n ApplicationName) String() string { return string(n) }

// ClusterName represents a name token used to identify a cluster.
type ClusterName string

func (n ClusterName) String() string { return string(n) }
