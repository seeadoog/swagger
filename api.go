package swagger

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io"
	"path"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Schema struct {
	Type        string             `json:"type"`
	Properties  map[string]*Schema `json:"properties,omitempty"`
	Enum        []string           `json:"enum,omitempty"`
	MaxLength   *int               `json:"maxLength,omitempty"`
	Items       *Schema            `json:"items,omitempty"`
	Location    string             `json:"location,omitempty"`
	Description string             `json:"description,omitempty"`
	Example     string             `json:"example,omitempty"`
	Default     string             `json:"default,omitempty"`
	Required    bool               `json:"required,omitempty"`
	Binding     string             `json:"binding,omitempty"`
}

type Api struct {
	Title          string                          `json:"title"`
	Request        any                             `json:"-"`
	Response       any                             `json:"-"`
	Route          string                          `json:"route"`
	Method         string                          `json:"method"`
	Description    string                          `json:"description"`
	RequestSchema  *Schema                         `json:"request_schema,omitempty"`
	ResponseSchema *Schema                         `json:"response_schema,omitempty"`
	ErrHandler     func(c *gin.Context, err error) `json:"-"`
}

type ApiGroup struct {
	apis []*Api

	vad *validator.Validate
}

func (a *ApiGroup) initValidator() {
	if a.vad != nil {
		return
	}
	a.vad = validator.New()
	a.vad.SetTagName("binding")
	a.vad.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return getFieldName(fld)
	})
}

func (a *ApiGroup) Add(api *Api) {
	a.apis = append(a.apis, api)
}

func (a *ApiGroup) GenerateMarkdown() string {
	return GenerateMarkdown(a.apis)
}

func (a *ApiGroup) GenerateHtml() string {
	return GenerateHtml(a.apis)
}
func boolOfStr(s string) bool {
	r, _ := strconv.ParseBool(s)
	return r
}
func (a *ApiGroup) schemaHandler(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		if boolOfStr(c.Query("get_schema")) {
			c.AbortWithStatusJSON(200, api)
			return
		}
		c.Next()
	}
}

func (a *ApiGroup) HandlerDocumentMd() gin.HandlerFunc {
	markdown := a.GenerateMarkdown()
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.Writer.WriteHeader(200)
		c.Writer.WriteString(markdown)
	}
}

func (a *ApiGroup) HandlerDocumentHtml() gin.HandlerFunc {
	markdown := a.GenerateHtml()
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		c.Writer.WriteHeader(200)
		c.Writer.WriteString(markdown)
	}
}

type OptFunc func(o *Api)

func WithDescription(desc string) OptFunc {
	return func(o *Api) {
		o.Description = desc
	}
}

func WithTitle(title string) OptFunc {
	return func(o *Api) {
		o.Title = title
	}
}

func WithErrHandler(f func(c *gin.Context, err error)) OptFunc {
	return func(o *Api) {
		o.ErrHandler = f
	}
}

func NewAPIGroup() *ApiGroup {
	a := &ApiGroup{}
	a.initValidator()
	return a
}

type BasicRouter interface {
	BasePath() string
	gin.IRouter
}

func RegisterAPI[Req, Resp any](r *ApiGroup, router BasicRouter, method, pth string, handler Handler[Req, Resp], opts ...OptFunc) {
	rsc := generateSchema(reflect.ValueOf(new(Req)), "")
	a := &Api{
		Request:  new(Req),
		Response: new(Resp),

		Method:         method,
		RequestSchema:  rsc,
		ResponseSchema: generateSchema(reflect.ValueOf(new(Resp)), ""),
	}
	for _, opt := range opts {
		opt(a)
	}
	rsc.Description = a.Description

	router.Handle(method, pth, r.schemaHandler(a), WrapHandler[Req, Resp](r, handler, a.ErrHandler))
	a.Route = path.Join(router.BasePath(), pth)
	r.apis = append(r.apis, a)
}

func bindRequest(r *ApiGroup, ctx *gin.Context, req any) error {

	bytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	if len(bytes) > 0 {
		err = json.Unmarshal(bytes, req)
		if err != nil {
			return fmt.Errorf("json unmarshal error: %w", err)
		}
	}

	err = bindPath(ctx, reflect.ValueOf(req))
	if err != nil {
		return err
	}
	vad, ok := req.(Validator)
	if ok {
		for _, f := range vad.Validate(ValidateCtx{}) {
			err = f()
			if err != nil {
				return err
			}
		}
	}
	return r.vad.Struct(req)
}

func getFieldName(f reflect.StructField) string {
	tag := f.Tag.Get("location")
	_, name, _ := strings.Cut(tag, ",")
	if name == "" {
		name = f.Tag.Get("json")
	}
	if name == "" {
		name = f.Name
	}
	return name
}

