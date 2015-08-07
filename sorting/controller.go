package sorting

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/qor/qor"
	"github.com/qor/qor/admin"
	"github.com/qor/qor/roles"
)

func updatePosition(context *admin.Context) {
	if result, err := context.FindOne(); err == nil {
		if position, ok := result.(sortingInterface); ok {
			if pos, err := strconv.Atoi(context.Request.Form.Get("to")); err == nil {
				var count int
				if _, ok := result.(sortingDescInterface); ok {
					context.GetDB().New().Model(modelValue(result)).Count(&count)
					pos = count - pos + 1
				}

				if MoveTo(context.GetDB(), position, pos) == nil {
					var pos = position.GetPosition()
					if _, ok := result.(sortingDescInterface); ok {
						pos = count - pos + 1
					}

					context.Writer.Write([]byte(fmt.Sprintf("%d", pos)))
					return
				}
			}
		}
	}
	context.Writer.Write([]byte("Error"))
	context.Writer.WriteHeader(422)
}

var injected bool

func (s *Sorting) InjectQorAdmin(res *admin.Resource) {
	Admin := res.GetAdmin()
	res.UseTheme("sorting")

	if res.Config.Permission == nil {
		res.Config.Permission = roles.NewPermission()
	}

	if !injected {
		injected = true
		for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
			admin.RegisterViewPath(path.Join(gopath, "src/github.com/qor/qor/sorting/views"))
		}
	}

	role := res.Config.Permission.Role
	if _, ok := role.Get("sorting_mode"); !ok {
		role.Register("sorting_mode", func(req *http.Request, currentUser qor.CurrentUser) bool {
			return req.URL.Query().Get("sorting") != ""
		})
	}

	if res.GetMeta("Position") == nil {
		res.Meta(&admin.Meta{
			Name: "Position",
			Valuer: func(value interface{}, ctx *qor.Context) interface{} {
				db := ctx.GetDB()
				var count int
				var pos = value.(sortingInterface).GetPosition()

				if _, ok := modelValue(value).(sortingDescInterface); ok {
					if total, ok := db.Get("sorting_total_count"); ok {
						count = total.(int)
					} else {
						db.New().Model(modelValue(value)).Count(&count)
					}
					pos = count - pos + 1
				}

				primaryKey := ctx.GetDB().NewScope(value).PrimaryKeyValue()
				url := path.Join(ctx.Request.URL.Path, fmt.Sprintf("%v", primaryKey), "sorting/update_position")
				return template.HTML(fmt.Sprintf("<input type=\"number\" class=\"qor-sorting-position\" value=\"%v\" data-sorting-url=\"%v\" data-position=\"%v\">", pos, url, pos))
			},
			Permission: roles.Allow(roles.Read, "sorting_mode"),
		})
	}

	var attrs []string
	for _, attr := range res.IndexAttrs() {
		if attr != "Position" {
			attrs = append(attrs, attr)
		}
	}
	res.IndexAttrs(append(attrs, "Position")...)

	router := Admin.GetRouter()
	router.Post(fmt.Sprintf("^/%v/\\d+/sorting/update_position$", res.ToParam()), updatePosition)
}