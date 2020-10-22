# lambda-go-rpc

AWS Lambda (golang runtime) RPC client for local environment

## Example

Run your lambda function server (with specific port).
```bash
_LAMBDA_SERVER_PORT=12345 ${YOUR_GOLANG_LAMBDA_SERVER_RUN_COMMAND}
```

Then you can send request with any payload and get response payload.
```bash
echo '{"any":"payload"}' | lambda-go-rpc -port 12345
```
