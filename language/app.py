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
iso_6393_to_model = {
    "ara": "xx_ent_wiki_sm",
    "chi_sim": "zh_core_web_lg",
    "chi_tra": "zh_core_web_lg",
    "deu": "de_core_news_lg",
    "eng": "en_core_web_lg",
    "fra": "fr_core_news_sm",
    "hin": "xx_ent_wiki_sm",
    "ita": "it_core_news_lg",
    "jpn": "ja_core_news_lg",
    "nld": "nl_core_news_lg",
    "por": "pt_core_news_lg",
    "rus": "ru_core_news_lg",
    "spa": "es_core_news_lg",
    "swe": "sv_core_news_lg",
    "nor": "nb_core_news_lg",
    "fin": "fi_core_news_lg",
    "dan": "da_core_news_lg",
}
nlp = {}
for key in iso_6393_to_model.keys():
    model = iso_6393_to_model[key]
    spacy.cli.download(model)
    nlp[key] = spacy.load(model)
    nlp[key].add_pipe("sentencizer")


@app.route("/v3/health", methods=["GET"])
def health():
    return "OK", 200


@app.route("/version", methods=["GET"])
def version():
    return {"version": "3.0.0"}


@app.route("/v3/entities", methods=["POST"])
def ner_entities():
    global nlp
    global iso_6393_to_6391

    content = request.json
    text = content["text"]
    language = content["language"]

    entities = [
        {"text": ent.text, "label": ent.label_}
        for doc in nlp[language].pipe([text], disable=["tagger"])
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
