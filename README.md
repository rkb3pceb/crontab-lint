# crontab-lint

Static analyzer and validator for crontab expressions with human-readable output

## Installation

```bash
go install github.com/yourusername/crontab-lint@latest
```

## Usage

```bash
# Lint a crontab file
crontab-lint /etc/cron.d/myjob

# Lint from stdin
echo "*/5 * * * * /usr/bin/mycommand" | crontab-lint -
```

## Example Output

```
line 3: invalid day-of-week value '8' (valid range: 0-7)
line 7: field count mismatch (got 4, expected 5)
```
