List font packages in SLE / openSUSE Leap:

```shell
zypper search --type package --not-installed-only '*-fonts' | awk '{ if ($2 ~ /.*-fonts$/) print $2 }'
```

List font packages in RHEL:

```shell
dnf search '*-fonts' --repo=rhel-9-for-x86_64-appstream-rpms | awk -F ':' '{gsub(/\..*/, "", $1); print $1}' | grep -v '^='
```
