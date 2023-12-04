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
	codec := cmdc.NewCodec()

	transport, err := tcp.NewServerTransport("localhost:8080")
	if err != nil {
		panic(err)
	}
	defer transport.Close()

	server := &mdd.Server{
		Codec:     codec,
		Transport: transport,
	}

	server.MessageHandler(func(containers *mdd.Containers) *mdd.Containers {
        ... 
		repsonse := mdd.Containers{}
        return response
	})

	err := transport.Listen()
	if err != nil {
		panic(err)
	}

```


## Client
```go
	codec := cmdc.NewCodec()

	transport, err := tcp.NewClientTransport("localhost:8080")
	if err != nil {
		panic(err)
	}
	defer transport.Close()

	client := &mdd.Client{
		Codec:     codec,
		Transport: transport,
	}

	request := mdd.Containers{}

	response, err := client.SendMessage(&request)
	if err != nil {
		panic(err)
	}
```