<!--
Copyright (C) 2018-2019 Hegemonie's AUTHORS
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
-->
{% include "header_map.tpl" %}
<script>

    var map = "{{ map }}"
    var cities = "{{ cities }}"

    function getTileStyle(tile) {
        tile--;
        var x = (tile % 20) * 5.25;
        var y = Math.floor(tile/20) * 5.25;
        return 'background-position:' + x + '% ' + y + '%'
    }

    function getCityInfo(id) {
        for (cid in cities) {
            if (cities[cid].Meta.Id == id) {
                return cities[cid]
            }
        }
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
        map = JSON.parse(map.replace(/&quot;/g, '\"'))
        cities = JSON.parse(cities.replace(/&quot;/g, '\"'));

        canvas = $("#canvas")
        for (idx in map.Cells) {
            toAppend = "<div class='tile' style='" + getTileStyle(map.Cells[idx].Biome) + "'>";
            if (map.Cells[idx].City > 0) {
                toAppend += "<div class='city' style='" + getTileStyle(121) + "'>" +
                    "<span class='cityName'>" + getCityInfo(map.Cells[idx].City).Meta.Name + "</span>"
                + "</div>"
            }
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
