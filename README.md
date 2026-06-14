# My Webpage

Nowadays everyone needs a webpage for lots of reasons. Mine is to have my curriculum and to add something every now and then, so I DIY it to _never_ need to write HTML, just Typst!

## Description 


I am not a fronted developer, and I don't aim to be one, but I needed a website. Before letting an LLM make a bloated JavaScript interactive website, i wanted to do something by myself and try to learn _something_ about doing a webpage.

The aim of this projects is that, once the layout is done, never to write any HTML nor CSS again, and to provide a natural way to edit all the potential writings as I do with every other text, using Typst in my terminal text editor.

The solution I landed on is to build an static template site using Go, which conveniently has both an `http.FileServer` and a `template` in its standard library. The program loads the `templates/layout.html` as a template when starting, and loads the contents of the previously compiled Typst contents with HTML, being the location of those provided as an argument to the program. Whenever a document is requested, opens the HTML output of the article in the file, extracts the body, completes the template and serves it to the user.

Fast, simple and reliable :)

## Updating the Contents

As explained in the past section, the article contents are not in this repository. The article contents are stored in another repository to back them up. To edit them or add new ones, I have a bare repository containing the articles in a different remote. That is:
- `origin`: a private GitHub repo that serves as back up.
- `production`: points to a bare repo on the VPS.

Every I push to production, the webpage will automatically update, as the next petition will already pick the correct document when building the template.

## TODO:
1. Document caching: once a file has been served, store it in a `cached/` files. Then, if requested again, just serve that file instead of merging both of those in memory _if the typst file_ has not been changed since last checked time.
2. Language Toggle: I would love to write in Catalan for certain things, or provide two versions of the same writing! I don't know how to do it, but probably involves having a `/ca/` or `/en/` and addressing the about page to one of those.
