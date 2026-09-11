# Promail
Email delivery application leveraging custom templates

## Rate limiting

The HTTP server applies a per-client fixed-window limit to all endpoints. The defaults are 100 requests per minute. Configure them with:

- `APP_RATE_LIMIT_REQUESTS`
- `APP_RATE_LIMIT_WINDOW_SECONDS`
