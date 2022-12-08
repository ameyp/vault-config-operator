/*
Copyright 2021.

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

package v1alpha1

import (
	"context"
	"reflect"

	"github.com/redhat-cop/operator-utils/pkg/util/apis"
	vaultutils "github.com/redhat-cop/vault-config-operator/api/v1alpha1/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DatabaseSecretEngineStaticRoleSpec defines the desired state of DatabaseSecretEngineStaticRole
type DatabaseSecretEngineStaticRoleSpec struct {
	// Authentication is the kube auth configuration to be used to execute this request
	// The authentication role must have the following capabilities = [ "create", "read", "update", "delete"] on that path.
	// +kubebuilder:validation:Required
	Authentication vaultutils.KubeAuthConfiguration `json:"authentication,omitempty"`

	// Path at which to create the role.
	// The final path will be {[spec.namespace]}/{spec.path}/static-roles/{metadata.name}.
	// +kubebuilder:validation:Required
	Path vaultutils.Path `json:"path,omitempty"`

	// Properties associated with the static role
	// +kubebuilder:validation:Required
	DBSESRole `json:",inline"`
}

var _ vaultutils.VaultObject = &DatabaseSecretEngineStaticRole{}

var _ apis.ConditionsAware = &DatabaseSecretEngineStaticRole{}

// DatabaseSecretEngineStaticRoleStatus defines the observed state of DatabaseSecretEngineStaticRole
type DatabaseSecretEngineStaticRoleStatus struct {
	// Important: Run "make" to regenerate code after modifying this file
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DatabaseSecretEngineStaticRole is the Schema for the databasesecretenginestaticroles API
type DatabaseSecretEngineStaticRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatabaseSecretEngineStaticRoleSpec   `json:"spec,omitempty"`
	Status DatabaseSecretEngineStaticRoleStatus `json:"status,omitempty"`
}

func (m *DatabaseSecretEngineStaticRole) GetConditions() []metav1.Condition {
	return m.Status.Conditions
}

func (m *DatabaseSecretEngineStaticRole) SetConditions(conditions []metav1.Condition) {
	m.Status.Conditions = conditions
}

func (d *DatabaseSecretEngineStaticRole) GetPath() string {
	return string(d.Spec.Path) + "/" + "static-roles" + "/" + d.Name
}
func (d *DatabaseSecretEngineStaticRole) GetPayload() map[string]interface{} {
	return d.Spec.toMap()
}
func (d *DatabaseSecretEngineStaticRole) IsEquivalentToDesiredState(payload map[string]interface{}) bool {
	desiredState := d.Spec.DBSESRole.toMap()
	return reflect.DeepEqual(desiredState, payload)
}

func (d *DatabaseSecretEngineStaticRole) IsInitialized() bool {
	return true
}

func (d *DatabaseSecretEngineStaticRole) PrepareInternalValues(context context.Context, object client.Object) error {
	return nil
}

func (r *DatabaseSecretEngineStaticRole) IsValid() (bool, error) {
	// TODO(user): fill in your validation logic upon object deletion.
	return true, nil
}

func (d *DatabaseSecretEngineStaticRole) GetKubeAuthConfiguration() *vaultutils.KubeAuthConfiguration {
	return &d.Spec.Authentication
}

//+kubebuilder:object:root=true

// DatabaseSecretEngineStaticRoleList contains a list of DatabaseSecretEngineStaticRole
type DatabaseSecretEngineStaticRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DatabaseSecretEngineStaticRole `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DatabaseSecretEngineStaticRole{}, &DatabaseSecretEngineStaticRoleList{})
}

type DBSESRole struct {
	// Username Specifies the database username that this Vault role corresponds to.
	// +kubebuilder:validation:Required
	Username string `json:"username,omitEmpty"`

	// RotationPeriod Specifies the amount of time Vault should wait before rotating the password. The minimum is 5 seconds.
	// +kubebuilder:validation:Required
	RotationPeriod metav1.Duration `json:"rotationPeriod,omitempty"`

	// DBName The name of the database connection to use for this role.
	// +kubebuilder:validation:Required
	DBName string `json:"dBName,omitempty"`

	// RotationStatements Specifies the database statements to be executed to rotate the password for the configured database user.Not every plugin type will support this functionality. See the plugin's API page for more information on support and formatting for this parameter.
	// +kubebuilder:validation:Optional
	// +listType=set
	// kubebuilder:validation:UniqueItems=true
	RotationStatements []string `json:"rotationStatements,omitempty"`
}

func (i *DBSESRole) toMap() map[string]interface{} {
	payload := map[string]interface{}{}
	payload["username"] = i.Username
	payload["rotation_period"] = i.RotationPeriod
	payload["db_name"] = i.DBName
	payload["rotation_statements"] = i.RotationStatements
	return payload
}
