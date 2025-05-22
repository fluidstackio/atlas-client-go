# Instance Lifecycle Example

This example demonstrates how to use `atlas-client-go` to create instances with the Atlas API.

## Usage

This example requires the following environment variables to be set:

| Variable           | Description                                                              |
| ------------------ | ------------------------------------------------------------------------ |
| `ATLAS_API_TOKEN`  | A valid JWT for the Atlas API.                                           |
| `ATLAS_REGION_URL` | The URL of an Atlas region. Example: https://sandbox.atlas.fluidstack.io |
| `ATLAS_PROJECT_ID` | UUID of an existing Atlas project in this region.                        |

Run:

```
ATLAS_API_TOKEN=<token> ATLAS_REGION_URL=<region> ATLAS_PROJECT_ID=<project-id> go run main.go
```
