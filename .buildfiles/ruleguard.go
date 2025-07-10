package buildfiles

import "github.com/quasilyte/go-ruleguard/dsl"

func logMessageUppercase(m dsl.Matcher) {
	m.Match(`$_.Log().$_($_).$_($s, $*_)`,
		`$_.Log().$_($_).$_($_).$_($s, $*_)`,
		`$_.Log().$_($s)`).
		Where(m["s"].Text.Matches(`^\"[a-z].*`)).
		Report(`$s message should start with uppercase letter`)
}
