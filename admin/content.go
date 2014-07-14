package admin

import (
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/rules"
)

type Content struct {
	Admin    *Admin
	Context  *qor.Context
	Resource *resource.Resource
	Result   interface{}
	Action   string
}

func (content *Content) AllowedMetas(mode rules.PermissionMode) []resource.Meta {
	var attrs []resource.Meta
	switch content.Action {
	case "index":
		attrs = content.Resource.IndexAttrs()
	case "show":
		attrs = content.Resource.ShowAttrs()
	case "edit":
		attrs = content.Resource.EditAttrs()
	case "new":
		attrs = content.Resource.NewAttrs()
	}

	var metas = []resource.Meta{}
	for _, meta := range attrs {
		if meta.HasPermission(mode, content.Context) {
			metas = append(metas, meta)
		}
	}
	return metas
}

func (content *Content) ValueOf(value interface{}, meta resource.Meta) interface{} {
	return meta.GetValue(value, content.Context)
}