module coke

go 1.16

require (
	github.com/gobuffalo/buffalo v0.17.5
	github.com/gobuffalo/buffalo-pop/v3 v3.0.0
	github.com/gobuffalo/envy v1.10.1
	github.com/gobuffalo/mw-csrf v1.0.0
	github.com/gobuffalo/mw-forcessl v0.0.0-20200131175327-94b2bd771862
	github.com/gobuffalo/mw-i18n/v2 v2.0.0
	github.com/gobuffalo/mw-paramlogger v1.0.0
	github.com/gobuffalo/pop/v6 v6.0.0
	github.com/gobuffalo/suite/v4 v4.0.0
	github.com/markbates/grift v1.5.0
	github.com/unrolled/secure v1.0.9
)

replace (
	github.com/gobuffalo/buffalo v0.17.5 => github.com/fasmat/buffalo v0.16.15-0.20211121195612-46c764b58057
	github.com/gobuffalo/buffalo-pop/v3 v3.0.0 => github.com/fasmat/buffalo-pop/v3 v3.0.0-20211121200722-a7fc8542fca5
	github.com/gobuffalo/mw-i18n/v2 v2.0.0 => github.com/fasmat/mw-i18n/v2 v2.0.0-20211121200253-969c01ee0cdf
	github.com/gobuffalo/suite/v4 v4.0.0 => github.com/fasmat/suite/v4 v4.0.0-20211121142907-1aafd50d9324
)
