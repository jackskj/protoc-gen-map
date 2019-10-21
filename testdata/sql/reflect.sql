{{ define "TypeCasting" }}
select
	1.1  as  double_cast,  -- float64 also testing name mapping
	1.1  as  FloatCast,    -- float32
	1    as  Int32Cast,    -- int32
	1    as  Int64Cast,    -- int64
	1    as  Uint32Cast,   -- uint32
	1    as  Uint64Cast,   -- uint64
	1    as  Sint32Cast,   -- int32
	1    as  Sint64Cast,   -- int64
	1    as  Fixed32Cast,  -- uint32
	1    as  Fixed64Cast,  -- uint64
	1    as  Sfixed32Cast, -- int32
	1    as  Sfixed64Cast, -- int64
	true as  BoolCast,     -- bool
	'1'  as  StringCast,   -- string
        P.created_on as TimestampCast --time.Time
from post P limit 1
{{ end }}

{{ define "IncorrectTypes" }}
select
        {{ .TypeValue }} as GoFloat64,
        {{ .TypeValue }} as GoFloat32,
        {{ .TypeValue }} as GoInt32,
        {{ .TypeValue }} as GoInt64,
        {{ .TypeValue }} as GoUint32,
        {{ .TypeValue }} as GoUint64,
        {{ .TypeValue }} as GoBool,
        {{ .TypeValue }} as GoString,
        {{ .TypeValue }} as GoTimestamp
        {{ if eq .TypeValue "created_on" }}
        from post limit 1
        {{ end }}
{{ end }}
