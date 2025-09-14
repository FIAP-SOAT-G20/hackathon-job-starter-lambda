# Hackathon Job Starter Lambda

A Go-based AWS Lambda function that processes S3 events and creates Kubernetes jobs to handle video processing tasks. This project demonstrates how to integrate AWS Lambda with Kubernetes for distributed job processing.

## 🏗️ Architecture

- **AWS Lambda**: Receives S3 events and triggers Kubernetes job creation
- **Kubernetes**: Executes video processing jobs using custom Docker images
- **S3**: Source of video files that trigger the processing pipeline
- **Docker**: Containerized job execution environment

## 🚀 Features

- S3 event-driven job creation
- Kubernetes job orchestration
- Configurable Docker images and commands
- Environment variable injection
- Local development support with minikube
- Comprehensive logging and monitoring

## 📋 Prerequisites

- Go 1.21+
- Docker Desktop
- minikube
- kubectl
- AWS CLI (for deployment)
- Make

## 🛠️ Installation & Setup

### 1. Clone the Repository
```bash
git clone <repository-url>
cd hackathon-job-starter-lambda
```

### 2. Install Dependencies
```bash
make install
```

### 3. Start Local Development Environment
```bash
# Start minikube
minikube start --driver=docker

# Start local database (if needed)
make compose-up
```

## 🏃‍♂️ Running the Project

### Local Development

1. **Build the Application**
   ```bash
   make build
   ```

2. **Start the Lambda Server**
   ```bash
   # Set environment variables for local testing
   export K8S_JOB_IMAGE="dummy-image:latest"
   export K8S_JOB_COMMAND="echo 'Hello from dummy image!'"
   export K8S_JOB_ENV_MY_BUILD_VAR="test build var"
   export K8S_JOB_ENV_MY_RUNTIME_VAR="test runtime var"
   
   # Start the lambda server
   make start-lambda
   ```

3. **Trigger the Lambda**
   ```bash
   # In a new terminal, trigger the lambda with test data
   make trigger-lambda
   ```

### Working with Local Docker Images

When developing with custom Docker images, you need to make them available to your Kubernetes cluster:

1. **Build Your Docker Image**
   ```bash
   # Build the dummy image (example)
   docker build -f Dockerfile.dummy -t dummy-image:latest .
   # Build the job checker (example)
   docker build -f Dockerfile.jobChecker -t job-checker .
   ```

2. **Load Image into minikube**
   ```bash
   # Make the image available to minikube's Docker environment
   minikube image load dummy-image:latest
   minikube image load job-checker:latest
   ```

3. **Configure Image Pull Policy**
   The application is configured to use `ImagePullPolicy: Never` for local images, which tells Kubernetes to use the local image instead of trying to pull from a registry.

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `K8S_NAMESPACE` | Kubernetes namespace for jobs | `default` |
| `K8S_JOB_IMAGE` | Docker image for job containers | `ghcr.io/fiap-soat-g20/hackathon-job-starter-lambda:latest` |
| `K8S_JOB_COMMAND` | Command to execute in job containers | `echo "Hello, World"` |
| `K8S_JOB_PREFIX` | Prefix for job names | `video-processor` |
| `K8S_JOB_ENV_*` | Environment variables with this format are set in the started job image and can contain any values as needed for your specific use case. | - |

## 📁 Project Structure

```
hackathon-job-starter-lambda/
├── bin/                          # Compiled binaries
├── internal/                     # Application code
│   ├── core/                     # Domain logic
│   │   ├── domain/               # Domain entities
│   │   ├── dto/                  # Data transfer objects
│   │   ├── port/                 # Interface definitions
│   │   └── usecase/              # Business logic
│   └── infrastructure/           # External integrations
│       ├── api/                  # Kubernetes API client
│       ├── aws/                  # AWS Lambda handlers
│       ├── config/               # Configuration management
│       ├── k8s/                  # Kubernetes client setup
│       └── logger/               # Logging utilities
├── test/                         # Test files and data
│   └── data/                     # Test event payloads
├── Dockerfile.dummy              # Example Docker image
├── main.go                       # Application entry point
├── Makefile                      # Build and deployment scripts
└── README.md                     # This file
```

## 🧪 Testing

### Run Tests
```bash
make test
```

### Run Tests with Coverage
```bash
make coverage
```

### Test Lambda Locally
```bash
# Start lambda server
make start-lambda

# In another terminal, trigger with test data
make trigger-lambda
```

## 🐳 Docker

### Building the Dummy Image
```bash
docker build -f Dockerfile.dummy -t dummy-image:latest .
```

### Using Custom Images
1. Build your custom image
2. Load it into minikube: `minikube image load your-image:latest`
3. Set `K8S_JOB_IMAGE="your-image:latest"` when starting the lambda

## 🚀 Deployment

### Build for AWS Lambda
```bash
make package
```

### Deploy to AWS
```bash
# Initialize Terraform (if using infrastructure as code)
make terraform-init

# Plan deployment
make terraform-plan

# Apply deployment
make terraform-apply
```

## 🔧 Development Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the application |
| `make start-lambda` | Start lambda server locally |
| `make trigger-lambda` | Trigger lambda with test data |
| `make test` | Run tests |
| `make coverage` | Run tests with coverage |
| `make clean` | Clean build artifacts |
| `make fmt` | Format code |
| `make lint` | Run linter |

## 🐛 Troubleshooting

### Common Issues

1. **Image Pull Errors**
   - Ensure your image is loaded into minikube: `minikube image load your-image:latest`
   - Check that `ImagePullPolicy: Never` is set for local images

2. **Port Already in Use**
   - Kill existing processes: `lsof -ti:3300 | xargs kill`

3. **Kubernetes Connection Issues**
   - Ensure minikube is running: `minikube status`
   - Start minikube: `minikube start --driver=docker`

4. **Job Creation Failures**
   - Check Kubernetes logs: `kubectl get pods` and `kubectl logs <pod-name>`
   - Verify environment variables are set correctly

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📞 Support

For questions or issues, please open an issue in the repository.