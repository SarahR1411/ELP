# Quick GitHub Command Reminders

## Cloning the Repository

To clone the repository, write this in the VS code terminal:

```bash
git clone https://github.com/SarahR1411/ELP.git
```

## 1. Pulling Changes

Before starting to work on your local repository, make sure that you have the latest changes from the remote repository using:
```bash
git pull origin main
```

## 2. Adding New Files

When you create a new file and want to add it to Git, use the following command to stage the file:

```bash
git add <file-name>
```

For example, to add a new go file called `main.go`, run:

```bash
git add main.go
```

To add all new and modified files at once, use:

```bash
git add .
```

## 3. Committing Changes

Once you've added your changes, you need to commit them. Make sure to add a short and clear commit message to explain what you did. For example:

```bash
git commit -a -m "Added new go file for ..."
```

## 4. Pushing Changes

After committing your changes, push them to the remote repository on GitHub using:

```bash
git push origin main
```

