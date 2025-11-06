resource "github_actions_variable" "dockerhub_repo_backend" {
  repository    = "SlotSwapper"
  variable_name = "DOCKERHUB_REPO_BACKEND"
  value         = "slotswapper-backend"
}

resource "github_actions_variable" "dockerhub_repo_frontend" {
  repository    = "SlotSwapper"
  variable_name = "DOCKERHUB_REPO_FRONTEND"
  value         = "slotswapper-frontend"
}