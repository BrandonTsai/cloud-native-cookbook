# This is python script to convert markdown2 to html for Darumatic blog

import markdown
from jinja2 import Template
import sys
from bs4 import BeautifulSoup as BSHTML
import datetime

# from xml.etree import ElementTree


if __name__ == "__main__":

    if len(sys.argv) < 2:
        print("Error! Usage: python " + sys.argv[0] + " <markdown filepath>")
        exit(1)

    with open (sys.argv[1], "r") as myfile:
        html = markdown.markdown(myfile.read(), extensions=['tables','codehilite', 'fenced_code'])


    # Replace image src path
    x = datetime.datetime.now()
    postdate = x.strftime("%Y_%m")
    soup = BSHTML(html, features="html.parser")
    images = soup.findAll('img')
    for image in images:
        img_name = image['src'].split("/")[1]
        html = html.replace(image['src'],"https://darumatic.com/media/blog_pics/"+ postdate + "/" +img_name)

    # Remove h1 title
    titles = soup.findAll('h1')
    for title in titles:
        html = html.replace(str(title), "")

    # write to new html file
    with open('tools/html_template.j2') as tmp:
        template = Template(tmp.read())
        rendered = template.render(content=html)

    output_path = "html/" + sys.argv[1].replace("./", "").replace("/", "_").replace(".md", ".html")
    print(output_path)
    # print(rendered)
    with open(output_path, "w") as f:
        f.write(rendered)