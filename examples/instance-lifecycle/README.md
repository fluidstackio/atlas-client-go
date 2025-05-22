# Instance Lifecycle Example

This example demonstrates how to use `atlas-client-go` to create instances with the Atlas API.

## Usage

This example requires the following environment variables to be set:

| Variable              | Description                                                              |
| --------------------- | ------------------------------------------------------------------------ |
| `ATLAS_CLIENT_ID`     | A valid OAuth client ID for the Atlas API.                               |
| `ATLAS_CLIENT_SECRET` | A valid OAuth client secret for the Atlas API.                           |
| `ATLAS_REGION_URL`    | The URL of an Atlas region. Example: https://sandbox.atlas.fluidstack.io |
| `ATLAS_PROJECT_ID`    | UUID of an existing Atlas project in this region.                        |

Run:

```
ATLAS_CLIENT_ID=<client-id> ATLAS_CLIENT_SECRET=<client-secret> ATLAS_REGION_URL=<region> ATLAS_PROJECT_ID=<project-id> go run main.go
```
