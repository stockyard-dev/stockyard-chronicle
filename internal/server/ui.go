package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Chronicle</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:.8rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:1rem;flex-wrap:wrap}
.hdr h1{font-size:.9rem;letter-spacing:2px}
.hdr h1 span{color:var(--rust)}
.hdr-r{display:flex;gap:.5rem;align-items:center;flex-wrap:wrap}
.stats-line{font-size:.6rem;color:var(--cm)}
.stats-line .num{color:var(--cream);font-weight:700}
.main{padding:1rem 1.5rem 1.5rem;max-width:1100px;margin:0 auto}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;flex-wrap:wrap;align-items:center}
.search{flex:1;min-width:200px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem;min-width:130px}
.type-bar{display:flex;gap:.3rem;margin-bottom:.8rem;flex-wrap:wrap}
.tbtn{font-size:.55rem;padding:.25rem .5rem;border:1px solid var(--bg3);background:var(--bg2);color:var(--cm);cursor:pointer;font-family:var(--mono)}
.tbtn:hover{border-color:var(--leather)}
.tbtn.active{border-color:var(--rust);color:var(--rust)}
.tbtn .count{color:var(--cm);margin-left:.3rem;font-size:.5rem}
.count-label{font-size:.6rem;color:var(--cm);margin-bottom:.5rem}
.events{display:flex;flex-direction:column}
.ev{display:grid;grid-template-columns:90px 90px 1fr;gap:.6rem;padding:.6rem .8rem;border:1px solid var(--bg3);border-bottom:none;background:var(--bg2);font-size:.7rem;align-items:flex-start;transition:border-color .15s}
.ev:last-child{border-bottom:1px solid var(--bg3)}
.ev:hover{border-left-color:var(--leather)}
.ev-time{font-size:.58rem;color:var(--cm);font-variant-numeric:tabular-nums;line-height:1.4}
.ev-time-rel{color:var(--cd);display:block}
.ev-time-abs{display:block;font-size:.5rem}
.ev-sev{font-size:.5rem;padding:.15rem .35rem;text-transform:uppercase;letter-spacing:1px;text-align:center;border:1px solid var(--bg3);align-self:flex-start;font-weight:700}
.ev-sev.info{border-color:var(--blue);color:var(--blue)}
.ev-sev.warning,.ev-sev.warn{border-color:var(--orange);color:var(--orange)}
.ev-sev.error{border-color:var(--red);color:var(--red)}
.ev-sev.debug{border-color:var(--cm);color:var(--cm)}
.ev-sev.success{border-color:var(--green);color:var(--green)}
.ev-body{min-width:0;line-height:1.4}
.ev-type{color:var(--gold);font-weight:700;font-size:.65rem}
.ev-source{color:var(--cm);font-size:.6rem;margin-left:.5rem}
.ev-subject{color:var(--cream);margin-top:.15rem;word-wrap:break-word}
.ev-tags{margin-top:.2rem;display:flex;gap:.3rem;flex-wrap:wrap}
.ev-tag{font-size:.5rem;background:var(--bg3);color:var(--cd);padding:.05rem .3rem}
.ev-data{font-size:.6rem;color:var(--cd);background:var(--bg);padding:.3rem .5rem;margin-top:.3rem;border-left:2px solid var(--bg3);max-height:80px;overflow-y:auto;white-space:pre-wrap;word-break:break-word;font-family:var(--mono)}
.btn{font-size:.6rem;padding:.3rem .6rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);font-family:var(--mono);transition:all .15s}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-p:hover{opacity:.85;color:#fff}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:480px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}
.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.fr-hint{font-size:.5rem;color:var(--cm);margin-top:.15rem}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem;background:var(--bg2);border:1px solid var(--bg3)}
.live-toggle{display:inline-flex;align-items:center;gap:.3rem;font-size:.55rem;color:var(--cm);cursor:pointer;user-select:none}
.live-toggle input{width:auto;cursor:pointer}
.live-dot{width:6px;height:6px;border-radius:50%;background:var(--cm);display:inline-block}
.live-dot.on{background:var(--green);box-shadow:0 0 4px var(--green)}
@media(max-width:700px){.ev{grid-template-columns:1fr;gap:.3rem}.ev-time,.ev-sev{align-self:flex-start}.toolbar{flex-direction:column;align-items:stretch}.search{min-width:100%}}
</style>
</head>
<body>

