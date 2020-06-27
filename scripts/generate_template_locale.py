# Author: KotRik
import re
import os
import io

PATH_FOR_SEARCH = "..\\templates"
T_SEARCH = r"\{\{(\snavbarItem\s\..*\(|\s)(\.T|\$global\.T|\$\.T)\s\"(.+?[^\\])\".+?\}\}"

paths = []
for i in os.walk(PATH_FOR_SEARCH):
    for filep in i[2]:
        paths.append(f"{i[0]}\\{filep}")

raw_strings = []
strings = {}
for path in paths:
    text = io.open(path, encoding="utf-8", mode="r").read()
    stringers = re.findall(T_SEARCH, text)
    for string in stringers:
        if string[2] in raw_strings: continue

        raw_strings.append(string[2])
        if path not in strings:
            strings[path] = []

        strings[path].append(string[2])

localeFiles = io.open("templates.pot", encoding="utf-8", mode="w")

for (k, v) in strings.items():
    localeFiles.write(f"\n\n# Translation for {k}")
    for string in v:
        localeFiles.write(f'''\n\nmsgid "{string}"
msgstr ""''')

localeFiles.close()
