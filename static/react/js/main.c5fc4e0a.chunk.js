(window.webpackJsonp=window.webpackJsonp||[]).push([[0],{39:function(e,t,a){e.exports=a(72)},44:function(e,t,a){},63:function(e,t,a){},68:function(e,t,a){},69:function(e,t,a){},72:function(e,t,a){"use strict";a.r(t);var n=a(0),s=a.n(n),i=a(35),c=a.n(i),r=(a(44),a(8)),o=a(9),l=a(11),m=a(10),d=a(12),u=a(13),h=a(1),p=a(18),v=a(6),f=a.n(v),g=(a(63),a(36)),b=a.n(g),w=function(e){function t(e){var a;return Object(r.a)(this,t),(a=Object(l.a)(this,Object(m.a)(t).call(this,e))).renderDummyData=function(){var e=a.state.beatmaps.map(function(e,t){var n=e.ChildrenBeatmaps.map(function(e,t){return s.a.createElement("div",{key:e.BeatmapID,className:"diff2"},s.a.createElement("div",{className:"faa fal fa-extra-mode-"+a.getGameMode(e.Mode),style:{color:a.getDiffColor(e.DifficultyRating)}}))});return s.a.createElement("div",{className:"eight wide column",key:e.SetID},s.a.createElement("div",{className:"map"},s.a.createElement("div",{className:"map-header"},s.a.createElement("a",{href:"https://kurikku.pw/b/"+e.ChildrenBeatmaps[e.ChildrenBeatmaps.length-1].BeatmapID},s.a.createElement("img",{src:"https://assets.ppy.sh/beatmaps/"+e.ChildrenBeatmaps[0].ParentSetID+"/covers/card.jpg",alt:""}))),s.a.createElement("div",{className:"status"},a.getRankStatus(e.RankedStatus)),s.a.createElement("div",{className:"name"},e.Title.substring(0,35)),s.a.createElement("div",{className:"artist"},e.Artist.substring(0,35)),s.a.createElement("div",{className:"creator"},"by ",e.Creator),s.a.createElement("a",{title:"Download beatmap",href:"https://storage.kurikku.pw/d/"+e.ChildrenBeatmaps[e.ChildrenBeatmaps.length-1].ParentSetID,className:"download"},s.a.createElement("i",{className:"download disk icon"})),s.a.createElement("div",{id:"alldiff"},n)))});return s.a.createElement("div",{className:"ui stackable two grid",onScrollCapture:a.scrollListener},e,s.a.createElement(b.a,{onBottom:a.scrollListener}))},a.state={offset:100,count:22,query:"",mode:0,loading:!0,status:1},a.renderDummyData.bind(Object(h.a)(Object(h.a)(a))),a.scrollListener=a.scrollListener.bind(Object(h.a)(Object(h.a)(a))),a.reCallApi=a.reCallApi.bind(Object(h.a)(Object(h.a)(a))),a.queryNew=a.queryNew.bind(Object(h.a)(Object(h.a)(a))),a}return Object(d.a)(t,e),Object(o.a)(t,[{key:"componentDidMount",value:function(){var e=this;f.a.get("/storage/api/search",{params:{offset:this.state.offset,amount:this.state.count,mode:this.state.mode,status:this.state.status,query:this.state.query}}).then(function(t){e.setState({loading:!1,beatmaps:t.data})})}},{key:"reCallApi",value:function(){var e=this;this.setState({loading:!0}),f.a.get("/storage/api/search",{params:{offset:this.state.offset,amount:this.state.count,mode:this.state.mode,status:this.state.status,query:this.state.query}}).then(function(t){e.setState({loading:!1,beatmaps:t.data})})}},{key:"getDiffColor",value:function(e){return e<=1.5?"#8AAE17":e>1.5&&e<=2.25?"#9AD4DF":e>2.25&&e<=3.75?"#DEB32A":e>3.75&&e<=5.25?"#EB69A4":e>5.25&&e<=6.75?"#7264B5":"#050505"}},{key:"getRankStatus",value:function(e){switch(e){case 0:return"PENDING";case 1:return"RANKED";case 3:return"APPROVED";case 4:return"Loved";default:return"UNKNOWN"}}},{key:"getGameMode",value:function(e){switch(e){case 0:return"osu";case 1:return"taiko";case 2:return"fruits";case 3:return"mania";default:return"osu"}}},{key:"scrollListener",value:function(){var e=this;this.setState({offset:this.state.offset+22}),f.a.get("/storage/api/search",{params:{offset:this.state.offset,amount:this.state.count,mode:this.state.mode,status:this.state.status,query:this.state.query}}).then(function(t){for(var a=e.state.beatmaps,n=0;n<t.data.length;n++)a.push(t.data[n]);e.setState({beatmaps:a})})}},{key:"switchMode",value:function(e){var t=this,a=Number(e.target.dataset.modeosu);this.setState({mode:a,offset:0},function(){t.reCallApi()})}},{key:"switchRank",value:function(e){var t=this,a=Number(e.target.dataset.rankstatus);this.setState({status:a,offset:0},function(){t.reCallApi()})}},{key:"queryNew",value:function(e){var t=this,a=e.target.value;this.setState({query:a,offset:0},function(){t.reCallApi()})}},{key:"render",value:function(){return s.a.createElement("div",null,s.a.createElement("div",{className:"ui segment"},s.a.createElement("div",{className:"ui one column stackable center aligned page grid"},s.a.createElement("div",{className:"column twelve wide"},s.a.createElement("center",null,s.a.createElement("h1",{className:"header"},"Beatmaps")),s.a.createElement("br",null),s.a.createElement("div",{class:"ui input",style:{width:"100%"}},s.a.createElement(p.DebounceInput,{minLength:0,debounceTimeout:350,onChange:this.queryNew})),s.a.createElement("div",{className:"ui segment wow-links"},s.a.createElement("a",{onClick:this.switchMode.bind(this),"data-modeosu":"-1",className:-1===this.state.mode?"clicked":"",href:"#"},"Any"),s.a.createElement("a",{onClick:this.switchMode.bind(this),"data-modeosu":"0",className:0===this.state.mode?"clicked":"",href:"#"},"osu!std"),s.a.createElement("a",{onClick:this.switchMode.bind(this),"data-modeosu":"1",className:1===this.state.mode?"clicked":"",href:"#"},"osu!taiko"),s.a.createElement("a",{onClick:this.switchMode.bind(this),"data-modeosu":"2",className:2===this.state.mode?"clicked":"",href:"#"},"osu!catch"),s.a.createElement("a",{onClick:this.switchMode.bind(this),"data-modeosu":"3",className:3===this.state.mode?"clicked":"",href:"#"},"osu!mania")),s.a.createElement("div",{className:"ui segment wow-links"},s.a.createElement("a",{onClick:this.switchRank.bind(this),className:-3===this.state.status?"clicked":"","data-rankstatus":"-3",href:"#"},"Any"),s.a.createElement("a",{onClick:this.switchRank.bind(this),className:1===this.state.status?"clicked":"","data-rankstatus":"1",href:"#"},"Ranked"),s.a.createElement("a",{onClick:this.switchRank.bind(this),className:3===this.state.status?"clicked":"","data-rankstatus":"3",href:"#"},"Qualified"),s.a.createElement("a",{onClick:this.switchRank.bind(this),className:4===this.state.status?"clicked":"","data-rankstatus":"4",href:"#"},"Loved"))))),s.a.createElement("div",{className:"bodySearching"},this.state.loading?s.a.createElement("div",{className:"ui active dimmer"},s.a.createElement("div",{className:"ui text loader"},"Loading")):this.renderDummyData()))}}]),t}(n.Component),E=a(16),y=a.n(E),k=a(20),N=(a(68),function(e){function t(e){var a;return Object(r.a)(this,t),(a=Object(l.a)(this,Object(m.a)(t).call(this,e))).state={loaded:!1,country11:[],country500:[],users:[],page:1,perPage:50,maxPage:1,mode:0,currentCountry:""},a.renderLeaderboard=a.renderLeaderboard.bind(Object(h.a)(Object(h.a)(a))),a.renderPagination=a.renderPagination.bind(Object(h.a)(Object(h.a)(a))),a}return Object(d.a)(t,e),Object(o.a)(t,[{key:"formatNumber",value:function(e){return e.toString().replace(/(\d)(?=(\d{3})+(?!\d))/g,"$1 ")}},{key:"showMessage",value:function(e,t){var a=void 0!==window.jQuery?window.jQuery:window.$,n=a('<div className="ui '+e+' message hidden"><i className="close icon"></i>'+t+"</div>");n.find(".close.icon").click(this.closeClosestMessage()),a("#messages-container").append(n),n.slideDown(300)}},{key:"componentDidMount",value:function(){var e=Object(k.a)(y.a.mark(function e(){var t,a;return y.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:return e.prev=0,e.next=3,f.a.get("/api/v2/leaderboardCountries");case 3:return t=e.sent,e.next=6,f.a.get("/api/v1/leaderboard",{params:{p:window.page,l:this.state.perPage,country:window.country,mode:window.favouriteMode}});case 6:a=e.sent,this.setState({country11:t.data.country11,country500:t.data.country500,maxPage:a.data.max_page,users:a.data.users,loaded:!0,mode:+window.favouriteMode,page:+window.page,currentCountry:window.country}),e.next=14;break;case 10:e.prev=10,e.t0=e.catch(0),console.log(e.t0),this.showMessage("error","Something happend when we tries to call server. Please, say to our developer that API is down.");case 14:window.jQuery(".ui.modal").modal({detachable:!1});case 15:case"end":return e.stop()}},e,this,[[0,10]])}));return function(){return e.apply(this,arguments)}}()},{key:"getNewLeaderboard",value:function(){var e=Object(k.a)(y.a.mark(function e(t,a,n){var s,i;return y.a.wrap(function(e){for(;;)switch(e.prev=e.next){case 0:if(!(!this.state.loaded||t>this.state.maxPage||t<=0)){e.next=2;break}return e.abrupt("return");case 2:return this.setState({loaded:!1}),e.prev=3,e.next=6,f.a.get("/api/v1/leaderboard",{params:{p:t,l:this.state.perPage,country:a.toLowerCase(),mode:n}});case 6:s=e.sent,this.setState({users:s.data.users,maxPage:s.data.max_page,page:+t,loaded:!0,mode:+n,currentCountry:a.toLowerCase()}),i=window.location,window.history.replaceState("",document.title,i.pathname+"?mode="+n+"&p="+t+(""!==a.toLowerCase()?"&country="+encodeURI(a.toLowerCase()):"")+i.hash),e.next=16;break;case 12:e.prev=12,e.t0=e.catch(3),console.log(e.t0),this.showMessage("error","Something happend when we tries to call server. Please, say to our developer that API is down.");case 16:case"end":return e.stop()}},e,this,[[3,12]])}));return function(t,a,n){return e.apply(this,arguments)}}()},{key:"renderLeaderboard",value:function(){var e=this;return this.state.users.map(function(t,a){return s.a.createElement("div",{className:"user-position",key:a},s.a.createElement("div",{className:"good"},s.a.createElement("div",{className:"columb col-1 center highlight ".concat(-1!==[1,2,3].indexOf(t.chosen_mode.global_leaderboard_rank)?"rat-pos":"")},s.a.createElement("span",null,"#",t.chosen_mode.global_leaderboard_rank," ",e.state.currentCountry.length>1?s.a.createElement("span",{className:"only-epic-color"},"(#".concat(t.chosen_mode.country_leaderboard_rank,")")):null)),s.a.createElement("div",{className:"columb full-size flexer highlight"},s.a.createElement("img",{src:"https://s.kurikku.pw/flags/".concat(t.country,".png"),alt:""}),s.a.createElement("a",{href:"/u/".concat(t.id)},t.username)),s.a.createElement("div",{className:"columb col-2 center highlight peppy-give-me-good-api"},s.a.createElement("span",null,e.formatNumber(t.chosen_mode.pp).replace("."," "),"pp"))),s.a.createElement("div",{className:"good2"},s.a.createElement("div",{className:"columb col-2 center highlight peppy-give-me-good-api"},s.a.createElement("span",null,e.formatNumber(t.chosen_mode.pp).replace("."," "),"pp")),s.a.createElement("div",{className:"columb col-3 center"},s.a.createElement("span",null,t.chosen_mode.accuracy.toFixed(2),"%")),s.a.createElement("div",{className:"columb col-4 center"},s.a.createElement("span",null,e.formatNumber(t.chosen_mode.playcount)))))})}},{key:"renderPagination",value:function(){var e=this;return s.a.createElement("div",{className:"pagination-by-trash"},s.a.createElement("div",{className:"pagination-prev",onClick:function(t){e.getNewLeaderboard(e.state.page-1,e.state.currentCountry,e.state.mode)}},s.a.createElement("span",null)),s.a.createElement("div",{className:"pagination-items"},[-2,-1,0,1,2].map(function(t,a){return 0===t?s.a.createElement("div",{className:"pagination-item active"},e.state.page):e.state.page+t<e.state.maxPage&&e.state.page+t>0?s.a.createElement("div",{className:"pagination-item",key:a,onClick:function(a){e.getNewLeaderboard(e.state.page+t,e.state.currentCountry,e.state.mode)}},e.state.page+t):null}),s.a.createElement("div",{className:"pagination-item dots"},s.a.createElement("span",null)),s.a.createElement("div",{className:"pagination-item",onClick:function(t){e.getNewLeaderboard(e.state.maxPage,e.state.currentCountry,e.state.mode)}},this.state.maxPage)),s.a.createElement("div",{className:"pagination-next",onClick:function(t){e.getNewLeaderboard(e.state.page+1,e.state.currentCountry,e.state.mode)}},s.a.createElement("span",null)))}},{key:"render",value:function(){var e=this;return s.a.createElement("div",{className:"ui container"},s.a.createElement("div",{className:"ui four item menu",id:"mode-menu"},["osu!standard","Taiko","Catch the Beat","osu!mania"].map(function(t,a){return s.a.createElement("a",{className:"".concat(a===+e.state.mode?"active ":"","item"),onClick:function(t){e.getNewLeaderboard(1,"",a)}},t)})),s.a.createElement("div",{className:"ui twelve item stackable menu"},this.state.country11.map(function(t,a){return s.a.createElement("a",{className:"item".concat(t===e.state.currentCountry?" active":""," lb-country"),key:a,onClick:function(a){e.getNewLeaderboard(1,t.toUpperCase(),e.state.mode)}},s.a.createElement("img",{src:"https://s.kurikku.pw/flags/".concat(t.toUpperCase(),".png"),className:t+" flag nopad"}))}),s.a.createElement("a",{className:"item",onClick:function(e){return window.jQuery(".ui.modal").modal("show")}},"...")),s.a.createElement("div",{className:"leader-table"},this.renderPagination(),s.a.createElement("div",{className:"table-leader"},s.a.createElement("div",{className:"table-head"},s.a.createElement("div",{className:"columb full-size"}),s.a.createElement("div",{className:"columb col-2 highlight"},"Peppy Points"),s.a.createElement("div",{className:"columb col-3"},"Accuracy"),s.a.createElement("div",{className:"columb col-4"},"Play Count")),s.a.createElement("div",{className:"table-body"},this.state.loaded?this.renderLeaderboard():s.a.createElement("div",{className:"ui active dimmer"},s.a.createElement("div",{className:"ui text loader only-white-color"},"Loading")))),this.renderPagination()),s.a.createElement("div",{class:"ui modal"},s.a.createElement("div",{class:"content"},s.a.createElement("div",{class:"ui sixteen stackable column grid"},this.state.country500.map(function(t,a){return s.a.createElement("div",{className:"ui clickable column lb-country simple-flex",key:a,onClick:function(a){window.jQuery(".ui.modal").modal("hide"),e.getNewLeaderboard(1,t.toUpperCase(),e.state.mode)}},s.a.createElement("img",{className:"ui mini image lb-country",src:"https://s.kurikku.pw/flags/".concat(t.toUpperCase(),".png"),alt:"".concat(t.toUpperCase()," flag")}))})))))}}]),t}(s.a.Component)),C=(a(69),function(e){function t(e){var a;return Object(r.a)(this,t),(a=Object(l.a)(this,Object(m.a)(t).call(this,e))).state={loading:!0,items:[],balance:0,isConfirmPopoutShowed:!1,selectedItem:0,selectedItemInd:-1},a.renderItems=a.renderItems.bind(Object(h.a)(Object(h.a)(a))),a.onCanceled=a.onCanceled.bind(Object(h.a)(Object(h.a)(a))),a.onConfirmed=a.onConfirmed.bind(Object(h.a)(Object(h.a)(a))),a}return Object(d.a)(t,e),Object(o.a)(t,[{key:"componentDidMount",value:function(){var e=this;f.a.get("/api/v1/shop/get_items").then(function(t){console.log(t.data),e.setState({loading:!1,items:t.data.items,balance:t.data.balance})})}},{key:"renderItems",value:function(){var e=this;return this.state.items.map(function(t,a){return s.a.createElement("div",{className:"column",key:a},s.a.createElement("div",{className:"ui left aligned fluid card"},s.a.createElement("div",{className:"image"},s.a.createElement("img",{src:t.Image,alt:"Item preview"})),s.a.createElement("div",{className:"content"},s.a.createElement("p",{className:"header"},t.Name),s.a.createElement("p",{className:"meta"},t.Description)),s.a.createElement("div",{className:"extra content"},s.a.createElement("div",{className:t.CanBuy?"ui animated fade positive button":"ui animated fade positive button disabled",onClick:function(){e.onClickItem(t.ID,a)},tabIndex:"0"},s.a.createElement("div",{className:"hidden content"},t.Cost," RUB"),s.a.createElement("div",{className:"visible content"},s.a.createElement("i",{className:"shop icon"}))))))})}},{key:"onClickItem",value:function(e,t){this.setState({isConfirmPopoutShowed:!0,selectedItem:e,selectedItemInd:t}),window.jQuery("#confirmModal").modal("show")}},{key:"closeClosestMessage",value:function(){window.jQuery(this).closest(".message").fadeOut(300,function(){window.jQuery(this).remove()})}},{key:"onConfirmed",value:function(){var e=this;console.log(this.state.items[this.state.selectedItemInd]),this.setState({isConfirmPopoutShowed:!1});var t=window.jQuery("<div id='loadingDimmer' class='ui active dimmer'><div class='ui text loader'>Loading</div></div>");window.jQuery("body").append(t),f.a.get("/api/v1/shop/buy_item",{params:{itemID:this.state.items[this.state.selectedItemInd].ID}}).then(function(t){e.setState({selectedItem:0,selectedItemInd:-1}),window.jQuery("#loadingDimmer").remove(),e.showMessage("success","Your buy has been complete. Now page will reload"),window.location.reload()}).catch(function(t){e.showMessage("error","Your buy has not been complete. Try later :("),window.jQuery("#loadingDimmer").remove()})}},{key:"onCanceled",value:function(){this.setState({isConfirmPopoutShowed:!1,selectedItem:0,selectedItemInd:-1})}},{key:"showMessage",value:function(e,t){var a=void 0!==window.jQuery?window.jQuery:window.$,n=a('<div class="ui '+e+' message hidden"><i class="close icon"></i>'+t+"</div>");n.find(".close.icon").click(this.closeClosestMessage()),a("#messages-container").append(n),n.slideDown(300)}},{key:"renderMain",value:function(){return s.a.createElement("div",null,s.a.createElement("div",{className:"ui segment"},s.a.createElement("div",{className:"Shop--up"},s.a.createElement("h3",null,"Kurikku!Shop"),s.a.createElement("h3",{className:"Shop--right"},this.state.balance," RUB")),s.a.createElement("div",{className:"ui divider"}),s.a.createElement("div",{className:"ui three column left aligned stackable grid"},this.renderItems(),this.renderItems(),this.renderItems())),s.a.createElement("div",{id:"confirmModal",className:"ui basic modal"},s.a.createElement("div",{className:"ui icon header",style:{color:"white"}},s.a.createElement("i",{className:"archive icon"}),"Do you Confirm your Payment"),s.a.createElement("div",{className:"actions"},s.a.createElement("div",{className:"ui red cancel button",onClick:this.onCanceled},s.a.createElement("i",{className:"remove icon"}),"Close"),s.a.createElement("div",{className:"ui green ok button",onClick:this.onConfirmed},s.a.createElement("i",{className:"checkmark icon"}),"Confirm"))))}},{key:"render",value:function(){return this.state.loading?s.a.createElement("div",{className:"ui active dimmer"},s.a.createElement("div",{className:"ui text loader"},"Loading")):this.renderMain()}}]),t}(s.a.Component)),j=function(e){function t(e){return Object(r.a)(this,t),Object(l.a)(this,Object(m.a)(t).call(this,e))}return Object(d.a)(t,e),Object(o.a)(t,[{key:"render",value:function(){return s.a.createElement("div",{className:"ui container"},s.a.createElement(u.c,null,s.a.createElement(u.a,{path:"/shop",component:C}),s.a.createElement(u.a,{path:"/beatmaps",component:w}),s.a.createElement(u.a,{path:"/leaderboard",component:N})))}}]),t}(s.a.Component),O=a(19);c.a.render(s.a.createElement(O.a,null,s.a.createElement(j,null)),document.getElementById("app"))}},[[39,1,2]]]);
//# sourceMappingURL=main.c5fc4e0a.chunk.js.map