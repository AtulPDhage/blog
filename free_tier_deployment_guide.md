# Postly Blog Platform: Free-Tier Deployment Guide

This guide describes how to deploy the entire Postly microservices stack (3 backend services + frontend SPA) using **free-tier** services, completely replacing the AWS EC2/ECR deployment.

---

## Architecture Topology (Free Tier)

```mermaid
graph TD
    subgraph Client ["Client Browser"]
        VercelSPA[Quasar Frontend SPA on Vercel]
    end

    subgraph Render ["Render (Free Web Services)"]
        UserService[User Service Go]
        AuthorService[Author Service Go]
        BlogService[Blog Service Go]
    end

    subgraph Datastores ["Free Cloud Datastores & Middleware"]
        Postgres[(PostgreSQL: Neon / Supabase)]
        MongoDB[(MongoDB: MongoDB Atlas)]
        Redis[(Redis: Upstash)]
        RabbitMQ[(RabbitMQ: CloudAMQP)]
        S3[(Object Storage: AWS S3 Free)]
    end

    VercelSPA -->|HTTP Requests| Render
    UserService --> MongoDB
    UserService --> S3
    AuthorService --> Postgres
    AuthorService --> RabbitMQ
    AuthorService --> S3
    BlogService --> Postgres
    BlogService --> Redis
    BlogService --> RabbitMQ
    BlogService --> UserService
```

---

## 1. Setup Free-Tier Datastores & Middleware

Before deploying the code, provision the necessary databases and messaging layers.

