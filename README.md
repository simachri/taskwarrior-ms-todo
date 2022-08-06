# Microsoft To-Do integration with Taskwarrior

## Install

  1. Register an application on _Microsoft Azure_:
     - Under _Authentication_ set _Allow public client flows_ to `Yes`.
     - Under _API permissions_ add `Tasks.Read`.

  1. Create a `$XDG_CONFIG_HOME/twtodo/credentials.yaml` file: 
     ```yaml
     # Tenant ID of the application on Azure. Set the value to 'consumers' if your 
     # Microsoft Account is a personal account.
     tenant_id: <tenantID>
     # Client ID of the application on Azure. 
     client_id: <clientID>
     ```

  1. Create a `$XDG_CONFIG_HOME/twtodo/config.yaml` file: 
     ```yaml
     server:
       port: 41001
     ```

  1. `go install github.com/simachri/taskwarrior-ms-todo/cmd/twtodo@latest` 

  1. The _CLI tool_ `grep` needs to be installed and available on path.
  

## Usage

### Start server

  Start the server that authenticates to Microsoft Azure and handles the commands from 
  the client:
  ```
  twtodo up
  ```

### Client: Pull tasks from a To-Do list

  When the server is started, execute from another terminal session:
  ```
  twtodo pull -l 'LIST_ID'
  ```
