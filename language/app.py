from flask import Flask, request, jsonify
import string
import spacy
import torch

if torch.backends.mps.is_available():
    mps_device = torch.device("mps")
    x = torch.ones(1, device=mps_device)
    print(f"🔥 MPS device is available: {x}")
else:
    print ("MPS device not found.")

app = Flask(__name__)
nlp = None
iso6393_to_model = {
    "ara": "xx",
    "chi_sim": "zh",
    "chi_tra": "zh",
    "deu": "de",
    "eng": "en",
    "fra": "fr",
    "hin": "xx",
    "ita": "it",
    "jpn": "ja",
    "nld": "nl",
    "por": "pt",
    "rus": "ru",
    "spa": "es",
    "swe": "sv",
}


@app.route("/v2/health", methods=["GET"])
def health():
    return "OK", 200


@app.route("/v2/entities", methods=["POST"])
def ner_entities():
    global nlp
    global iso6393_to_model

    content = request.json
    text = content["text"]
    language = content["language"]

    entities = []
    if nlp is None:
        nlp = {
            "xx": spacy.load("xx_ent_wiki_sm"),
            "zh": spacy.load("zh_core_web_trf"),
            "de": spacy.load("de_core_news_lg"),
            "en": spacy.load("en_core_web_trf"),
            "fr": spacy.load("fr_core_news_lg"),
            "it": spacy.load("it_core_news_lg"),
            "ja": spacy.load("ja_core_news_trf"),
            "nl": spacy.load("nl_core_news_lg"),
            "pt": spacy.load("pt_core_news_lg"),
            "ru": spacy.load("ru_core_news_lg"),
            "es": spacy.load("es_core_news_lg"),
            "sv": spacy.load("sv_core_news_lg"),
        }
        for key in nlp:
            nlp[key].add_pipe("sentencizer")

    for doc in nlp[iso6393_to_model[language]].pipe([text], disable=["tagger"]):
        for sent in doc.sents:
            for ent in sent.ents:
                entities.append({"text": ent.text, "label": ent.label_})

    # Filter out entities with less than 3 characters and CARDINAL
    entities = [
        entity
        for entity in entities
        if len(entity["text"]) >= 3 and entity["label"] != "CARDINAL"
    ]

    # Group by text and count frequency
    result = {}
    whitespace_and_nonprintable = string.whitespace + "".join(chr(i) for i in range(32))
    for entity in entities:
        entity["text"] = entity["text"].strip(whitespace_and_nonprintable)
        key = entity["text"].lower()
        if key in result:
            result[key]["frequency"] += 1
        else:
            result[key] = {"text": entity["text"], "frequency": 1}
    
    # Convert the dictionary back to a list of entities with the "frequency" field
    result = [{"text": value["text"], "frequency": value["frequency"]} for value in result.values()]

    # Sort by descending order of frequency
    result.sort(key=lambda x: x["frequency"], reverse=True)

    return jsonify(result)
