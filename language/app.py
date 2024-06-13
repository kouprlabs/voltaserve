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
iso_6393_to_model = {
    "ara": "xx_ent_wiki_sm",
    "chi_sim": "xx_ent_wiki_sm",
    "chi_tra": "xx_ent_wiki_sm",
    "deu": "xx_ent_wiki_sm",
    "eng": "xx_ent_wiki_sm",
    "fra": "xx_ent_wiki_sm",
    "hin": "xx_ent_wiki_sm",
    "ita": "xx_ent_wiki_sm",
    "jpn": "xx_ent_wiki_sm",
    "nld": "xx_ent_wiki_sm",
    "por": "xx_ent_wiki_sm",
    "rus": "xx_ent_wiki_sm",
    "spa": "xx_ent_wiki_sm",
    "swe": "xx_ent_wiki_sm",
}


@app.route("/v2/health", methods=["GET"])
def health():
    return "OK", 200


@app.route("/v2/entities", methods=["POST"])
def ner_entities():
    global nlp
    global iso_6393_to_6391

    content = request.json
    text = content["text"]
    language = content["language"]

    entities = []
    if nlp is None:
        nlp = {}
        for key in iso_6393_to_model.keys():
            nlp[key] = None
    if nlp[language] is None:
        nlp[language] = spacy.load(iso_6393_to_model[language])
        nlp[language].add_pipe("sentencizer")

    for doc in nlp[language].pipe([text], disable=["tagger"]):
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
