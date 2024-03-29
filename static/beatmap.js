(function () {
  var mapset = {};
  setData.ChildrenBeatmaps.forEach(function (diff) {
    mapset[diff.BeatmapID] = diff;
  });
  console.log(mapset);
  function loadLeaderboard(b, m, mode) {
    var wl = window.location;
    window.history.replaceState('', document.title,
      "/b/" + b + "?mode=" + m + "&mod=" + mode + wl.hash);

    api("scores?sort=score,desc&sort=id,asc", {
      mode: m,
      gamemode: mode,
      b: b,
      p: 1,
      l: 50,
    },
      function (data) {
        console.log(data);
        var tb = $(".ui.table tbody");
        tb.find("tr").remove();
        if (data.scores == null) {
          data.scores = [];
        }
        var i = 0;
        data.scores.sort(function (a, b) { return b.score - a.score; });
        data.scores.forEach(function (score) {
          var user = score.user;
          tb.append($("<tr />").append(
            $("<td data-sort-value=" + (++i) + " />")
              .text("#" + ((page - 1) * 50 + i)),
            $("<td />").html("<a href='/u/" + user.id +
              "' title='View profile'><i class='" +
              user.country.toLowerCase() + " flag'></i>" +
              escapeHTML(user.username) + "</a>"),
            $("<td data-sort-value=" + score.score + " />")
              .html(addCommas(score.score)),
            $("<td />").html(getScoreMods(score.mods, true)),
            $("<td data-sort-value=" + score.accuracy + " />")
              .text(score.accuracy.toFixed(2) + "%"),
            $("<td data-sort-value=" + score.max_combo + " />")
              .text(addCommas(score.max_combo)),
            $("<td data-sort-value=" + score.pp + " />")
              .html(score.pp.toFixed(2))));
        });
      });
  }

  function changeDifficulty(bid) {
    // load info
    var diff = mapset[bid];

    // column 2
    $("#cs").html(diff.CS);
    $("#hp").html(diff.HP);
    $("#od").html(diff.OD);
    $("#passcount").html(addCommas(diff.Passcount));
    $("#playcount").html(addCommas(diff.Playcount));

    // column 3
    $("#ar").html(diff.AR);
    $.getJSON(`https://osu-pp-calc-api.glitch.me/?id=${bid}`, function (res) {
      $("#stars").html(diff.DifficultyRating.toFixed(2) + ` (${res.pp_100}pp)`);
    });
    $("#length").html(timeFormat(diff.TotalLength));
    $("#drainLength").html(timeFormat(diff.HitLength));
    $("#bpm").html(diff.BPM);

    // hide mode for non-std maps
    console.log("favMode", favMode);
    if (diff.Mode != 0) {
      currentMode = (currentModeChanged ? currentMode : favMode);
      $("#mode-menu").hide();
    } else {
      currentMode = diff.Mode;
      $("#mode-menu").show();
    }

    const urlParams = new URLSearchParams(window.location.search);
    const modParam = urlParams.get('mod');
    currentMod = (modParam !== null && modParam.length > 0) ? modParam : "nomod"

    // update mode menu
    $("#mode-menu .active.item").removeClass("active");
    $("#mode-" + currentMode).addClass("active");

    $("#sm-menu .active.item").removeClass("active");
    $("#mode-" + currentMod).addClass("active");

    loadLeaderboard(bid, currentMode, currentMod);
  }
  window.loadLeaderboard = loadLeaderboard;
  window.changeDifficulty = changeDifficulty;
  changeDifficulty(beatmapID);
  // loadLeaderboard(beatmapID, currentMode);

  function smMenuHandler(e) {
    e.preventDefault();
    $("#sm-menu .active.item").removeClass("active");
    $(this).addClass("active");
    console.log("clicked");
    console.log($(this).data)
    currentMod = $(this).data("mod");
    loadLeaderboard(beatmapID, currentMode, currentMod);
  }

  $("#diff-menu .item")
    .click(function (e) {
      e.preventDefault();
      $(this).addClass("active");
      beatmapID = $(this).data("bid");
      changeDifficulty(beatmapID);
    });
  $("#mode-menu .item")
    .click(function (e) {
      e.preventDefault();
      $("#mode-menu .active.item").removeClass("active");
      $(this).addClass("active");
      currentMode = $(this).data("mode");
      currentModeChanged = true;
      if (currentMode !== 0) {
        // Удаляем все менюшки
        currentMod = "nomod";
        $("#sm-menu .item").remove()
        $("#sm-menu").hide();
      } else {
        // Пытаемся вернуть всё на место(тесто)
        $("#sm-menu").append(
          $(`
          <a class="item active" id="mode-nomod" data-mod="nomod" href="">Standart</a>
          <a class="item" id="mode-rx" data-mod="rx" href="#">Relax</a>
          <a class="item" id="mode-ap" data-mod="ap" href="#">Autopilot</a>
          `)
        )
        $("#sm-menu .item").click(smMenuHandler)
        $("#sm-menu").show()
      }
      loadLeaderboard(beatmapID, currentMode, currentMod);
    });
  $("#sm-menu .item")
    .click(smMenuHandler)

  $("table.sortable").tablesort();
})();
