{% include "header.tpl" %}

<h2>The characters managed by {{User.Meta.Name}}</h2>
<ul>{% for c in User.Characters %}
    <li><a href="/game/character?cid={{c.Id}}">{{c.Name}}</a></li>{% endfor %}
</ul>

<h2>Admin</h2>
<p>Logged as {{User.Meta.Name}}.</p>
<p>Your email is {{User.Meta.Email}}</p>
<p>Check <a href="/game/user">your profile</a></p>
<form action="/action/logout" method="post"><input type="submit" value="Log Out"/></form>

{% include "footer.tpl" %}