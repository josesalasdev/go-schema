# Go-Schema

Go-Schema is a lightweight data validation library for Go, inspired by [Cerberus from Python](https://github.com/pyeve/cerberus). It allows for easy schema definition and provides clear and detailed validation errors.

## Features

- **Simple yet Powerful Validation**: Define clear and concise schemas for structured data validation.
- **Support for Multiple Data Types**:
  - `string`: With length validation and regex patterns
  - `int`: With numeric ranges (min/max)
  - `float`: With numeric ranges (min/max) 
  - `bool`: For boolean values
  - `list`: For validating collections of elements
  - `map`: For nested data structures
- **Comprehensive Validation Rules**:
  - `Required`: Mandatory fields
  - `Min`/`Max`: Ranges for numeric values
  - `MinLength`/`MaxLength`: String length constraints
  - `Regex`: Pattern validation
  - `Default`: Default values
- **Nested Structure Validation**: Validate lists and maps with complex structures.
- **Customizable Error Messages**: Define specific messages for each error type.
- **Schema Self-Validation**: Verify the integrity of your own schemas.
- **Detailed Error Reporting**: Provides precise information about location and nature of errors.
- **Easy Integration**: Simple interface and clear results for any Go application.

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
        "name": {Type: "string", MinLength: 2, Required: true},
        "age":   {Type: "int", Min: 18, Max: 99, Required: true}
    }
    
    data := map[string]interface{}{
        "name": "A",
        "age":   17
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

