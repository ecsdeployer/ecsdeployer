import logging, re
import mkdocs.plugins
import json

log = logging.getLogger('mkdocs')

schemaDefaultsJson = json.load(open('data/defaults.json'))

@mkdocs.plugins.event_priority(-50)
def on_page_markdown(markdown, page, **kwargs):
    path = page.file.src_uri
    # {{schema:default:NameTemplates.PROPNAME}}
    for m in re.finditer(r'\{\{schema:default:([a-zA-Z_]+)\.([a-zA-Z_]+)\}\}', markdown):
        # value = schemaDefaultsJson["$defs"][m[1]]["properties"][m[2]]["default"]
        value = schemaDefaultsJson[m[1]][m[2]]
        markdown = markdown.replace(f"{m[0]}", f"{value}")
    
    return markdown