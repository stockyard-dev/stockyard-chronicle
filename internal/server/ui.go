package server
import "net/http"
func(s *Server)dashboard(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","text/html");w.Write([]byte(dashHTML))}
const dashHTML=`<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Chronicle</title>
<style>:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--blue:#4a7ec9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}
.main{padding:1rem 1.5rem;max-width:1000px;margin:0 auto}
.filters{display:flex;gap:.3rem;margin-bottom:.8rem;flex-wrap:wrap}
.fbtn{font-size:.55rem;padding:.2rem .4rem;border:1px solid var(--bg3);background:var(--bg);color:var(--cm);cursor:pointer}.fbtn:hover{border-color:var(--leather)}.fbtn.active{border-color:var(--rust);color:var(--rust)}
.ev{display:flex;align-items:flex-start;gap:.6rem;padding:.4rem 0;border-bottom:1px solid var(--bg3);font-size:.72rem}
.ev-time{font-size:.6rem;color:var(--cm);width:70px;flex-shrink:0}
.ev-type{font-size:.55rem;padding:.1rem .3rem;min-width:60px;text-align:center;text-transform:uppercase}
.ev-type-info{background:#4a7ec922;color:var(--blue);border:1px solid #4a7ec944}
.ev-type-warning{background:#d4843a22;color:var(--orange);border:1px solid #d4843a44}
.ev-type-error{background:#c9444422;color:var(--red);border:1px solid #c9444444}
.ev-type-debug{background:var(--bg3);color:var(--cm)}
.ev-type-custom{background:#4a9e5c22;color:var(--green);border:1px solid #4a9e5c44}
.ev-body{flex:1;color:var(--cd)}
.ev-source{font-size:.6rem;color:var(--cm)}
.ev-data{font-size:.6rem;color:var(--cm);background:var(--bg);padding:.2rem .4rem;margin-top:.2rem;max-height:40px;overflow:hidden;cursor:pointer}
.ev-data:hover{max-height:none}
.stats{font-size:.65rem;color:var(--cm);margin-bottom:.8rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}
</style></head><body>
<div class="hdr"><h1>CHRONICLE</h1><div class="stats" id="stats"></div></div>
<div class="main">
<div class="filters" id="filters"></div>
<div id="events"></div>
</div>
<script>
const A='/api';let events=[],types=[],filterType='';
async function load(){const[e,s]=await Promise.all([fetch(A+'/events'+(filterType?'?type='+encodeURIComponent(filterType):'')).then(r=>r.json()),fetch(A+'/stats').then(r=>r.json())]);
events=e.events||[];
document.getElementById('stats').textContent=(s.total||0)+' events';
const tc=s.by_type||{};types=Object.keys(tc);
let fh='<button class="fbtn'+(filterType===''?' active':'')+'" onclick="setType(\'\')">All</button>';
types.forEach(t=>{fh+='<button class="fbtn'+(filterType===t?' active':'')+'" onclick="setType(\''+t+'\')">'+t+' ('+tc[t]+')</button>';});
document.getElementById('filters').innerHTML=fh;render();}
function setType(t){filterType=t;load();}
function render(){if(!events.length){document.getElementById('events').innerHTML='<div class="empty">No events logged. POST to /api/events to start.</div>';return;}
let h='';events.forEach(e=>{
const typeClass=(['info','warning','error','debug'].includes(e.type)?e.type:'custom');
h+='<div class="ev"><span class="ev-time">'+ft(e.created_at)+'</span><span class="ev-type ev-type-'+typeClass+'">'+esc(e.type||'event')+'</span><div class="ev-body">'+esc(e.message);
if(e.source)h+=' <span class="ev-source">['+esc(e.source)+']</span>';
if(e.data&&e.data!=='{}')h+='<div class="ev-data">'+esc(e.data)+'</div>';
h+='</div></div>';});
document.getElementById('events').innerHTML=h;}
function ft(t){if(!t)return'';const d=new Date(t);return d.toLocaleTimeString([],{hour:'2-digit',minute:'2-digit',second:'2-digit'});}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
load();setInterval(load,5000);
</script></body></html>`
