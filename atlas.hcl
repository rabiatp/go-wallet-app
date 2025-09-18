env "local_host" {
  url = "postgres://postgres:254798@db:5432/wallet_db?sslmode=disable"
  dev = "postgres://postgres:254798@db:5432/wallet_db?sslmode=disable"
  migration { dir = "file://migrations" }
}
