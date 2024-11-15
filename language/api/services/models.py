import spacy.cli
import pip
import pkg_resources

models = {
    "eng": {
        "package": "en_core_web_sm",
        "url": "https://huggingface.co/spacy/en_core_web_sm/resolve/main/en_core_web_sm-any-py3-none-any.whl",
    },
    "deu": {
        "package": "de_core_news_sm",
        "url": "https://huggingface.co/spacy/de_core_news_sm/resolve/main/de_core_news_sm-any-py3-none-any.whl",
    },
    "fra": {
        "package": "fi_core_news_sm",
        "url": "https://huggingface.co/spacy/fi_core_news_sm/resolve/main/fi_core_news_sm-any-py3-none-any.whl",
    },
    "ita": {
        "package": "it_core_news_sm",
        "url": "https://huggingface.co/spacy/it_core_news_sm/resolve/main/it_core_news_sm-any-py3-none-any.whl",
    },
    "jpn": {
        "package": "ja_core_news_sm",
        "url": "https://huggingface.co/spacy/ja_core_news_sm/resolve/main/ja_core_news_sm-any-py3-none-any.whl",
    },
    "nld": {
        "package": "nl_core_news_sm",
        "url": "https://huggingface.co/spacy/nl_core_news_sm/resolve/main/nl_core_news_sm-any-py3-none-any.whl",
    },
    "por": {
        "package": "pt_core_news_sm",
        "url": "https://huggingface.co/spacy/pt_core_news_sm/resolve/main/pt_core_news_sm-any-py3-none-any.whl",
    },
    "spa": {
        "package": "es_core_news_sm",
        "url": "https://huggingface.co/spacy/es_core_news_sm/resolve/main/es_core_news_sm-any-py3-none-any.whl",
    },
    "swe": {
        "package": "sv_core_news_sm",
        "url": "https://huggingface.co/spacy/sv_core_news_sm/resolve/main/sv_core_news_sm-any-py3-none-any.whl",
    },
    "nor": {
        "package": "nb_core_news_sm",
        "url": "https://huggingface.co/spacy/nb_core_news_sm/resolve/main/nb_core_news_sm-any-py3-none-any.whl",
    },
    "fin": {
        "package": "fi_core_news_sm",
        "url": "https://huggingface.co/spacy/fi_core_news_sm/resolve/main/fi_core_news_sm-any-py3-none-any.whl",
    },
    "dan": {
        "package": "da_core_news_sm",
        "url": "https://huggingface.co/spacy/da_core_news_sm/resolve/main/da_core_news_sm-any-py3-none-any.whl",
    },
    "chi_sim": {
        "package": "zh_core_web_sm",
        "url": "https://huggingface.co/spacy/zh_core_web_sm/resolve/main/zh_core_web_sm-any-py3-none-any.whl",
    },
    "chi_tra": {
        "package": "zh_core_web_sm",
        "url": "https://huggingface.co/spacy/zh_core_web_sm/resolve/main/zh_core_web_sm-any-py3-none-any.whl",
    },
    "rus": {
        "package": "ru_core_news_sm",
        "url": "https://huggingface.co/spacy/ru_core_news_sm/resolve/main/ru_core_news_sm-any-py3-none-any.whl",
    },
    "hin": {
        "package": "xx_ent_wiki_sm",
        "url": "https://huggingface.co/spacy/xx_ent_wiki_sm/resolve/main/xx_ent_wiki_sm-any-py3-none-any.whl",
    },
    "ara": {
        "package": "xx_ent_wiki_sm",
        "url": "https://huggingface.co/spacy/xx_ent_wiki_sm/resolve/main/xx_ent_wiki_sm-any-py3-none-any.whl",
    },
}


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
