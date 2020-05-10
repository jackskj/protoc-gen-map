package mapper

import (
	"bytes"
	"encoding/base64"
	"errors"
	// "log"
	"regexp"
	"strconv"
)

var paramRegexp *regexp.Regexp

var dialectVars = map[string]string{
	"postgres":         "$",
	"cloudsqlpostgres": "$",
	"sqlite3":          "?",
	"common":           "?",
	"mssql":            "?",
	"mysql":            "?",
}

func init() {
	// matches "map_param_ followed by base64 encoded string"
	re := "map_param_(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?"
	paramRegexp = regexp.MustCompile(re)
}

func PrepareQuery(dialect string, rawSql []byte) (string, []interface{}, error) {

	var (
		preparedSql *bytes.Buffer
		arg64       []byte
		argBuff     []byte
		sqlPart     []byte
		err         error
	)

	if dialect == "" {
		return "", nil, errors.New("Parameterized query detected, but dialect is unknown. " +
			"Specify DB connection along with dialect name, (postgres, mysql, mssql, sqlite3) in the gprc service server.")
	} else if _, found := dialectVars[dialect]; found == false {
		return "", nil, errors.New("Parameterized queries are not supported for specified dialect: " + dialect +
			". Known dialects include postgres, mysql, mssql, and sqlite3.")
	}

	// param values in plain bytes,
	// it is an empty interface as db/sql requires args to be interface{}
	args := []interface{}{}

	if paramLoc := paramRegexp.FindAllIndex(rawSql, -1); paramLoc != nil {
		// copy over the sql from the begginning to the first parameter
		preparedSql = &bytes.Buffer{}

		// variable used to copy over sql query and while ommitting parametarized values
		lastLoc := 0

		// deepcopy, prevents changes to the original rawSql
		rawSql = append(make([]byte, 0, len(rawSql)), rawSql...)
		for i, loc := range paramLoc {
			arg64 = rawSql[loc[0]+10 : loc[1]] // +10 removes "map_param_"
			argBuff = make([]byte, len(arg64))
			n, err := base64.StdEncoding.Decode(argBuff, arg64)
			if err != nil {
				return "", nil, err
			}
			args = append(args, string(argBuff[0:n]))
			sqlPart = append(rawSql[lastLoc:loc[0]], bindVar(dialect, i+1)...)
			_, err = preparedSql.Write(sqlPart)
			if err != nil {
				return "", nil, err
			}
			lastLoc = loc[1]
		}
		sqlPart = rawSql[lastLoc:len(rawSql)]
		_, err = preparedSql.Write(sqlPart)
		if err != nil {
			return "", nil, err
		}
	} else {
		preparedSql = bytes.NewBuffer(rawSql)
	}
	return preparedSql.String(), args, nil
}

func bindVar(dialect string, i int) []byte {
	if varSymbol, _ := dialectVars[dialect]; varSymbol == "$" {
		return []byte("$" + strconv.FormatInt(int64(i), 10))
	} else if varSymbol, _ := dialectVars[dialect]; varSymbol == "?" {
		return []byte("?")
	}
	// this should not happen
	return []byte("?")
}
