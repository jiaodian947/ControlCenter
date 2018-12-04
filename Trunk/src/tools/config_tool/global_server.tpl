<?xml version="1.0" encoding="GB2312"?>
<Globals>{{range $_, $kv := .KV}}{{if $kv.IsServer}}
    <Global {{$kv.Name}}="{{$kv.Value}}" />{{end}}{{end}}
</Globals>