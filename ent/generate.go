package ent

// go generate bu dosyayı ent/ klasörü baz alarak çalıştırır.
// Bu yüzden burada YOL: ./schema olmalı (./ent/schema değil!)

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/lock ./schema
