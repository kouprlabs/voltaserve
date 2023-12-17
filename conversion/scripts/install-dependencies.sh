#!/bin/bash

zypper install -y ffmpeg-4

zypper install -y gawk ghostscript ImageMagick

# https://www.suse.com/support/kb/doc/?id=000019384
policy_path=$(identify -list policy | awk '/Path:/ {print $2}' | sed 's/\[built-in\]//') \
  && awk '/pattern="(PS|PS2|PS3|PDF|XPS|EPS|PCL)"/ { sub(/rights="write"/, "rights=\"read|write\"") } { print }' "$policy_path" > "$policy_path.tmp" \
  && mv "$policy_path.tmp" "$policy_path"

zypper install -y poppler-tools

zypper install -y \
  libreoffice \
  libreoffice-writer \
  libreoffice-calc \
  libreoffice-impress \
  libreoffice-draw \
  libreoffice-math \
  libreoffice-base
