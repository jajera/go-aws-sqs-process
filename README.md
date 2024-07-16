# go-aws-sqs-process

export AWS_VAULT_FILE_PASSPHRASE="$(cat /root/.awsvaultk)"

aws-vault exec dev -- terraform -chdir=./terraform init
aws-vault exec dev -- terraform -chdir=./terraform apply --auto-approve

bash ./terraform/terraform.tmp

go mod init go-aws-sqs-process
go mod tidy

go run cmd/main.go --task send --queueUrl https://sqs.ap-southeast-1.amazonaws.com/000000000000/myqueue --messageBody "Your message body here"

go run cmd/main.go --task receive --queueUrl https://sqs.ap-southeast-1.amazonaws.com/000000000000/myqueue

go run cmd/main.go --task delete --queueUrl https://sqs.ap-southeast-1.amazonaws.com/000000000000/myqueue --receiptHandle AQEBczEjCX2xWELfAd0svfHHalBgouth5gGiqHIf6EjIn1j1BJQj8pJcM7v1JZfDQmxHALmrXQ2es28xAlw4XrQA65K6RW2sfWhBrpwFfsGuSfvLenzwtq2IAYUNowQ8HAjcdod5dYl8YaguEeoBGG9d6aWxJ79qsqc7HclG7ia9FUw5K0Vnxl99WxsQYvO9HYeuZC28ITtlNLNrV2WAC54BdECJkhaf8XqBJ97uhvVheyh/amiYDv5C+NaakEGUcHK2K6v0zZMTYYUvXMMZblShpl0q7UvXwjVJnfZaEGTsq28Tmf1o6SqgGjdo5GJACMVH0dCnAJX/B1/ojL5H0XGaZ9pO4JgUL9fKNdIq+DqIapObVDMPGVEOwsRA3dJsmtinVBRYvxdLJ9gfXcO31OT9Ng==
