![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/hdb-api)
[![Actions Status](https://github.com/tommzn/hdb-api/actions/workflows/go.image.build.yml/badge.svg)](https://github.com/tommzn/hdb-api/actions)

# HomeDashboard API
External services can use this API to publish HomeDashboard data.

## Supported Datasources
- IndoorClimate

## API Description
```swagger:
openapi: 3.0.0
info:
  version: 1.0.0
  title: HomeDashboard API
  description: Provides endpoints for external services to publish events with HomeDashboard data.
  license:
    name: MIT

paths:
  //api/v1/indoorclimate:
    post: 
      summary: Publish new indoor climate data.
      requestBody:
        required: true
        content:
          application/json:
            schema: 
              $ref: '#/components/schemas/IndoorClimateData'
      responses:
        '204':
          description: An event has been published.
        '400':
          description: Invalid indoor climate data passed.
        '500':
          description: Unable to publish event.
    
  /health:
    get: 
      summary: Health status check endpoint.
      responses:
        '204':
          description: Server is still healty.

components:
  schemas:
    IndoorClimateData:
      type: object
      required:
        - deviceid
        - measurementtype
        - value
      properties:
        timestamp:
          description: Timestamp for indoor climate measurement. RFC3339
          type: string
        deviceid:
          description: ID of a device.
          type: string
        measurementtype:
          description: Type of an indoor climate measurement.
          type: string
          enum:
            - temperature
            - humidity
            - battery
        value:
          description: Measurement value.
          type: string
```

# Links
- [HomeDashboard Documentation](https://github.com/tommzn/hdb-docs/wiki)
