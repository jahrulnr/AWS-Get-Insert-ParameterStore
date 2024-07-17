# AWS Parameter Store Management

This package provides a set of functions to manage parameters in the AWS Parameter Store. It includes functionality to retrieve parameters, generate a new list of parameters, and insert parameters into the store.

## Package Structure

The package is divided into three files:

1. **config.go**: Defines the data structures used throughout the package, including `InsertPayload` and `GetResponse`.
2. **core.go**: Implements the core functionality of the package, including retrieving parameters from the AWS Parameter Store, generating a new list of parameters, and inserting parameters into the store.
3. **main.go**: Provides the entry point for the package, allowing users to interact with the package through command-line arguments.

## Functions

### GetParameterStore()

Retrieves parameters from the AWS Parameter Store and stores them in a JSON file. This function is used to fetch parameters from the store and update the local cache.

### GenerateList()

Generates a new list of parameters by modifying the existing parameters in the local cache. This function is used to transform the parameters for a different environment or use case.

### InsertParameterStore()

Inserts parameters into the AWS Parameter Store. This function is used to update the parameters in the store with new values or to add new parameters.

### getPayloadParameterStore()

Retrieves parameters from a JSON file and transforms them into a slice of `InsertPayload` structures. This function is used by `GenerateList()` and `InsertParameterStore()` to work with the parameters.

### getParameterStore()

Reads the contents of a JSON file containing parameters. This function is used by `getPayloadParameterStore()` to read the parameters from a file.

## Usage

To use this package, you can run the following commands:

* `go run main.go getlist`: Retrieves parameters from the AWS Parameter Store and stores them in a JSON file.
* `go run main.go generatelist`: Generates a new list of parameters based on the existing parameters in the local cache.
* `go run main.go insertlist`: Inserts parameters into the AWS Parameter Store.

Note: Ensure you have the necessary AWS credentials set up on your system to use this package.
