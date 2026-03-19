module github.com/my-username/server-todo

go 1.23.0

replace github.com/Rishabhgoswami0/shared-go => ../shared-go

require (
	github.com/Rishabhgoswami0/shared-go v0.0.0
	github.com/joho/godotenv v1.5.1
)

require github.com/lib/pq v1.12.0 // indirect
