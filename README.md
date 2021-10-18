# Task App with etcd store and Github OAuth 

## Build & Run

**Build**
``` bash
make build
```

**Run**

1. Create a .env file and set values for the following environment variables.
```
GITHUB_CLIENT_ID=<YOUR_GITHUB_CLIENT_ID>
GITHUB_CLIENT_SECRET=<YOUR_GITHUB_CLIENT_SECRET>
JWT_SECRET_KEY=<JWT_SECRET_KEY>
```

2. Run the service
``` bash 
make run
```
