# Go-Schema

Go-Schema is a lightweight data validation library for Go, inspired by [Cerberus from Python](https://github.com/pyeve/cerberus). It allows for easy schema definition and provides clear and detailed validation errors.

## Features
- Simple schema definition for validation.
- Supports data types: `string`, `int`, `bool`, `list`, `map`.
- Validation rules such as `required`, `min`, `max`, `regex`, `allowed_values`.
- Detailed error handling.
- Extensible with custom validation rules.

## Installation
```sh
 go get github.com/josesalasdev/go-schema/validator
```

## Usage
### Define a Schema
```go
package main

import (
    "fmt"
    "github.com/josesalasdev/go-schema/validator"
)

func main() {
    schema := validator.Schema{
        "name": {Type: "string", MinLength: intPtr(2), Required: true},
        "age":   {Type: "int", Min: intPtr(18), Max: intPtr(99), Required: true},
        "email":  {Type: "string", Regex: strPtr(`^[^@]+@[^@]+\.[^@]+$`), Required: true},
    }
    
    data := map[string]interface{}{
        "name": "A",
        "age":   17,
        "email":  "invalid_email",
    }
    
    result := validator.Validate(data, schema)
    
    if !result.IsValid {
        fmt.Println("Validation errors:", result.Errors)
    }
}
```

## Testing
```sh
go test ./...
```

## License
This project is licensed under the MIT License.

