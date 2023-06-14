from flask import Flask, request, jsonify
from iso639 import Lang
import spacy
import spacy_fastlang  # noqa: F401 # pylint: disable=unused-import

app = Flask(__name__)

nlp = spacy.blank("xx")
nlp.add_pipe("language_detector")


@app.route("/v1/health", methods=["GET"])
def health():
    return "OK"


@app.route("/v1/detect", methods=["POST"])
def detect():
    content = request.json
    text = content["text"]
    doc = nlp(text)
    result = {"language": Lang(doc._.language).pt3, "score": doc._.language_score}
    return jsonify(result)
