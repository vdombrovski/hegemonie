{% include "header.tpl" %}
<h1>Hegemonie</h1>

<h2>Hello, {{Character.Name}}</h2>
<p>Check <a href="/game/character?pid={{Character.Id}}">your profile</a></p>

<h2>You are in {{Land.Name}}</h2>
<p><a href="/game/land?lid={{Land.Id}}">City</a></p>

<h2>Your Production</h2>
{% for r in Land.Production %}
    <li>{{r.Amount}} {{r.Name}}</li>
{% endfor %}

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

<h2>Admin</h2>
<p>Logged as {{userid}}. Check <a href="/game/user">your profile</a></p>
<form action="/action/logout" method="post"><input type="submit" value="Log Out"/></form>
{% include "footer.tpl" %}
