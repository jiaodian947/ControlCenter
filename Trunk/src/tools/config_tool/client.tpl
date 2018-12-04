<?xml version="1.0" encoding="utf-8"?>{{$ps := .Propertys}}
<{{.Name}}sConfig>
    <{{.Name}}s>{{range $_, $obj := .Objects}}
        <{{$.Name}}{{range $_, $val := $obj.Values}} {{$t := index $ps $val.Index}}{{if $t.IsClient}}{{$t.Name}}="{{$val.Value}}"{{end}}{{end}} />{{end}}     
    </{{.Name}}s>
</{{.Name}}sConfig>