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

## Cluster

See aks-engine examples on how to `secrets/kubernetes.json`. Generate arm templates:

```
aks-engine generate secrets/kubernetes.json
```

Deploy:

```
az group deployment create \
--name "ronenfreemanmshomework" \
--resource-group "microsoft-homework" \
--template-file "./_output/ronenfreemanmshomework/azuredeploy.json" \
--parameters "./_output/ronenfreemanmshomework/azuredeploy.parameters.json"
```

## DNS

Get ingress controller lb ip (13.77.181.241) ID:

```
PUBLIC_IP_ID=$(az network public-ip list --query "[?ipAddress=='13.77.181.241].id" -o tsv)
```
Create zone:

```
az network dns zone create \
  --resource-group microsoft-homework \
  --name ronenfreemanmshomework.io
```
Create record:
```
az network dns record-set a add-record \
    --resource-group microsoft-homework \
    --record-set-name "@" \
    --zone-name ronenfreemanmshomework.io \
    --ipv4-address 13.77.181.241
```
Update recode:
```
az network dns record-set a update --name @ \
  --resource-group microsoft-homework \
  --zone-name ronenfreemanmshomework.io \
  --target-resource $PUBLIC_IP_ID
```
```
az network dns record-set a add-record \
    --resource-group microsoft-homework \
    --zone-name ronenfreemanmshomework.io \
    --record-set-name www \
    --ipv4-address 13.77.181.241
```

Query nameservers:
```
az network dns zone show \
  --resource-group microsoft-homework  \
  --name ronenfreemanmshomework.io \
  --query nameServers
```

## Application

The application, given three environment variables `BITCOIN_ENDPOINT`, `MINUTES_TO_SLEEP` and `MINUTES_TO_GET_AVERAGE` gets the rate of Bitcoin from `BITCOIN_ENDPOINT` every `MINUTES_TO_SLEEP` minute(s) and averages the rate every MINUTES_TO_GET_AVERAGE minutes.

### Assumptions

- The app is not statefull. The app will never die.
- External monitoring exists to capture if the app would throw an error.
- The first minute to print the Bitcoin rate begins as the app starts up.
- There is zero delay in getting the rate.

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


## Future Work

- Add horizontal pod autoscalling
- Take into account the delay to fetch the Bitcoin rate and other app processing (to decide how long a minute is between checking for new rates)
- Add persistence in case the app crashes
- Ensure quotas are available for node autoscaling


## Problems:
- quotas. Had to use cloud shell to see
- Network policy didnt allow routing from nginx controller. label namespace