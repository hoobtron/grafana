package timeinterval

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TimeInterval struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              Spec   `json:"spec"`
	Status            Status `json:"status"`
}

func (o *TimeInterval) GetSpec() any {
	return o.Spec
}

func (o *TimeInterval) SetSpec(spec any) error {
	cast, ok := spec.(Spec)
	if !ok {
		return fmt.Errorf("cannot set spec type %#v, not of type Spec", spec)
	}
	o.Spec = cast
	return nil
}

func (o *TimeInterval) GetSubresources() map[string]any {
	return map[string]any{
		"status": o.Status,
	}
}

func (o *TimeInterval) GetSubresource(name string) (any, bool) {
	switch name {
	case "status":
		return o.Status, true
	default:
		return nil, false
	}
}

func (o *TimeInterval) SetSubresource(name string, value any) error {
	switch name {
	case "status":
		cast, ok := value.(Status)
		if !ok {
			return fmt.Errorf("cannot set status type %#v, not of type Status", value)
		}
		o.Status = cast
		return nil
	default:
		return fmt.Errorf("subresource '%s' does not exist", name)
	}
}

func (o *TimeInterval) GetCreatedBy() string {
	if o.ObjectMeta.Annotations == nil {
		o.ObjectMeta.Annotations = make(map[string]string)
	}

	return o.ObjectMeta.Annotations["grafana.com/createdBy"]
}

func (o *TimeInterval) SetCreatedBy(createdBy string) {
	if o.ObjectMeta.Annotations == nil {
		o.ObjectMeta.Annotations = make(map[string]string)
	}

	o.ObjectMeta.Annotations["grafana.com/createdBy"] = createdBy
}

func (o *TimeInterval) GetUpdateTimestamp() time.Time {
	if o.ObjectMeta.Annotations == nil {
		o.ObjectMeta.Annotations = make(map[string]string)
	}

	parsed, _ := time.Parse(time.RFC3339, o.ObjectMeta.Annotations["grafana.com/updateTimestamp"])
	return parsed
}

func (o *TimeInterval) SetUpdateTimestamp(updateTimestamp time.Time) {
	if o.ObjectMeta.Annotations == nil {
		o.ObjectMeta.Annotations = make(map[string]string)
	}

	o.ObjectMeta.Annotations["grafana.com/updateTimestamp"] = updateTimestamp.Format(time.RFC3339)
}

func (o *TimeInterval) GetUpdatedBy() string {
	if o.ObjectMeta.Annotations == nil {
		o.ObjectMeta.Annotations = make(map[string]string)
	}

	return o.ObjectMeta.Annotations["grafana.com/updatedBy"]
}

func (o *TimeInterval) SetUpdatedBy(updatedBy string) {
	if o.ObjectMeta.Annotations == nil {
		o.ObjectMeta.Annotations = make(map[string]string)
	}

	o.ObjectMeta.Annotations["grafana.com/updatedBy"] = updatedBy
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type TimeIntervalList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []TimeInterval `json:"items"`
}
