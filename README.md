#  End-to-end proxy authentication service to unauthenticated resources
The service allows authenticating requests with JWT tokens supported by the authentication service (https://github.com/trotill/auth-service)

Use cases:
- adding native JWT authentication to a non-secured (open) service, e.g. grafana with authentication disabled
- reducing requests to the authentication service
- parallel, multi-threaded authentication of requests,

Docker image of the service takes about 25MB and is in request idle mode, it does not consume memory and CPU time.

# How to use
```
docker pull monkeyhouse1/auth_proxy:<version>, e.g. docker pull monkeyhouse1/auth_proxy:0.1.0
```

Copy .env.example to .env, specify your settings, put the file in docker compose or docker.  
If the .env file is missing, the application will take the variables from the environment variables.

Example of running under docker
```
docker run -p 9080:9080 --pid host --privileged --env-file ./.env --name "auth_proxy" monkeyhouse1/auth_proxy:0.1.0
```

Example of use in docker-compose
```
   auth_proxy:
     container_name: auth_proxy
     image: monkeyhouse1/auth_proxy:0.1.0
     env_file:
       - .env
     expose:
       - "9080"
     restart: unless-stopped
```
OR
```
  auth_proxy:
     container_name: auth_proxy
     image: monkeyhouse1/auth_proxy:0.1.0
     privileged: true
     pid: host
     environment:
        PORT: 9180
        TARGET_URL: http://127.0.0.1:3000
        ACCESS_TOKEN_NAME: access
        DB_PATH: /app/db/auth.db
        PUBLIC_KEY_PATH: /app/db/public.key
        ROLE_GUEST_BLOCK: 1
        ROLE_OPERATOR_BLOCK: 1
        ROLE_ADMIN_BLOCK: 0

     expose:
       - "9080"
     restart: unless-stopped
```

# .env file options
PORT. Port number (9080 default)  
```
PORT=9080
```

TARGET_URL. URL of the application to which requests will be proxied
```
TARGET_URL=http://grafana:3000
```

ACCESS_TOKEN_NAME. Access token name
```
ACCESS_TOKEN_NAME=access
```

DB_PATH. Database path
```
DB_PATH=/app/db/auth.db
```

PUBLIC_KEY_PATH. Path to the public key
```
PUBLIC_KEY_PATH=/opt/scrips/
```

ROLE_GUEST_BLOCK. Block requests to a user with the role - guest
```
ROLE_GUEST_BLOCK=1
```

ROLE_OPERATOR_BLOCK.  Block requests to the user with the role - operator 
```
ROLE_OPERATOR_BLOCK=1
```

ROLE_ADMIN_BLOCK.  Block requests to a user with the admin role
```
ROLE_ADMIN_BLOCK=1
```