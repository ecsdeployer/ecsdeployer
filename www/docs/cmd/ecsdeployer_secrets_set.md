---
hide:
  toc: true
---
# ecsdeployer secrets set

Set (or delete) a secret

```
ecsdeployer secrets set VARIABLE_NAME {VALUE | --file filename | --stdin | --unset} [flags]
```

## Options

```
  -c, --config file   Configuration file to check
      --file file     Read value from file
      --force         Do not ask for confirmation
  -h, --help          help for set
      --keyid KeyID   Specify the KMS KeyID to use. If not provided, will use the default.
      --stdin         Get value from stdin
      --unset         Removes a variable entirely
```

## Global Options

```
      --debug   Enable debug mode
```

## See also

* [`ecsdeployer secrets`](ecsdeployer_secrets.md)	 - Manage application secrets in SSM

