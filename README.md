# Matrixx Domain Data

## Installation
To install the MDD library, add the following import statement to your Go code:

```go
import "github.com/matrixxsoftware/go-mdd"
```

Then, run the following command to fetch the library:
```bash
go get -u github.com/matrixxsoftware/go-mdd@latest
```

# Example

## Server
```go
	server, err := mdd.NewServer("localhost:8080", codec)
	if err != nil {
		panic(err)
	}
	defer server.Close()

	server.Handler(func(containers *mdd.Containers) *mdd.Containers {
        ... 
        // return response containers
	})

```


## Client
```go
	client, err := mdd.NewClient("localhost:8080", codec)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	request := mdd.Containers{}

	response, err := client.SendMessage(&request)
	if err != nil {
		panic(err)
	}
```