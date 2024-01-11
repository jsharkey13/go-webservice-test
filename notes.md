# Notes

Some very brief notes on things I found building this proof-of-concept single Isaac endpoint and cookie auth in Go.
This is only ~1 day of work, with quite a lot of internet searching; there are likely approaches that could help with
the issues encountered and questions raised that I just did not find in time.


## Questions and issues

### Error handling

The `func doSomething(...) (ReturnType, error)` pattern seems common in the Go standard library and popular libraries.
In this code it led to a lot of duplicated code for handling errors at every step of a method in exactly the same way
in each place. If the function calls were independent, then `errors.Join(...)` could combine the error messages and 
handling. This does not help when the result of the first call is passed into the second, etc.

It also meant there were lots of instances of `EmptyStruct{}, errors.New(...)`, where the empty struct
seems unnecessary. It leads to the potential for control flow issues if the error is ignored by accident and the empty
struct (which is by necessity a valid value) is mistakenly used. I made this mistake a lot here.

The common `something, err := someMethodCall(...)` sometimes seems to allow `err` to be redefined and sometimes warns
about no new variable being defined.


### Shared database access

This was a useful comparison of database access ideas: https://www.alexedwards.net/blog/organising-database-access.

Matching the Java object-with-own-database-instance pattern seems difficult in Go. Is there a way without using global
variables to pass the database many layers down the call stack without every method requiring a database argument?

This also applies to things like a logging framework or other global properties management.


### Representing optional properties

Our data model makes extensive use of optional properties of objects in Java and TypeScript. Structs don't support this,
at least when using primitive types. There seem to be two approaches; pointers everywhere (and all the associated
referencing and dereferencing *s and &s) or using a null-library (and doing `x.Valid && x.String == "y"` everywhere).


### Code splitting

All files in a Go package seem to need to live in the same directory, and all methods and top-level variables are
globally visible in all other files of that package. Packaging things seems to require manual package indexing inside
the `go.mod` file to ensure the remote version of the local subpackage isn't used.
Go Workspaces might be an option here, but I didn't learn enough to use them in time.


### Automatic JSON marshalling

The automatic JSON marshalling of structs is pretty neat, but very opinionated. It is not possible to exclude `null`
valued keys without also excluding `0`, `""` and `[]`, because Go has no way of distinguishing `null` from these
zero-like values. There are null libraries which do this, and one is tested here.

The other issue encountered was with datetimes; `time.Time` is serialised in a specific format and there is no way to
alter this. One fix is to use a custom type that marshals to JSON in the desired way (in this case Unix milliseconds), 
but this has the downside of then requiring type casts whenever you want to compare it to an ordinary `time.Time` 
object. The other approach would have been to wrap it in a struct that marshals as desired, and then get the `time.Time`
value from the struct each time you need it.
