# Katt Backend


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