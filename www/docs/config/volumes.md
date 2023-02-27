# Volumes and Mounts

Volumes are specified at the task level. The recommended place to specify all your volumes is within the [`task_defaults`](defaults.md) block. Any volumes that are not referenced within a task's containers will be ignored and not included in the task definition. This means you can specify all possible volumes within the task defaults and just reference them in various containers using mount points.

## Example Usage

=== "Using Elastic File System"

    ```yaml
    task_defaults:
      volumes:
        - name: archive
          efs:
            file_system_id: fs-1234567
            access_point_id: fsap-abc12345

      mounts:
        - path: /mnt/archive
          source: archive
    ```

=== "Using Bind Volumes"

    ```yaml
    task_defaults:
      volumes:
        - name: bindvol
        - other-bind-vol # short hand

      mounts:
        - path: /mnt/shared
          source: bindvol

        - path: /mnt/roshared
          source: other-bind-vol
          readonly: true
    ```



## Specifying Volumes

Volumes are specified as an array of objects with the following fields:

[`name`](#volumes.name){ #volumes.name } - **(required)**

:   The name of the volume.

[`efs`](#volumes.efs){ #volumes.efs }

:   Adding this block marks this as an EFS Volume reference.
    
    See [EFS Volumes](#efs-volumes) below for details.

### EFS Volumes

[`file_system_id`](#efs.file_system_id){ #efs.file_system_id } - **(required)**

:   The ID of the EFS file system you want to reference.

[`access_point_id`](#efs.access_point_id){ #efs.access_point_id } - **(recommended)**

:   The [EFS Access Point ID](https://docs.aws.amazon.com/efs/latest/ug/efs-access-points.html). This is the recommended way of using an EFS volume.

[`root`](#efs.root){ #efs.root }

:   The directory within the EFS file system to mount as the root directory inside the host.

    Do not specify this if you are using an AccessPoint.

[`use_iam`](#efs.use_iam){ #efs.use_iam }

:   Whether to use the task's role when mounting the EFS Volume.

    _Default_: `{{schema:default:VolumeEFSConfig.use_iam}}`

[`transit_encryption`](#efs.transit_encryption){ #efs.transit_encryption }

:   Whether or not to enable encryption for data in transit. This will be forcibly enabled if [`use_iam`](#volumes.use_iam) or [`access_point_id`](#volumes.access_point_id) are set.

    You should leave this enabled. There is no reason to ever disable it.

    _Default_: `{{schema:default:VolumeEFSConfig.transit_encryption}}`


## Specifying Mount Points

Mounts specified at the task level will be applied to the primary container only. If you need to use a mount in a sidecar container, you will need to specify it within the sidecar object.


[`path`](#mounts.path){ #mounts.path } - **(required)**

:   The path inside the container where this volume will be mounted.

[`source`](#mounts.source){ #mounts.source } - **(required)**

:   The name of the volume that this mount references.
    
    Must match the [`name`](#volumes.name) used when declaring the volume.

[`readonly`](#mounts.readonly){ #mounts.readonly }

:   Whether to mount the volume as read-only.

    _Default_: `{{schema:default:Mount.readonly}}`

## See Also

* [Using Data Volumes](https://docs.aws.amazon.com/AmazonECS/latest/userguide/using_data_volumes.html)
* [EFS Volumes](https://docs.aws.amazon.com/AmazonECS/latest/userguide/efs-volumes.html)
* [Bind Mounts](https://docs.aws.amazon.com/AmazonECS/latest/userguide/bind-mounts.html)
* [Working with Access Points](https://docs.aws.amazon.com/efs/latest/ug/efs-access-points.html)