package resources

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/lib/pq"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"

	"github.com/Depado/articles/code/qor/models"
)

// AddProduct adds the product model to the admin interface
func AddProduct(adm *admin.Admin) {
	p := adm.AddResource(&models.Product{})
	p.Meta(&admin.Meta{
		Name: "Tags",
		Type: "string",
		Valuer: func(record interface{}, context *qor.Context) (result interface{}) { // nolint: unparam
			s := record.(*models.Product)
			var out string
			for _, t := range s.Tags {
				if out == "" {
					out = t
					continue
				}
				out = fmt.Sprintf("%s, %s", out, t)
			}
			return out
		},
		Setter: func(record interface{}, metaValue *resource.MetaValue, context *qor.Context) { // nolint: unparam
			s := record.(*models.Product)
			values := metaValue.Value.([]string)
			s.Tags = pq.StringArray{}
			if len(values) > 0 && values[0] != "" {
				tags := strings.Split(values[0], ",")
				for _, t := range tags {
					s.Tags = append(s.Tags, strings.TrimSpace(t))
				}
			}
		},
	})

	p.IndexAttrs("ID", "Name", "Price", "Tag")
	p.EditAttrs("Name", "Price", "Tags")
	p.NewAttrs("Name", "Price", "Tags")
	p.Meta(&admin.Meta{
		Name: "Tag",
		Valuer: func(record interface{}, context *qor.Context) interface{} { // nolint: unparam
			p := record.(*models.Product)
			var out string
			var tot int
			badge := `<span style="background-color: rgb(33,150,243); padding: 5px; border-radius: 25%%; color: white; margin-right: 5px;">%s</span>`
			for _, t := range p.Tags {
				if tot+len(t) >= 30 {
					more := fmt.Sprintf(badge, "…")
					out = fmt.Sprintf("%s%s", out, more)
					break
				}
				tot += len(t)
				tt := fmt.Sprintf(badge, t)
				if out == "" {
					out = tt
					continue
				}
				out = fmt.Sprintf("%s%s", out, tt)
			}
			return template.HTML(out) // nolint: gosec
		},
	})
}
