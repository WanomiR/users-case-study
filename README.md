# Users: case study
Simple REST API with Postgres and Redis.

## Instructions
- To generate swagger docs, build and start the app run:
    ```bash
    make all
    ```
- To build the app without generating docs run:
    ```bash
    make build
    ```
- Use other make commands to conduct database operations.
- Pre-existing users (for tests):
  - `admin@example.com password`
  -	`john.doe@gmail.com password`
  -	`jen.star@gmail.com password`


> [!note]  
> It is important to use `make` as it loads necessary environment variables.
