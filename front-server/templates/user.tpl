{% include "header.tpl" %}
<h1>Hegemonie</h1>

<h2>Your characters</h2>
<a href="/game/character">Character</a>

<h2>Your cities</h2>
<a href="/game/land">City</a>

<h2>Admin</h2>
<p>Logged as {{userid}}. Check <a href="/game/user">your profile</a></p>
<form action="/action/logout" method="post"><input type="submit" value="Log Out"/></form>
{% include "footer.tpl" %}
