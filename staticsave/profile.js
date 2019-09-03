// code that is executed on every user profile
$(document).ready(function() {
  var wl = window.location;
  var newPathName = wl.pathname;
  // userID is defined in profile.html
  if (newPathName.split("/")[2] != userID) {
    newPathName = "/u/" + userID;
  }
  // if there's no mode parameter in the querystring, add it
  if (wl.search.indexOf("mode=") === -1)
    window.history.replaceState('', document.title, newPathName + "?mode=" + favouriteMode + wl.hash);
  else if (wl.pathname != newPathName)
    window.history.replaceState('', document.title, newPathName + wl.search + wl.hash);
  setDefaultScoreTable();
  // when an item in the mode menu is clicked, it means we should change the mode.
  $("#mode-menu>.item").click(function(e) {
    e.preventDefault();
    if ($(this).hasClass("active"))
      return;
    var m = $(this).data("mode");
    $("[data-mode]:not(.item):not([hidden])").attr("hidden", "");
    $("[data-mode=" + m + "]:not(.item)").removeAttr("hidden");
    $("#mode-menu>.active.item").removeClass("active");
    var needsLoad = $("#scores-zone>[data-mode=" + m + "][data-loaded=0]");
    if (needsLoad.length > 0)
      initialiseScores(needsLoad, m);
    $(this).addClass("active");
    window.history.replaceState('', document.title, wl.pathname + "?mode=" + m + wl.hash);
  });
  initialiseFriends();
  // load scores page for the current favourite mode
  var i = function(){initialiseScores($("#scores-zone>div[data-mode=" + favouriteMode + "]"), favouriteMode)};
  if (i18nLoaded)
    i();
  else
    i18next.on("loaded", function() {
      i();
    });
});

function initialiseFriends() {
  var b = $("#add-friend-button");
  if (b.length == 0) return;
  api('friends/with', {id: userID}, setFriendOnResponse);
  b.click(friendClick);
}
function setFriendOnResponse(r) {
  var x = 0;
  if (r.friend) x++;
  if (r.mutual) x++;
  setFriend(x);
}
function setFriend(i) {
  var b = $("#add-friend-button");
  b.removeClass("loading green blue red");
  switch (i) {
  case 0:
    b
      .addClass("blue")
      .attr("title", T("Add friend"))
      .html("<i class='plus icon'></i>");
    break;
  case 1:
    b
      .addClass("green")
      .attr("title", T("Remove friend"))
      .html("<i class='minus icon'></i>");
    break;
  case 2:
    b
      .addClass("red")
      .attr("title", T("Unmutual friend"))
      .html("<i class='heart icon'></i>");
    break;
  }
  b.attr("data-friends", i > 0 ? 1 : 0)
}

function humanizeNumber(number) {
  return number.toString().replace(/(\d)(?=(\d{3})+(?!\d))/g, '$1 ')
}

function friendClick() {
  var t = $(this);
  if (t.hasClass("loading")) return;
  t.addClass("loading");
  api("friends/" + (t.attr("data-friends") == 1 ? "del" : "add"), {user: +userID}, setFriendOnResponse, true);
}