<div class="hdr">
<h1 id="dash-title"><span>&#9670;</span> CHRONICLE</h1>
<div class="hdr-r">
<span class="stats-line" id="stats"></span>
<label class="live-toggle"><input type="checkbox" id="live-cb" checked> <span class="live-dot on" id="live-dot"></span>Live</label>
<button class="btn btn-p" onclick="openEmit()">+ Log Event</button>
</div>
</div>

<div class="main">
<div class="toolbar">
<input class="search" id="search" placeholder="Search subject, data, tags..." oninput="debouncedRender()">
<select class="filter-sel" id="severity-filter" onchange="render()">
<option value="">All Severities</option>
<option value="info">Info</option>
<option value="success">Success</option>
<option value="warning">Warning</option>
<option value="error">Error</option>
<option value="debug">Debug</option>
</select>
<select class="filter-sel" id="source-filter" onchange="render()">
<option value="">All Sources</option>
</select>
</div>

<div class="type-bar" id="typeBar"></div>
<div class="count-label" id="count"></div>
<div class="events" id="events"></div>
</div>

<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()">
<div class="modal" id="mdl"></div>
</div>

<script>
var A='/api';
var events=[];
var stats={};
var curType='';
var liveOn=true;
var liveTimer=null;
var searchTimer=null;
var presetTypes=[];
var presetSources=[];

// ─── Helpers ──────────────────────────────────────────────────────

function fmtRelative(s){
if(!s)return'';
try{
var d=new Date(s);
if(isNaN(d.getTime()))return s;
var now=new Date();
var diff=Math.floor((now-d)/1000);
if(diff<10)return'just now';
if(diff<60)return diff+'s ago';
if(diff<3600)return Math.floor(diff/60)+'m ago';
if(diff<86400)return Math.floor(diff/3600)+'h ago';
return Math.floor(diff/86400)+'d ago';
}catch(e){return s}
}

function fmtAbsolute(s){
if(!s)return'';
try{
var d=new Date(s);
if(isNaN(d.getTime()))return s;
return d.toLocaleTimeString('en-US',{hour:'2-digit',minute:'2-digit',second:'2-digit'});
}catch(e){return s}
}

function debouncedRender(){
clearTimeout(searchTimer);
searchTimer=setTimeout(render,200);
}

// ─── Loading ──────────────────────────────────────────────────────

async function load(){
try{
var qs=[];
if(curType)qs.push('type='+encodeURIComponent(curType));
qs.push('limit=200');
var url=A+'/events'+(qs.length?'?'+qs.join('&'):'');
var resp=await Promise.all([
fetch(url).then(function(r){return r.json()}),
fetch(A+'/stats').then(function(r){return r.json()})
]);
events=resp[0].events||[];
stats=resp[1]||{};
}catch(e){
console.error('load failed',e);
events=[];
stats={};
}
renderStats();
renderTypeBar();
renderSourceFilter();
render();
}

function renderStats(){
var total=stats.total||0;
var today=stats.today||0;
var types=stats.types||0;
var sources=stats.sources||0;
document.getElementById('stats').innerHTML=
'<span class="num">'+total+'</span> events &middot; '+
'<span class="num">'+today+'</span> today &middot; '+
'<span class="num">'+types+'</span> types &middot; '+
'<span class="num">'+sources+'</span> sources';
}

function renderTypeBar(){
var byType=stats.by_type||{};
var typeKeys=Object.keys(byType).sort(function(a,b){return byType[b]-byType[a]});
// Merge in preset types from personalization, even if no events yet
presetTypes.forEach(function(t){if(typeKeys.indexOf(t)===-1)typeKeys.push(t)});

var h='<button class="tbtn'+(curType===''?' active':'')+'" onclick="setType(\'\')">All</button>';
typeKeys.forEach(function(t){
var c=byType[t]||0;
h+='<button class="tbtn'+(curType===t?' active':'')+'" onclick="setType(\''+esc(t)+'\')">';
h+=esc(t);
if(c>0)h+='<span class="count">'+c+'</span>';
h+='</button>';
});
document.getElementById('typeBar').innerHTML=h;
}