func bindPath(ctx *gin.Context, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Ptr:
		return bindPath(ctx, v.Elem())
	case reflect.Struct:
		t := v.Type()
	conn:
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("location")

			fv := v.Field(i)
			if field.Anonymous {
				err := bindPath(ctx, fv)
				if err != nil {
					return err
				}
			}
			location, name, _ := strings.Cut(tag, ",")

			if name == "" {
				name = field.Tag.Get("json")
			}
			if name == "" {
				name = field.Name
			}
			var val string
			switch location {
			case "path":
				val = ctx.Param(name)
			case "query":
				val = ctx.Query(name)
			case "", "json":
				if (fv.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct) || field.Type.Kind() == reflect.Struct {
					err := bindPath(ctx, fv)
					if err != nil {
						return fmt.Errorf("bind field: %w", err)
					}
					continue conn
				}
				// body 中数据只有为指针类型时且为nil，才会走下去，去设置默认值
				if fv.Kind() == reflect.Ptr && fv.IsNil() {

				} else {
					continue conn
				}
			default:
				return fmt.Errorf("inner error: invalid location: %s", location)
			}
			if val == "" {
				def := field.Tag.Get("default")
				if def != "" {
					val = def
				}
			}
			if val == "" {
				continue
			}
			err := bindValue(ctx, fv, val)
			if err != nil {
				return fmt.Errorf("bind value: %w :%v", err, val)
			}
		}
	}
	return nil
}

func bindValue(ctx *gin.Context, v reflect.Value, str string) error {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			elem := reflect.New(v.Type().Elem())
			v.Set(elem)
			return bindValue(ctx, elem, str)
		}
		return bindValue(ctx, v.Elem(), str)
	case reflect.String:
		v.SetString(str)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		i, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}
		v.SetFloat(f)
	case reflect.Bool:
		b, err := strconv.ParseBool(str)
		if err != nil {
			return err
		}
		v.SetBool(b)

	case reflect.Struct:

	default:
		return fmt.Errorf("unsupported type in path binding: %s", v.Type())
	}
	return nil
}

func getEnumsFromTag(tag reflect.StructTag) []string {
	enus := tag.Get("enum")
	if enus == "" {
		return nil
	}
	return strings.Split(enus, ",")
}
func getLocationFromTag(tag reflect.StructTag) string {
	location := tag.Get("location")
	if location == "" {
		return "json"
	}
	lo, _, _ := strings.Cut(location, ",")
	return lo
}

func getNameFromTag(f reflect.StructField) string {
	tag := f.Tag
	location := tag.Get("location")
	if location == "" {
		v := tag.Get("json")
		if v == "" {
			v = f.Name
		}
		return v
	}
	_, name, _ := strings.Cut(location, ",")
	if name == "" {
		name = f.Name
	}
	return name
}
func generateSchema(v reflect.Value, tags reflect.StructTag) *Schema {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return generateSchema(reflect.New(v.Type().Elem()), tags)
		}
		return generateSchema(v.Elem(), tags)
	case reflect.Struct:
		t := v.Type()

		sc := &Schema{
			Type:       "object",
			Properties: map[string]*Schema{},
		}
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if field.Anonymous {
				jst := field.Tag.Get("json")
				inline := jst == "" || strings.Contains(jst, ",inline")
				if inline {
					for key, schema := range generateSchema(v.Field(i), field.Tag).Properties {
						sc.Properties[key] = schema
					}
					continue
				}

			}
			jsonTag := getNameFromTag(field)
			fsc := generateSchema(v.Field(i), field.Tag)
			fsc.Location = getLocationFromTag(field.Tag)
			fsc.Description = field.Tag.Get("desc")
			fsc.Example = field.Tag.Get("example")
			fsc.Default = field.Tag.Get("default")

			fsc.Required, _ = strconv.ParseBool(field.Tag.Get("required"))
			bind, ok := field.Tag.Lookup("binding")
			if ok {
				if strings.Contains(bind, ",required") {
					fsc.Required = true
				}
			}
			fsc.Binding = bind
			sc.Properties[jsonTag] = fsc

		}
		return sc
	case reflect.Slice:
		sc := Schema{
			Type:  "array",
			Items: generateSchema(reflect.New(v.Type().Elem()).Elem(), tags),
		}
		return &sc
	case reflect.Map:
		return &Schema{
			Type: "object",
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return &Schema{
			Type: "integer",
			Enum: getEnumsFromTag(tags),
		}
	case reflect.Float32, reflect.Float64:
		return &Schema{
			Type: "number",
			Enum: getEnumsFromTag(tags),
		}
	case reflect.Bool:
		return &Schema{
			Type: "boolean",
		}
	case reflect.String:
		return &Schema{
			Type: "string",
			Enum: getEnumsFromTag(tags),
		}
	}
	panic(fmt.Errorf("unsupported type when gen schema: %s", v.Type()))
	return nil
}

func (s *Schema) Doc() []*FiledDoc {
	docs := s.toDoc("")

	sort.Slice(docs, func(i, j int) bool {
		if docs[i].Location == "path" && docs[j].Location == "path" {
			return docs[i].Field < docs[j].Field
		}
		if docs[i].Location == "path" {
			return true
		}
		if docs[j].Location == "path" {
			return false
		}
		if docs[i].Location == "query" && docs[j].Location == "query" {
			return docs[i].Field < docs[j].Field
		}
		if docs[i].Location == "query" {
			return true
		}
		if docs[j].Location == "query" {
			return false
		}
		return docs[i].Field < docs[j].Field
	})
	return docs
}

