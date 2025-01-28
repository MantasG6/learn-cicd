# Google Cloud Platform

GCP is one of the "big three" cloud providers, along with AWS and Azure. We're going to use GCP to host our Notely application!

Everything we do in this course falls under the free tier of GCP, at the time of writing. That said, you will need to provide a credit card to sign up, and you should be careful to not exceed the free tier and free trial limits if you don't want to be charged.

## Create a GCP Account

First, you'll need to create a GCP account. You can do that [here](https://cloud.google.com/?hl=en).


## Create a Project

Once you've created an account, you'll need to [create a project](https://developers.google.com/workspace/guides/create-project).

Name the project `notely`.

One of my favorite aspects of GCP is how it groups resources by project. We'll keep everything for Notely in a single project, and when you're done with this course you can simply delete the project to clean everything up in one place.

## Create a Billing Account

Next, you'll need to create a billing account. This is where you'll provide your credit card information. You can find the billing section in the GCP console by clicking the hamburger menu in the top left, then "Billing".

Ensure your billing account is linked to your project, and you are able to see the billing information for your project in the GCP console.

# Google Cloud SDK

For some tasks, it makes sense to use the `gcloud` CLI instead of the GCP web console. For example, to run tasks from a GitHub Actions workflow, we'll need to use the `gcloud` CLI.

1. Install the gcloud CLI tool [here](https://cloud.google.com/sdk/docs/install).
2. Initialize it by running gcloud init in your terminal.
    - It will prompt you to login by opening a browser window. Login with the same account you used to create your GCP project.
    - Select your `notely` project.

## Assignment

Run `gcloud config list` and make sure your authenticated account and project are set correctly.

*Note: You should already be authenticated after running `gcloud init`. If not, run `gcloud auth login`.*

# Google Artifact Registry

We'll be using [Google Artifact Registry](https://cloud.google.com/artifact-registry/docs/overview) to store our Docker images. It's similar to Docker Hub, but it's private and hosted on GCP.

Whenever we create a new version of Notely, we'll build it into a new Docker image version and push that to Artifact Registry.

## Assignment

1. Search for and enable the `Cloud Build API`.
2. Within Artifact Registry in the GCP console, enable the Artifact Registry API.
3. Click `Create Repository`:

- Name: `notely-ar-repo`
- Format: `Docker`
- Mode: `Standard`
- Location type: `Region`
- Region for deployment: `us-central1`
- Leave "Google-managed encryption key" selected

Note that the image hosting region from earlier, and service deployment region we are targeting now, may not necessarily be the same region. Cloud providers provide flexibility with [availability zones](https://cloud.google.com/compute/docs/regions-zones), so that engineers can pick and choose the most optimal regions for your system.

4. Build and push the Docker image to Artifact Registry:

```
gcloud builds submit --tag REGION-docker.pkg.dev/PROJECT_ID/REPOSITORY/IMAGE:TAG .
```

*You can copy/paste your actual value for `REGION-docker.pkg.dev/PROJECT_ID/REPOSITORY` from your Artifact Registry repo page in the GCP console.*

## Tips

Copy and paste `REGION-docker.pkg.dev/PROJECT_ID/REPOSITORY` from the repository in the Artifact Registry! It should look like this at the top:

`us-central1-docker.pkg.dev > your-project-123456 > notely`

Example build command:

```
gcloud builds submit --tag us-central1-docker.pkg.dev/your-project-123456/notely-ar-repo/notely:latest .
```

# Automate Builds

Now that we've built the Docker image locally, let's build it automatically on every push to the `main` branch.

## Assignment

Use the [setup-gcloud](https://github.com/google-github-actions/setup-gcloud) action to authenticate with GCP.

I recommend using the simple [service account key JSON](https://github.com/google-github-actions/setup-gcloud#service-account-key-json) setup.

## Creating a Service Account

1. Go to the [IAM & Admin Service Accounts](https://console.cloud.google.com/iam-admin/serviceaccounts) section of the GCP console.
2. Create a service account and name it "Cloud Run Deployer" with these permissions:

- `Cloud Build Editor`
- `Cloud Build Service Account`
- `Cloud Run Admin`
- `Service Account User`
- `Viewer`

3. Create a JSON key for that service account and download it to your computer.

## Add the Key As a Secret in GitHub Actions

4. Go to your GitHub Repo > Repository Settings > Secrets and variables > Actions > New repository secret (**not** "environment secret" - use a repository secret)

- Name: `GCP_CREDENTIALS`
- Secret: Paste the entire JSON key from the file you downloaded from GCP

5. Save the secret

## Update Your GitHub Action Workflow

After the `buildprod` script runs, add the [setup-gcloud](https://github.com/google-github-actions/setup-gcloud) steps to setup the `gcloud` CLI and authenticate with GCP.

Finally, add a step to build the Docker image and push it to Google Artifact Registry.

```
gcloud builds submit --tag REGION-docker.pkg.dev/PROJECT_ID/REPOSITORY/IMAGE:TAG .
```

*You can copy/paste the actual value for `REGION-docker.pkg.dev/PROJECT_ID/REPOSITORY` from the repo's page in the GCP console.*

Commit and push your changes to GitHub. You should see the GitHub Action run and successfully build and push the Docker image to Google Artifact Registry.

# Google Cloud Run

Cloud Run is a serverless container hosting service. It's a great fit for Notely because it has a generous free tier, scales automatically, and automagically configures pesky infrastructure like:

- Load balancing
- DNS
- HTTPS

In a nutshell, we give Cloud Run a Docker image, and it runs it for us.

## Cloud Run Services

Cloud run has 2 types of applications:
- Services
- Jobs

A "service" is a Cloud Run application that listens and responds to web requests, as opposed to a "job" which is simply a task that runs to completion.

Notely makes sense as a service because it's a web application that needs to respond to web requests. It needs to serve the frontend, and it needs to respond to API requests.

## Assignment

For now, let's create a test service. Once we're comfortable with the process, we'll deploy Notely.

1. Navigate to the Cloud Run section of the GCP console.
2. Click "Create Service"
3. Container image URL: Click `SELECT` -> Find your image in Artifact Registry
4. Service name: `YOUR SERVICE NAME`
5. Region: us-central1
6. Check the "Allow unauthenticated invocations" checkbox. This will allow anyone to access the webpage the service serves without having to log in to GCP
7. Ingress control: "All" checked (allow direct access to the service)
8. Open "Container(s), Volumes, Networking, Security" and scroll to the bottom. Set the Maximum number of instances: `4`
9. Create the service
10. Wait for the service to deploy
11. Click the service's URL to see the webpage it serves

# Cloud Run Updates

Now that we have a service configured and running, let's update our GitHub actions workflow to automatically deploy changes to the app when we push to the `main` branch.

## Assignment

1. Add a new step to the `deploy` job in the GitHub actions workflow to deploy the app to Cloud Run. Because we want to allow unauthenticated access to the app, we'll also need to add a new security setting to the Cloud Run service. I recommend simply adding this to the same step in CD:

```
- name: Deploy to Cloud Run
  run: gcloud run deploy notely --image REGION-docker.pkg.dev/PROJECT_ID/REPO_NAME/IMAGE:TAG --region REGION --allow-unauthenticated --project PROJECT_ID --max-instances=4
```

2. Target the us-central1 region for service deployment and the latest image.

## Make a Change and Push to GitHub

3. Change the `/static/index.html` file so that the `h1` tag says "Welcome to Notely" instead of just "Notely".
4. Commit and push your changes to GitHub. You should see the GitHub Action run and successfully deploy the new version of the app to Cloud Run.
5. Open the URL in your browser and you should see the new version of the app.