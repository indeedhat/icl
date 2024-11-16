# ICL: Indeedhats config language

Inspired by hcl and nginx's config languages this is pretty much just a stripped down version that does the data format
part without any of the complex logic

It attempts to be somewhat tollerent of invalid illegal tokens where possible

## Literals
### String
```hcl
"string literal"
'string literal'
```

### Integer
```hcl
102836 # basic ints only atm
```

### Float
```hcl
102836.112
```

### Boolean
```hcl
# bools are case sensetive
true
false
```

### Null
```hcl
null # null is also case sensetive
```

## Constructs

### assignment
```hcl
myvar = "thing"
```

### array
```hcl
[1, 2, 3]
```

### map
```hcl
{
    # identifier key
    key1: "value",
    # string key
    "key 2": "value"
}
```

### block
```hcl
# blocks can have multiple values assigned as such, values can only be string literals
my_block "with value" {

}
```

### environment variable interpolation
```hcl
my_key = env(HOME)
```

### comment
```hcl
# comments start with a hash
# they are single line
```

## Under Consideration
### variables
```hcl
# assignment
$var = "value"

# Usage
my_key = $var
```


## Known issues
- [ ] parser is probably too tolerant of issues
- [ ] slices do not support pointers
