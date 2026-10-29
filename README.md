# Running The test application

## As a Docker Container
After cloning the repo exectute:
```bash
docker build -t rate-limiter-app:test .
```

and then, for instance you can run the following command
```bash
docker run -d --name rl-test -p "127.0.0.1:8080:8000" rate-limiter-app:test
```
This will run the application and will be accessible in http://localhost:8000

## As an executable
