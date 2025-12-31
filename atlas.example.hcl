env "local" {
  src = "ent://internal/ent/schema"

  dev = "postgres://YOUR_USERNAME:YOUR_PASSWORD@localhost:13579/setetes?search_path=public&sslmode=disable"
  url = "postgres://YOUR_USERNAME:YOUR_PASSWORD@localhost:13579/setetes?search_path=public&sslmode=disable"

  migration {
    dir = "file://migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
