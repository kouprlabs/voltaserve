from flask import Flask, request, jsonify
import string
import spacy

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
        if nlp is None:
            nlp = {
                "xx": spacy.load("xx_ent_wiki_sm"),
                "zh": spacy.load("zh_core_web_trf"),
                "de": spacy.load("de_dep_news_trf"),
                "en": spacy.load("en_core_web_trf"),
                "fr": spacy.load("fr_dep_news_trf"),
                "it": spacy.load("it_core_news_lg"),
                "ja": spacy.load("ja_core_news_trf"),
                "nl": spacy.load("nl_core_news_lg"),
                "pt": spacy.load("pt_core_news_lg"),
                "ru": spacy.load("ru_core_news_lg"),
                "es": spacy.load("es_dep_news_trf"),
                "sv": spacy.load("sv_core_news_lg"),
            }
        nlp[iso6393_to_model[language]].add_pipe("sentencizer")

    for doc in nlp[iso6393_to_model[language]].pipe([text], disable=["tagger"]):
        for sent in doc.sents:
            for ent in sent.ents:
                entities.append({"text": ent.text, "label": ent.label_})

    grouped_entities = {}
    for entity in entities:
        entity["text"] = entity["text"].strip(
            string.whitespace + "".join(chr(i) for i in range(32))
        )
        key = (entity["text"], entity["label"])
        if key in grouped_entities:
            grouped_entities[key] += 1
        else:
            grouped_entities[key] = 1

    result = [
        {"text": text, "label": label, "frequency": frequency}
        for (text, label), frequency in grouped_entities.items()
    ]

    filtered_result = [
        entity
        for entity in result
        if len(entity["text"]) >= 3 and entity["label"] != "CARDINAL"
    ]

    return jsonify(filtered_result)
