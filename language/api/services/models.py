# Copyright (c) 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

import spacy.cli
import pip
import pkg_resources
import yaml

with open("models.yaml", "r") as f:
    models = yaml.safe_load(f)


nlp = {}
package_max_length = max(len(model["package"]) for model in models.values())
for key in models.keys():
    package = models[key]["package"]
    url = models[key]["url"]

    try:
        pkg_resources.get_distribution(package)
    except pkg_resources.DistributionNotFound:
        pip.main(["install", f"{package} @ {url}"])

    nlp[key] = spacy.load(package)
    nlp[key].add_pipe("sentencizer")

    highlighted_package = f"\033[1m{package.ljust(package_max_length)}\033[0m"
    print(f"🧠 Model {highlighted_package} is ready.")
