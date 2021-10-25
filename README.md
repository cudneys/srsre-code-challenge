# Sr SRE Code Challenge

This repo contains the source code and binaries for a REST API that generates and validates papsswords.  The server needs to be packaged in a Docker container and deployed in Kubernetes >= 1.20 with Istio >= 1.10.

# Steps
1. Clone this repo
2. Write an appropriate Dockerfile 
3. Build a Docker container
4. Write all of the deployment manifests
5. Check your code in and send us the URL to the repo

# Details
## App
The server is a REST API server that includes a Swagger UI (/swagger/index.html) and a Prometheus exporter (/metrics).  App APIs are versioned and served from /api/v1/.

The application has two simple functions:
* Password Generator (/generate)
* Password Validator (/validate)

The server accepts three inputs:
* **Bind Host**: 
  * CLI Flag: --host
  * Env Var: BIND_HOST
* **Bind Port**:
  * CLI Flag: --port
  * Env Var: BIND_PORT
* **Release Mode**:
  * Env Var: GIN_MODE
  * Values: "release" or "debug" (Defaults to "debug")

The application generally uses less than 75M of RAM, but you likely want to give it a bit more to give the garbage collector time to catch up when the app is under load.

Pre-compiled binaries are available in the [dist/](https://github.com/cudneys/srsre-code-challenge/tree/main/dist) directory.  

## Deployment
Your deployment manifests must be placed into a directory named "kubernetes" in the root of your repo.  

You may deliver the manifests as YAML files or a helm chart.  Please choose the method that you're most comfortable with. 

Your manifests should include a deployment, HPA, service, and virtual service (Istio).  The app will be deployed in a namespace that already contains a gateway named "myhappygatewway"  

## Extra Credit
* Since we don't want to expose the Swagger or Prometheus endpoints, it would be helpful to return a response greater than a 300 for requests from originating outside the cluster.
* We set our imagePullPolicy to Always, so it's helpful to have the smallest possible image.  Try to build a container that's less than 40M. 





