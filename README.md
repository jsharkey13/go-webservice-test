# go-webservice-test

A quick hackathon project to learn some [Go](https://go.dev/), and investigate its use as a replacement for Java for [`isaac-api`](https://github.com/isaacphysics/isaac-api).

The only existing endpoint it offers is `/users/current_user`, where it will parse an existing session cookie and load user information from the database. It does not (quite) support all properties of the Java `RegisteredUserDTO`.

A second endpoint, `/_/experiment`, explores some Go functionality around the user object and handling values etc.

Start with [`main.go`](./main.go), if you want to explore the code.
