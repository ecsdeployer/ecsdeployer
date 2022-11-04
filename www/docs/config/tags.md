# Tags

Tags are defined globally, as a list of objects under the `tags` key.

Each tag has 2 required fields: `name` and `value`.

These 2 fields are always interpreted as strings. You can also use [template values](../templating.md) in both fields.


!!! tip 
    If either the `name` or the `value` evaluates to an empty string, then that tag will not be added. You can use this to condition your tags.


```yaml title="Example Tag Declaration"
tags:
  - name: cloud87/billable
    value: true
  
  - name: cloud87/application
    value: "{{ .ProjectName }}"

  - name: "cloud87/{{.ProjectName}}/tag"
    value: "{{ .ImageTag }}"
```

## Fields

[`name`](#tags.name){ #tags.name } - **(required)**

:   The name or key for the tag. You may use [template clauses](../templating.md) in this field.


[`value`](#tags.value){ #tags.value } - **(required)**

:   The value for the tag. You may use [template clauses](../templating.md) in this field.
