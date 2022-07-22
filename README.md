# Microsoft To-Do integration with Taskwarrior

## Install

  1. Register an application on _Microsoft Azure_:
     - Under _Authentication_ set _Allow public client flows_ to `Yes`.
     - Under _API permissions_ add `Tasks.Read`.

  1. Create a `$XDG_CONFIG_HOME/twtodo/credentials.env` file: 
     ```env
     # Tenant ID of the application on Azure. Set the value to 'consumers' if your 
     # Microsoft Account is a personal account.
     TENANT_ID=<tenantID>
     # Client ID of the application on Azure. 
     CLIENT_ID=<clientID>
     ```

  1. `go install github.com/simachri/taskwarrior-ms-todo/cmd/twtodo@latest` 

  
## Usage

### Pull tasks from a to-do list

  1. `twtodo pull -l 'LIST_ID'`
