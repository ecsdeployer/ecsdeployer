---
hide:
  toc: true
---
# ecsdeployer clean

Runs the cleanup step only. Skips actual deployment

## Synopsis

Use this command to purge any unused services, cronjobs, task definitions, etc 
from your environment that are no longer being referenced in your configuration file.


```
ecsdeployer clean [flags]
```

## Options

```
      --app-version string   Set the application version. Useful for templates
  -c, --config string        Configuration file to check
  -h, --help                 help for clean
      --image string         Specify a container image URI.
  -q, --quiet                Quiet mode: no output
      --tag string           Specify a custom image tag to use.
      --timeout duration     Timeout for the entire cleanup process (default 30m0s)
```

## Global Options

```
      --debug   Enable debug mode
```

## See also

* [`ecsdeployer`](ecsdeployer.md)	 - Deploy applications to Fargate

