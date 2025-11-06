provider "github" {
  token = var.github_token
}

resource "github_actions_secret" "docker_password" {
  repository      = "SlotSwapper"
  secret_name     = "DOCKER_PASSWORD"
  plaintext_value = var.dockerhub_token
}

resource "github_actions_secret" "docker_username" {
  repository      = "SlotSwapper"
  secret_name     = "DOCKER_USERNAME"
  plaintext_value = var.dockerhub_username
}