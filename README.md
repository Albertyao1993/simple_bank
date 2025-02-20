# Simple Bank

Simple Bank is a simple banking system that provides basic account management and transaction functionalities. This project utilizes a modern tech stack to achieve efficient and secure backend services.

## Tech Stack

- **Go**: The backend service is written in Go, offering high performance and concurrency handling.
- **Gin**: The Gin framework is used to build RESTful APIs, providing fast and flexible routing.
- **GORM**: GORM is used as the ORM framework to simplify database operations.
- **PostgreSQL**: PostgreSQL is used as the database to store user and transaction data.
- **Docker**: Docker is used to containerize the application, simplifying deployment and environment configuration.
- **JWT**: JSON Web Tokens are used for user authentication and authorization.
- **GoMock**: GoMock is used for unit testing to ensure code reliability and maintainability.

## Features

- **User Registration and Login**: Provides user registration and login functionalities with JWT for authentication.
- **Account Management**: Users can create and manage multiple bank accounts.
- **Transfer Functionality**: Supports transfers between accounts, ensuring transaction security and consistency.
- **Transaction History**: Offers detailed transaction history query functionality for users to view past transactions.
- **Security**: Uses encryption to protect user data and transaction information.

## Deployment

### AWS ECR and EKS

This project supports deployment on AWS's ECR (Elastic Container Registry) and EKS (Elastic Kubernetes Service).

#### Deployment Steps

1. **Build Docker Image**:
   ```bash
   docker build -t simple_bank .
   ```

2. **Push to AWS ECR**:
   - Log in to ECR:
     ```bash
     aws ecr get-login-password --region <your-region> | docker login --username AWS --password-stdin <your-account-id>.dkr.ecr.<your-region>.amazonaws.com
     ```
   - Create ECR repository:
     ```bash
     aws ecr create-repository --repository-name simple_bank
     ```
   - Tag and push the image:
     ```bash
     docker tag simple_bank:latest <your-account-id>.dkr.ecr.<your-region>.amazonaws.com/simple_bank:latest
     docker push <your-account-id>.dkr.ecr.<your-region>.amazonaws.com/simple_bank:latest
     ```

3. **Deploy to EKS**:
   - Configure `kubectl` to connect to your EKS cluster.
   - Create Kubernetes deployment and service:
     ```yaml
     apiVersion: apps/v1
     kind: Deployment
     metadata:
       name: simple-bank
     spec:
       replicas: 2
       selector:
         matchLabels:
           app: simple-bank
       template:
         metadata:
           labels:
             app: simple-bank
         spec:
           containers:
           - name: simple-bank
             image: <your-account-id>.dkr.ecr.<your-region>.amazonaws.com/simple_bank:latest
             ports:
             - containerPort: 8080
     ---
     apiVersion: v1
     kind: Service
     metadata:
       name: simple-bank
     spec:
       type: LoadBalancer
       ports:
       - port: 80
         targetPort: 8080
       selector:
         app: simple-bank
     ```
   - Apply the configuration:
     ```bash
     kubectl apply -f simple_bank_deployment.yaml
     ```

## Quick Start

1. **Clone the Project**:
   ```bash
   git clone https://github.com/yourusername/simple_bank.git
   cd simple_bank
   ```

2. **Build and Run Docker Container**:
   ```bash
   docker-compose up --build
   ```

3. **Access the API**:
   The API will run on `http://localhost:8080`, and you can test it using Postman or other tools.

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) to learn how to participate in the project.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
