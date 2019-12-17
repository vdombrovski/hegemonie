<!DOCTYPE html>
<html lang="fr">
<head><meta charset="UTF-8"><title>Hegemonie</title>
<meta name="author" content="Jean-Francois Smigielski"/>
<meta name="theme-color" content="#FFF"/>
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
<meta name="description" content="${description}"/>
<link rel="stylesheet" href="/static/style.css"/>
</head>
<body>
<header><h1>Hegemonie</h1></header>
<nav>
    {% if Character.Meta.Id %}<a href="/game/character?cid={{ Character.Meta.Id }}">My Character</a>{% endif %}
    {% if userid %}<a href="/game/user">Character</a>{% endif %}
    {% if userid %}<a href="/action/logout">Log-Out</a>{% endif %}
</nav>
<main>