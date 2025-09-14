#!/bin/bash

echo "Testing logging setup..."

# Create a test pod that generates logs
kubectl run test-logger --image=busybox --rm -it --restart=Never -- sh -c "
echo 'This is a test log message from Minikube'
echo 'Another test message with timestamp: $(date)'
echo 'JSON test: {\"level\":\"info\",\"message\":\"test json log\"}'
sleep 5
"

echo "Test pod completed. Check Elasticsearch for logs..."
echo "You can check logs with: make check-logs"
echo "Or access Kibana with: make port-forward-kibana"
```

Make it executable:
```bash
chmod +x scripts/test-logging.sh
```

## 8. Deployment Steps

Now you can deploy your complete logging stack:

```bash
# Deploy the complete logging stack
make deploy-logging

# Test the setup
./scripts/test-logging.sh

# Check if logs are being indexed
make check-logs

# Access Kibana (in a separate terminal)
make port-forward-kibana
```

## 9. Verify the Setup

After deployment, you can verify everything is working:

```bash
# Check if all pods are running
kubectl get pods -n observability
kubectl get pods -n kube-system | grep fluent-bit

# Check Elasticsearch health
make check-elasticsearch

# Check if logs are being indexed
make check-logs
```

## 10. Access Kibana

Once everything is deployed, you can access Kibana:

```bash
<code_block_to_apply_changes_from>
```

Then open your browser to `http://localhost:5601`

In Kibana, you can:
1. Go to "Stack Management" → "Index Patterns"
2. Create an index pattern for `minikube-*`
3. Go to "Discover" to view your logs
4. Create dashboards to monitor your video processing jobs

## 11. Monitor Your Application Logs

Your application will continue to log to stdout/stderr as usual, and Fluent Bit will automatically collect these logs and send them to Elasticsearch. You can view them in Kibana by:

1. Filtering by namespace: `kubernetes.namespace_name: default`
2. Filtering by job name: `kubernetes.labels.job-name: video-processor-*`
3. Searching for specific log messages

This setup gives you a complete local logging solution where:
- Your application logs normally to console
- Fluent Bit collects all container logs
- Elasticsearch stores and indexes the logs
- Kibana provides a web interface to search and visualize logs
- Everything runs locally in your Minikube cluster

The logs will be automatically collected from all your video processing jobs and other containers, making it easy to monitor and debug your application.