function renderSourceFilter(){
var sel=document.getElementById('source-filter');
if(!sel)return;
var current=sel.value;
// Build the source list from current events + preset sources
var seen={};
var sources=[];
events.forEach(function(e){if(e.source&&!seen[e.source]){seen[e.source]=true;sources.push(e.source)}});
presetSources.forEach(function(s){if(!seen[s]){seen[s]=true;sources.push(s)}});
sources.sort();
sel.innerHTML='<option value="">All Sources</option>'+sources.map(function(s){return'<option value="'+esc(s)+'"'+(s===current?' selected':'')+'>'+esc(s)+'</option>'}).join('');
}

function setType(t){
curType=t;
load();
}

// ─── Rendering ────────────────────────────────────────────────────

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var sf=document.getElementById('severity-filter').value;
var srcF=document.getElementById('source-filter').value;

var f=events;
if(sf)f=f.filter(function(e){return e.severity===sf});
if(srcF)f=f.filter(function(e){return e.source===srcF});
if(q)f=f.filter(function(e){
return(e.subject||'').toLowerCase().includes(q)||
       (e.data||'').toLowerCase().includes(q)||
       (e.tags||'').toLowerCase().includes(q);
});

document.getElementById('count').textContent=f.length+' event'+(f.length!==1?'s':'')+(events.length!==f.length?' (of '+events.length+')':'');

if(!f.length){
var msg=window._emptyMsg||'No events to show. Click "Log Event" to add one, or POST to /api/events.';
document.getElementById('events').innerHTML='<div class="empty">'+esc(msg)+'</div>';
return;
}

var h='';
f.forEach(function(e){h+=eventHTML(e)});
document.getElementById('events').innerHTML=h;
}

function eventHTML(e){
var sev=(e.severity||'info').toLowerCase();
var h='<div class="ev">';

// Time column
h+='<div class="ev-time">';
h+='<span class="ev-time-rel">'+esc(fmtRelative(e.created_at))+'</span>';
h+='<span class="ev-time-abs">'+esc(fmtAbsolute(e.created_at))+'</span>';
h+='</div>';

// Severity column
h+='<div class="ev-sev '+esc(sev)+'">'+esc(sev)+'</div>';

// Body column
h+='<div class="ev-body">';
h+='<span class="ev-type">'+esc(e.type||'event')+'</span>';
if(e.source)h+='<span class="ev-source">['+esc(e.source)+']</span>';
if(e.subject)h+='<div class="ev-subject">'+esc(e.subject)+'</div>';
if(e.tags){
var tagList=String(e.tags).split(',').map(function(t){return t.trim()}).filter(function(t){return t});
if(tagList.length){
h+='<div class="ev-tags">';
tagList.forEach(function(t){h+='<span class="ev-tag">#'+esc(t)+'</span>'});
h+='</div>';
}
}
if(e.data&&e.data!=='{}'&&e.data.trim()!==''){
h+='<div class="ev-data">'+esc(e.data)+'</div>';
}
h+='</div>';

h+='</div>';
return h;
}

// ─── Emit form ────────────────────────────────────────────────────

