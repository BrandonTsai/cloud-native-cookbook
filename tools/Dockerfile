FROM python

RUN pip install markdown Jinja2 beautifulsoup4

WORKDIR /opt

CMD find ./blogs/ -name "*.md" | xargs -n 1 python tools/md2html.py