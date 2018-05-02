package stacks

import (
	"encoding/json"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
	"time"
)
// CreatedStack represents the object extracted from a Create operation.
type CreatedStack struct {
	ID    string             `json:"id"`
	Links []golangsdk.Link `json:"links"`
}

// CreateResult represents the result of a Create operation.
type CreateResult struct {
	golangsdk.Result
}

// Extract returns a pointer to a CreatedStack object and is called after a
// Create operation.
func (r CreateResult) Extract() (*CreatedStack, error) {
	var s struct {
		CreatedStack *CreatedStack `json:"stack"`
	}
	err := r.ExtractInto(&s)
	return s.CreatedStack, err
}

// StackPage is a pagination.Pager that is returned from a call to the List function.
type StackPage struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a ListResult contains no Stacks.
func (r StackPage) IsEmpty() (bool, error) {
	stacks, err := ExtractStacks(r)
	return len(stacks) == 0, err
}

// ListedStack represents an element in the slice extracted from a List operation.
type ListedStack struct {
	CreationTime time.Time        `json:"-"`
	Description  string           `json:"description"`
	ID           string           `json:"id"`
	Links        []golangsdk.Link `json:"links"`
	Name         string           `json:"stack_name"`
	Status       string           `json:"stack_status"`
	StatusReason string           `json:"stack_status_reason"`
	UpdatedTime  time.Time        `json:"-"`
}

func (r *ListedStack) UnmarshalJSON(b []byte) error {
	type tmp ListedStack
	var s struct {
		tmp
		CreationTime golangsdk.JSONRFC3339NoZ `json:"creation_time"`
		UpdatedTime  golangsdk.JSONRFC3339NoZ `json:"updated_time"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = ListedStack(s.tmp)

	r.CreationTime = time.Time(s.CreationTime)
	r.UpdatedTime = time.Time(s.UpdatedTime)

	return nil
}

// ExtractStacks extracts and returns a slice of ListedStack. It is used while iterating
// over a stacks.List call.
func ExtractStacks(r pagination.Page) ([]ListedStack, error) {
	var s struct {
		ListedStacks []ListedStack `json:"stacks"`
	}
	err := (r.(StackPage)).ExtractInto(&s)
	return s.ListedStacks, err
}

// RetrievedStack represents the object extracted from a Get operation.
type RetrievedStack struct {
	Capabilities        []interface{}            `json:"capabilities"`
	CreationTime        time.Time                `json:"-"`
	Description         string                   `json:"description"`
	DisableRollback     bool                     `json:"disable_rollback"`
	ID                  string                   `json:"id"`
	Links               []golangsdk.Link         `json:"links"`
	NotificationTopics  []interface{}            `json:"notification_topics"`
	Outputs             []map[string]interface{} `json:"outputs"`
	Parameters          map[string]string        `json:"parameters"`
	Name                string                   `json:"stack_name"`
	Status              string                   `json:"stack_status"`
	StatusReason        string                   `json:"stack_status_reason"`
	Tags                []string                 `json:"tags"`
	TemplateDescription string                   `json:"template_description"`
	Timeout             int                      `json:"timeout_mins"`
	UpdatedTime         time.Time                `json:"-"`
}

func (r *RetrievedStack) UnmarshalJSON(b []byte) error {
	type tmp RetrievedStack
	var s struct {
		tmp
		CreationTime golangsdk.JSONRFC3339NoZ `json:"creation_time"`
		UpdatedTime  golangsdk.JSONRFC3339NoZ `json:"updated_time"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = RetrievedStack(s.tmp)

	r.CreationTime = time.Time(s.CreationTime)
	r.UpdatedTime = time.Time(s.UpdatedTime)

	return nil
}

// GetResult represents the result of a Get operation.
type GetResult struct {
	golangsdk.Result
}

// Extract returns a pointer to a CreatedStack object and is called after a
// Create operation.
func (r GetResult) Extract() (*RetrievedStack, error) {
	var s struct {
		Stack *RetrievedStack `json:"stack"`
	}
	err := r.ExtractInto(&s)
	return s.Stack, err
}

// UpdateResult represents the result of a Update operation.
type UpdateResult struct {
	golangsdk.ErrResult
}

// DeleteResult represents the result of a Delete operation.
type DeleteResult struct {
	golangsdk.ErrResult
}