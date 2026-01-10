package seed

import "github.com/samber/do/v2"

var Packages = do.Package(
	do.Lazy[Seeder](NewSeeder),
)
