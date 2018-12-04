<?xml version="1.0" encoding="utf-8"?>
<Globals>{{range $_, $kv := .KV}}{{if $kv.IsClient}}
    <Global {{$kv.Name}}="{{$kv.Value}}" />{{end}}{{end}}
</Globals>