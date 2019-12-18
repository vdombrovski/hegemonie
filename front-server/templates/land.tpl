{% include "header.tpl" %}

<h2>You are in {{Land.Meta.Name}}</h2>
<p><a href="/game/land?lid={{Land.Meta.Id}}">City</a></p>

<h2>Production</h2>
<ul>
    {% for r in Land.Meta.Production %}
    <li>{{r}}</li>
    {% endfor %}
</ul>

<h2>Stock</h2>
<ul>
    {% for r in Land.Meta.Stock %}
    <li>{{r}}</li>
    {% endfor %}
</ul>

<h2>Your buildings</h2>
<ul>
    {% for b in Land.Buildings %}
    <li>{{b.Name}} (id {{b.Id}})</li>
    {% endfor %}
</ul>

<h2>Troops defending</h2>
<ul>
    {% for u in Land.Units %}
    <li>{{u.Name}} (id {{u.Id}})</li>
    {% endfor %}
</ul>

{% include "footer.tpl" %}