function openEmit(){
var typeOpts='';
if(presetTypes.length){
typeOpts+='<option value="">— Select or type below —</option>';
presetTypes.forEach(function(t){typeOpts+='<option value="'+esc(t)+'">'+esc(t)+'</option>'});
}
var sourceOpts='';
if(presetSources.length){
sourceOpts+='<option value="">— Select or type below —</option>';
presetSources.forEach(function(s){sourceOpts+='<option value="'+esc(s)+'">'+esc(s)+'</option>'});
}

var h='<h2>NEW EVENT</h2>';

if(presetTypes.length){
h+='<div class="fr"><label>Event Type *</label><select id="f-type-select" onchange="document.getElementById(\'f-type\').value=this.value">'+typeOpts+'</select></div>';
h+='<div class="fr"><label>Or custom type</label><input type="text" id="f-type" placeholder="custom_type"></div>';
}else{
h+='<div class="fr"><label>Event Type *</label><input type="text" id="f-type" placeholder="custom_type"><div class="fr-hint">e.g. login, order_placed, fermentation_started</div></div>';
}

h+='<div class="row2">';
h+='<div class="fr"><label>Severity</label><select id="f-severity"><option value="info">Info</option><option value="success">Success</option><option value="warning">Warning</option><option value="error">Error</option><option value="debug">Debug</option></select></div>';
if(presetSources.length){
h+='<div class="fr"><label>Source</label><select id="f-source"><option value="">None</option>'+presetSources.map(function(s){return'<option value="'+esc(s)+'">'+esc(s)+'</option>'}).join('')+'</select></div>';
}else{
h+='<div class="fr"><label>Source</label><input type="text" id="f-source" placeholder="api / cron / web"></div>';
}
h+='</div>';

h+='<div class="fr"><label>Subject</label><input type="text" id="f-subject" placeholder="Short description"></div>';
h+='<div class="fr"><label>Data</label><textarea id="f-data" rows="3" placeholder="Optional details, JSON, or context"></textarea></div>';
h+='<div class="fr"><label>Tags</label><input type="text" id="f-tags" placeholder="comma separated"></div>';

h+='<div class="acts">';
h+='<button class="btn" onclick="closeModal()">Cancel</button>';
h+='<button class="btn btn-p" onclick="submitEvent()">Log Event</button>';
h+='</div>';

document.getElementById('mdl').innerHTML=h;
document.getElementById('mbg').classList.add('open');
var t=document.getElementById('f-type');
if(t)t.focus();
}

async function submitEvent(){
var typeEl=document.getElementById('f-type');
if(!typeEl||!typeEl.value.trim()){
alert('Event type is required');
return;
}
var body={
type:typeEl.value.trim(),
severity:document.getElementById('f-severity').value,
source:document.getElementById('f-source').value.trim(),
subject:document.getElementById('f-subject').value.trim(),
data:document.getElementById('f-data').value.trim(),
tags:document.getElementById('f-tags').value.trim()
};
try{
var r=await fetch(A+'/events',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r.ok){
var err=await r.json().catch(function(){return{}});
alert(err.error||'Failed to log event');
return;
}
}catch(e){
alert('Network error: '+e.message);
return;
}
closeModal();
load();
}

function closeModal(){
document.getElementById('mbg').classList.remove('open');
}

function esc(s){
if(s===undefined||s===null)return'';
var d=document.createElement('div');
d.textContent=String(s);
return d.innerHTML;
}

document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal()});

// ─── Live polling ─────────────────────────────────────────────────
// Auto-refresh every 5s when the toggle is on. Pauses when modal is open
// so the user doesn't lose form input mid-edit.

function startLive(){
if(liveTimer)return;
liveTimer=setInterval(function(){
if(!document.getElementById('mbg').classList.contains('open')){
load();
}
},5000);
}

function stopLive(){
if(liveTimer){clearInterval(liveTimer);liveTimer=null}
}

document.getElementById('live-cb').addEventListener('change',function(e){
liveOn=e.target.checked;
var dot=document.getElementById('live-dot');
if(liveOn){
dot.classList.add('on');
startLive();
}else{
dot.classList.remove('on');
stopLive();
}
});

// ─── Personalization ──────────────────────────────────────────────

(function loadPersonalization(){
fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
if(!cfg||typeof cfg!=='object')return;

if(cfg.dashboard_title){
var h1=document.getElementById('dash-title');
if(h1)h1.innerHTML='<span>&#9670;</span> '+esc(cfg.dashboard_title);
document.title=cfg.dashboard_title;
}

if(cfg.empty_state_message)window._emptyMsg=cfg.empty_state_message;

// Chronicle uses 'event_types' and 'sources' from the config to pre-populate
// the type filter bar and the emit form's type/source selects. This lets a
// brewery's chronicle suggest 'fermentation_started', 'kegging_complete', etc.
if(Array.isArray(cfg.event_types)){
presetTypes=cfg.event_types.filter(function(t){return typeof t==='string'&&t});
}
if(Array.isArray(cfg.sources)){
presetSources=cfg.sources.filter(function(s){return typeof s==='string'&&s});
}
}).catch(function(){
}).finally(function(){
load();
if(liveOn)startLive();
});
})();
</script>
</body>
</html>`
