module example

go 1.23.1

require (
	github.com/JohnPlummer/llm-client v0.0.0
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_golang v1.20.5 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/sashabaranov/go-openai v1.37.0 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	github.com/sony/gobreaker/v2 v2.0.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)

replace (
	github.com/JohnPlummer/llm-client => ../../
	github.com/JohnPlummer/reddit-client => ../../../reddit-client
)
