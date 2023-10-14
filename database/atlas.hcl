env "local" {
  src = "file://database/schema.hcl"
  url = "postgres://root:root@localhost:5432/app?search_path=public&sslmode=disable"
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir    = "file://database/migrations"
    format = "golang-migrate"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
