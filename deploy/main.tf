provider "aws" {
  region = "ap-southeast-1"  # Change to your preferred region
}

# Lambda function
resource "aws_lambda_function" "katt_api" {
  function_name = "katt-api"
  filename      = "function.zip"
  handler       = "bootstrap"
  runtime       = "provided.al2"
  role          = aws_iam_role.lambda_role.arn

  environment {
    variables = {
      JWT_SECRET      = var.jwt_secret
      AUTH0_AUDIENCE  = var.auth0_audience
      AUTH0_DOMAIN    = var.auth0_domain
      DB_HOST         = var.db_host
      DB_USER         = var.db_user
      DB_PASS         = var.db_pass
      DB_NAME         = var.db_name
    }
  }

  timeout     = 30
  memory_size = 1024
}

# IAM role for Lambda
resource "aws_iam_role" "lambda_role" {
  name = "katt_lambda_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

# Attach basic Lambda execution policy
resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# API Gateway
resource "aws_apigatewayv2_api" "katt_api" {
  name          = "katt-api-gateway"
  protocol_type = "HTTP"
  cors_configuration {
    allow_origins = ["*"]
    allow_methods = ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    allow_headers = ["*"]
  }
}

# API Gateway stage - change from "prod" to "$default"
resource "aws_apigatewayv2_stage" "katt_api" {
  api_id      = aws_apigatewayv2_api.katt_api.id
  name        = "$default"  # This removes the stage prefix from the URL
  auto_deploy = true
}

# API Gateway integration with Lambda
resource "aws_apigatewayv2_integration" "katt_api" {
  api_id           = aws_apigatewayv2_api.katt_api.id
  integration_type = "AWS_PROXY"

  integration_method = "POST"
  integration_uri    = aws_lambda_function.katt_api.invoke_arn
}

# API Gateway route (catch-all)
resource "aws_apigatewayv2_route" "katt_api" {
  api_id    = aws_apigatewayv2_api.katt_api.id
  route_key = "ANY /{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.katt_api.id}"
}

# Fix Lambda permission for API Gateway
resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.katt_api.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.katt_api.execution_arn}/*/*"
}

# Add a default route for the root path
resource "aws_apigatewayv2_route" "root" {
  api_id    = aws_apigatewayv2_api.katt_api.id
  route_key = "ANY /"
  target    = "integrations/${aws_apigatewayv2_integration.katt_api.id}"
}

# Output the API endpoint URL
output "api_endpoint" {
  value = "${aws_apigatewayv2_stage.katt_api.invoke_url}"
}