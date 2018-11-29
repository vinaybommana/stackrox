package generator

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/proto"
)

const schemaTemplate = `
schema {
	query: Query
}
{{range $td := .Entries}}
{{if isEnum .Data.Type -}}
enum {{.Data.Name}} {
{{- range $k, $v := enumValues .Data}}
	{{$k -}}
{{end}}
}
{{else -}}
type {{.Data.Name}} {
{{- range .Data.FieldData -}}
{{$t := schemaType .}}{{if $t }}
	{{lower .Name}}: {{ $t -}}
{{end -}}
{{end -}}
{{range .Data.UnionData }}
	{{lower .Name}}: {{$td.Data.Name}}{{.Name -}}
{{end -}}
{{range .ExtraResolvers }}
	{{ . -}}
{{end}}
}
{{ range $u := .Data.UnionData}}
union {{$td.Data.Name}}{{$u.Name}} = 
{{- range $i, $f := $u.Entries}}{{if $i}} |{{end}} {{schemaType $f}}{{end}}
{{ end }}
{{end}}{{end}}

type Label {
	key: String!
	value: String!
}

scalar Time
`

type schemaEntry struct {
	Data           typeData
	ListData       map[string]bool
	ExtraResolvers []string
}

func isListType(p reflect.Type) bool {
	if p == nil {
		return false
	}
	name := p.Name()
	if p.Kind() == reflect.Ptr {
		name = p.Elem().Name()
	}
	return isProto(p) && len(name) > 4 && name[0:4] == "List"
}

func makeSchemaEntries(data []typeData, extraResolvers map[string][]string) []schemaEntry {
	output := make([]schemaEntry, 0)
	listRef := make(map[string]map[string]bool)
	for _, td := range data {
		if isListType(td.Type) {
			fm := make(map[string]bool)
			for _, f := range td.FieldData {
				fm[f.Name] = true
			}
			listRef[td.Name[4:]] = fm
		}
	}

	for _, td := range data {
		if (td.Name == "Query" || isProto(td.Type) || isEnum(td.Type)) && !isListType(td.Type) {
			ldt, ok := listRef[td.Name]
			ldp := ldt
			if !ok {
				ldp = nil
			}
			se := schemaEntry{
				Data:           td,
				ListData:       ldp,
				ExtraResolvers: extraResolvers[td.Name],
			}
			output = append(output, se)
		}
	}
	return output
}

func schemaType(fd fieldData) string {
	if strings.HasPrefix(fd.Name, "XXX_") {
		return ""
	}
	if fd.Name == "Id" && fd.Type.Kind() == reflect.String {
		return "ID!"
	}
	return schemaExpand(fd.Type)
}

func schemaExpand(p reflect.Type) string {
	switch p.Kind() {
	case reflect.String:
		return "String!"
	case reflect.Int32:
		if isEnum(p) {
			return p.Name() + "!"
		}
		return "Int!"
	case reflect.Int64:
		return "Int!"
	case reflect.Uint32:
		return "Int!"
	case reflect.Float32:
		return "Float!"
	case reflect.Float64:
		return "Float!"
	case reflect.Bool:
		return "Boolean!"
	case reflect.Slice:
		inner := schemaExpand(p.Elem())
		if inner != "" {
			return fmt.Sprintf("[%s]!", inner)
		}
		return ""
	case reflect.Map:
		if p.Elem().Kind() == reflect.String &&
			p.Elem().Kind() == reflect.String {
			return "[Label!]!"
		}
	case reflect.Ptr:
		if p == timestampType {
			return "Time!"
		}
		if isProto(p) {
			return p.Elem().Name()
		}
		inner := schemaExpand(p.Elem())
		if inner == "" {
			return ""
		}
		if strings.HasSuffix(inner, "!") {
			return inner[:len(inner)-1]
		}
		return inner
	}
	return ""
}

func enumValues(data typeData) map[string]int32 {
	return proto.EnumValueMap(importedName(data.Type))
}

// GenerateSchema produces a valid GraphQL schema based on the results of a type walk and including
// the extra resolvers provided. A Query object is automatically added, so entry point resolver methods
// should be exposed on that name in extraResolvers.
func GenerateSchema(parameters TypeWalkParameters, extraResolvers map[string][]string) string {
	walkResults := typeWalk(
		parameters.IncludedTypes,
		[]reflect.Type{
			reflect.TypeOf((*types.Timestamp)(nil)),
		},
	)
	data := make([]typeData, 1, len(walkResults)+1)
	data[0] = typeData{
		Name:      "Query",
		Type:      nil,
		FieldData: nil,
		Package:   "",
	}
	data = append(data, walkResults...)
	buf := &bytes.Buffer{}
	schemaEntries := makeSchemaEntries(data, extraResolvers)
	t, err := template.New("schema").Funcs(template.FuncMap{
		"lower":      lower,
		"plural":     plural,
		"schemaType": schemaType,
		"isEnum":     isEnum,
		"enumValues": enumValues,
	}).Parse(schemaTemplate)
	if err != nil {
		panic(err)
	}
	err = t.Execute(buf, struct{ Entries []schemaEntry }{schemaEntries})
	if err != nil {
		panic(err)
	}
	return buf.String()
}
