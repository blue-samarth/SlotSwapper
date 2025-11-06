variable "github_token" {
  description = "GitHub token with repo and admin:repo_hook permissions"
  type        = string
  sensitive   = true
}
variable "dockerhub_username" {
  description = "Docker Hub username"
  type        = string
  sensitive   = true
}
variable "dockerhub_token" {
  description = "Docker Hub token"
  type        = string
  sensitive   = true
}