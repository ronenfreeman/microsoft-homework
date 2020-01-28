# Home Work

The following repo contains the solution to a Microsoft home work problem.

## Problem Definition

1. Cluster details:
    a. Use Azure and AKS-Engine  (On GitHub not pre-defined managed AKS solutions)
    b. Set up K8s 1.15+ cluster, with RBAC enabled
    c. Cluster should have 2 micro-services – A and B
    d. Cluster should have Ingress controller, redirecting traffic by URL: xxx/serviceA or xxx/serviceB
    e. ServiceA should not be able to talk with ServiceB (policy disabled by RBAC). Network policy.
    f. For Service A:write a script\application which retrieves the bitcoin value in dollar from an API on the web (you should find one), every minute and prints it, Every 10 minutes it should print the average value of the last 10 minutes.

2. General Guideline
    a. Please, consider this as process for setting up “production-ready” cluster by all meaning, 1 from 100 that will be created later
    b. Please, share cluster templates and yaml files as GitHub repo / zip file

## Application

The application, given three environment variables `BITCOIN_ENDPOINT`, `MINUTES_TO_SLEEP` and `MINUTES_TO_GET_AVERAGE` gets the rate of Bitcoin from `BITCOIN_ENDPOINT` every `MINUTES_TO_SLEEP` minute(s) and averages the rate every MINUTES_TO_GET_AVERAGE minutes.

### Assumptions

### Build Locally

The application itself can be tested locally by running the following command in the root directory.
```
export BITCOIN_ENDPOINT=https://api.coindesk.com/v1/bpi/currentprice/usd.json
export MINUTES_TO_SLEEP=1
export MINUTES_TO_GET_AVERAGE=10

go run cmd/main.go
```
To Build the Docker image run:

`docker build -t <repo-owner>/<repo-name> .` 

To test the built Docker image run:

`docker run -e BITCOIN_ENDPOINT=https://api.coindesk.com/v1/bpi/currentprice/usd.json -e MINUTES_TO_SLEEP=1 -e MINUTES_TO_GET_AVERAGE=10 <repo-owner>/<repo-name>` 
