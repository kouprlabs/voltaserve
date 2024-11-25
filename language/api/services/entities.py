# Copyright (c) 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

import string
import re
from spacy.language import Language


class EntityExtractor:
    def __init__(self, nlp: Language):
        self.nlp = nlp

    def run(self, text):
        """
        Run the full extraction pipeline.

        Args
            text (string): A string containing text to extract entities from.

        Returns:
            list: A list of entities
        """
        entities = self.extract_entities(text)
        entities = self.filter_entities(entities)
        groups = self.group_and_count_frequency(entities)
        return self.convert_dict_to_list(groups)

    def extract_entities(self, text):
        """
        Extracts entities from text.

        Args:
            text (string): A string containing text to extract entities from.

        Returns:
            list: A list of entities
        """
        return [
            {"text": re.sub(r"\n+", " ", ent.text).strip(), "label": ent.label_}
            for doc in self.nlp.pipe([text], disable=["tagger"])
            for sent in doc.sents
            for ent in sent.ents
        ]

    @staticmethod
    def filter_entities(entities):
        """
        Filters out entities with less than 3 characters,
        and those with the label 'CARDINAL'.

        Args:
            entities (list): A list of entities to be filtered.

        Returns:
            list: A filtered list of entities.
        """
        return [
            entity
            for entity in entities
            if len(entity["text"]) >= 3 and entity["label"] != "CARDINAL"
        ]

    @staticmethod
    def group_and_count_frequency(entities):
        """
        Groups items by their text representation, and counts
        the frequency of each group.

        Args:
            entities (list): A list of items to be grouped and counted.

        Returns:
            dict: A dictionary where keys are the unique text representations
                and values are the corresponding counts.
        """
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
                result[key] = {
                    "text": entity["text"],
                    "label": entity["label"],
                    "frequency": 1,
                }
        return result

    @staticmethod
    def convert_dict_to_list(groups):
        """
        Convert the dictionary back to a list of entities, with "text" and
        "frequency" fields, then sort by descending order of frequency.

        Args:
            groups (dict): A dictionary where keys are entity names,
                and values are frequencies.

        Returns:
            list: A list of tuples sorted by frequency in descending order.
        """
        result = [
            {
                "text": value["text"],
                "label": value["label"],
                "frequency": value["frequency"],
            }
            for value in groups.values()
        ]
        result.sort(key=lambda x: x["frequency"], reverse=True)
        return result
