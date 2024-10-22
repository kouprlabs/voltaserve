# Copyright 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

from flask import Flask, request, jsonify
import string
import spacy.cli

app = Flask(__name__)
spacy.cli.download("xx_ent_wiki_sm")
nlp = spacy.load("xx_ent_wiki_sm")
nlp.add_pipe("sentencizer")


@app.route("/v3/health", methods=["GET"])
def health():
    return "OK", 200


@app.route("/version", methods=["GET"])
def version():
    return {"version": "3.0.0"}


@app.route("/v3/entities", methods=["POST"])
def ner_entities():
    global nlp

    content = request.json
    text = content["text"]

    entities = [
        {"text": ent.text, "label": ent.label_}
        for doc in nlp.pipe([text], disable=["tagger"])
        for sent in doc.sents
        for ent in sent.ents
    ]

    # Filter out entities with less than 3 characters and CARDINAL
    entities = [
        entity
        for entity in entities
        if len(entity["text"]) >= 3 and entity["label"] != "CARDINAL"
    ]

    # Group by text and count frequency
    result = {}
    whitespace_and_non_printable = string.whitespace + "".join(
        chr(i) for i in range(32)
    )
    for entity in entities:
        entity["text"] = entity["text"].strip(whitespace_and_non_printable)
        key = entity["text"].lower()
        if key in result:
            result[key]["frequency"] += 1
        else:
            result[key] = {"text": entity["text"], "frequency": 1}

    # Convert the dictionary back to a list of entities with the "frequency" field
    result = [
        {"text": value["text"], "frequency": value["frequency"]}
        for value in result.values()
    ]

    # Sort by descending order of frequency
    result.sort(key=lambda x: x["frequency"], reverse=True)

    return jsonify(result)