### A. PostgreSQL (Neon / Supabase)
Used by **Author Service** and **Blog Service**.
1. Sign up at [Neon](https://neon.tech/) or [Supabase](https://supabase.com/).
2. Create a new project and select PostgreSQL database.
3. Copy the database connection URL (e.g. `postgresql://[user]:[password]@[host]/[dbname]?sslmode=require`).
   > [!IMPORTANT]
   > Make sure the connection string contains `sslmode=require` or `sslmode=disable` appropriate to the service. Neon requires `sslmode=require` for secure connection over the internet.

### B. MongoDB (MongoDB Atlas)
Used by **User Service**.
1. Sign up at [MongoDB Atlas](https://www.mongodb.com/products/platform/atlas-database).
2. Deploy a free **M0 Sandbox** cluster (choose your closest AWS/GCP region).
3. Under **Database Access**, create a user with read/write permissions.
4. Under **Network Access**, add IP `0.0.0.0/0` (allow connections from anywhere since Render IPs rotate).
5. Copy the connection string (e.g. `mongodb+srv://[user]:[password]@[host]/?retryWrites=true&w=majority`).

### C. Redis (Upstash)
Used by **Blog Service** for caching.
1. Sign up at [Upstash](https://upstash.com/).
2. Create a serverless Redis database.
3. Copy the **Redis URL** (e.g. `redis://default:[password]@[host]:[port]`).

### D. RabbitMQ (CloudAMQP)
Used for asynchronous message queuing between services.
1. Sign up at [CloudAMQP](https://www.cloudamqp.com/).
2. Create a new instance on the **Lemur (Free)** plan.
3. Copy the **Host** name, **Username**, and **Password** from the dashboard.

### E. AWS S3 (AWS Free Tier)
Used by **User** and **Author** services for media uploads.
- If you already have an AWS account, you can continue using the S3 bucket under the AWS Free Tier (5GB free for 12 months). Make sure to create an IAM user with S3 read/write policy, and copy the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`.

---

## 2. Docker Hub Repositories Setup

Your GitHub Actions will build Docker images and push them here. Render will pull from here to run them.

1. Sign up at [Docker Hub](https://hub.docker.com/).
2. Create three new **public** repositories:
   - `[your-dockerhub-username]/user-service`
   - `[your-dockerhub-username]/author-service`
   - `[your-dockerhub-username]/blog-service`
3. Generate a Personal Access Token (PAT) for GitHub Actions:
   - Go to **Account Settings > Security > Personal access tokens**.
   - Click **Generate new token**, name it (e.g., `GitHub Actions CI/CD`), and copy it.

---

## 3. Render Web Services Setup (Backend)

For each Go service, create a new Web Service in Render that pulls from Docker Hub.

1. Sign up at [Render](https://render.com/).
2. Click **New +** > **Web Service**.
3. Select **Deploy an existing image** and click **Next**.
4. Enter the Docker Hub image path and click **Next**:
   - User Service: `docker.io/[your-dockerhub-username]/user-service:latest`
   - Author Service: `docker.io/[your-dockerhub-username]/author-service:latest`
   - Blog Service: `docker.io/[your-dockerhub-username]/blog-service:latest`
5. Choose the **Free** instance type.
6. Under **Environment Variables**, add the respective configuration variables:

   > [!TIP]
   > For RabbitMQ, you can either provide the host, username, and password separately in the environment variables, OR you can simply paste the full AMQP connection URI (e.g. `amqps://user:pass@host/vhost`) into `Rabbitmq_Host` and leave `Rabbitmq_Username` and `Rabbitmq_Password` blank. Both formats are supported!

| Service | Environment Variable Key | Suggested Value / Example |
| :--- | :--- | :--- |
| **User Service** | `PORT` | `5002` |
| | `BIND_ADDR` | `0.0.0.0` |
| | `MONGO_URI` | *Your MongoDB Atlas connection string* |
| | `DB_NAME` | `MasterDB` |
| | `JWT_SECRET` | *Your long random secret key* |
| | `Google_Client_Id` | *Your Google client ID* |
| | `Google_client_Secret` | *Your Google client secret* |
| | `AWS_ACCESS_KEY_ID` | *Your AWS access key* |
| | `AWS_SECRET_ACCESS_KEY` | *Your AWS secret key* |
| | `AWS_REGION` | `us-east-1` |
| | `AWS_S3_BUCKET` | *Your S3 bucket name* |
| **Author Service** | `PORT` | `5000` |
| | `BIND_ADDR` | `0.0.0.0` |
| | `DB_URL` | *Your Neon / Supabase Postgres connection string* |
| | `Rabbitmq_Host` | *Your CloudAMQP Host (e.g. whale.rmq.cloudamqp.com)* |
| | `Rabbitmq_Username` | *Your CloudAMQP Username* |
| | `Rabbitmq_Password` | *Your CloudAMQP Password* |
| | `JWT_SECRET` | *Your long random secret key* |
| | `GEMINI_API_KEY` | *Your Gemini API key* |
| | `AWS_ACCESS_KEY_ID` | *Your AWS access key* |
| | `AWS_SECRET_ACCESS_KEY` | *Your AWS secret key* |
| | `AWS_REGION` | `us-east-1` |
| | `AWS_S3_BUCKET` | *Your S3 bucket name* |
| **Blog Service** | `PORT` | `5001` |
| | `BIND_ADDR` | `0.0.0.0` |
| | `DB_URL` | *Your Neon / Supabase Postgres connection string* |
| | `REDIS_URL` | *Your Upstash Redis URL* |
| | `Rabbitmq_Host` | *Your CloudAMQP Host* |
| | `Rabbitmq_Username` | *Your CloudAMQP Username* |
| | `Rabbitmq_Password` | *Your CloudAMQP Password* |
| | `USER_SERVICE` | *Your Render User Service URL (e.g. `https://user-service-xxxx.onrender.com`)* |
| | `JWT_SECRET` | *Your long random secret key* |

7. Save and let the services run. Once deployed, locate the **Deploy Hook** for each service:
   - Go to the service's **Settings** tab.
   - Scroll down to **Deploy Hook**.
   - Copy the URL (e.g. `https://api.render.com/deploy/srv-xxxxxxxxxxxxx?key=yyyyyyyyyy`).

---

## 4. Vercel Project Setup (Frontend)

1. Sign up at [Vercel](https://vercel.com/).
2. Create a Vercel Personal Access Token:
   - Go to **Account Settings > Tokens**.
   - Click **Create**, name it `github-actions-token`, select Scope, and copy the token.
3. Install Vercel CLI locally (or check dashboard URLs) to find Organization and Project IDs:
   - Alternatively, link your GitHub repository to Vercel via Vercel's UI first (which auto-creates the project). Go to project settings to grab the IDs:
     - **Project ID**: Found in Project Settings > General.
     - **Organization ID / Team ID**: Found in Team Settings > General.
4. Set the following environment variables under **Project Settings > Environment Variables** on Vercel:
   - `VITE_USER_SERVICE`: `https://[your-user-service-name].onrender.com`
   - `VITE_AUTHOR_SERVICE`: `https://[your-author-service-name].onrender.com`
   - `VITE_BLOG_SERVICE`: `https://[your-blog-service-name].onrender.com`
   - `VITE_GOOGLE_CLIENT_ID`: *Your Google OAuth Client ID*
5. Configure the build parameters in Vercel settings:
   - **Root Directory**: `frontend`
   - **Build Command**: `npm run build`
   - **Output Directory**: `dist/spa`

---

## 5. Configure GitHub Repository Secrets

Go to your repository on GitHub, then navigate to **Settings > Secrets and variables > Actions** and add the following repository secrets:

| Secret Name | Description |
| :--- | :--- |
| `DOCKERHUB_USERNAME` | Your Docker Hub account username. |
| `DOCKERHUB_TOKEN` | Your Docker Hub Personal Access Token (PAT). |
| `RENDER_DEPLOY_HOOK_USER_SERVICE` | Deploy hook URL for `user-service`. |
| `RENDER_DEPLOY_HOOK_AUTHOR_SERVICE` | Deploy hook URL for `author-service`. |
| `RENDER_DEPLOY_HOOK_BLOG_SERVICE` | Deploy hook URL for `blog-service`. |
| `VERCEL_TOKEN` | Your Vercel Personal Access Token. |
| `VERCEL_ORG_ID` | Your Vercel Organization/Team ID. |
| `VERCEL_PROJECT_ID` | Your Vercel Project ID. |

Once these secrets are set and you push your changes to the `main` branch, GitHub Actions will build, push, and trigger automatic deployments for all services!