var defaultScoreTable;
function setDefaultScoreTable() {
  defaultScoreTable = $("<table class='ui table score-table' />")
    .append(
      $("<thead />").append(
        $("<tr />").append(
          $("<th class=\"text-kr-center\">" + T("General info") + "</th>"),
          $("<th class=\"text-kr-center\">"+ T("Score") + "</th>")
        )
      )
    )
    .append(
      $("<tbody />")
    )
    .append(
      $("<tfoot />").append(
        $("<tr />").append(
          $("<th colspan=2 />").append(
            $("<div class='ui right floated pagination menu' />").append(
              $("<a class='disabled item load-more-button'>" + T("Load more") + "</a>").click(loadMoreClick)
            )
          )
        )
      )
    )
  ;
}
i18next.on('loaded', function(loaded) {
  setDefaultScoreTable();
});
function initialiseScores(el, mode) {
  el.attr("data-loaded", "1");
  var best = defaultScoreTable.clone(true).addClass("orange");
  var recent = defaultScoreTable.clone(true).addClass("blue");
  best.attr("data-type", "best");
  recent.attr("data-type", "recent");
  recent.addClass("no bottom margin");
  el.append($("<div class='ui segments no bottom margin' />").append(
    $("<div class='ui segment' />").append("<h2 class='ui header'>" + T("Best scores") + "</h2>", best),
    $("<div class='ui segment' />").append("<h2 class='ui header'>" + T("Recent scores") + "</h2>", recent)
  ));
  loadScoresPage("best", mode);
  loadScoresPage("recent", mode);
};
function loadMoreClick() {
  var t = $(this);
  if (t.hasClass("disabled"))
    return;
  t.addClass("disabled");
  var type = t.parents("table[data-type]").data("type");
  var mode = t.parents("div[data-mode]").data("mode");
  loadScoresPage(type, mode);
}
// currentPage for each mode
var currentPage = {
  0: {best: 0, recent: 0},
  1: {best: 0, recent: 0},
  2: {best: 0, recent: 0},
  3: {best: 0, recent: 0},
};
var scoreStore = {};
function loadScoresPage(type, mode) {
  var table = $("#scores-zone div[data-mode=" + mode + "] table[data-type=" + type + "] tbody");
  var page = ++currentPage[mode][type];
  console.log("loadScoresPage with", {
    page: page,
    type: type,
    mode: mode,
  });
  api("users/scores/" + type, {
    mode: mode,
    p: page,
    l: 20,
    id: userID,
  }, function(r) {
    if (r.scores == null) {
      disableLoadMoreButton(type, mode);
      return;
    }
    r.scores.forEach(function(v, idx){
      scoreStore[v.id] = v;
      var scoreRank = getRank(mode, v.mods, v.accuracy, v.count_300, v.count_100, v.count_50, v.count_miss);
      table.append(
        $("<tr class='new score-row' data-scoreid='" + v.id + "' />").append(
          $(`
            <div class="scores-table">
              <div class="scores-table-left">
                <img src='/static/ranking-icons/${scoreRank}.svg' class='score rank' alt='${scoreRank}'> 
              </div>
              <div class="scores-table-left-info">
                ${escapeHTML(v.beatmap.song_name)} <b>${getScoreMods(v.mods)}</b><br/>
                
                <div class="subtitle">
                  ${v.accuracy.toFixed(2)}% / ${humanizeNumber(v.score)} / ${v.max_combo}x (${v.beatmap.max_combo}x) { ${v.count_300} / ${v.count_100} / ${v.count_50} / ${v.count_miss} }
                </div>
                <div class="subtitle">
                  <time class='new timeago' datetime='${v.time}'>${v.time}</time>
                </div>
              </div>
            </div>
          `),
          $("<td class=\"text-kr-center\"><b>" + ppOrScore(v.pp, v.score) + "</b> " + weightedPP(type, page, idx, v.pp) +  (v.completed == 3 ? "<br>" + downloadStar(v.id) : "") +  "</td>")
      ));
    });
    $(".new.timeago").timeago().removeClass("new");
    $(".new.score-row").click(viewScoreInfo).removeClass("new");
    $(".new.downloadstar").click(function(e) {
      e.stopPropagation();
    }).removeClass("new");
    var enable = true;
    if (r.scores.length != 20)
      enable = false;
    disableLoadMoreButton(type, mode, enable);
  });
}
function downloadStar(id) {
  return "<a href='/web/replays/" + id + "' class='new downloadstar'><i class='star icon'></i>" + T("Download") + "</a>";
}
function weightedPP(type, page, idx, pp) {
  if (type != "best" || pp == 0)
    return "";
  var perc = Math.pow(0.95, ((page - 1) * 20) + idx);
  var wpp = pp * perc;
  return "<i title='Weighted PP, " + Math.round(perc*100) + "%'>(" + wpp.toFixed(2) + "pp)</i>";
}
function disableLoadMoreButton(type, mode, enable) {
  var button = $("#scores-zone div[data-mode=" + mode + "] table[data-type=" + type + "] .load-more-button");
  if (enable) button.removeClass("disabled");
  else button.addClass("disabled");
}
function viewScoreInfo() {
  var scoreid = $(this).data("scoreid");
  if (!scoreid && scoreid !== 0) return;
  var s = scoreStore[scoreid];
  if (s === undefined) return;

  // data to be displayed in the table.
  // Data by KotRik, because i want!
  var ultimateStupidDatav2 = {
    // key, value, shouldinsert?
    'bg': {
      val: `https://assets.ppy.sh/beatmaps/${s.beatmap.beatmap_id}/covers/cover.jpg`
    },
    'Score': {
      name: 'score.png',
      val:  humanizeNumber(s.score)
    },
    'PP': {
      name: 'pp.png',
      val: addCommas(s.pp)
    },
    'Starrate': {
      name: 'starrate.png',
      val: Math.round(s.beatmap.difficulty2[modesShort[s.play_mode]])
    },
    'Accuracy': {
      name: 'accuracy.png',
      val: s.accuracy + "%"
    },
    'Mods': {
      name: 'mods.png',
      val: getScoreMods(s.mods, true)
    },
    'Combo': {
      name: 'combo.png',
      val: s.max_combo
    }
  }

  // hits data
  var hd = {};
  var trans = modeTranslations[s.play_mode];
  [
    s.count_300,
    s.count_100,
    s.count_50,
    s.count_geki,
    s.count_katu,
    s.count_miss,
  ].forEach(function(val, i) {
    hd[trans[i][0]] = {
        name: trans[i][1],
        val: val
    };
  });

  data = $.extend(ultimateStupidDatav2, hd);

  var els = [];
  $.each(data, function(key, value) {
    if (key === "bg") continue;

    els.push(
      $(`
      <div class="four wide column">
        <div>
          <div class="ui grid score-info-grid">
            <div class="four wide column">
              <img src="/static/images/scoreic/${value.name}" class="ui centered icon-score"/>
            </div>
            <div class="twelve wide column">
              <p class="status-head-score">${value.val}</p>
              <p class="status-footer">${T(key)}</p>
            </div>
          </div>
        </div>
      </div>
      `)
    );
  });

  $("#scores-header").css("background-image", data['bg']['val']); // Update header image
  $("#scores-body div").remove(); // Remove old stats
  $("#scores-body").append(els); // Add new stats imho ;d

  $(".ui.modal").modal("show");
}

