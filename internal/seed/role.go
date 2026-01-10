package seed

func (s *seedBuilder) Role() {
	tx, err := s.client.Tx(s.ctx)
	if err != nil {
		panic(err)
	}

	role, err := tx.Role.Create().
		SetName("Donor").
		SetKey("donor").
		SetActivated(true).
		SetDomain("region:*").
		SetDescription("General blood donor role with limited access to donation features.").
		Save(s.ctx)
	if err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			panic(rerr)
		}

		panic(err)
	}

	permission, err := tx.Permission.Create().
		SetName("Get self profile").
		SetKey("get-self-profile").
		SetDomain("*").
		SetDescription("Allow donor to view their own profile details.").
		SetResource("/account/v1/self").SetAction("GET").Save(s.ctx)
	if err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			panic(rerr)
		}

		panic(err)
	}

	if err = tx.Commit(); err != nil {
		panic(err)
	}

	if err = s.rbac.AddPolicy(
		role.Key,
		permission.Domain,
		permission.Resource,
		permission.Action,
	); err != nil {
		panic(err)
	}
}
