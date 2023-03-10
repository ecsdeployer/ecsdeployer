# Templating

The templating engine is powered by [Go Templates](https://pkg.go.dev/text/template#hdr-Actions). You can use any normal Go template operator within your template strings.


## Fields


### Common Fields

These fields are available to all templates, everywhere. Note they must be wrapped in double curly braces.

<div class="tbl-nowrap-key" markdown>

Key             | Description
----------------|---------------
`.Project`      | the project name
`.Env.VARNAME`  | a Map with all the current environment variables
`.Date`         | current UTC date in RFC 3339 format
`.Timestamp`    | current UTC time in Unix format
`.AppVersion`   | the value you provided to `--app-version` (if you provided it)
`.Version`      | an alias for `.AppVersion`
`.ImageTag`     | the value you provided to `--tag`
`.Tag`          | alias for `.ImageTag`
`.Image`        | the value you provided to `--image` (a container image URI)
`.Cluster`      | the name of the ECS Cluster the app will be deployed on
`.Stage`        | the stage name for the application (i.e. "production", "staging", etc). You can specify this with `stage: VALUE` in the config file.<br><br>**IMPORTANT:** Specifying a `stage` in your file _WILL_ modify the naming conventions for your application. If you have already deployed, do not add this.<br>See [Naming](config/naming.md#fields) to see how this will change the names of resources.
`AwsAccountId`  | numeric AWS Account number (note the lack of `.` at the start)
`AwsRegion`     | current AWS region (note the lack of `.` at the start)

</div>

The following sections denote fields that are only available in certain contexts.

### Task Related

<small>These are only available within sections related to individual tasks (CronJobs, PreDeploy, Services).</small>

<div class="tbl-nowrap-key" markdown>

Key          | Description
-------------|---------------
`.Arch`      | The architecture of the task. `amd64` or `arm64`
`.Name`      | the task name you are referencing
`.Container` | the name of the individual container (will be the same as `.Name` for the primary container)

</div>

## Functions
For all fields, you can use the following functions:

<div class="tbl-nowrap-key" markdown>

Usage                   |Description
------------------------|-----------------
`join "sep" "x" "y" "z"`|concatenates the 2nd thru last parameter using the 1st as a separator
`prefix "value" 4`      |only returns the first N characters of a string
`replace "v1.2" "v" ""` |replaces all matches. See [ReplaceAll](https://golang.org/pkg/strings/#ReplaceAll)
`split "1.2" "."`       |split string at separator. See [Split](https://golang.org/pkg/strings/#Split)
`time "01/02/2006"`     |current UTC time in the specified format (this is not deterministic, a new time will be returned for every call)
`tolower "V1.2"`        |makes input string lowercase. See [ToLower](https://golang.org/pkg/strings/#ToLower)
`toupper "v1.2"`        |makes input string uppercase. See [ToUpper](https://golang.org/pkg/strings/#ToUpper)
`trim " v1.2  "`        |removes all leading and trailing white space. See [TrimSpace](https://golang.org/pkg/strings/#TrimSpace)
`trimprefix "v1.2" "v"` |removes provided leading prefix string, if present. See [TrimPrefix](https://golang.org/pkg/strings/#TrimPrefix)
`trimsuffix "1.2v" "v"` |removes provided trailing suffix string, if present. See [TrimSuffix](https://pkg.go.dev/strings#TrimSuffix)

</div>