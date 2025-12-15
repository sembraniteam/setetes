package handler

import "github.com/samber/do/v2"

var Packages = do.Package(
	do.Lazy[Account](NewAccount),
)
