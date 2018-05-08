package stacktemplates

import (
	"github.com/huaweicloud/golangsdk"
)

func getTemplateURL(c *golangsdk.ServiceClient, stackName, stackID string) string {
	return c.ServiceURL("stacks", stackName, stackID, "template")
}


