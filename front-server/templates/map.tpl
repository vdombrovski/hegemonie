<!--
Copyright (C) 2018-2019 Hegemonie's AUTHORS
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
-->
{% include "header_map.tpl" %}
<script>
    function getTileStyle(tile) {
        tile--;
        var x = (tile % 20) * 5.25;
        var y = Math.floor(tile/20) * 5.25;
        return 'background-position:' + x + '% ' + y + '%'
    }

    function enablePanDrag() {
      var current = [0, 0];
      var canvas = $("#canvas");
      const ROT = "rotateX(64deg) rotateY(0deg) rotateZ(-45deg)";
      canvas.panzoom({'cursor':'default', 'easing': null, 'disableZoom': true, transition: false, onPan: function(e, panzoom) {
        var matrix = panzoom.getMatrix();
        canvas.css(
            {'transform': 'translate('+ (parseInt(matrix[4]) + parseInt(current[0])) + 'px,'+ (parseInt(matrix[5]) + parseInt(current[1])) +'px) ' + ROT})
      }})

      canvas.on('panzoomend', function(e, panzoom, matrix, changed) {
        current = [parseInt(matrix[4]) + parseInt(current[0]),   parseInt(matrix[5]) + parseInt(current[1])]
        canvas.css({'transform': 'translate('+ current[0] + 'px,'+ current[1] +'px) ' + ROT})
      });
    }

    $(document).ready(function() {
        var map = "{{ map }}"
        map = JSON.parse(map.replace(/&quot;/g, '\"'))
        var overlay = "{{ overlay }}"
        overlay = JSON.parse(overlay.replace(/&quot;/g, '\"'))

        console.log(overlay)

        canvas = $("#canvas")
        for (idx in map) {
            toAppend = "<div class='tile' style='" + getTileStyle(map[idx].Terrain) + "'>";
            if (overlay[idx].Terrain > 0)
                toAppend += "<div class='overlay' style='" + getTileStyle(overlay[idx].Terrain) + "'></div>"

            toAppend += "</div>"
            canvas.append(toAppend)
        }


        enablePanDrag()
    })
</script>
<div class="frame">
    <div id="canvas" class="canvas"></div>
    <div class="sidebar">

    </div>
</div>
{% include "footer.tpl" %}
