{% include "header.tpl" %}

<p>{{Flash.InfoMsg}}{{Flash.WarningMsg}}{{Flash.ErrorMsg}}</p>
<p>Online management RPG game</p>
<form action="/action/login" method="post">
    <input type="text" name="email" value=""/>
    <input type="password" name="password" value=""/>
    <input type="submit" value="Enter"/>
</form>

{% include "footer.tpl" %}