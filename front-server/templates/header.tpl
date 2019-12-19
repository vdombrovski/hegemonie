<!--
Copyright (C) 2018-2019 Hegemonie's AUTHORS
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
-->
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
    {% if userid %}<a href="/game/user">User</a>{% endif %}
    {% if cid %}<a href="/game/character?uid={{ userid }}&cid={{ cid }}">Character</a>{% endif %}
    {% if lid %}<a href="/game/land?uid={{ userid }}&cid={{ cid }}&lid={{ lid }}">Land</a>{% endif %}
    {% if userid %}<a href="/action/logout">Log-Out</a>{% endif %}
</nav>
<main>
