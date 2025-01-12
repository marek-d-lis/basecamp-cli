# BaseCamp CLI

BaseCamp CLI is a command-line tool designed to simplify and automate the process of bootstrapping environments. It allows you to run Ansible playbooks from multiple repositories in a structured and efficient way.

## **Purpose**
The main purpose of BaseCamp CLI is to:

- Automate the setup of your development environment.
- Support multiple repositories containing different configurations (e.g., PHP, Node.js, Docker, IDE setups).
- Provide a clean and modular approach to managing and running Ansible playbooks.
- Automatically clean up temporary files after execution.

This tool is particularly useful for developers who frequently set up new machines or want a repeatable, consistent setup across multiple projects and devices.

---

## **Key Features**

- **Run Multiple Repositories:** Specify multiple repositories with Ansible playbooks to run sequentially.
- **YAML Configuration Support:** Use a configuration file (`basecamp.yml`) to list repositories.
- **Command-line Options:** Add repositories directly via `--repo` flags.
- **Automatic Cleanup:** Automatically deletes cloned repositories after execution to keep the system clean.

---

## **Installation**

To build the BaseCamp CLI from source:

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/basecamp-cli.git
   cd basecamp-cli
   ```

2. Build the binary:
   ```bash
   go build -o basecamp
   ```

3. Run the CLI tool:
   ```bash
   ./basecamp run --help
   ```

---

## **Usage**

### **Run with Repositories via `--repo` Flags**
You can specify repositories directly:
```bash
./basecamp run --repo https://github.com/user/vue-setup.git --repo https://github.com/user/python-setup.git
```

### **Run with a YAML Configuration File**
Create a `basecamp.yml` file:
```yaml
repos:
  - https://github.com/your-username/BaseCamp.git
  - https://github.com/user/vue-setup.git
```

Run the command with the configuration file:
```bash
./basecamp run --config basecamp.yml
```

### **Combination of Config File and Flags**
Add more repositories on top of those listed in the configuration file:
```bash
./basecamp run --config basecamp.yml --repo https://github.com/extra/setup.git
```

---

## **Example Workflow**

1. Clone your repository.
2. Create or edit `basecamp.yml`.
3. Run the command:
   ```bash
   ./basecamp run --config basecamp.yml
   ```
4. After execution, all temporary files will be cleaned up automatically.


## **Contributing**

Contributions are welcome! Please open an issue or submit a pull request with improvements or bug fixes.

---

## **License**
This project is licensed under the MIT License.

---

With BaseCamp CLI, you can easily bootstrap and extend your development environments with a consistent, repeatable process.

