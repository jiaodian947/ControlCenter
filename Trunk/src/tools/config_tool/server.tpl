<?xml version="1.0" encoding="GB2312"?>{{$ps := .Propertys}}
<Object>{{range $_, $obj := .Objects}}
    <Property{{range $_, $val := $obj.Values}} {{$t := index $ps $val.Index}}{{if $t.IsServer}}{{$t.Name}}="{{$val.Value}}"{{end}}{{end}} />{{end}}     
</Object>