func (s *Schema) toDoc(path string) []*FiledDoc {
	docs := []*FiledDoc{}
	prefix := path
	if prefix != "" {
		prefix = path + "."
	}

	if s.Type == "object" && (strings.HasSuffix(path, "[]") || path == "") {

	} else {
		docs = append(docs, &FiledDoc{
			Field:       path,
			Type:        s.Type,
			Enum:        strings.Join(s.Enum, ","),
			Required:    s.Required,
			Description: s.Description,
			Location:    s.Location,
			Default:     s.Default,
			Binding:     s.Binding,
		})
	}

	for name, schema := range s.Properties {
		docs = append(docs, schema.toDoc(prefix+name)...)
	}
	if s.Items != nil {
		docs = append(docs, s.Items.toDoc(path+"[]")...)
	}
	return docs
}

func (s *Schema) getExample() string {
	if s.Example != "" {
		return s.Example
	}
	if s.Default != "" {
		return s.Default
	}

	if len(s.Enum) > 0 {
		return s.Enum[0]
	}
	return ""
}

func (s *Schema) GenExampleJson() string {
	bs, _ := json.MarshalIndent(s.GenExample(), "", "   ")
	return string(bs)
}

func (s *Schema) genExampleQuery() []string {
	res := []string{}
	for name, schema := range s.Properties {
		if schema.Location == "query" {
			ex := schema.getExample()
			if ex != "" && ex != "-" {
				res = append(res, name+"="+ex)

			}
		}
	}
	sort.Strings(res)
	return res
}

func (s *Schema) GenExample() any {

	switch s.Type {
	case "object":
		m := make(map[string]any)
		for name, schema := range s.Properties {
			if schema.Location == "path" || schema.Location == "query" {
				continue
			}
			exp := schema.GenExample()
			es, ok := exp.(string)
			if ok && es == "-" {
			} else {
				m[name] = exp
			}
		}
		return m
	case "integer":
		exm := s.getExample()
		if len(exm) > 0 {
			v, _ := strconv.Atoi(exm)
			return v
		}
		return 0
	case "string":
		return s.getExample()
	case "array":
		if s.Items != nil {

			if s.getExample() != "" {
				return formatByType(s.Type, s.getExample(), s.Items.Type)
			}
			exp := s.Items.GenExample()
			return []any{exp}
		}
	case "number":
		exm := s.getExample()
		if len(exm) > 0 {
			v, _ := strconv.ParseFloat(exm, 64)
			return v
		}
		return 0
	case "boolean":
		exm := s.getExample()
		if len(exm) > 0 {
			v, _ := strconv.ParseBool(exm)
			return v
		}
		return false
	}
	return nil
}

func (s *Schema) generateExamplePath(path string) string {

	params := parsePathParams(path)
	values := make([]string, len(params))
	for i, param := range params {
		s.Walk(func(name string, node *Schema) bool {
			if name == param[1:] {
				exp := node.getExample()
				if exp == "" {
					exp = name
				}
				values[i] = exp
				return false
			}
			return true
		})
	}

	rps := make([]string, 0, len(values)+len(params))
	for i, value := range values {
		rps = append(rps, params[i], value)
	}

	return strings.NewReplacer(rps...).Replace(path)
}

func (s *Schema) Walk(f func(name string, node *Schema) bool) {
	for name, schema := range s.Properties {
		if !f(name, schema) {
			return
		}
		schema.Walk(f)
	}
	if s.Items != nil {
		s.Items.Walk(f)
	}
}

// /abc/:key
func parsePathParams(path string) []string {
	tkn := []byte{}
	tokens := []string{}
	state := 0
	for i := 0; i < len(path); i++ {
		c := path[i]
		switch state {
		case 0:
			switch c {
			case ':', '*':
				tkn = append(tkn, c)
				state = 1
			}
		case 1:
			switch c {
			case '/':
				state = 0
				tokens = append(tokens, string(tkn))
				tkn = tkn[:0]
			default:
				tkn = append(tkn, c)
			}
		}
	}
	if len(tkn) > 0 && state == 1 {
		tokens = append(tokens, string(tkn))
	}
	return tokens
}

func formatByType(typ string, str string, sbtype string) any {
	switch typ {
	case "string":
		return str
	case "integer":
		i, _ := strconv.Atoi(str)
		return i
	case "number":
		f, _ := strconv.ParseFloat(str, 64)
		return f
	case "boolean":
		b, _ := strconv.ParseBool(str)
		return b
	case "array":
		dst := []any{}
		for _, s := range strings.Split(str, ",") {
			dst = append(dst, formatByType(sbtype, s, ""))
		}
		return dst
	}
	return nil
}
