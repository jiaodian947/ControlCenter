// create by tools.
// source file: {{.Path}}
#ifndef __{{.Name}}_wrapper_h__
#define __{{.Name}}_wrapper_h__
#include "game_common.h"

class T{{setName .Name}} : public BaseWrapper
{
public:
	DECLARE_GAMEOBJECT_INHERIT_CLASS(T{{setName .Name}}, BaseWrapper);
public:
    {{range .Propertys}}
    //{{.Desc}}
    Decleration_Property({{.Name}}, {{getType .Type}}){{end}}
    {{if .Childs}}  
    {{range $k, $v := .Childs}}
    //{{$v.Path}}{{range $v.Propertys}}
    //{{.Desc}}
    Decleration_Property({{.Name}}, {{getType .Type}}){{end}}
    //{{$v.Path}} end
    {{end}}
    {{end}}
};

{{range .Records}}
//{{.Desc}}
class TR{{setName .Name}} : public BaseRecord
{
public:
    DECLARE_GAMEOBJECT_INHERIT_CLASS1(TR{{setName .Name}}, BaseRecord, {{.Name}});
public:
    {{range $k, $v := .ColTypes}}
    //{{.Desc}}
    Decleration_Column({{$k}}, {{$v.Name}}, {{getType $v.Type}}){{end}}
};
{{end}}
{{if .Childs}}{{range $k, $v := .Childs}}{{if $v.Records}}//{{$v.Path}}{{range $v.Records}}
//{{.Desc}}
class TR{{setName .Name}} : public BaseRecord
{
public:
    DECLARE_GAMEOBJECT_INHERIT_CLASS1(TR{{setName .Name}}, BaseRecord, {{.Name}});
public:
    {{range $k, $v := .ColTypes}}
    //{{.Desc}}
    Decleration_Column({{$k}}, {{$v.Name}}, {{getType $v.Type}}){{end}}
};
{{end}}{{end}}{{end}}{{end}}
#endif
