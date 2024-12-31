# Copyright (c) 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

import pip
import pkg_resources
import spacy.cli
import yaml

with open("models.yaml", "r") as file:
    models = yaml.safe_load(file)


def highlight(text):
    return f"\033[1m{text}\033[0m"


nlp = {}
padding = max(len(model["package"]) for model in models.values())
for key in models.keys():
    package = models[key]["package"]
    url = models[key]["url"]

    try:
        pkg_resources.get_distribution(package)
    except pkg_resources.DistributionNotFound:
        pip.main(["install", f"{package} @ {url}"])

    nlp[key] = spacy.load(package)
    nlp[key].add_pipe("sentencizer")

    print(f"🧠 Loaded model {highlight(package.ljust(padding))} for language {highlight(key)}.")
