{% include "header.tpl" %}

<h2>The cities managed by {{Character.Meta.Name}}</h2>
<ul>{% for c in Character.OwnerOf %}
    <li><a href="/game/land?cid={{Character.Meta.Id}}&lid={{c.Id}}">{{c.Name}}</a></li>{% endfor %}
</ul>
<ul>{% for c in Character.DeputyOf %}
    <li><a href="/game/land?cid={{Character.Meta.Id}}&lid={{c.Id}}">{{c.Name}}</a></li>{% endfor %}
</ul>

{% include "footer.tpl" %}
