env "local" {
    src = "file://schema.hcl"
    url = "postgres://postgres:postgres@localhost:5432/postgres?search_path=public&sslmode=disable"
    dev = "docker://postgres/16/dev?search_path=public"
}