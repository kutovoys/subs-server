# Subs Server

**Subs Server** is a lightweight HTTP server for serving base64-encoded subscription files from the local filesystem. It watches files for changes and provides customizable HTTP headers for clients.

## ✨ Features

- 📁 Serves subscription data from the local filesystem.
- 🔄 Automatically watches for file changes (create, update, delete).
- 🌐 Supports direct protocols (`vmess://`, `vless://`, `ss://`, `trojan://`, etc.).
- 🌍 Supports remote base64-encoded subscriptions via `http(s)://...` links.
- ➕ Combines local and remote subscription entries into a single output.
- 🛠 Fully configurable via CLI flags or environment variables.
- 📦 Custom HTTP headers for profile metadata.
- 🐞 Debug mode with listing of all available endpoints.

---

## 🚀 Usage

### CLI

```bash
subs-server \
  --location ./subs \
  --port 2115 \
  --debug
```

### Docker

```bash
docker run -d -p 2115:2115 \
  -e LOCATION=/subs \
  -v $(pwd)/subs:/subs \
  kutovoys/subs-server
```

### Docker Compose

```yaml
services:
  subs-server:
    image: kutovoys/subs-server
    container_name: subs-server
    restart: always
    ports:
      - "2115:2115"
    volumes:
      - ./subs:/subs
    environment:
      LOCATION: /subs
      DEBUG: "true"
```

## ⚙️ Configuration

Configuration can be set via CLI flags or environment variables.

| Parameter             | CLI Flag                    | Environment Variable      | Default                                   | Description                                   |
| --------------------- | --------------------------- | ------------------------- | ----------------------------------------- | --------------------------------------------- |
| Files source location | `--location`, `-l`          | `LOCATION`                | —                                         | Path to the directory with subscription files |
| Source type           | `--source`, `-s`            | `SOURCE`                  | `filesystem`                              | Currently only `filesystem` is supported      |
| Host                  | `--host`, `-h`              | `HOST`                    | `0.0.0.0`                                 | Host for the HTTP server                      |
| Port                  | `--port`, `-p`              | `PORT`                    | `2115`                                    | Port for the HTTP server                      |
| Debug mode            | `--debug`, `-d`             | `DEBUG`                   | `false`                                   | Enables debug output and endpoint listing     |
| Profile title         | `--profile-title`           | `PROFILE_TITLE`           | `Subs-Server`                             | Base64-encoded and sent as a header           |
| Update interval       | `--profile-update-interval` | `PROFILE_UPDATE_INTERVAL` | `12`                                      | Sent as a header, in hours                    |
| Profile page URL      | `--profile-web-page-url`    | `PROFILE_WEB_PAGE_URL`    | `https://github.com/kutovoys/subs-server` | Sent as a header                              |
| Support URL           | `--support-url`             | `SUPPORT_URL`             | `https://github.com/kutovoys/subs-server` | Sent as a header                              |
| Version               | `--version`                 | —                         | —                                         | Print version and exit                        |

## 🛠 Example

If you place a file `test.txt` inside the `./subs` directory with the content:

```
vmess://example
https://some-remote-url.com/encoded.txt
```

and run the server, visiting `http://localhost:2115/test` will return a single base64-encoded response that includes:

- The `vmess://example` line.
- All decoded entries fetched from `https://some-remote-url.com/encoded.txt`.

## 🤝 Contributing

Contributions to Subs Server are warmly welcomed. Whether it's bug fixes, new features, or documentation improvements, your input helps make this project better. Here's a quick guide to contributing:

1. **Fork & Branch**: Fork this repository and create a branch for your work.
2. **Implement Changes**: Work on your feature or fix, keeping code clean and well-documented.
3. **Test**: Ensure your changes maintain or improve current functionality, adding tests for new features.
4. **Commit & PR**: Commit your changes with clear messages, then open a pull request detailing your work.
5. **Feedback**: Be prepared to engage with feedback and further refine your contribution.

Happy contributing! If you're new to this, GitHub's guide on [Creating a pull request](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request) is an excellent resource.

## VPN Recommendation

For secure and reliable internet access, we recommend [BlancVPN](https://getblancvpn.com/?ref=subs). Use promo code `TRYBLANCVPN` for 15% off your subscription.
