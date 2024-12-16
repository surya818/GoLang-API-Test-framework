# Candidate Take-Home Exercise - SDET

Welcome to the candidate take-home exercise for the SDET position. This challenge is designed to evaluate
your technical skills, specifically focusing on testing, documentation, and evaluating an existing API.
Below, you'll find all the relevant information to help you get started, including requirements, setup
instructions, and guidelines for submission.

## Table of Contents
1. [Requirements for Candidates](#requirements-for-candidates)
2. [Endpoints and Operations](#endpoints-and-operations)
3. [OpenAPI Specification](#openapi-specification)
4. [How to Run the Project](#how-to-run-the-project)
5. [Configuration Options](#configuration-options)
6. [Dependencies](#dependencies)
7. [Submission Guidelines](#submission-guidelines)
8. [Extra Credit - CI/CD Integration](#extra-credit---cicd-integration)

## Project Overview
This is a service catalog API written in Go, using SQLite for persistence. It provides CRUD operations
for services and their versions. The application is set up to run using Docker, allowing you to easily
start the API and interact with it.

## Requirements for Candidates

1. **Create a README File**:
   - Include instructions on how to execute and configure your **end-to-end tests**.
     - Simplify this by using scripting mechanisms (e.g. bash, Makefile, ...etc).

2. **Write Tests**:
   - Implement tests for all endpoints, including both normal and edge cases.
   - Include **integration** and **end-to-end tests** as appropriate.
   - Use `go test` or an equivalent test framework to execute your tests.

3. **Documentation**:
   - Document your approach to testing, including a comprehensive **test plan** using the README.
   - Summarize your **findings**, including any identified issues.
   - Outline the **steps you would take to address these issues**.
   - Document any **shortcuts** taken due to time constraints and the rationale behind them.

4. **Extra Credit (Optional)**:
   - Integrate CI/CD tools (e.g., GitHub Actions, Jenkins) to automate the testing process on each
     commit.

## Endpoints and Operations

### Services Endpoints
1. **`POST /v1/token`** - Generate a JWT token for authentication.
2. **`GET /v1/services`** - List all services.
3. **`POST /v1/services`** - Create a new service.
4. **`GET /v1/services/{serviceId}`** - Retrieve details of a specific service.
5. **`PATCH /v1/services/{serviceId}`** - Update an existing service.
6. **`DELETE /v1/services/{serviceId}`** - Delete a specific service.

### Service Versions Endpoints
1. **`GET /v1/services/{serviceId}/versions`** - List all versions of a specific service.
2. **`POST /v1/services/{serviceId}/versions`** - Create a new version for a specific service.
3. **`GET /v1/services/{serviceId}/versions/{versionId}`** - Retrieve details of a specific service
   version.
4. **`PATCH /v1/services/{serviceId}/versions/{versionId}`** - Update an existing version.
5. **`DELETE /v1/services/{serviceId}/versions/{versionId}`** - Delete a specific service version.

## OpenAPI Specification
The OpenAPI Specification for this project is located at
[./docs/openapi.html](./docs/openapi.html). It provides a detailed description of all the API
endpoints, request/response structures, and authentication requirements.

## How to Run the Project

To run the project using Docker:

1. **Run the Docker Container**:
   ```
   make docker-run
   ```
   This will build and start the application on port **18080**. You can interact with the API at
   `http://localhost:18080`.

### Example `curl` Commands

1. **Generate a JWT Token**:
   ```sh
   curl -X POST http://localhost:18080/v1/token \
     -H "Content-Type: application/json" \
     -d '{"username": "kong", "password": "onward"}'
   ```
   Expected Response:
   ```json
   {
     "token": "<JWT_TOKEN>"
   }
   ```

2. **List All Services** (Requires the Token from Step 1):
   ```sh
   curl -X GET http://localhost:18080/v1/services \
     -H "Authorization: Bearer <JWT_TOKEN>"
   ```
   Replace `<JWT_TOKEN>` with the token obtained in the previous step.

## Configuration Options

The project uses a `config.yml` file with the following default options:

- **JWT Secret** (`jwt_secret`): Default is `"kong"`.
  - **Note**: This is used to sign JWT tokens.
- **Username** (`username`): Default is `"kong"`.
- **Password** (`password`): Default is `"onward"`.
- **JWT Token Timeout** (`jwt_token_timeout`): Default is `5m`.
- **Request Timeout** (`request_timeout`): Default is `5s`.

These can be configured by updating the `config.yml` file.

## Dependencies

To complete this exercise, you will need the following:

1. **Docker**: To build and run the Docker container.
2. **Go**: To run tests and work with Go modules.

## Submission Guidelines

1. **Create a Private GitHub Repository**:
   - Create a private GitHub repository and push your changes.
   - Add the GitHub handles provided by the recruiter as collaborators to give them access.

2. **Submit Your Solution**:
   - Include a README file with your approach, test plan, and findings.
   - Your submission should include tests and any CI/CD integrations for extra credit.

---

We hope you enjoy the challenge and look forward to reviewing your work. Good luck!
# kong-takehometask
