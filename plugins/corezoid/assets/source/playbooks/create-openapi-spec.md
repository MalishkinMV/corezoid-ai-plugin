# Creating a process API Wrapper from Template

## 1. General Concept

This playbook describes how to create an API wrapper for any process (so that the process can be called via API). Below are the step-by-step instructions on how to create a wrapper for the process.

The user must specify the following parameters:
- API_URL
- APIGW_URL
- API_LOGIN
- API_USER_ID
- API_TOKEN
- WORKSPACE_ID
- PROC_ID
- PROJECT_ID
- ROOT_FOLDER_ID
- FOLDER_ID
- HTTP_METHOD (optional, default: POST)
- HTTP_PATH (optional, default: generated based on process name)
- LEAN_URL (optional, default empty string "")

All required parameters must be passed when calling the convctl.sh command. Optional parameters can be omitted.

## 2. Workflow
 

### 2.1 Study the Input and Output Parameters of the Process
The process has already been exported previously and is located at `.processes/target_process.json`. No upload is needed.

 - Input parameters of the process are described in the `params` array.
 - All variants of output parameters are described in the `api_rpc_reply` logic. Example 

```json
 {
  "type": "api_rpc_reply",
  "mode": "key_value",
  "res_data": {
    "description": "failed to download",    
    "result": "error"
  },
  "res_data_type": {
    "description": "string",    // here describes the data type of the description field
    "result": "string"  // here describes the data type of the result field
  },
  "throw_exception": true
}
```
### 2.2 Create OpenAPI Specification
You need to create an OpenAPI specification for the process. It describes all input and output parameters of the process. There should be only one request.

**HTTP Method:**
- If the **HTTP_METHOD** parameter is provided, use it (e.g., POST, GET, PUT, DELETE, PATCH).
- If NOT provided, use **POST** by default.

**HTTP Path:**
- If the **HTTP_PATH** parameter is provided, use it (e.g., `/get_user_by_id`, `/users/{id}`).
- If NOT provided, generate the path automatically based on the process name:
  - Convert the process name to snake_case
  - Add a leading slash
  - Example: if process name is "Get User By Id" → generate path `/get_user_by_id`
  - Example: if process name is "CreateOrder" → generate path `/create_order`

Parameters are passed in the request body in JSON format (for POST, PUT, PATCH methods) or as query parameters (for GET, DELETE methods). The response is returned in JSON format.
In the Swagger file, put the descriptions of all entities under `paths` (do not use the `components` object).
For all input and output parameters, add default values (`default` field) based on their data types and descriptions. Use meaningful example values that illustrate the expected format and content.
The specification needs to be placed in the file `.processes/openapi_spec.json` ( version 3.0.0 ). 
Template:
```json
{
  "info": {
    "description": "",
    "title": "API",
    "version": "1.0.0"
  },
  "openapi": "3.0.0",
  "components": {
    "securitySchemes": {
      "BearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    }
  },
  "security": [
    {
      "BearerAuth": []
    }
  ],
  "paths": {}, // here you need to describe the request
  "servers": []
}

```


### 2.3 Creating Corezoid Wrapper Process
You need to create a process using the template `templates/api_wrapper.json`. To do this:
  - Copy the template content to `.processes/api_wrapper.json`
  - give it the same name as the target process, but prefixed with "API:" + title
  - Add the list of input parameters. Here you simply need to copy the `params` content from `target_process.json` to `api_wrapper.json`
  - Replace parameters. Replace the following placeholders from the template with the actual values that were passed in the parameters:
    - `{{WORKSPACE_ID}}` with WORKSPACE_ID value
    - `{{PROJECT_ID}}` with PROJECT_ID value
    - `{{ROOT_FOLDER_ID}}` with ROOT_FOLDER_ID value
    - `{{PROC_ID}}` with PROC_ID value
    - `{{LEAN_URL}}` with LEAN_URL value

**Example:**
```
if
  WORKSPACE_ID=111
  PROJECT_ID=222
  ROOT_FOLDER_ID=333
  PROC_ID=444
  LEAN_URL=https://lean.example.com
then
  {{WORKSPACE_ID}} -> 111
  {{PROJECT_ID}} -> 222
  {{ROOT_FOLDER_ID}} -> 333
  {{PROC_ID}} -> 444
  {{LEAN_URL}} -> https://lean.example.com
```

--- 


### 2.4 Save the Process
Save the modified process to:
```bash
~/repos/corezoid-ai-doc/.processes/api_wrapper.json
```

### 6.2 Test the Process
Test the Process
```bash
API_URL=<API_URL> APIGW_URL=<APIGW_URL> API_TOKEN=<API_TOKEN> WORKSPACE_ID=<WORKSPACE_ID> PROJECT_ID=<PROJECT_ID> ./convctl.sh make-api-wrapper ~/repos/corezoid-ai-doc/.processes/api_wrapper.json ~/repos/corezoid-ai-doc/.processes/openapi_spec.json  <ROOT_FOLDER_ID> <API_LOGIN> <API_USER_ID> <FOLDER_ID>
```

**Example:**
```bash
API_URL=https://admin.corezoid.com APIGW_URL=https://apigw.corezoid.com API_TOKEN=123 WORKSPACE_ID=123 PROJECT_ID=123 ./convctl.sh make-api-wrapper ~/repos/corezoid-ai-doc/.processes/api_wrapper.json ~/repos/corezoid-ai-doc/.processes/openapi_spec.json 123 123 123 123
```
If something is wrong, fix it.
