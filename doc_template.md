
{{range $index,$api := .}}
#### {{$index}} {{ $api.Api.Title }}
{{$api.Api.Description}}

````
{{$api.Api.Method}} {{$api.Api.Route}}
````
**请求说明**

|参数名称|参数类型|取值范围|必要性|参数位置|默认值|描述|
|-------|-------|------|-----|-------|-----|----|{{ range $_,$f := $api.Req }}
|{{$f.Field}}|{{$f.Type}}|{{$f.Enum}}|{{$f.Required}}|{{$f.Location}}|{{$f.Default}}|{{$f.Description}}|{{end}}

**请求URL示例**

```` 
{{$api.ReqRequestLineExample}}
````
{{if $api.ReqBodyExample}}
**请求header示例**

````
{{ $api.ReqHeaderExample }}
````
{{end}}

{{if $api.ReqBodyExample}}
**请求body示例**

````
{{ $api.ReqBodyExample }}
````
{{end}}
**响应示例**

````
{{$api.ResExample}}
````
**响应说明**
|参数名称|参数类型|取值范围|描述|
|-------|-------|------|----|{{ range $_,$f := $api.Res }}
|{{$f.Field}}|{{$f.Type}}|{{$f.Enum}}|{{$f.Description}}|{{end}}

{{end}}