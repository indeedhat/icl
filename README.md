# ICL: Indeedhat's config language

Inspired by hcl and nginx's config languages this is pretty much just a stripped down version that does the data format
part without any of the complex logic

ICL is designed to be a configuration data language not a data interchange language

It attempts to be somewhat tollerent of invalid illegal tokens where possible

## Syntax
<table>
    <tr>
        <th>Type</th>
        <th>Example</th>
    </tr>
    <tr>
        <td>String</td>
        <td>

```hcl
"string literal"
'string literal'
```
</td>
    </tr>
    <tr>
        <td>Integer</td>
        <td>

```hcl
102836
```
</td>
    </tr>
    <tr>
        <td>Float</td>
        <td>

```hcl
15.3
```
</td>
    </tr>
    <tr>
        <td>Boolean</td>
        <td>

```hcl
true
false
```
</td>
    </tr>
    <tr>
        <td>Null</td>
        <td>

```hcl
null
```
</td>
    </tr>
    <tr>
        <td>Comemnts</td>
        <td>

```hcl
# comments start with a hash
# they are single line
```
</td>
    </tr>
    <tr>
        <td>Assignment</td>
        <td>

```hcl
mavar = "value"
```
</td>
    </tr>
    <tr>
        <td>Slice</td>
        <td>

```hcl
[1, 2, 3]
```
</td>
    </tr>
    <tr>
        <td>Map</td>
        <td>

```hcl
{
    # identifier key
    key1: "value",
    # string key
    "key 2": "value"
}
```
</td>
    </tr>
    <tr>
        <td>Block</td>
        <td>

```hcl
# blocks can have multiple values assigned as such, values can only be string literals
my_block "with value" "another value" {
    inner_var = "some value"
    inner_block val1 {
        has_data = true
    }
    inner_block val2 {
        has_data = false
    }
}
```
</td>
    </tr>
    <tr>
        <td>Environment variables</td>
        <td>

```hcl
# the env() macro is only available on primitives, they cannot be used in slices or maps
my_key = env(HOME)
```
</td>
    </tr>
</table>

## Supported types
- string
- bool
- int, int8, int16, int32, int64
- uint, uint8, uint16, uint32, uint64
- float32, float64
- slice
- map (string keys only)
- struct

## Marshaling data
Data can be marshaled directly from a struct into an ICL document

- Only structs can be marshaled
- ICL is opt-in, only fields marked with an icl tag will be included in the document
- there is no support for marshaling comments, any comments in an existing file will be lost

```go
type MyConfig struct {
    Version    int `icl:"version"`
    MyVar      string `icl:"my_var"`
    Unexported string
}

c := MyConfig{
    Version: 2,
    MyVar: "data",
    Unexported: "Other data",
}

fmt.Printf(icl.MarshalString(c))

/* Output:
version = 2
my_var = "data"
*/
```

### Marshaling envars
In cases where you want certain fields to be filled via environment variables ICL provides the env(ENVAR_KEY) macro.  
in order to maintain this macro when marshaling into an ICL document you must suffix the struct tag with the macro.
```go
type MyConfig struct {
    Version    int `icl:"version"`
    MyEnvar    string `icl:"my_envar,env(MY_ENVAR)"`
}

c := MyConfig{
    Version: 2,
    MyEnvar: "data",
}

fmt.Printf(icl.MarshalString(c))

/* Output:
version = 2
my_var = env(MY_ENVAR)
*/
```

## Marshaling data
Data can be unmalhaled directly into a struct from an ICL document

the same rule apply as unmarshaling
```go
type MyConfig struct {
    Version    int `icl:"version"`
    MyVar      string `icl:"my_var"`
    Unexported string
}

document = `
version = 2
my_var = "data"
`

c := MyConfig{Unexported: "private"}

_ = icl.UnmarshalString(document, &c)

/* Output:
c.Version == 2
c.MyVar == "data"
c.Unexported == "private"
*/

```

## ICL struct tags
- "my_var" the icl struct tag is used to define the identifier for a variable/block in the ICL document
- "my_float.2" the /.\n/ suffix is used to define the precision level of a float when marshaled into an ICL document
- ".param" is used to define a field as a param on its parent block, params will get marshaled/unmarshaled in the order they appear
- "my_key,env(ENVAR_KEY)" the `env(ENVAR_KEY)` macro tells the encoder to set the variable value to be a env macro when building the ICL document

## Version assignment
ICL provides some versioning support out of the box, this is accomplished via the `version assigment`

The version assignment is as follows:
```hcl
version = 1
```

- It MUST be the first non comment, non whitespace node in the file.
- It MUST be an integer value

### Checking the version number
```go
document = `
# this is an ICL document
version = 1

my_var = "data"
`
a, _ := icl.Parse(document)
fmt.Print(a.Version())
```
> should no version assignment be found as the first node Version() will return 0

### Marshaling Version number
There is no special support for marshaling version numbers, it is however recommended that you set the first
exported struct field as the version in your root struct:
```go
type MyConfig struct {
    Version    int `icl:"version"`
    MyVar      string `icl:"my_var"`
    Unexported string
}

c := MyConfig{
    Version: 2,
    MyVar: "data",
    Unexported: "Other data",
}

_ = icl.MarshalFile(c, "my/save/path.icl")
```

### Unmarshaling versions
ICL provides a helper method `UnmarshalVersion(data []byte, versions map[int]any)` for handling multi version
unmarshaling

```go
type MyConfigV1 struct {
    Version    int `icl:"version"`
    ...
}
type MyConfigV2 struct {
    Version    int `icl:"version"`
    ...
}

document = `
version = 2
...
`

foundVersion, config, err := icl.UnmarshalVersion(
    []byte(document),
    map[int]any{
        1: &MyConfigV1{},
        2: &MyConfigV2{},
    }
)

// foundVersion == 2
// config == *MyConfigV2
// err == nil
```

## LImitations
- Slices do not support map values

## Known issues
- [ ] parser is probably too tolerant of issues
- [ ] error messages still need some work

## Under consideration
### Marshalling/Unmarshaling comments
```go
document = `
version = 1
...

# These Comments will get
# Unmarshaled into the structs PreMyVar1 field
my_var_1 = "data"

# these comments will be ignored
my_var_2 = "data"

# This comment will get Unmarshaled into the structs PostMyVar2 field
`

type config struct {
    Version int `icl:"version"`
    ...
    PreMyVar1 []string `icl:".comments"`
    MyVar1 string `icl:"my_var_1"`
    MyVar2 string `icl:"my_var_2"`
    PostMyVar2 []string `icl:".comments"`
}

type configAlt struct {
    Version int `icl:"version"`
    ...
    // the key being the line number the comment resides on
    Comments map[int]string `icl:".comments"`
    MyVar1 string `icl:"my_var_1"`
    MyVar2 string `icl:"my_var_2"`
}
```

### default values
Allow setting a default value in the struct tag for if there is no assignemt in the ICL document

```go
document = `
version = 1
my_var_1 = "data"
`

type config struct {
    Version int `icl:"version"`
    MyVar1 string `icl:"my_var_1"`
    MyVar2 string `icl:"my_var_2,default=some default value"`
}
```
