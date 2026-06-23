# AWS EC2 Single Instance Deployment Guide

This guide details the step-by-step process of deploying the Postly Blog microservices platform onto a single AWS EC2 instance. This all-in-one architecture hosts the frontend (Nginx/Quasar), Go microservices, database layers (PostgreSQL, MongoDB), caching (Redis), and message broker (RabbitMQ) on a single virtual server using Docker Compose.

---

## 1. Instance Selection & Security Groups

### A. EC2 Instance Sizing
> [!IMPORTANT]
> The Go microservices and the Vue/Quasar frontend are compiled inside Docker containers from source. **Compilation is resource-intensive.**
> - **Recommended size:** At least `t3.medium` or `t2.medium` (4GB RAM, 2 vCPUs).
> - **Minimum size:** `t3.micro` or `t2.micro` (1GB RAM) **ONLY IF** you configure **Swap Space** (see Step 3) before running builds. Otherwise, the compiler will freeze or crash with an Out-of-Memory (OOM) error.

### B. Security Group Inbound Rules
Create or edit your EC2 instance's Security Group to allow inbound traffic to the following ports:

| Service | Port | Protocol | Source | Description |
| :--- | :--- | :--- | :--- | :--- |
| **SSH** | `22` | TCP | `My IP` | SSH Terminal access |
| **Frontend** | `80` | TCP | `0.0.0.0/0` | Public Web traffic |
| **Author Service** | `5000` | TCP | `0.0.0.0/0` | API requests for writing blogs |
| **Blog Service** | `5001` | TCP | `0.0.0.0/0` | API requests for reading blogs |
| **User Service** | `5002` | TCP | `0.0.0.0/0` | API requests for user profiles/auth |
| **RabbitMQ Web** | `15672` | TCP | `My IP` | Optional: RabbitMQ Management Panel |

---

## 2. Install Docker & Compose on EC2

SSH into your Ubuntu EC2 instance (`ssh -i key.pem ubuntu@YOUR_EC2_IP`) and execute the following commands:

```bash
# Update package database
sudo apt-get update -y && sudo apt-get upgrade -y

# Install Docker dependencies
sudo apt-get install -y ca-certificates curl gnupg lsb-release

# Add Docker's official GPG key
sudo fold -s -w 80 <<'EOF'
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
EOF

# Set up the repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine and Docker Compose
sudo apt-get update -y
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Run Docker without sudo (optional, requires logging back in)
sudo usermod -aG docker $USER
```

---

## 3. Configure Swap Space (Crucial for Micro Instances)
If you chose a `t3.micro` or `t2.micro` instance, enable swap space to avoid build crashes:

```bash
# Allocate a 4GB swap file
sudo fallocate -l 4G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# Make the swap file persistent across reboots
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/etc/fstab
```

---

## 4. Get the Code & Prepare the Environment

### A. Clone the Repository
Clone the codebase onto the EC2 instance:
```bash
git clone <YOUR_GIT_REPOSITORY_URL> ~/blog-platform
cd ~/blog-platform
```

### B. Create the `.env` File
Create a `.env` file in the root of `~/blog-platform` containing the external parameters required by the application. 

Replace `YOUR_EC2_PUBLIC_IP` with your instance's public IP address (or domain name, e.g. `blog.mydomain.com`).

```env
# ── Authentication ──
JWT_SECRET=supersecretjwttokenhere

# ── API endpoints (Used by browser clients to talk to services) ──
VITE_USER_SERVICE=http://YOUR_EC2_PUBLIC_IP:5002
VITE_AUTHOR_SERVICE=http://YOUR_EC2_PUBLIC_IP:5000
VITE_BLOG_SERVICE=http://YOUR_EC2_PUBLIC_IP:5001

# ── Google OAuth Credentials (Required by user-service) ──
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=GOCSPX-your-google-client-secret

# ── AWS S3 / MinIO Configuration (Required by User & Author services) ──
# Option 1: Use actual AWS S3 buckets (Production)
AWS_ACCESS_KEY_ID=your-aws-access-key-id
AWS_SECRET_ACCESS_KEY=your-aws-secret-access-key
AWS_REGION=us-east-1
AWS_S3_BUCKET=your-production-bucket-name

# Option 2: Mock S3 configurations (if you want to bypass actual S3 validation)
# Leave these as placeholders if you don't require file uploads immediately:
# AWS_ACCESS_KEY_ID=mockkey
# AWS_SECRET_ACCESS_KEY=mocksecret
# AWS_REGION=us-east-1
# AWS_S3_BUCKET=mock-bucket

# ── Gemini API (Required by author-service for AI capabilities) ──
GEMINI_API_KEY=AIzaSyYourGeminiApiKey
```

---

## 5. Build and Deploy

Run the docker compose build command. This builds the static Quasar frontend (injecting the correct public endpoint URLs at build time) and compiles the Go microservice binaries:

```bash
# Build the Docker images
docker compose build

# Start the services in detached mode
docker compose up -d
```

---

## 6. Verification and Health Checks

### A. Check Container Status
Verify that all 8 containers are running and healthy:
```bash
docker compose ps
```

### B. Check Logs
Check logs of the services if any fail to start:
```bash
docker compose logs -f user-service
docker compose logs -f author-service
docker compose logs -f blog-service
```

### C. Access the Site
Open your web browser and navigate to:
- Frontend: `http://YOUR_EC2_PUBLIC_IP`
- Health check endpoints:
  - User Service: `http://YOUR_EC2_PUBLIC_IP:5002/health`
  - Author Service: `http://YOUR_EC2_PUBLIC_IP:5000/health`
  - Blog Service: `http://YOUR_EC2_PUBLIC_IP:5001/health`

---

## 7. Production Enhancements (Next Steps)

Once verified, consider the following best practices for a production-ready application:

### A. SSL Support via Domain Name + Nginx Reverse Proxy
To secure your frontend and API endpoints with HTTPS (SSL) on port `443` without opening ports `5000-5002` to the public:
1. Register a domain name (e.g. `example.com`) and point its A-records to your EC2 IP.
2. Install `nginx` directly on the EC2 host: `sudo apt-get install -y nginx`
3. Configure Nginx to reverse proxy port 80/443 traffic to the docker containers:
   - Proxy `/` to `http://127.0.0.1:80` (Frontend)
   - Proxy `/api/v1/user` (or map specific routes) to `http://127.0.0.1:5002`
   - Proxy `/api/v1/author` to `http://127.0.0.1:5000`
   - Proxy `/api/v1/blog` to `http://127.0.0.1:5001`
4. Use Certbot to generate Free SSL Certificates:
   ```bash
   sudo apt-get install -y certbot python3-certbot-nginx
   sudo certbot --nginx -d example.com -d www.example.com
   ```
5. Update your frontend build arguments (`VITE_USER_SERVICE` etc.) in the `.env` to use the secure domain `https://example.com` and run `docker compose up --build -d` to rebuild the frontend assets.
