# yum2npm

This service transforms metadata from YUM repositories to a
format compatible with NPM registries to be used with tools
like [Renovate](https://github.com/renovatebot/renovate).

## API

| Path | Description |
|---|---|
| `/repos` | Get available repositories |
| `/repos/{repo}/packages` | Get packages available in a repository |
| `/repos/{repo}/packages/{package}` | Get available versions of a package in a format like a NPM registry |
| `/repos/{repo}/modules` | Get modules in a repository |
| `/repos/{repo}/modules/{module}/packages` | Get packages in a module. The module parameter must be in the format "`{name}:{stream}`" |
| `/repos/{repo}/modules/{module}/packages/{package}` | Get available versions of a package in a module in a format like a NPM registry |

## Config

The config is stored in a yaml file. The path to the config defaults to `/etc/yum2npm/config.yaml`,
but can be changed via the command line argument `-c` or `--config`.

The default config is in this repo under `configs/config.yaml`.

```yaml
---

http:
  host: "0.0.0.0"
  port: "8080"

refreshInterval: "1h"

repos:
  - name: rocky-9-baseos-x86_64
    url: https://dl.rockylinux.org/pub/rocky/9/BaseOS/x86_64/os/
  - name: rocky-9-appstream-x86_64
    url: https://dl.rockylinux.org/pub/rocky/9/AppStream/x86_64/os/
  - name: rocky-9-extras-x86_64
    url: https://dl.rockylinux.org/pub/rocky/9/extras/x86_64/os/
```

## Usage example

Below is an example on how to use Renovate to update yum packages.

### renovate.json

```json
{
    "$schema": "https://docs.renovatebot.com/renovate-schema.json",
    "regexManagers": [
        {
            "fileMatch": ["Dockerfile$"],
            "matchStrings": [
                "# renovate: datasource=yum repo=(?<registryUrl>[^\\s]+)\\s+(?<depName>[^\\s]+)-(?<currentValue>[^\\s-]+-[^\\s-]+)"
            ],
            "datasourceTemplate": "npm",
            "versioningTemplate": "loose",
            "registryUrlTemplate": "https://yum2npm.io/repos/{{#if (containsString registryUrl '/')}}{{{replace '/' '/modules/' registryUrl}}}{{else}}{{registryUrl}}{{/if}}/packages"
        }
    ]
}
```

### Dockerfile

```Dockerfile
FROM rockylinux:9-minimal

RUN microdnf install -y \
    # renovate: datasource=yum repo=rocky-9-appstream-x86_64/nodejs:18
    nodejs-18.12.1-1.module+el9.1.0+13234+90e40c60 \
    # renovate: datasource=yum repo=rocky-9-appstream-x86_64
    java-17-openjdk-headless-17.0.5.0.8-2.el9_0
```
