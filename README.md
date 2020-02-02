# Cloud Native Cookbook

Export markdown 2.0 to HTML
---------------------------

- build image: `docker build -f tools/Dockerfile -t md2html .`
- run container `./docker_run.sh`
- in container: `python tools/md2html.py <markdown file>`

This script will convert markdown file to HTML file in *./html/* folder.


Tool for darwing:
-----------------

- https://www.draw.io/


Tips:
------------------

- No more than 2 images in one article.