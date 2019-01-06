# sjson

Streaming json parser implementation

## Project Status

I am not sure where i left this two years ago...

## Usage

Once you create the decoder, you can iterate
over the tokens:

```go
r := strings.NewReader(`{"hello":"world"}`)
dec := sjson.NewDecoder(r)

for t, err := dec.Next(); err != nil && t != nil; t, err = dec.Next() {
    // iterate over every token
}
```

## Supported tokens

 * Simple Tokens:
   * NumberType
   * StringType
   * BoolType
   * NullType

 * Complex Tokens:
   * ObjectType
   * ArrayType

 * Special Tokens
   * MemberType
   * EndType

NOTE: whitespace is currently ignored
