# Use an Nginx image
FROM nginx:alpine

# Copy the frontend files to Nginx's default content directory
COPY . /usr/share/nginx/html

# Expose the HTTP port
EXPOSE 80

# Start Nginx server
CMD ["nginx", "-g", "daemon off;"]
