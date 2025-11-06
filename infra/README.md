# Terraform GitHub Secrets & Variables Management

This Terraform configuration automates the management of GitHub Actions secrets and variables for the SlotSwapper project's CI/CD pipelines.

## 📋 Overview

This setup automatically configures:
- **GitHub Actions Secrets**: Sensitive credentials for Docker Hub authentication
- **GitHub Actions Variables**: Non-sensitive configuration values for Docker repositories

## 🏗️ Infrastructure Components

### Providers
- **GitHub Provider** (`versions.tf`): Manages GitHub repository resources
  - Version: `6.6.0`
  - Source: `integrations/github`

### Resources Created

#### Secrets (Encrypted)
1. **DOCKER_USERNAME** - Docker Hub username for authentication
2. **DOCKER_HUB_TOKEN** - Docker Hub access token for pushing images

#### Variables (Public)
1. **DOCKERHUB_REPO_BACKEND** - Backend Docker repository name (`slotswapper-backend`)
2. **DOCKERHUB_REPO_FRONTEND** - Frontend Docker repository name (`slotswapper-frontend`)

## 🚀 Setup Instructions

### Prerequisites
1. **Terraform** installed (v1.0+)
   ```bash
   brew install terraform  # macOS
   ```

2. **GitHub Personal Access Token** with permissions:
   - `repo` (Full control of private repositories)
   - `admin:repo_hook` (Admin access to repository hooks)
   - `workflow` (Update GitHub Action workflows)
   
   Create token at: https://github.com/settings/tokens/new

3. **Docker Hub Access Token**
   - Go to: https://hub.docker.com/settings/security
   - Click "New Access Token"
   - Give it a name (e.g., "SlotSwapper CI/CD")
   - Copy the token immediately (you won't see it again!)

### Installation Steps

1. **Create secrets file** (never commit this!)
   ```bash
   cd Terraform
   touch secrets.tfvars
   ```

2. **Add your credentials to `secrets.tfvars`**
   ```hcl
   github_token        = "ghp_your_github_personal_access_token_here"
   dockerhub_username  = "your-dockerhub-username"
   dockerhub_token     = "dckr_pat_your_dockerhub_token_here"
   ```

3. **Ensure `.gitignore` excludes secrets**
   ```bash
   # Already in .gitignore (verify):
   *.tfvars
   *.tfstate
   *.tfstate.backup
   .terraform/
   ```

4. **Initialize Terraform**
   ```bash
   terraform init
   ```

5. **Review the execution plan**
   ```bash
   terraform plan -var-file="secrets.tfvars"
   ```

6. **Apply the configuration**
   ```bash
   terraform apply -var-file="secrets.tfvars"
   ```

7. **Confirm** by typing `yes` when prompted

## 📁 File Structure

```
Terraform/
├── README.md                    # This file
├── versions.tf                  # Provider configuration
├── variables.tf                 # Variable definitions
├── github_secrets.tf           # Secret resources
├── github_variables.tf         # Variable resources
├── secrets.tfvars              # Your actual values (NEVER COMMIT!)
├── terraform.tfstate           # State file (managed by Terraform)
└── .terraform/                 # Terraform cache directory
```

## 🔐 Security Best Practices

### ✅ DO:
- Store `secrets.tfvars` securely (password manager, vault)
- Use access tokens instead of passwords
- Rotate tokens regularly (every 90 days)
- Use minimal required permissions for tokens
- Review `.gitignore` before committing

### ❌ DON'T:
- Commit `secrets.tfvars` to version control
- Share tokens via insecure channels (email, Slack)
- Use admin-level tokens if not needed
- Leave old tokens active after rotation

## 🔄 Updating Secrets/Variables

### Update a Secret Value
1. Edit `secrets.tfvars` with new value
2. Run: `terraform apply -var-file="secrets.tfvars"`
3. Terraform will detect the change and update GitHub

### Add New Secret/Variable
1. Add variable definition to `variables.tf`:
   ```hcl
   variable "new_secret" {
     description = "Description of new secret"
     type        = string
     sensitive   = true
   }
   ```

2. Add resource to `github_secrets.tf` or `github_variables.tf`:
   ```hcl
   resource "github_actions_secret" "new_secret" {
     repository      = "SlotSwapper"
     secret_name     = "NEW_SECRET"
     plaintext_value = var.new_secret
   }
   ```

3. Add value to `secrets.tfvars`:
   ```hcl
   new_secret = "actual_value_here"
   ```

4. Apply changes:
   ```bash
   terraform apply -var-file="secrets.tfvars"
   ```

## 🧹 Cleanup

To remove all managed secrets and variables:

```bash
terraform destroy -var-file="secrets.tfvars"
```

⚠️ **Warning**: This will delete all secrets and variables from GitHub Actions!

## 🔍 Verification

After applying, verify in GitHub:

1. Go to: `https://github.com/blue-samarth/SlotSwapper/settings/secrets/actions`
2. Check secrets exist:
   - ✅ DOCKER_USERNAME
   - ✅ DOCKER_HUB_TOKEN

3. Go to: `https://github.com/blue-samarth/SlotSwapper/settings/variables/actions`
4. Check variables exist:
   - ✅ DOCKERHUB_REPO_BACKEND = `slotswapper-backend`
   - ✅ DOCKERHUB_REPO_FRONTEND = `slotswapper-frontend`

## 🔗 Integration with CI/CD

These secrets/variables are used by:

- **`.github/workflows/cd.yml`** - Continuous Deployment
  - Uses `DOCKER_USERNAME` and `DOCKER_HUB_TOKEN` to authenticate
  - Uses `DOCKERHUB_REPO_BACKEND` and `DOCKERHUB_REPO_FRONTEND` for image names
  - Triggers on push to `main`, `master`, or `development` branches

## 📚 Terraform Commands Reference

| Command | Description |
|---------|-------------|
| `terraform init` | Initialize Terraform and download providers |
| `terraform plan -var-file="secrets.tfvars"` | Preview changes without applying |
| `terraform apply -var-file="secrets.tfvars"` | Apply changes to GitHub |
| `terraform destroy -var-file="secrets.tfvars"` | Remove all managed resources |
| `terraform show` | Display current state |
| `terraform refresh` | Update state from actual infrastructure |

## 🐛 Troubleshooting

### Error: "403 Resource not accessible by personal access token"
**Solution**: Ensure your GitHub token has `repo` and `workflow` scopes.

### Error: "404 Not Found"
**Solution**: Verify repository name is correct (`SlotSwapper`) and token has access.

### Error: "401 Unauthorized" with Docker Hub
**Solution**: Check Docker Hub token is valid and not expired.

### Secrets not updating in GitHub
**Solution**: Run `terraform refresh` then `terraform apply` again.

### State file conflicts
**Solution**: If working in a team, consider using [Terraform Cloud](https://cloud.hashicorp.com/products/terraform) or S3 backend for shared state.

## 📖 Additional Resources

- [Terraform GitHub Provider Docs](https://registry.terraform.io/providers/integrations/github/latest/docs)
- [GitHub Actions Secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
- [Docker Hub Access Tokens](https://docs.docker.com/docker-hub/access-tokens/)
- [Terraform Best Practices](https://www.terraform-best-practices.com/)

## 🤝 Contributing

When adding new secrets/variables:
1. Update `variables.tf` with variable definition
2. Update relevant resource file (`github_secrets.tf` or `github_variables.tf`)
3. Update this README with the new secret/variable
4. **Never commit actual secret values!**

---

**Last Updated**: November 7, 2025  
**Maintained by**: SlotSwapper Team
