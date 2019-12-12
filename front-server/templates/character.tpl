{% include "header.tpl" %}
<h1>Hegemonie</h1>

<h2>Your are...</h2>
<p><a href="/game/character">Character</a></p>

<h2>Your cities</h2>
<li>
    <ul><a href="/game/land">City</a></ul>
</li>

<h2>Admin</h2>
<p>Logged as {{userid}}. Check <a href="/game/user">your profile</a></p>
<form action="/action/logout" method="post"><input type="submit" value="Log Out"/></form>
{% include "footer.tpl" %}
