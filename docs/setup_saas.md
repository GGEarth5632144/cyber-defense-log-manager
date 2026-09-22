# SaaS Setup Guide

To deploy this log management system in a public cloud, follow these steps:

1. **Provision a VM:** Spin up an Ubuntu 22.04 VM (e.g., AWS EC2 t3.medium, DigitalOcean Droplet) with at least 4 vCPUs and 8 GB RAM.
2. **Open Ports:** Ensure ports 80 (HTTP), 443 (HTTPS), and 514 (UDP for Syslog) are open in the firewall.
3. **Install Docker:**
   ```bash
   sudo apt-get update
   sudo apt-get install -y docker.io docker-compose
   ```
4. **Clone the repository:**
   ```bash
   git clone <YOUR_GIT_REPO_URL>
   cd log-management-system
   ```
5. **Set up HTTPS (Reverse Proxy):**
   - Install Caddy (recommended for easy Let's Encrypt HTTPS):
     ```bash
     sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
     curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
     curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
     sudo apt update && sudo apt install caddy
     ```
   - Update `/etc/caddy/Caddyfile`:
     ```
     yourdomain.com {
         reverse_proxy localhost:3000
     }
     api.yourdomain.com {
         reverse_proxy localhost:8080
     }
     ```
   - Restart Caddy: `sudo systemctl restart caddy`
6. **Start the application:**
   ```bash
   docker-compose up -d
   ```
7. Access the application at `https://yourdomain.com`.
