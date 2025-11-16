# GitHub Setup Guide

Follow these steps to push your RBAC & ReBAC API project to GitHub.

## Step 1: Create a GitHub Repository

1. Go to [GitHub](https://github.com) and sign in
2. Click the **+** icon in the top right corner
3. Select **New repository**
4. Fill in the repository details:
   - **Repository name**: `go-rebac-rbac-postgres` (or your preferred name)
   - **Description**: `RBAC & ReBAC API demonstration with Go, Gin, GORM, and PostgreSQL`
   - **Visibility**: Choose Public or Private
   - **DO NOT** initialize with README, .gitignore, or license (we already have these)
5. Click **Create repository**

## Step 2: Add Remote and Push

After creating the repository on GitHub, you'll see a page with setup instructions. Use these commands:

### Option A: If you haven't committed yet (First time)

```bash
# Make sure you're in the project directory
cd /Users/amoako/Downloads/RBAC_ReBAC/go-rebac-rbac-postgres

# Initialize git (if not already done)
git init

# Add all files
git add .

# Create initial commit
git commit -m "Initial commit: RBAC & ReBAC API with Go, Gin, GORM, and PostgreSQL"

# Add your GitHub repository as remote (replace YOUR_USERNAME with your GitHub username)
git remote add origin https://github.com/YOUR_USERNAME/go-rebac-rbac-postgres.git

# Or if using SSH:
# git remote add origin git@github.com:YOUR_USERNAME/go-rebac-rbac-postgres.git

# Push to GitHub
git branch -M main
git push -u origin main
```

### Option B: If you've already committed locally

```bash
# Add your GitHub repository as remote (replace YOUR_USERNAME with your GitHub username)
git remote add origin https://github.com/YOUR_USERNAME/go-rebac-rbac-postgres.git

# Rename branch to main (if needed)
git branch -M main

# Push to GitHub
git push -u origin main
```

## Step 3: Verify

1. Go to your GitHub repository page
2. You should see all your files there
3. The README.md will be displayed on the repository homepage

## Troubleshooting

### Authentication Issues

If you get authentication errors:

**For HTTPS:**
- Use a Personal Access Token instead of password
- Generate token: GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
- Use the token as your password when pushing

**For SSH:**
- Set up SSH keys: [GitHub SSH Guide](https://docs.github.com/en/authentication/connecting-to-github-with-ssh)

### Push Rejected

If push is rejected:
```bash
# Pull first (if repository was initialized with files)
git pull origin main --allow-unrelated-histories

# Then push again
git push -u origin main
```

### Change Remote URL

If you need to change the remote URL:
```bash
# Remove existing remote
git remote remove origin

# Add new remote
git remote add origin https://github.com/YOUR_USERNAME/REPO_NAME.git
```

## Next Steps After Pushing

1. **Add repository topics**: Go to repository → ⚙️ → Topics → Add: `go`, `golang`, `rbac`, `rebac`, `postgresql`, `gin`, `gorm`, `api`
2. **Add description**: Update repository description
3. **Add license**: Consider adding a LICENSE file
4. **Enable GitHub Actions** (if you want CI/CD)
5. **Add badges** to README (optional)

## Quick Reference Commands

```bash
# Check remote
git remote -v

# Check status
git status

# View commits
git log --oneline

# Push changes (after initial setup)
git add .
git commit -m "Your commit message"
git push
```

## Repository URL Format

Replace `YOUR_USERNAME` and `REPO_NAME` with your actual values:

- **HTTPS**: `https://github.com/YOUR_USERNAME/REPO_NAME.git`
- **SSH**: `git@github.com:YOUR_USERNAME/REPO_NAME.git`

Example:
- `https://github.com/johndoe/go-rebac-rbac-postgres.git`
- `git@github.com:johndoe/go-rebac-rbac-postgres.git`

