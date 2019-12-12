{% include "header.tpl" %}
<h1>Hegemonie</h1>
<p>Online management RPG game</p>
<form action="/action/login" method="post">
    <input type="text" name="userid" value=""/>
    <input type="password" name="passwd" value=""/>
    <input type="submit" value="Log-in"/>
</form>
{% include "footer.tpl" %}
