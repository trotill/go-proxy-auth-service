# Docker Black Hole (DBH)
The service provides a secure way for services in the docker-compose environment to communicate with the external environment. The service has its own REST API and swagger file, with the help of which neighbouring services can execute scripts and applications outside the docker-compose environment.  

Common application examples:
- rebooting or shutting down a VDS
- obtaining complete information about the host system,
- manipulating the graphical environment with utilities,
- configuring the network on the host system and firewall,
- manipulating docker with the docker command

Docker image of the service takes about 20MB and is in request idle mode, it does not consume memory and CPU time.

# How to use
```
docker pull monkeyhouse1/black_hole:<version>, e.g. docker pull monkeyhouse1/black_hole:0.1.0
```

Copy .env.example to .env, specify your settings, put the file in docker compose or docker.  
If the .env file is missing, the application will take the variables from the environment variables.

Example of running under docker
```
docker run -p 9080:9080 --pid host --privileged --env-file ./.env --name "blh-service" monkeyhouse1/black_hole:0.1.0
```

Example of use in docker-compose
```
   blh-service:
     container_name: blh-service
     image: monkeyhouse1/black_hole:0.1.0
     privileged: true
     pid: host
     env_file:
       - .env
     expose:
       - "9080"
     restart: unless-stopped
```
OR
```
   blh-service:
     container_name: blh-service
     image: monkeyhouse1/black_hole:0.1.0
     privileged: true
     pid: host
     environment:
       PORT=9080
       DOCKER=0
       EXECUTE_MAX_TIMEOUT_SEC=600
       SCRIPT_PATH=/home/ivan/work/golang/docker-black-hole/scripts/
       ALLOW_ABSOLUTE_MODE=1
       EXECUTE_FROM_USER=ivan
       SHELL_PATH=/usr/bin/sh
     expose:
       - "9080"
     restart: unless-stopped
```
# Safety
Safety is ensured by the following:
- only one service is running in privileged mode, which is minimalistic and has only a few endpoints,
- with the SCRIPT_PATH and ALLOW_ABSOLUTE_MODE=false options you can restrict the directory with scripts,
- using the EXECUTE_FROM_USER option you can specify the user on behalf of which commands will be executed in the external environment

# SWAGGER
There is a swagger for easy integration, use it to evaluate the service and generate code

# .env file options
PORT. Port number (9080 default)  
```
PORT=9080
```

DOCKER. 1 - the service is run in a docker-compose environment, 0 - the service is run as a separate application in the host system  (1 default)
```
DOCKER=1
```

EXECUTE_FROM_USER. User in the host system on behalf of which applications/scripts will be run  (root default)
```
EXECUTE_FROM_USER=develinux
```

EXECUTE_MAX_TIMEOUT_SEC. Max. time (in sec.) of application/script execution. The setting protects against hangs of the running application/script. (600 default)
```
EXECUTE_MAX_TIMEOUT_SEC=600
```

SCRIPT_PATH. Path to scripts. If ALLOW_ABSOLUTE_MODE option is disabled, the service can execute scripts only from the SCRIPT_PATH folder (scripts default)   
```
SCRIPT_PATH=/opt/scrips/
```

ALLOW_ABSOLUTE_MODE. Allows execution of programmes/scripts from any folder of the host system. If this option is disabled, execution is allowed only from the SCRIPT_PATH folder (0 default)
```
ALLOW_ABSOLUTE_MODE=1
```

DISABLE_LOGS. Will disable all logs (0 default).    
```
DISABLE_LOGS=0
```

SHELL_PATH. Shell path, bash for default (bash default)   
```
SHELL_PATH=/usr/bin/sh
```

# REST

The app provides 2 endpoints for interaction.  
POST -> Job to run the task (RunJob)  
GET -> Job/{ID} to get the state of the task (GetJob).

The endpoints fields are described in detail in the swagger file.

Interaction logic.  
A RunJob execution request is sent, then GetJob is executed at intervals until the state is set to finish or error  
As soon as the finish status is received, the result is analysed