# Use the official Node.js image as the base image
FROM node:16-alpine

# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json to the working directory
COPY package*.json ./

# Forcefully install dependencies (replace with the appropriate command)
RUN npm install --force

# Copy the remaining frontend files to the working directory
COPY . .

# Expose the port that the React app listens on
EXPOSE 3000

# Run the React app
CMD ["npm", "start"]
