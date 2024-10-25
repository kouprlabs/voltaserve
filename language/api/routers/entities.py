# Copyright 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

from flask import Blueprint, request, jsonify
import spacy.cli
from ..services.entities import EntityExtractor

bp = Blueprint("entities", __name__)

multi_language_model = "xx_ent_wiki_sm"
spacy.cli.download(multi_language_model)

nlp = spacy.load(multi_language_model)
nlp.add_pipe("sentencizer")


@bp.route("/v3/entities", methods=["POST"])
def get_entities():
    global nlp

    content = request.json
    text = content["text"]

    entity_extractor = EntityExtractor(nlp)
    dtos = entity_extractor.run(text)

    return jsonify(dtos)
