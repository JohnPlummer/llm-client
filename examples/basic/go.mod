module example

go 1.23.1

require (
	github.com/JohnPlummer/post-scorer v0.0.0
	github.com/JohnPlummer/reddit-client v0.9.0
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/sashabaranov/go-openai v1.37.0 // indirect
	golang.org/x/time v0.5.0 // indirect
)

replace (
	github.com/JohnPlummer/post-scorer => ../../
	github.com/JohnPlummer/reddit-client => ../../../reddit-client
)
