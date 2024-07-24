# DL(Default Loader)

## Overview

DL (Default Loader) is a tool designed to generate and assign default values to fields within Go structs based on tags.

This utility allows you to specify default values for your struct fields using a simple tag syntax, 
making it easier to initialize structs with predefined values without having to explicitly set them in your code.

---
`defaults.go` is forked from [creasty/defaults](https://github.com/creasty/defaults) and modified to meet the requirements of this project.

## Features

- **Tag-based Default Values**: Use the `default` tag to set default values for struct fields.
- **Easy Integration**: Quickly generate methods to load default values into your structs.
- **Efficient Initialization**: Simplifies the initialization process of complex structs by automatically setting
  default values.

## Installation

To use DL (Default Loader), first install the tool via:

```shell
go install github.com/godcong/dl/cmd@latest
```

## Usage

### Step 1: Define Struct with Default Tags

Add the `default` tag to your struct fields to specify their default values:

```go
// example: demo.go
type Demo struct {
Name string `default:"demo"`
}
```

### Step 2: Generate Default Value Loading Method

Run DL to generate the necessary loading method for your struct:

```
dl -f ./demo.go
```

### Step 3: Load Default Values

In your code, use `dl.Load()` to populate your struct with the default values:

```go
func main() {
demo := &Demo{}
if err := dl.Load(demo); err != nil {
panic(err)
} 
// Now 'demo' has its fields initialized with default values. 
}
```
