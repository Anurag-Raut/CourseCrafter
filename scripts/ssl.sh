#!/bin/bash

PASSWORD="Anurag12345"


# Step 1: Update and install Certbot
echo "Updating package lists and installing Certbot..."
echo "$PASSWORD" | sudo -S apt update
echo "$PASSWORD" | sudo -S apt install -y certbot python3-certbot-nginx

# Step 2: Set variables for domain and email
DOMAIN="backend.coursecrafter.site"
EMAIL="rautab0@gmail.com"

# Step 3: Add basic Nginx server block for Certbot HTTP challenge
echo "Configuring Nginx for HTTP challenge..."
echo "$PASSWORD" | sudo -S tee /etc/nginx/sites-available/default > /dev/null <<EOT
server {
    listen 80;
    server_name $DOMAIN www.$DOMAIN;

    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    location / {
        return 404;
    }
}
EOT

# Reload Nginx to apply the changes
echo "Testing and reloading Nginx for Certbot challenge..."
echo "$PASSWORD" | sudo -S nginx -t
if [ $? -ne 0 ]; then
    echo "Error: Nginx configuration test failed."
    exit 1
fi
echo "$PASSWORD" | sudo -S systemctl reload nginx

# Step 4: Obtain SSL certificate using Certbot
echo "Obtaining SSL certificate for $DOMAIN..."
echo "$PASSWORD" | sudo -S certbot --nginx -d $DOMAIN --email $EMAIL --agree-tos --non-interactive --debug

# Check if Certbot succeeded
if [ $? -ne 0 ]; then
    echo "Error: Certbot failed to obtain SSL certificate."
    exit 1
fi

# Step 5: (Optional) Generate Diffie-Hellman parameters for extra security
echo "Generating Diffie-Hel lman parameters..."
echo "$PASSWORD" | sudo -S openssl dhparam -out /etc/ssl/certs/dhparam.pem 2048

# Step 6: Configure Nginx for HTTPS redirection and reverse proxy
echo "Configuring Nginx for HTTPS and reverse proxy..."
echo "$PASSWORD" | sudo -S tee /etc/nginx/sites-available/default > /dev/null <<EOT
# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name $DOMAIN;

    location / {
        return 301 https://\$host\$request_uri;
    }
}

# SSL Server Block (HTTPS)
server {
    listen 443 ssl;
    server_name $DOMAIN;

    ssl_certificate /etc/letsencrypt/live/$DOMAIN/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/$DOMAIN/privkey.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384';
    ssl_prefer_server_ciphers on;
    ssl_dhparam /etc/ssl/certs/dhparam.pem;

    location / {
        proxy_pass http://127.0.0.1:8080;  # Forward requests to port 8003
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
    }
}
EOT

# Step 7: Test and reload Nginx
echo "Testing Nginx configuration..."
echo "$PASSWORD" | sudo -S nginx -t
if [ $? -ne 0 ]; then
    echo "Error: Nginx configuration test failed."
    exit 1
fi

echo "Reloading Nginx..."
echo "$PASSWORD" | sudo -S systemctl reload nginx
if [ $? -ne 0 ]; then
    echo "Error: Failed to reload Nginx."
    exit 1
fi

# Step 8: Open firewall ports for HTTP/HTTPS (if UFW is enabled)
echo "Opening ports 80 and 443 in the firewall..."
echo "$PASSWORD" | sudo -S ufw allow 80/tcp
echo "$PASSWORD" | sudo -S ufw allow 443/tcp
echo "$PASSWORD" | sudo -S ufw reload

echo "Script execution completed successfully. SSL certificate is active and Nginx is configured."
