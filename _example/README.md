## sjsontest
 
This is an example of how a JSON file can be
iterated and processed.

The file is opened and `sjson.ReadAll` sends
the json tokens into the channel. `sjson.ReadAll`
blocks until it is finished so we do this in a goroutine:

```go
	f, _ := os.Open(os.Args[1])
	ch := make(chan sjson.Token)
	go func() {
		defer close(ch)
		err = sjson.ReadAll(f, ch)
		if err != nil {
			panic(err)
		}
	}()
```

A for select loop is used to read each token and print
its details, using an offset for whitespace. 

```go
	var offset int
	for {
		select {
		case t, more := <-ch:
			if !more {
				return
			}
```

We decrease the offset on an EndType before printing so that
the start token and the end tokens line up:

```go
			if t.Type() == sjson.EndType {
				offset--
			}
```

We print our whitespace followed by the item. Each token type should
implement a debug-useful `String()` method:

```go
			for i := 0; i <= offset; i++ {
				fmt.Printf(" ")
			}
			fmt.Printf("%s\n", t)
```

We increase the offset on JSON objects and arrays. The next tokens, until
`EndType`, will be the members of this current token and need to be 
offset:

```go
			if t.Type() == sjson.ObjectType || t.Type() == sjson.ArrayType {
				offset++
			}
		}
	}
```