var modeTranslations = [
  [
    ["300s", '300.png'], 
    ["100s", '100.png'],
    ["50s", '50.png'],
    ["Gekis", 'gekis.png'],
    ["Katus", 'katus.png'],
    ["Misses", 'misses.png']
  ],
  [
    ["GREATs", '300.png'],
    ["GOODs", '100.png'],
    ["50s", '50.png'],
    ["GREATs (Gekis)", 'gekis.png'],
    ["GOODs (Katus)", 'katus.png'],
    ["Misses", 'misses.png']
  ],
  [
    ["Fruits (300s)", '300.png'],
    ["Ticks (100s)", '100.png'],
    ["Droplets", '50.png'],
    ["Gekis", 'gekis.png'],
    ["Droplet misses", 'katus.png'],
    ["Misses", 'misses.png']
  ],
  [
    ["300s", '300.png'],
    ["200s", '100.png'],
    ["50s", '50.png'],
    ["Max 300s", 'gekis.png'],
    ["100s", 'katus.png'],
    ["Misses", 'misses.png']
  ]
];

// helper functions copied from user.js in old-frontend
function getScoreMods(m, noplus) {
	var r = [];
  // has nc => remove dt
  if ((m & 512) == 512)
    m = m & ~64;
  // has pf => remove sd
  if ((m & 16384) == 16384)
    m = m & ~32;
  modsString.forEach(function(v, idx) {
    var val = 1 << idx;
    if ((m & val) > 0)
      r.push(v);
  });
	if (r.length > 0) {
		return (noplus ? "" : "+ ") + r.join(", ");
	} else {
		return (noplus ? T('None') : '');
	}
}

var modsString = [
  "NF",
	"EZ",
	"NV",
	"HD",
	"HR",
	"SD",
	"DT",
	"RX",
	"HT",
	"NC",
	"FL",
	"AU", // Auto.
	"SO",
	"AP", // Autopilot.
	"PF",
	"K4",
	"K5",
	"K6",
	"K7",
	"K8",
	"K9",
	"RN", // Random
	"LM", // LastMod. Cinema?
	"K9",
	"K0",
	"K1",
	"K3",
	"K2",
];

function getRank(gameMode, mods, acc, c300, c100, c50, cmiss) {
	var total = c300+c100+c50+cmiss;

  // Hidden | Flashlight | FadeIn
	var hdfl = (mods & (1049608)) > 0;

	var ss = hdfl ? "SSHD" : "SS";
	var s = hdfl ? "SHD" : "S";

	switch(gameMode) {
		case 0:
		case 1:
			var ratio300 = c300 / total;
			var ratio50 = c50 / total;

			if (ratio300 == 1)
				return ss;

			if (ratio300 > 0.9 && ratio50 <= 0.01 && cmiss == 0)
				return s;

			if ((ratio300 > 0.8 && cmiss == 0) || (ratio300 > 0.9))
				return "A";

			if ((ratio300 > 0.7 && cmiss == 0) || (ratio300 > 0.8))
				return "B";

			if (ratio300 > 0.6)
				return "C";

			return "D";

		case 2:
			if (acc == 100)
				return ss;

			if (acc > 98)
				return s;

			if (acc > 94)
				return "A";

			if (acc > 90)
				return "B";

			if (acc > 85)
				return "C";

			return "D";

		case 3:
			if (acc == 100)
				return ss;

			if (acc > 95)
				return s;

			if (acc > 90)
				return "A";

			if (acc > 80)
				return "B";

			if (acc > 70)
				return "C";

			return "D";
	}
}

function ppOrScore(pp, score) {
  if (pp != 0)
    return addCommas(pp.toFixed(2)) + "pp";
  return addCommas(score);
}

function beatmapLink(type, id) {
  if (type == "s")
    return "<a href='/s/" + id + "'>" + id + '</a>';
  return "<a href='/b/" + id + "'>" + id + '</a>';  
}
