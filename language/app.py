from flask import Flask, request, jsonify
import spacy

app = Flask(__name__)
nlp = None


@app.route("/v1/health", methods=["GET"])
def health():
    return "OK", 200


@app.route("/v1/entities", methods=["POST"])
def ner_entities():
    global nlp
    content = request.json
    text = content["text"]

    entities = []
    if nlp is None:
        nlp = spacy.load("xx_ent_wiki_sm")
        nlp.add_pipe("sentencizer")

    for doc in nlp.pipe([text], disable=["tagger"]):
        for sent in doc.sents:
            for ent in sent.ents:
                entities.append({"text": ent.text, "label": ent.label_})

    return jsonify(entities)
