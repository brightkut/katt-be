variable "jwt_secret" {
  description = "JWT Secret Key"
  type        = string
  sensitive   = true
}

variable "auth0_audience" {
  description = "Auth0 API Audience"
  type        = string
}

variable "auth0_domain" {
  description = "Auth0 Domain"
  type        = string
}

variable "db_host" {
  description = "Database Host"
  type        = string
}

variable "db_user" {
  description = "Database User"
  type        = string
  sensitive   = true
}

variable "db_pass" {
  description = "Database Password"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "Database Name"
  type        = string
}