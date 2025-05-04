# Katt Backend


### Tech Stack
1. Golang
2. Terraform
3. AWS Lambda
4. AWS APIGateway
5. Supabase for database

### Step To Deploy
1. run build command in the root directory

```bash
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
zip function.zip bootstrap
mv function.zip deploy/
```

2. access to folder deploy

```bash
terraform init
terraform plan
terraform apply
```